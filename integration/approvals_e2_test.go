package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/routes"
	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	authServices "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/services"
	authorityModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/authority/models"
	stakeholderModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/stakeholders/models"
	approvalRepositories "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/workflow/repositories"
	eosModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/eos/models"
	uuid "github.com/satori/go.uuid"
)

// E2 test users. They are seeded by the E2 harness and removed by its cleanup,
// independently of the E1 fixture's user list.
const (
	e2BoardActive  = "ems-e2-board-active"
	e2BoardExpired = "ems-e2-board-expired"
	e2BoardFuture  = "ems-e2-board-future"

	// e2StakeholderDecision is a stakeholder-anchored decision type used to prove
	// the ratified D10 validity-window rule (snapshot and action time).
	e2StakeholderDecision = "e2_stakeholder_request"
)

// e2Harness wraps the E1 fixture with HTTP helpers for the E2 lifecycle tests.
type e2Harness struct {
	fixture *eFixture
	server  *httptest.Server
	client  *http.Client
	t       *testing.T
}

func newE2Harness(t *testing.T) *e2Harness {
	t.Helper()

	fixture := setupEFixture(t)
	server := httptest.NewServer(routes.SetupRoutes())
	t.Cleanup(server.Close)

	// E2-owned users, seeded after the E1 fixture so the LIFO cleanup below runs
	// before the fixture's own cleanup.
	users := []string{e2BoardActive, e2BoardExpired, e2BoardFuture}
	db := fixture.db

	cleanup := func() {
		db.Unscoped().Where("user_number IN ?", users).Delete(&authModel.Contact{})
		db.Unscoped().Where("user_number IN ?", users).Delete(&authModel.UserPassword{})
		db.Unscoped().Where("user_number IN ?", users).Delete(&eosModel.EntityUser{})
		db.Unscoped().Where("user_number IN ?", users).Delete(&eosModel.ActiveUserEntity{})
		db.Unscoped().Where("number IN ?", users).Delete(&authModel.User{})
	}
	cleanup()
	t.Cleanup(cleanup)

	for _, number := range users {
		user := authModel.User{
			Number: number, FirstName: "E2", SecondName: "Board",
			Gender: "M", Type: "Admin", Status: "Active",
		}
		if err := db.Create(&user).Error; err != nil {
			t.Fatalf("seeding E2 user %s: %v", number, err)
		}

		active := eosModel.ActiveUserEntity{UserNumber: number, EntityNumber: eOrgA}
		if err := db.Omit("Entity", "User").Create(&active).Error; err != nil {
			t.Fatalf("setting active entity for %s: %v", number, err)
		}

		// Foundation C2 eligibility: a candidate must be able to sign in.
		contact := authModel.Contact{
			Contact: number + "@example.com", ContactType: "Email", UserNumber: number,
		}
		if err := db.Omit("User").Create(&contact).Error; err != nil {
			t.Fatalf("seeding contact for %s: %v", number, err)
		}
		password := authModel.UserPassword{UserNumber: number, UserPassword: "seeded-hash"}
		if err := db.Omit("User").Create(&password).Error; err != nil {
			t.Fatalf("seeding password for %s: %v", number, err)
		}
	}

	// A stakeholder-anchored decision type at level 1, so the D10 window rule can
	// be observed end to end.
	authority := authorityModel.DecisionAuthority{
		Number: uuid.NewV4().String(), OrganisationNumber: eOrgA,
		DecisionType: e2StakeholderDecision, ApprovalLevel: 1,
		StakeholderType: stakeholderModel.StakeholderTypeBoardMember,
		Status:          "ACTIVE", AddedByNumber: eAdminA,
	}
	if err := db.Create(&authority).Error; err != nil {
		t.Fatalf("seeding E2 stakeholder authority: %v", err)
	}

	// Windows relative to a fixed reference so the semantics are unambiguous.
	reference := time.Now()
	window := func(from, to *time.Time) stakeholderModel.Stakeholder {
		return stakeholderModel.Stakeholder{
			Number: uuid.NewV4().String(), OrganisationNumber: eOrgA,
			StakeholderType: stakeholderModel.StakeholderTypeBoardMember,
			Status:          "ACTIVE", ValidFrom: from, ValidTo: to, AddedByNumber: eAdminA,
		}
	}

	// e2-board-active: bounded but currently inside the window.
	past := reference.Add(-24 * time.Hour)
	future := reference.Add(24 * time.Hour)
	for _, spec := range []struct {
		userNumber string
		from, to   *time.Time
	}{
		{e2BoardActive, &past, &future},
		// e2-board-expired: ACTIVE but the window already closed.
		{e2BoardExpired, &past, &past},
		// e2-board-future: ACTIVE but the window has not opened.
		{e2BoardFuture, &future, &future},
	} {
		stakeholder := window(spec.from, spec.to)
		stakeholder.UserNumber = spec.userNumber
		if err := db.Create(&stakeholder).Error; err != nil {
			t.Fatalf("seeding stakeholder for %s: %v", spec.userNumber, err)
		}
	}

	h := &e2Harness{fixture: fixture, server: server, client: server.Client(), t: t}
	return h
}

// request issues one HTTP call. It is safe to call from a goroutine when token
// and payload are read-only.
func (h *e2Harness) request(method, path, token string, payload any) (int, []byte) {
	h.t.Helper()

	var reader io.Reader
	if payload != nil {
		raw, err := json.Marshal(payload)
		if err != nil {
			h.t.Fatalf("encoding payload: %v", err)
		}
		reader = bytes.NewReader(raw)
	}

	req, err := http.NewRequest(method, h.server.URL+path, reader)
	if err != nil {
		h.t.Fatalf("building request: %v", err)
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := h.client.Do(req)
	if err != nil {
		h.t.Fatalf("request %s %s failed: %v", method, path, err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	return resp.StatusCode, body
}

// create issues a create with an optional previous_request_number.
func (h *e2Harness) create(token, decisionType, resourceNumber string, amount *float64, previous string) (int, eApprovalEnvelope) {
	h.t.Helper()

	payload := map[string]any{
		"module": "finance", "resource_type": "purchase_order",
		"resource_number": resourceNumber, "decision_type": decisionType,
	}
	if amount != nil {
		payload["amount"] = *amount
	}
	if previous != "" {
		payload["previous_request_number"] = previous
	}

	status, body := h.request(http.MethodPost, eApprovalsBase, token, payload)

	var envelope eApprovalEnvelope
	if status == http.StatusCreated {
		if err := json.Unmarshal(body, &envelope); err != nil {
			h.t.Fatalf("decoding created request: %v (%s)", err, string(body))
		}
	}
	return status, envelope
}

// decide drives approve / reject / return / cancel.
func (h *e2Harness) decide(token, number, action, comment string) (int, eApprovalEnvelope, []byte) {
	h.t.Helper()

	status, body := h.request(http.MethodPost, eApprovalsBase+"/"+number+"/"+action, token,
		map[string]any{"comment": comment})

	var envelope eApprovalEnvelope
	if status == http.StatusOK {
		if err := json.Unmarshal(body, &envelope); err != nil {
			h.t.Fatalf("decoding %s response: %v (%s)", action, err, string(body))
		}
	}
	return status, envelope, body
}

// fetch reads a request back through the API.
func (h *e2Harness) fetch(token, number string) eApprovalEnvelope {
	h.t.Helper()

	status, body := h.request(http.MethodGet, eApprovalsBase+"/"+number, token, nil)
	if status != http.StatusOK {
		h.t.Fatalf("fetch %s = %d (%s)", number, status, string(body))
	}
	var envelope eApprovalEnvelope
	if err := json.Unmarshal(body, &envelope); err != nil {
		h.t.Fatalf("decoding fetched request: %v", err)
	}
	return envelope
}

func (h *e2Harness) amount(v float64) *float64 { return &v }

// shape renders the parts of a request that must never change once a decision
// has been recorded, so immutability can be asserted by comparison.
func shape(envelope eApprovalEnvelope) string {
	var builder strings.Builder
	fmt.Fprintf(&builder, "%s|%s|%d|%s", envelope.Data.Number, envelope.Data.Status,
		envelope.Data.CurrentLevel, envelope.Data.PreviousRequestNumber)
	for _, step := range envelope.Data.Steps {
		fmt.Fprintf(&builder, "||%d:%s:%v", step.ApprovalLevel, step.Status, step.Actionable)
		for _, assignment := range step.Assignments {
			fmt.Fprintf(&builder, "|%s:%s:%s", assignment.UserNumber, assignment.Status, assignment.Comment)
		}
	}
	return builder.String()
}

// ============================================================================
// Return-for-correction (D1, D2, D3, D4)
// ============================================================================

func TestWorkflowE2Return(t *testing.T) {
	h := newE2Harness(t)
	fx := h.fixture

	t.Run("eligible_approver_returns_and_requests_becomes_inert", func(t *testing.T) {
		status, created := h.create(fx.tokenAdminA, "purchase_request", "e2-return-1", h.amount(100000), "")
		if status != http.StatusCreated {
			t.Fatalf("create = %d", status)
		}
		number := created.Data.Number

		status, returned, body := h.decide(fx.tokenAppr1, number, "return", "please attach the quote")
		if status != http.StatusOK {
			t.Fatalf("return = %d (%s)", status, string(body))
		}
		if returned.Data.Status != "RETURNED" {
			t.Fatalf("request status = %s, want RETURNED", returned.Data.Status)
		}

		step := stepAtLevel(returned, 1)
		if step == nil || step.Status != "RETURNED" || step.Actionable {
			t.Fatalf("step must be RETURNED and not actionable: %+v", step)
		}

		// The returning approver's own row carries the immutable decision.
		var actor, sibling *eAssignment
		for i := range step.Assignments {
			switch step.Assignments[i].UserNumber {
			case eApprover1:
				actor = &step.Assignments[i]
			case eApprover2:
				sibling = &step.Assignments[i]
			}
		}
		if actor == nil || actor.Status != "RETURNED" || actor.Comment != "please attach the quote" || actor.DecidedAt == nil {
			t.Fatalf("returner's assignment not recorded correctly: %+v", actor)
		}
		if sibling == nil || sibling.Status != "UNDECIDED" || sibling.DecidedAt != nil {
			t.Fatalf("a sibling candidate must stay undecided: %+v", sibling)
		}

		// Only statuses moved: the request identity, the snapshotted path, the
		// anchors and the candidate set are all unchanged, so nothing was rebuilt.
		if returned.Data.Number != number {
			t.Fatal("return must not change the request identity")
		}
		if len(returned.Data.Steps) != len(created.Data.Steps) {
			t.Fatalf("the snapshotted path must not be rebuilt (%d steps, want %d)",
				len(returned.Data.Steps), len(created.Data.Steps))
		}
		for _, originalStep := range created.Data.Steps {
			movedStep := stepAtLevel(returned, originalStep.ApprovalLevel)
			if movedStep == nil {
				t.Fatalf("level %d disappeared after return", originalStep.ApprovalLevel)
			}
			if len(movedStep.Anchors) != len(originalStep.Anchors) {
				t.Fatalf("snapshotted anchors changed on level %d", originalStep.ApprovalLevel)
			}
			if len(movedStep.Assignments) != len(originalStep.Assignments) {
				t.Fatalf("candidate set changed on level %d", originalStep.ApprovalLevel)
			}
			for i := range originalStep.Assignments {
				if movedStep.Assignments[i].UserNumber != originalStep.Assignments[i].UserNumber {
					t.Fatalf("candidate identity changed on level %d", originalStep.ApprovalLevel)
				}
			}
		}

		// A RETURNED request is terminal: no further decision of any kind.
		for _, action := range []string{"approve", "reject", "return"} {
			if status, _, _ := h.decide(fx.tokenAppr1, number, action, "x"); status != http.StatusConflict {
				t.Fatalf("%s on a RETURNED request = %d, want %d", action, status, http.StatusConflict)
			}
		}
		// The requester cannot cancel it either: RETURNED is not a cancellable
		// state (only PENDING/BLOCKED are).
		if status, _, _ := h.decide(fx.tokenAdminA, number, "cancel", ""); status != http.StatusConflict {
			t.Fatalf("requester cancel of a RETURNED request = %d, want %d", status, http.StatusConflict)
		}
		// A non-requester is refused by the requester-only cancellation rule.
		if status, _, _ := h.decide(fx.tokenAppr1, number, "cancel", ""); status != http.StatusForbidden {
			t.Fatalf("non-requester cancel = %d, want %d", status, http.StatusForbidden)
		}
	})

	t.Run("later_steps_are_never_actionable_after_return", func(t *testing.T) {
		status, created := h.create(fx.tokenAdminA, "purchase_request", "e2-return-2", h.amount(600000), "")
		if status != http.StatusCreated {
			t.Fatalf("create = %d", status)
		}
		number := created.Data.Number
		if len(created.Data.Steps) != 2 {
			t.Fatalf("expected a two-level path, got %d", len(created.Data.Steps))
		}

		// Level 1 is approved first, then level 2 is returned: the earlier, already
		// recorded decision must remain untouched.
		if status, _, _ := h.decide(fx.tokenAppr1, number, "approve", "level one ok"); status != http.StatusOK {
			t.Fatalf("level 1 approve = %d", status)
		}
		status, returned, _ := h.decide(fx.tokenOwnerA, number, "return", "needs owner revision")
		if status != http.StatusOK {
			t.Fatalf("level 2 return = %d", status)
		}

		level1 := stepAtLevel(returned, 1)
		if level1.Status != "APPROVED" {
			t.Fatalf("a completed prior step must stay APPROVED, got %s", level1.Status)
		}
		decided := false
		for _, assignment := range level1.Assignments {
			if assignment.UserNumber == eApprover1 && assignment.Status == "APPROVED" && assignment.Comment == "level one ok" {
				decided = true
			}
		}
		if !decided {
			t.Fatalf("the prior recorded approval must remain immutable: %+v", level1.Assignments)
		}
		if level2 := stepAtLevel(returned, 2); level2.Status != "RETURNED" || level2.Actionable {
			t.Fatalf("the returned level must be inert: %+v", level2)
		}
		if returned.Data.Status != "RETURNED" {
			t.Fatalf("request must be RETURNED, got %s", returned.Data.Status)
		}
	})

	t.Run("reason_is_mandatory", func(t *testing.T) {
		status, created := h.create(fx.tokenAdminA, "purchase_request", "e2-return-3", h.amount(100000), "")
		if status != http.StatusCreated {
			t.Fatalf("create = %d", status)
		}
		if status, _, _ := h.decide(fx.tokenAppr1, created.Data.Number, "return", "   "); status != http.StatusBadRequest {
			t.Fatalf("blank reason = %d, want %d", status, http.StatusBadRequest)
		}
	})

	t.Run("requester_and_unassigned_users_cannot_return", func(t *testing.T) {
		status, created := h.create(fx.tokenAdminA, "purchase_request", "e2-return-4", h.amount(100000), "")
		if status != http.StatusCreated {
			t.Fatalf("create = %d", status)
		}
		// The requester is not an assignment of the step.
		if status, _, _ := h.decide(fx.tokenAdminA, created.Data.Number, "return", "mine"); status != http.StatusForbidden {
			t.Fatalf("requester return = %d, want %d", status, http.StatusForbidden)
		}
		// An admin with the permission but no assignment is likewise refused.
		if status, _, _ := h.decide(fx.tokenAdminB, created.Data.Number, "return", "x"); status != http.StatusNotFound {
			t.Fatalf("cross-organisation return = %d, want %d", status, http.StatusNotFound)
		}
	})

	t.Run("cross_organisation_return_is_not_found_and_changes_nothing", func(t *testing.T) {
		status, created := h.create(fx.tokenAdminA, "purchase_request", "e2-return-5", h.amount(100000), "")
		if status != http.StatusCreated {
			t.Fatalf("create = %d", status)
		}
		before := shape(h.fetch(fx.tokenAdminA, created.Data.Number))

		if status, _, _ := h.decide(fx.tokenAdminB, created.Data.Number, "return", "x"); status != http.StatusNotFound {
			t.Fatalf("cross-org return = %d, want %d", status, http.StatusNotFound)
		}
		if after := shape(h.fetch(fx.tokenAdminA, created.Data.Number)); after != before {
			t.Fatalf("a cross-organisation attempt must not change the request")
		}
	})
}

// ============================================================================
// D3 comment limit: 1,000 Unicode characters, uniformly, counted as runes
// ============================================================================

func TestWorkflowE2CommentLimit(t *testing.T) {
	h := newE2Harness(t)
	fx := h.fixture

	ascii := func(n int) string { return strings.Repeat("a", n) }
	// Each 'é' is a single rune but two UTF-8 bytes, so a 1,000-rune string is
	// 2,000 bytes: accepting it proves rune counting rather than len().
	multibyte := func(n int) string { return strings.Repeat("é", n) }

	t.Run("exactly_1000_unicode_characters_accepted_1001_rejected", func(t *testing.T) {
		status, created := h.create(fx.tokenAdminA, "purchase_request", "e2-limit-1", h.amount(100000), "")
		if status != http.StatusCreated {
			t.Fatalf("create = %d", status)
		}
		if status, _, _ := h.decide(fx.tokenAppr1, created.Data.Number, "return", ascii(1001)); status != http.StatusBadRequest {
			t.Fatalf("1001 chars = %d, want %d", status, http.StatusBadRequest)
		}

		status, returned, body := h.decide(fx.tokenAppr1, created.Data.Number, "return", ascii(1000))
		if status != http.StatusOK {
			t.Fatalf("1000 chars = %d (%s)", status, string(body))
		}
		if got := len(stepAtLevel(returned, 1).Assignments); got == 0 {
			t.Fatal("expected assignments")
		}
	})

	t.Run("1000_multibyte_runes_accepted_and_stored_intact", func(t *testing.T) {
		status, created := h.create(fx.tokenAdminA, "purchase_request", "e2-limit-2", h.amount(100000), "")
		if status != http.StatusCreated {
			t.Fatalf("create = %d", status)
		}

		reason := multibyte(1000)
		if len(reason) <= 1000 {
			t.Fatalf("test fixture must exceed 1000 bytes to be meaningful (got %d)", len(reason))
		}

		status, returned, body := h.decide(fx.tokenAppr1, created.Data.Number, "return", reason)
		if status != http.StatusOK {
			t.Fatalf("1000 multibyte runes = %d (%s)", status, string(body))
		}

		// Round-trips through MySQL varchar(1000) unchanged.
		fetched := h.fetch(fx.tokenAdminA, returned.Data.Number)
		var stored string
		for _, assignment := range stepAtLevel(fetched, 1).Assignments {
			if assignment.UserNumber == eApprover1 {
				stored = assignment.Comment
			}
		}
		if stored != reason {
			t.Fatalf("stored reason mismatch: got %d runes, want %d", len([]rune(stored)), len([]rune(reason)))
		}

		status, created2 := h.create(fx.tokenAdminA, "purchase_request", "e2-limit-3", h.amount(100000), "")
		if status != http.StatusCreated {
			t.Fatalf("create = %d", status)
		}
		if status, _, _ := h.decide(fx.tokenAppr2, created2.Data.Number, "return", multibyte(1001)); status != http.StatusBadRequest {
			t.Fatalf("1001 multibyte runes = %d, want %d", status, http.StatusBadRequest)
		}
	})

	t.Run("limit_applies_to_approval_comments_and_rejection_reasons", func(t *testing.T) {
		status, created := h.create(fx.tokenAdminA, "purchase_request", "e2-limit-4", h.amount(100000), "")
		if status != http.StatusCreated {
			t.Fatalf("create = %d", status)
		}
		if status, _, _ := h.decide(fx.tokenAppr1, created.Data.Number, "approve", multibyte(1001)); status != http.StatusBadRequest {
			t.Fatalf("approval comment over limit = %d, want %d", status, http.StatusBadRequest)
		}
		if status, _, _ := h.decide(fx.tokenAppr1, created.Data.Number, "approve", multibyte(1000)); status != http.StatusOK {
			t.Fatalf("approval comment at limit = %d", status)
		}

		status, created2 := h.create(fx.tokenAdminA, "purchase_request", "e2-limit-5", h.amount(100000), "")
		if status != http.StatusCreated {
			t.Fatalf("create = %d", status)
		}
		if status, _, _ := h.decide(fx.tokenAppr1, created2.Data.Number, "reject", multibyte(1001)); status != http.StatusBadRequest {
			t.Fatalf("rejection reason over limit = %d, want %d", status, http.StatusBadRequest)
		}
		if status, _, _ := h.decide(fx.tokenAppr1, created2.Data.Number, "reject", multibyte(1000)); status != http.StatusOK {
			t.Fatalf("rejection reason at limit = %d", status)
		}
	})
}

// ============================================================================
// Return concurrency: exactly one winner, exactly one decision
// ============================================================================

func TestWorkflowE2ReturnConcurrency(t *testing.T) {
	h := newE2Harness(t)
	fx := h.fixture

	race := func(number string, a, b eDecisionSpec) (int, int) {
		t.Helper()

		start := make(chan struct{})
		statuses := make([]int, 2)
		var wait sync.WaitGroup

		fire := func(index int, spec eDecisionSpec) {
			defer wait.Done()
			<-start
			status, _ := h.request(http.MethodPost, eApprovalsBase+"/"+number+"/"+spec.action, spec.token,
				map[string]any{"comment": spec.comment})
			statuses[index] = status
		}

		wait.Add(2)
		go fire(0, a)
		go fire(1, b)
		close(start)
		wait.Wait()
		return statuses[0], statuses[1]
	}

	for attempt := 0; attempt < 3; attempt++ {
		status, created := h.create(fx.tokenAdminA, "purchase_request", fmt.Sprintf("e2-race-%d", attempt), h.amount(100000), "")
		if status != http.StatusCreated {
			t.Fatalf("create = %d", status)
		}
		number := created.Data.Number

		approveStatus, returnStatus := race(number,
			eDecisionSpec{token: fx.tokenAppr1, action: "approve", comment: "yes"},
			eDecisionSpec{token: fx.tokenAppr2, action: "return", comment: "needs work"})

		successes := 0
		for _, status := range []int{approveStatus, returnStatus} {
			if status == http.StatusOK {
				successes++
			} else if status != http.StatusConflict && status != http.StatusForbidden {
				t.Fatalf("attempt %d: loser status = %d, want 409 or 403", attempt, status)
			}
		}
		if successes != 1 {
			t.Fatalf("attempt %d: %d decisions succeeded, want exactly 1", attempt, successes)
		}

		final := h.fetch(fx.tokenAdminA, number)
		t.Logf("return race attempt %d: approve=%d return=%d final=%s", attempt, approveStatus, returnStatus, final.Data.Status)

		decided := 0
		for _, step := range final.Data.Steps {
			for _, assignment := range step.Assignments {
				if assignment.Status == "APPROVED" || assignment.Status == "RETURNED" {
					decided++
				}
			}
		}
		if decided != 1 {
			t.Fatalf("attempt %d: %d decisions recorded, want exactly 1", attempt, decided)
		}

		if approveStatus == http.StatusOK && final.Data.Status != "APPROVED" {
			t.Fatalf("attempt %d: approve won but status = %s", attempt, final.Data.Status)
		}
		if returnStatus == http.StatusOK && final.Data.Status != "RETURNED" {
			t.Fatalf("attempt %d: return won but status = %s", attempt, final.Data.Status)
		}
	}
}

// ============================================================================
// Resubmission and previous-request linkage (D5, D6, D7, D8)
// ============================================================================

func TestWorkflowE2Resubmission(t *testing.T) {
	h := newE2Harness(t)
	fx := h.fixture

	// returnedRequest creates a request and returns it, yielding a linkable
	// predecessor plus its immutable shape.
	returnedRequest := func(t *testing.T, resource string) (eApprovalEnvelope, string) {
		t.Helper()
		status, created := h.create(fx.tokenAdminA, "purchase_request", resource, h.amount(100000), "")
		if status != http.StatusCreated {
			t.Fatalf("create = %d", status)
		}
		status, returned, body := h.decide(fx.tokenAppr1, created.Data.Number, "return", "fix the totals")
		if status != http.StatusOK {
			t.Fatalf("return = %d (%s)", status, string(body))
		}
		return returned, shape(returned)
	}

	t.Run("requester_resubmits_as_a_new_linked_request", func(t *testing.T) {
		predecessor, predecessorShape := returnedRequest(t, "e2-resub-1")

		status, created := h.create(fx.tokenAdminA, "purchase_request", "e2-resub-1", h.amount(100000), predecessor.Data.Number)
		if status != http.StatusCreated {
			t.Fatalf("resubmit = %d", status)
		}
		if created.Data.Number == predecessor.Data.Number {
			t.Fatal("resubmission must create a NEW request")
		}
		if created.Data.PreviousRequestNumber != predecessor.Data.Number {
			t.Fatalf("previous_request_number = %q, want %q", created.Data.PreviousRequestNumber, predecessor.Data.Number)
		}
		if created.Data.Status != "PENDING" {
			t.Fatalf("a new request must be PENDING/BLOCKED, got %s", created.Data.Status)
		}

		// The new request has its own fresh path: all assignments undecided.
		for _, step := range created.Data.Steps {
			for _, assignment := range step.Assignments {
				if assignment.Status != "UNDECIDED" {
					t.Fatalf("a fresh snapshot must not reuse prior decisions: %+v", assignment)
				}
			}
		}

		// The predecessor is untouched.
		if after := shape(h.fetch(fx.tokenAdminA, predecessor.Data.Number)); after != predecessorShape {
			t.Fatalf("the previous request must remain intact:\nbefore %s\nafter  %s", predecessorShape, after)
		}
	})

	t.Run("fresh_authority_resolution_and_snapshot", func(t *testing.T) {
		predecessor, _ := returnedRequest(t, "e2-resub-2")
		if anchors := stepAtLevel(predecessor, 1).Anchors; len(anchors) != 1 {
			t.Fatalf("expected one anchor on the predecessor, got %d", len(anchors))
		}

		// Authority changes while the request is returned. The resubmission must
		// re-resolve D and snapshot the new configuration; the predecessor must not.
		extra := authorityModel.DecisionAuthority{
			Number: uuid.NewV4().String(), OrganisationNumber: eOrgA,
			DecisionType: "purchase_request", ApprovalLevel: 1,
			MaxAmount: h.amount(500000), RoleNumber: eRoleSolo,
			Status: "ACTIVE", AddedByNumber: eAdminA,
		}
		if err := fx.db.Create(&extra).Error; err != nil {
			t.Fatalf("adding authority: %v", err)
		}
		t.Cleanup(func() {
			fx.db.Unscoped().Where("number = ?", extra.Number).Delete(&authorityModel.DecisionAuthority{})
		})

		status, created := h.create(fx.tokenAdminA, "purchase_request", "e2-resub-2", h.amount(100000), predecessor.Data.Number)
		if status != http.StatusCreated {
			t.Fatalf("resubmit = %d", status)
		}

		if got := len(stepAtLevel(created, 1).Anchors); got != 2 {
			t.Fatalf("the resubmission must re-resolve authority (expected 2 anchors, got %d)", got)
		}
		if got := len(stepAtLevel(h.fetch(fx.tokenAdminA, predecessor.Data.Number), 1).Anchors); got != 1 {
			t.Fatalf("the predecessor's snapshot must not change (got %d anchors)", got)
		}
	})

	t.Run("linkage_rules_reject_invalid_previous_requests", func(t *testing.T) {
		predecessor, _ := returnedRequest(t, "e2-resub-3")

		// Non-requester may not resubmit another user's request.
		if status, _ := h.create(fx.tokenAppr2, "purchase_request", "e2-resub-3", h.amount(100000), predecessor.Data.Number); status != http.StatusBadRequest {
			t.Fatalf("non-requester resubmission = %d, want %d", status, http.StatusBadRequest)
		}

		// Unknown number.
		if status, _ := h.create(fx.tokenAdminA, "purchase_request", "e2-resub-3", h.amount(100000), uuid.NewV4().String()); status != http.StatusBadRequest {
			t.Fatalf("unknown previous request = %d, want %d", status, http.StatusBadRequest)
		}

		// Cross-organisation: the predecessor belongs to organisation A.
		if status, _ := h.create(fx.tokenAdminB, "purchase_request", "e2-resub-x", h.amount(100000), predecessor.Data.Number); status != http.StatusBadRequest {
			t.Fatalf("cross-organisation link = %d, want %d", status, http.StatusBadRequest)
		}

		// A PENDING request is not linkable (only RETURNED is).
		status, pending := h.create(fx.tokenAdminA, "purchase_request", "e2-resub-4", h.amount(100000), "")
		if status != http.StatusCreated {
			t.Fatalf("create = %d", status)
		}
		if status, _ := h.create(fx.tokenAdminA, "purchase_request", "e2-resub-4", h.amount(100000), pending.Data.Number); status != http.StatusBadRequest {
			t.Fatalf("linking a PENDING request = %d, want %d", status, http.StatusBadRequest)
		}

		// A REJECTED request is terminal and must never be linkable (BD-E9 / D6).
		status, rejected := h.create(fx.tokenAdminA, "purchase_request", "e2-resub-5", h.amount(100000), "")
		if status != http.StatusCreated {
			t.Fatalf("create = %d", status)
		}
		if status, _, _ := h.decide(fx.tokenAppr1, rejected.Data.Number, "reject", "no"); status != http.StatusOK {
			t.Fatalf("reject = %d", status)
		}
		if status, _ := h.create(fx.tokenAdminA, "purchase_request", "e2-resub-5", h.amount(100000), rejected.Data.Number); status != http.StatusBadRequest {
			t.Fatalf("linking a REJECTED request = %d, want %d", status, http.StatusBadRequest)
		}
		// A brand-new unlinked request after a rejection is still allowed.
		if status, _ := h.create(fx.tokenAdminA, "purchase_request", "e2-resub-6", h.amount(100000), ""); status != http.StatusCreated {
			t.Fatalf("unlinked request after rejection = %d", status)
		}

		// None of the rejected attempts created a request.
		var count int64
		if err := fx.db.Table("ems_approval_requests").
			Where("organisation_number = ? AND resource_number = ?", eOrgA, "e2-resub-3").Count(&count).Error; err != nil {
			t.Fatalf("counting: %v", err)
		}
		if count != 1 {
			t.Fatalf("invalid linkage must create nothing (found %d requests for e2-resub-3)", count)
		}
	})

	t.Run("chains_link_to_the_immediate_predecessor", func(t *testing.T) {
		first, firstShape := returnedRequest(t, "e2-chain-1")

		status, second := h.create(fx.tokenAdminA, "purchase_request", "e2-chain-1", h.amount(100000), first.Data.Number)
		if status != http.StatusCreated {
			t.Fatalf("second create = %d", status)
		}
		if second.Data.PreviousRequestNumber != first.Data.Number {
			t.Fatalf("B must link to A")
		}

		// Return B, then resubmit C from B.
		if status, _, _ := h.decide(fx.tokenAppr1, second.Data.Number, "return", "still wrong"); status != http.StatusOK {
			t.Fatalf("return B = %d", status)
		}
		status, third := h.create(fx.tokenAdminA, "purchase_request", "e2-chain-1", h.amount(100000), second.Data.Number)
		if status != http.StatusCreated {
			t.Fatalf("third create = %d", status)
		}
		if third.Data.PreviousRequestNumber != second.Data.Number {
			t.Fatalf("C must link to B, got %q", third.Data.PreviousRequestNumber)
		}
		if third.Data.PreviousRequestNumber == first.Data.Number {
			t.Fatal("C must link only to its immediate predecessor, not to A")
		}

		// A and B remain untouched.
		if after := shape(h.fetch(fx.tokenAdminA, first.Data.Number)); after != firstShape {
			t.Fatal("A must remain intact across the chain")
		}
		if after := h.fetch(fx.tokenAdminA, second.Data.Number); after.Data.Status != "RETURNED" {
			t.Fatal("B must remain RETURNED")
		}
	})
}

// ============================================================================
// Stakeholder validity windows (D10)
// ============================================================================

func TestWorkflowE2StakeholderValidity(t *testing.T) {
	h := newE2Harness(t)
	fx := h.fixture

	t.Run("validity_window_boundaries_are_inclusive", func(t *testing.T) {
		// Repository-level check with a controlled clock, so the boundary is exact.
		reference := time.Now().Truncate(time.Second)
		boundary := reference

		number := "e2-repo-boundary"
		stakeholder := stakeholderModel.Stakeholder{
			Number: uuid.NewV4().String(), OrganisationNumber: eOrgA,
			UserNumber: number, StakeholderType: "consultant", Status: "ACTIVE",
			ValidFrom: &boundary, ValidTo: &boundary, AddedByNumber: eAdminA,
		}
		if err := fx.db.Create(&stakeholder).Error; err != nil {
			t.Fatalf("seeding boundary stakeholder: %v", err)
		}

		cases := []struct {
			name string
			now  time.Time
			want bool
		}{
			{"exactly at the window", reference, true},
			{"after valid_to", reference.Add(time.Second), false},
			{"before valid_from", reference.Add(-time.Second), false},
		}
		for _, testCase := range cases {
			numbers, err := approvalRepositories.FindUserNumbersByStakeholderType(eOrgA, "consultant", testCase.now)
			if err != nil {
				t.Fatalf("%s: %v", testCase.name, err)
			}
			found := false
			for _, candidate := range numbers {
				if candidate == number {
					found = true
				}
			}
			if found != testCase.want {
				t.Fatalf("%s (now=%s): eligible=%v, want %v", testCase.name, testCase.now, found, testCase.want)
			}
		}

		// Null bounds are unbounded on that side.
		unbounded := stakeholderModel.Stakeholder{
			Number: uuid.NewV4().String(), OrganisationNumber: eOrgA,
			UserNumber: "e2-repo-unbounded", StakeholderType: "consultant",
			Status: "ACTIVE", AddedByNumber: eAdminA,
		}
		if err := fx.db.Create(&unbounded).Error; err != nil {
			t.Fatalf("seeding unbounded stakeholder: %v", err)
		}
		numbers, err := approvalRepositories.FindUserNumbersByStakeholderType(eOrgA, "consultant", reference)
		if err != nil {
			t.Fatalf("unbounded: %v", err)
		}
		found := false
		for _, candidate := range numbers {
			if candidate == "e2-repo-unbounded" {
				found = true
			}
		}
		if !found {
			t.Fatal("a null validity window must remain eligible")
		}
	})

	t.Run("snapshot_honours_the_window", func(t *testing.T) {
		status, created := h.create(fx.tokenAdminA, e2StakeholderDecision, "e2-window-1", h.amount(10), "")
		if status != http.StatusCreated {
			t.Fatalf("create = %d", status)
		}

		step := stepAtLevel(created, 1)
		has := map[string]bool{}
		for _, assignment := range step.Assignments {
			has[assignment.UserNumber] = true
		}

		if !has[e2BoardActive] {
			t.Fatalf("a stakeholder inside the window must be a candidate: %+v", step.Assignments)
		}
		if has[e2BoardExpired] {
			t.Fatalf("an ACTIVE but time-expired stakeholder must not be a candidate")
		}
		if has[e2BoardFuture] {
			t.Fatalf("a not-yet-valid stakeholder must not be a candidate")
		}
	})

	t.Run("action_time_recheck_honours_the_window", func(t *testing.T) {
		status, created := h.create(fx.tokenAdminA, e2StakeholderDecision, "e2-window-2", h.amount(10), "")
		if status != http.StatusCreated {
			t.Fatalf("create = %d", status)
		}
		number := created.Data.Number

		// The window closes after the snapshot: the person may not act, and the
		// snapshot itself is not rebuilt.
		past := time.Now().Add(-time.Hour)
		if err := fx.db.Model(&stakeholderModel.Stakeholder{}).
			Where("organisation_number = ? AND user_number = ? AND stakeholder_type = ?",
				eOrgA, e2BoardActive, stakeholderModel.StakeholderTypeBoardMember).
			Update("valid_to", past).Error; err != nil {
			t.Fatalf("expiring the stakeholder: %v", err)
		}

		token, err := authServices.JWTAuthService().GenerateToken(e2BoardActive)
		if err != nil {
			t.Fatalf("token: %v", err)
		}

		if status, _, _ := h.decide(token, number, "approve", "should fail"); status != http.StatusForbidden {
			t.Fatalf("acting after the window closed = %d, want %d", status, http.StatusForbidden)
		}
		// The assignment row is untouched: no silent reassignment or rebuild.
		if after := h.fetch(fx.tokenAdminA, number); stepAtLevel(after, 1).Status != "PENDING" {
			t.Fatal("a failed eligibility check must not change the step")
		}
	})
}
