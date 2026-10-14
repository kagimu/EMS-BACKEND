package integration

import (
	"encoding/json"
	"net/http"
	"sort"
	"strconv"
	"testing"
	"time"

	assignmentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/assignments/models"
	authorityModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/authority/models"
	workflowModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/workflow/models"
	permissionsModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/permissions/models"
	uuid "github.com/satori/go.uuid"
)

// Foundation E3 (ratified D2 + D3) — read-side behaviour:
//   - stalled-request exposure, derived at READ TIME and never stored;
//   - stale / ineligible approval-task suppression.
//
// These tests use the Foundation D unit-narrowing fixture, which already seeds a
// unit hierarchy, placements with validity windows/statuses, and role holders.

const unDecMulti = "un_multi_request"

// unRoleL2 anchors the future level of the multi-level fixture and is held only by
// unErin, so it can be lapsed independently of the current level's candidates.
const unRoleL2 = "un-role-l2"

// createRequestWithAmount creates a request for the requester with an explicit amount.
func (h *unHarness) createRequestWithAmount(decisionType, unitNumber string, amount float64) eApprovalEnvelope {
	h.t.Helper()

	payload := map[string]any{
		"module": "finance", "resource_type": "purchase_order",
		"resource_number": uuid.NewV4().String(), "decision_type": decisionType,
		"amount": amount,
	}
	if unitNumber != "" {
		payload["organisation_unit_number"] = unitNumber
	}

	status, body := h.do(http.MethodPost, eApprovalsBase, h.token(unRequester), payload)
	if status != http.StatusCreated {
		h.t.Fatalf("create %s (%s) = %d (%s)", decisionType, unitNumber, status, string(body))
	}

	var envelope eApprovalEnvelope
	if err := json.Unmarshal(body, &envelope); err != nil {
		h.t.Fatalf("decoding created request: %v", err)
	}
	return envelope
}

// taskRequestNumbers returns the set of request numbers in a user's actionable inbox.
func (h *unHarness) taskRequestNumbers(userNumber string) map[string]bool {
	h.t.Helper()

	status, body := h.do(http.MethodGet, eTasksURL, h.token(userNumber), nil)
	if status != http.StatusOK {
		h.t.Fatalf("tasks for %s = %d (%s)", userNumber, status, string(body))
	}

	var tasks eTaskEnvelope
	if err := json.Unmarshal(body, &tasks); err != nil {
		h.t.Fatalf("decoding tasks: %v", err)
	}

	result := map[string]bool{}
	for _, task := range tasks.Data {
		result[task.RequestNumber] = true
	}
	return result
}

// assignmentUsersOf returns the sorted snapshotted candidate user numbers of a level.
func assignmentUsersOf(envelope eApprovalEnvelope, level int) []string {
	step := stepAtLevel(envelope, level)
	if step == nil {
		return nil
	}
	result := make([]string, 0, len(step.Assignments))
	for _, assignment := range step.Assignments {
		result = append(result, assignment.UserNumber)
	}
	sort.Strings(result)
	return result
}

// expirePlacement sets a user's placement end date into the past.
func (h *unHarness) expirePlacement(userNumber string) {
	h.t.Helper()
	past := time.Now().Add(-time.Hour)
	if err := h.db.Model(&assignmentModel.UserAssignment{}).
		Where("organisation_number = ? AND user_number = ?", unOrg, userNumber).
		Update("valid_to", past).Error; err != nil {
		h.t.Fatalf("expiring placement for %s: %v", userNumber, err)
	}
}

// restorePlacementWindow clears a user's placement validity window.
func (h *unHarness) restorePlacementWindow(userNumber string) {
	h.t.Helper()
	if err := h.db.Model(&assignmentModel.UserAssignment{}).
		Where("organisation_number = ? AND user_number = ?", unOrg, userNumber).
		Updates(map[string]any{"valid_from": nil, "valid_to": nil}).Error; err != nil {
		h.t.Fatalf("restoring placement window for %s: %v", userNumber, err)
	}
}

// ============================================================================
// Stalled-request exposure
// ============================================================================

func TestStalledExposure(t *testing.T) {
	h := newUnHarness(t)

	t.Run("not_stalled_when_a_candidate_is_currently_eligible", func(t *testing.T) {
		envelope := h.createRequest(unDecKampala, unKampala)
		assertCandidates(t, envelope, 1, unAlice, unCarol)

		if read := h.fetch(envelope.Data.Number); read.Data.Stalled {
			t.Fatal("a request with an eligible snapshotted candidate must not be stalled")
		}
	})

	t.Run("not_stalled_when_one_of_several_candidates_remains_eligible", func(t *testing.T) {
		envelope := h.createRequest(unDecKampala, unKampala)
		h.expirePlacement(unCarol)
		t.Cleanup(func() { h.restorePlacementWindow(unCarol) })

		read := h.fetch(envelope.Data.Number)
		if read.Data.Stalled {
			t.Fatal("one eligible candidate is enough: the request must not be stalled")
		}
	})

	t.Run("stalled_when_every_snapshotted_candidate_is_ineligible", func(t *testing.T) {
		envelope := h.createRequest(unDecKampala, unKampala)
		assertCandidates(t, envelope, 1, unAlice, unCarol)

		h.expirePlacement(unAlice)
		h.expirePlacement(unCarol)
		t.Cleanup(func() {
			h.restorePlacementWindow(unAlice)
			h.restorePlacementWindow(unCarol)
		})

		read := h.fetch(envelope.Data.Number)
		if !read.Data.Stalled {
			t.Fatalf("with no currently eligible snapshotted candidate the request must read as stalled: %+v", read.Data)
		}
		// The snapshot itself is untouched by the derivation.
		assertCandidates(t, read, 1, unAlice, unCarol)
		if read.Data.Status != "PENDING" {
			t.Fatalf("a stalled request stays PENDING: got %s", read.Data.Status)
		}
	})

	t.Run("stalled_is_derived_at_read_time_and_clears_when_eligibility_returns", func(t *testing.T) {
		envelope := h.createRequest(unDecKampala, unKampala)
		h.expirePlacement(unAlice)
		h.expirePlacement(unCarol)
		t.Cleanup(func() {
			h.restorePlacementWindow(unAlice)
			h.restorePlacementWindow(unCarol)
		})

		if !h.fetch(envelope.Data.Number).Data.Stalled {
			t.Fatal("expected stalled after both candidates lapsed")
		}

		// Nothing but the placement window changes, then the next read reports the
		// opposite: the condition is derived, not stored.
		h.restorePlacementWindow(unAlice)

		if h.fetch(envelope.Data.Number).Data.Stalled {
			t.Fatal("restoring one candidate's eligibility must clear the derived stalled state")
		}
	})

	t.Run("computing_stalled_never_writes", func(t *testing.T) {
		envelope := h.createRequest(unDecKampala, unKampala)
		step := stepAtLevel(envelope, 1)

		before, err := approvalSnapshotSignature(h, envelope.Data.Number, step.Number)
		if err != nil {
			t.Fatalf("reading snapshot: %v", err)
		}

		for i := 0; i < 3; i++ {
			h.fetch(envelope.Data.Number)
		}

		after, err := approvalSnapshotSignature(h, envelope.Data.Number, step.Number)
		if err != nil {
			t.Fatalf("reading snapshot: %v", err)
		}
		if before != after {
			t.Fatalf("the read-side derivation must not change any stored row:\nbefore %s\nafter  %s", before, after)
		}
	})

	t.Run("no_candidate_is_substituted_and_no_authority_re_resolution_occurs", func(t *testing.T) {
		envelope := h.createRequest(unDecKampala, unKampala)
		h.expirePlacement(unAlice)
		h.expirePlacement(unCarol)
		t.Cleanup(func() {
			h.restorePlacementWindow(unAlice)
			h.restorePlacementWindow(unCarol)
		})

		read := h.fetch(envelope.Data.Number)
		if !read.Data.Stalled {
			t.Fatal("expected a stalled request")
		}
		// Nobody is added and the frozen snapshot is unchanged.
		assertCandidates(t, read, 1, unAlice, unCarol)

		// Deactivating the authority rule must not re-resolve the frozen level, so the
		// stalled state is unchanged and the same candidate set is reported.
		if err := h.db.Model(&authorityModel.DecisionAuthority{}).
			Where("organisation_number = ? AND decision_type = ?", unOrg, unDecKampala).
			Update("status", "INACTIVE").Error; err != nil {
			t.Fatalf("deactivating authority: %v", err)
		}
		t.Cleanup(func() {
			h.db.Model(&authorityModel.DecisionAuthority{}).
				Where("organisation_number = ? AND decision_type = ?", unOrg, unDecKampala).
				Update("status", "ACTIVE")
		})

		after := h.fetch(envelope.Data.Number)
		if !after.Data.Stalled {
			t.Fatal("authority changes must not re-resolve a frozen request")
		}
		assertCandidates(t, after, 1, unAlice, unCarol)
	})

	t.Run("no_persisted_stalled_state_exists", func(t *testing.T) {
		var stallColumns int64
		// Match stalled-status column names precisely. The earlier substring
		// check ('%stall%') also matched legitimate columns such as
		// max_installations introduced by Foundation F Increment 1, so the
		// guard is narrowed to names that denote a stalled state.
		if err := h.db.Raw(
			"SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND (COLUMN_NAME LIKE 'stall%' OR COLUMN_NAME LIKE '%stall' OR COLUMN_NAME LIKE '%stalled%')",
		).Scan(&stallColumns).Error; err != nil {
			t.Fatalf("checking stalled columns: %v", err)
		}
		if stallColumns != 0 {
			t.Fatalf("E3 is read-time only: no stalled-status column may exist, found %d", stallColumns)
		}
	})
}

// approvalSnapshotSignature renders the stored rows a stalled read must never modify.
func approvalSnapshotSignature(h *unHarness, requestNumber, stepNumber string) (string, error) {
	request, err := approvalSnapshotRequest(h, requestNumber)
	if err != nil {
		return "", err
	}
	step, err := approvalSnapshotStep(h, stepNumber)
	if err != nil {
		return "", err
	}
	assignments, err := approvalSnapshotAssignments(h, stepNumber)
	if err != nil {
		return "", err
	}
	return request + "|" + step + "|" + assignments, nil
}

func approvalSnapshotRequest(h *unHarness, requestNumber string) (string, error) {
	var row workflowModel.ApprovalRequest
	if err := h.db.Where("number = ?", requestNumber).First(&row).Error; err != nil {
		return "", err
	}
	return row.Number + ":" + row.Status + ":" + strconv.Itoa(row.CurrentLevel) + ":" + row.UpdatedAt.String(), nil
}

func approvalSnapshotStep(h *unHarness, stepNumber string) (string, error) {
	var row workflowModel.ApprovalStep
	if err := h.db.Where("number = ?", stepNumber).First(&row).Error; err != nil {
		return "", err
	}
	return row.Number + ":" + row.Status + ":" + row.Anchors, nil
}

func approvalSnapshotAssignments(h *unHarness, stepNumber string) (string, error) {
	var rows []workflowModel.ApprovalAssignment
	if err := h.db.Where("step_number = ?", stepNumber).Order("user_number asc").Find(&rows).Error; err != nil {
		return "", err
	}
	signature := ""
	for _, row := range rows {
		signature += row.UserNumber + ":" + row.Status + ":" + row.Comment + ";"
	}
	return signature, nil
}

// ============================================================================
// Stale task suppression
// ============================================================================
// Stale task suppression
// ============================================================================

func TestStaleTaskSuppression(t *testing.T) {
	h := newUnHarness(t)

	t.Run("eligible_candidate_sees_the_task", func(t *testing.T) {
		envelope := h.createRequest(unDecKampala, unKampala)
		if !h.taskRequestNumbers(unAlice)[envelope.Data.Number] {
			t.Fatal("a currently eligible candidate must see the task")
		}
	})

	t.Run("placement_becoming_inactive_suppresses_the_task", func(t *testing.T) {
		envelope := h.createRequest(unDecKampala, unKampala)
		if err := h.db.Model(&assignmentModel.UserAssignment{}).
			Where("organisation_number = ? AND user_number = ?", unOrg, unAlice).
			Update("status", "INACTIVE").Error; err != nil {
			t.Fatalf("inactivating placement: %v", err)
		}
		t.Cleanup(func() {
			h.db.Model(&assignmentModel.UserAssignment{}).
				Where("organisation_number = ? AND user_number = ?", unOrg, unAlice).
				Update("status", "ACTIVE")
		})

		if h.taskRequestNumbers(unAlice)[envelope.Data.Number] {
			t.Fatal("an INACTIVE placement must suppress the actionable task")
		}
	})

	t.Run("expired_placement_suppresses_the_task", func(t *testing.T) {
		envelope := h.createRequest(unDecKampala, unKampala)
		h.expirePlacement(unAlice)
		t.Cleanup(func() { h.restorePlacementWindow(unAlice) })

		if h.taskRequestNumbers(unAlice)[envelope.Data.Number] {
			t.Fatal("an expired placement must suppress the actionable task")
		}
	})

	t.Run("future_placement_suppresses_the_task", func(t *testing.T) {
		envelope := h.createRequest(unDecKampala, unKampala)
		future := time.Now().Add(48 * time.Hour)
		if err := h.db.Model(&assignmentModel.UserAssignment{}).
			Where("organisation_number = ? AND user_number = ?", unOrg, unAlice).
			Update("valid_from", future).Error; err != nil {
			t.Fatalf("setting future valid_from: %v", err)
		}
		t.Cleanup(func() { h.restorePlacementWindow(unAlice) })

		if h.taskRequestNumbers(unAlice)[envelope.Data.Number] {
			t.Fatal("a not-yet-valid placement must suppress the actionable task")
		}
	})

	t.Run("ended_placement_suppresses_the_task", func(t *testing.T) {
		envelope := h.createRequest(unDecKampala, unKampala)
		if err := h.db.Model(&assignmentModel.UserAssignment{}).
			Where("organisation_number = ? AND user_number = ?", unOrg, unAlice).
			Update("status", "ENDED").Error; err != nil {
			t.Fatalf("ending placement: %v", err)
		}
		t.Cleanup(func() {
			h.db.Model(&assignmentModel.UserAssignment{}).
				Where("organisation_number = ? AND user_number = ?", unOrg, unAlice).
				Update("status", "ACTIVE")
		})

		if h.taskRequestNumbers(unAlice)[envelope.Data.Number] {
			t.Fatal("an ENDED placement must suppress the actionable task")
		}
	})

	t.Run("candidate_in_the_wrong_unit_has_no_task", func(t *testing.T) {
		envelope := h.createRequest(unDecKampala, unKampala)
		// Bob is a Finance Manager, but placed in Jinja: never a candidate here.
		if h.taskRequestNumbers(unBob)[envelope.Data.Number] {
			t.Fatal("a Jinja-placed approver must not receive a Kampala task")
		}
	})

	t.Run("another_eligible_snapshotted_candidate_still_sees_the_task", func(t *testing.T) {
		envelope := h.createRequest(unDecKampala, unKampala)
		h.expirePlacement(unAlice)

		if h.taskRequestNumbers(unAlice)[envelope.Data.Number] {
			t.Fatal("the lapsed candidate must not see the task")
		}
		if !h.taskRequestNumbers(unCarol)[envelope.Data.Number] {
			t.Fatal("the remaining eligible candidate must still see the task")
		}
	})

	t.Run("cross_organisation_inbox_never_contains_another_orgs_task", func(t *testing.T) {
		envelope := h.createRequest(unDecKampala, unKampala)
		if h.taskRequestNumbers(unOtherOrg)[envelope.Data.Number] {
			t.Fatal("organisation isolation must hold for the task inbox")
		}
	})

	t.Run("position_title_never_creates_a_task", func(t *testing.T) {
		envelope := h.createRequest(unDecKampala, unKampala)
		// unTitleOnly is placed at Kampala with the title "Finance Manager" but holds
		// no role, so no task may exist for them.
		if h.taskRequestNumbers(unTitleOnly)[envelope.Data.Number] {
			t.Fatal("a placement or job title must never create an approval task")
		}
	})
}

// ============================================================================
// Action security (the inbox filter is visibility only)
// ============================================================================

func TestStalledActionSecurity(t *testing.T) {
	h := newUnHarness(t)

	t.Run("stale_candidate_cannot_approve_through_the_api", func(t *testing.T) {
		envelope := h.createRequest(unDecKampala, unKampala)
		number := envelope.Data.Number
		h.expirePlacement(unAlice)
		t.Cleanup(func() { h.restorePlacementWindow(unAlice) })

		// The task is suppressed …
		if h.taskRequestNumbers(unAlice)[number] {
			t.Fatal("the stale task must not be listed")
		}
		// … but even calling the action endpoint directly is refused by Foundation D.
		if status, _ := h.decide(h.token(unAlice), number, "approve", "stale"); status != http.StatusForbidden {
			t.Fatalf("direct action by an ineligible snapshotted candidate = %d, want %d", status, http.StatusForbidden)
		}
	})

	t.Run("an_eligible_candidate_can_still_act", func(t *testing.T) {
		envelope := h.createRequest(unDecKampala, unKampala)
		number := envelope.Data.Number

		if status, body := h.decide(h.token(unAlice), number, "approve", "ok"); status != http.StatusOK {
			t.Fatalf("eligible candidate approve = %d (%s)", status, string(body))
		}
		if final := h.fetch(number); final.Data.Status != "APPROVED" {
			t.Fatalf("request status = %s, want APPROVED", final.Data.Status)
		}
	})

	t.Run("no_automatic_replacement_when_a_candidate_lapses", func(t *testing.T) {
		envelope := h.createRequest(unDecKampala, unKampala)
		h.expirePlacement(unAlice)

		read := h.fetch(envelope.Data.Number)
		assertCandidates(t, read, 1, unAlice, unCarol)
		for _, userNumber := range assignmentUsersOf(read, 1) {
			if userNumber != unAlice && userNumber != unCarol {
				t.Fatalf("nobody may be substituted into the snapshot, found %s", userNumber)
			}
		}
	})
}

// ============================================================================
// Legacy compatibility
// ============================================================================

func TestStalledLegacyCompatibility(t *testing.T) {
	h := newUnHarness(t)

	t.Run("legacy_request_reads_and_is_not_stalled_while_a_candidate_is_eligible", func(t *testing.T) {
		envelope := h.createRequest(unDecKampala, unKampala)
		step := stepAtLevel(envelope, 1)

		// Emulate a pre-amendment snapshot: no normalised anchor rows.
		if err := h.db.Unscoped().
			Where("step_number = ?", step.Number).
			Delete(&workflowModel.ApprovalStepAnchor{}).Error; err != nil {
			t.Fatalf("removing anchor rows: %v", err)
		}

		read := h.fetch(envelope.Data.Number)
		if keys := anchorKeysOf(read, 1); len(keys) != 1 || keys[0] != "role:"+unRoleFin+"@" {
			t.Fatalf("a legacy step must read as organisation-wide: %v", keys)
		}
		if read.Data.Stalled {
			t.Fatal("with a legacy organisation-wide anchor and an eligible candidate the request is not stalled")
		}
		// The snapshot is not retroactively narrowed.
		assertCandidates(t, read, 1, unAlice, unCarol)
	})

	t.Run("legacy_task_behaviour_remains_compatible", func(t *testing.T) {
		envelope := h.createRequest(unDecKampala, unKampala)
		step := stepAtLevel(envelope, 1)
		if err := h.db.Unscoped().
			Where("step_number = ?", step.Number).
			Delete(&workflowModel.ApprovalStepAnchor{}).Error; err != nil {
			t.Fatalf("removing anchor rows: %v", err)
		}

		// A snapshotted candidate that is a role holder remains eligible under the
		// legacy organisation-wide anchor, so the task remains visible.
		if !h.taskRequestNumbers(unAlice)[envelope.Data.Number] {
			t.Fatal("a legacy request's task must remain visible to an eligible candidate")
		}
		read := h.fetch(envelope.Data.Number)
		assertCandidates(t, read, 1, unAlice, unCarol)
	})
}

// ============================================================================
// Multi-level approvals
// ============================================================================

func TestStalledMultiLevel(t *testing.T) {
	h := newUnHarness(t)

	// A second role, held only by Erin, anchors the FUTURE level. It is
	// organisation-wide, because a rule scoped to Jinja would not be applicable to a
	// Kampala request at all (Foundation D: the rule's unit must be the request's
	// unit or an ancestor of it).
	levelTwoRole := permissionsModel.Role{
		Number: unRoleL2, Name: "Level Two Approver", EntityNumber: unOrg,
		Status: "ACTIVE", UserNumber: unRequester,
	}
	if err := h.db.Omit("User", "Entity", "Permissions").Create(&levelTwoRole).Error; err != nil {
		t.Fatalf("seeding level-two role: %v", err)
	}
	levelTwoGrant := permissionsModel.UserRole{
		UserNumber: unErin, RoleNumber: unRoleL2,
		ExpiryDate: time.Now().Add(365 * 24 * time.Hour), AddedBy: unRequester,
	}
	if err := h.db.Omit("User", "Role", "AddedByUser").Create(&levelTwoGrant).Error; err != nil {
		t.Fatalf("seeding level-two role grant: %v", err)
	}

	// Level 1 does not cover 5000, so the path is [1, 2].
	amountCap := 1000.0
	authorities := []authorityModel.DecisionAuthority{
		{Number: uuid.NewV4().String(), OrganisationNumber: unOrg, DecisionType: unDecMulti,
			ApprovalLevel: 1, MaxAmount: &amountCap, RoleNumber: unRoleFin,
			OrganisationUnitNumber: unKampala, Status: "ACTIVE", AddedByNumber: unRequester},
		{Number: uuid.NewV4().String(), OrganisationNumber: unOrg, DecisionType: unDecMulti,
			ApprovalLevel: 2, RoleNumber: unRoleL2, Status: "ACTIVE", AddedByNumber: unRequester},
	}
	for _, authority := range authorities {
		if err := h.db.Create(&authority).Error; err != nil {
			t.Fatalf("seeding multi-level authority: %v", err)
		}
	}

	t.Run("a_future_level_with_no_eligible_candidate_does_not_stall_the_request", func(t *testing.T) {
		envelope := h.createRequestWithAmount(unDecMulti, unKampala, 5000)
		if len(envelope.Data.Steps) != 2 {
			t.Fatalf("expected a two-level path, got %d", len(envelope.Data.Steps))
		}
		assertCandidates(t, envelope, 1, unAlice, unCarol)
		assertCandidates(t, envelope, 2, unErin)

		// The FUTURE level's only candidate lapses (their role grant expires); the
		// current level is unaffected.
		if err := h.db.Model(&permissionsModel.UserRole{}).
			Where("user_number = ? AND role_number = ?", unErin, unRoleL2).
			Update("expiry_date", time.Now().Add(-time.Hour)).Error; err != nil {
			t.Fatalf("lapsing level-two role grant: %v", err)
		}
		t.Cleanup(func() {
			h.db.Model(&permissionsModel.UserRole{}).
				Where("user_number = ? AND role_number = ?", unErin, unRoleL2).
				Update("expiry_date", time.Now().Add(365*24*time.Hour))
		})

		read := h.fetch(envelope.Data.Number)
		if read.Data.CurrentLevel != 1 {
			t.Fatalf("current level = %d, want 1", read.Data.CurrentLevel)
		}
		if read.Data.Stalled {
			t.Fatal("a lapsed FUTURE level must not mark the request stalled while the current level is actionable")
		}
	})

	t.Run("current_level_with_no_eligible_candidate_is_stalled", func(t *testing.T) {
		envelope := h.createRequestWithAmount(unDecMulti, unKampala, 5000)
		h.expirePlacement(unAlice)
		h.expirePlacement(unCarol)
		t.Cleanup(func() {
			h.restorePlacementWindow(unAlice)
			h.restorePlacementWindow(unCarol)
		})

		read := h.fetch(envelope.Data.Number)
		if !read.Data.Stalled {
			t.Fatal("the current level has no eligible snapshotted candidate: the request must read as stalled")
		}
		assertCandidates(t, read, 1, unAlice, unCarol)
	})

	t.Run("completed_level_does_not_stall_and_progression_is_unchanged", func(t *testing.T) {
		envelope := h.createRequestWithAmount(unDecMulti, unKampala, 5000)
		number := envelope.Data.Number

		if status, body := h.decide(h.token(unAlice), number, "approve", "level one"); status != http.StatusOK {
			t.Fatalf("level 1 approve = %d (%s)", status, string(body))
		}

		read := h.fetch(number)
		if read.Data.CurrentLevel != 2 || read.Data.Status != "PENDING" {
			t.Fatalf("sequencing changed: level %d status %s", read.Data.CurrentLevel, read.Data.Status)
		}
		// Level 1 is completed (its candidate lapsed or not, it is not the current
		// level) and level 2 has an eligible candidate, so the request is not stalled.
		level1 := stepAtLevel(read, 1)
		if level1.Status != "APPROVED" {
			t.Fatalf("completed level status = %s", level1.Status)
		}
		if read.Data.Stalled {
			t.Fatal("a completed level must never make the request read as stalled")
		}
	})
}
