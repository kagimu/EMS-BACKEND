package integration

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/routes"
	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	authServices "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/services"
	assignmentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/assignments/models"
	authorityModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/authority/models"
	organisationModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/organisation/models"
	workflowModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/workflow/models"
	eosModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/eos/models"
	permissionsModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/permissions/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/internal/testutil"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/migrations"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// Foundation D (unit narrowing) test fixture identities.
const (
	unRequester = "un-requester"
	unAlice     = "un-alice"    // Finance Manager placed at Kampala
	unBob       = "un-bob"      // Finance Manager placed at Jinja
	unCarol     = "un-carol"    // Finance Manager placed at a child unit of Kampala
	unDave      = "un-dave"     // Finance Manager placed at the PARENT unit (region)
	unErin      = "un-erin"     // Finance Manager placed at an unrelated branch
	unExpired   = "un-expired"  // Finance Manager, Kampala placement already expired
	unFuture    = "un-future"   // Finance Manager, Kampala placement not yet valid
	unInactive  = "un-inactive" // Finance Manager, Kampala placement INACTIVE
	unNoPlace   = "un-noplace"  // Finance Manager with NO placement at all
	unTitleOnly = "un-title"    // NO role, Kampala placement, position_title "Finance Manager"
	unOtherOrg  = "un-other-org-user"

	unOrg  = "un-org"
	unOrgB = "un-org-b"

	unRoleFin = "un-role-fin"

	unRegion      = "un-unit-region"
	unKampala     = "un-unit-kampala"
	unJinja       = "un-unit-jinja"
	unKampalaFin  = "un-unit-kampala-fin" // child of Kampala
	unOtherBranch = "un-unit-other"

	unDecOrgWide  = "un_orgwide_request"
	unDecKampala  = "un_kampala_request"
	unDecJinja    = "un_jinja_request"
	unDecRegion   = "un_region_request"
	unDecCollapse = "un_collapse_request"
)

type unFixture struct {
	db     *gorm.DB
	tokens map[string]string
}

func (f *unFixture) token(userNumber string) string { return f.tokens[userNumber] }

func setupUnitNarrowingFixture(t *testing.T) *unFixture {
	t.Helper()

	db := testutil.RequireTestDB(t)
	if err := migrations.Init(); err != nil {
		t.Fatalf("running migrations: %v", err)
	}
	t.Setenv("JWT_SECRET", "ems-foundation-d-test-secret")

	users := []string{unRequester, unAlice, unBob, unCarol, unDave, unErin,
		unExpired, unFuture, unInactive, unNoPlace, unTitleOnly, unOtherOrg}
	organisations := []string{unOrg, unOrgB}
	roles := []string{unRoleFin}

	cleanup := func() {
		db.Unscoped().Where("organisation_number IN ?", organisations).Delete(&workflowModel.ApprovalStepAnchor{})
		db.Unscoped().Where("organisation_number IN ?", organisations).Delete(&workflowModel.ApprovalAssignment{})
		db.Unscoped().Where("organisation_number IN ?", organisations).Delete(&workflowModel.ApprovalStep{})
		db.Unscoped().Where("organisation_number IN ?", organisations).Delete(&workflowModel.ApprovalRequest{})
		db.Unscoped().Where("organisation_number IN ?", organisations).Delete(&authorityModel.DecisionAuthority{})
		db.Unscoped().Where("organisation_number IN ?", organisations).Delete(&assignmentModel.UserAssignment{})
		db.Unscoped().Where("organisation_number IN ?", organisations).Delete(&organisationModel.OrganisationUnit{})
		db.Unscoped().Where("role_number IN ?", roles).Delete(&permissionsModel.UserRole{})
		db.Unscoped().Where("user_number IN ?", users).Delete(&permissionsModel.UserRole{})
		db.Unscoped().Where("user_number IN ?", users).Delete(&authModel.Contact{})
		db.Unscoped().Where("user_number IN ?", users).Delete(&authModel.UserPassword{})
		db.Unscoped().Where("entity_number IN ?", organisations).Delete(&permissionsModel.Role{})
		db.Unscoped().Where("user_number IN ?", users).Delete(&eosModel.EntityUser{})
		db.Unscoped().Where("user_number IN ?", users).Delete(&eosModel.ActiveUserEntity{})
		db.Unscoped().Where("number IN ?", organisations).Delete(&eosModel.Entity{})
		db.Unscoped().Where("number IN ?", users).Delete(&authModel.User{})
	}
	cleanup()
	t.Cleanup(cleanup)

	for _, number := range users {
		user := authModel.User{
			Number: number, FirstName: "Unit", SecondName: "Narrowing",
			Gender: "M", Type: "Admin", Status: "Active",
		}
		if err := db.Create(&user).Error; err != nil {
			t.Fatalf("seeding user %s: %v", number, err)
		}
	}

	for _, spec := range []struct{ number, name string }{
		{unOrg, "Unit Narrowing Org"}, {unOrgB, "Unit Narrowing Other Org"},
	} {
		entity := eosModel.Entity{Number: spec.number, Name: spec.name, AddedBy: unRequester}
		if err := db.Create(&entity).Error; err != nil {
			t.Fatalf("seeding entity %s: %v", spec.number, err)
		}
	}

	for _, userNumber := range users {
		entityNumber := unOrg
		if userNumber == unOtherOrg {
			entityNumber = unOrgB
		}
		active := eosModel.ActiveUserEntity{UserNumber: userNumber, EntityNumber: entityNumber}
		if err := db.Omit("Entity", "User").Create(&active).Error; err != nil {
			t.Fatalf("setting active entity for %s: %v", userNumber, err)
		}
	}

	role := permissionsModel.Role{
		Number: unRoleFin, Name: "Finance Manager", EntityNumber: unOrg,
		Status: "ACTIVE", UserNumber: unRequester,
	}
	if err := db.Omit("User", "Entity", "Permissions").Create(&role).Error; err != nil {
		t.Fatalf("seeding role: %v", err)
	}

	expiry := time.Now().Add(365 * 24 * time.Hour)
	for _, userNumber := range []string{unAlice, unBob, unCarol, unDave, unErin,
		unExpired, unFuture, unInactive, unNoPlace} {
		userRole := permissionsModel.UserRole{
			UserNumber: userNumber, RoleNumber: unRoleFin, ExpiryDate: expiry, AddedBy: unRequester,
		}
		if err := db.Omit("User", "Role", "AddedByUser").Create(&userRole).Error; err != nil {
			t.Fatalf("seeding role for %s: %v", userNumber, err)
		}
	}

	// Foundation C2 eligibility: everyone is able to sign in, so exclusion from a
	// candidate set can only come from the anchor or the placement rules.
	for _, userNumber := range users {
		contact := authModel.Contact{
			Contact: userNumber + "@example.com", ContactType: "Email", UserNumber: userNumber,
		}
		if err := db.Omit("User").Create(&contact).Error; err != nil {
			t.Fatalf("seeding contact for %s: %v", userNumber, err)
		}
		password := authModel.UserPassword{UserNumber: userNumber, UserPassword: "seeded-hash"}
		if err := db.Omit("User").Create(&password).Error; err != nil {
			t.Fatalf("seeding password for %s: %v", userNumber, err)
		}
	}

	units := []organisationModel.OrganisationUnit{
		{Number: unRegion, OrganisationNumber: unOrg, ParentNumber: "", UnitType: "region",
			Name: "Region", Status: "ACTIVE", AddedByNumber: unRequester},
		{Number: unKampala, OrganisationNumber: unOrg, ParentNumber: unRegion, UnitType: "branch",
			Name: "Kampala", Status: "ACTIVE", AddedByNumber: unRequester},
		{Number: unJinja, OrganisationNumber: unOrg, ParentNumber: unRegion, UnitType: "branch",
			Name: "Jinja", Status: "ACTIVE", AddedByNumber: unRequester},
		{Number: unKampalaFin, OrganisationNumber: unOrg, ParentNumber: unKampala, UnitType: "department",
			Name: "Kampala Finance", Status: "ACTIVE", AddedByNumber: unRequester},
		{Number: unOtherBranch, OrganisationNumber: unOrg, ParentNumber: "", UnitType: "branch",
			Name: "Other", Status: "ACTIVE", AddedByNumber: unRequester},
	}
	for _, unit := range units {
		if err := db.Create(&unit).Error; err != nil {
			t.Fatalf("seeding unit %s: %v", unit.Number, err)
		}
	}

	past := time.Now().Add(-48 * time.Hour)
	future := time.Now().Add(48 * time.Hour)
	placements := []assignmentModel.UserAssignment{
		{OrganisationNumber: unOrg, UserNumber: unAlice, OrganisationUnitNumber: unKampala,
			AssignmentType: "manager", Status: "ACTIVE", AddedByNumber: unRequester},
		{OrganisationNumber: unOrg, UserNumber: unBob, OrganisationUnitNumber: unJinja,
			AssignmentType: "manager", Status: "ACTIVE", AddedByNumber: unRequester},
		// Carol sits in a CHILD unit of Kampala and holds a non-manager assignment
		// type: D3=S2 descendant placement and D4=any placement type.
		{OrganisationNumber: unOrg, UserNumber: unCarol, OrganisationUnitNumber: unKampalaFin,
			AssignmentType: "staff", Status: "ACTIVE", AddedByNumber: unRequester},
		// Dave sits ABOVE the governed unit: rejected by S2.
		{OrganisationNumber: unOrg, UserNumber: unDave, OrganisationUnitNumber: unRegion,
			AssignmentType: "head", Status: "ACTIVE", AddedByNumber: unRequester},
		{OrganisationNumber: unOrg, UserNumber: unErin, OrganisationUnitNumber: unOtherBranch,
			AssignmentType: "member", Status: "ACTIVE", AddedByNumber: unRequester},
		{OrganisationNumber: unOrg, UserNumber: unExpired, OrganisationUnitNumber: unKampala,
			AssignmentType: "manager", Status: "ACTIVE", ValidTo: &past, AddedByNumber: unRequester},
		{OrganisationNumber: unOrg, UserNumber: unFuture, OrganisationUnitNumber: unKampala,
			AssignmentType: "manager", Status: "ACTIVE", ValidFrom: &future, AddedByNumber: unRequester},
		{OrganisationNumber: unOrg, UserNumber: unInactive, OrganisationUnitNumber: unKampala,
			AssignmentType: "manager", Status: "INACTIVE", AddedByNumber: unRequester},
		// A placement is NOT authority: this person carries the job title but holds
		// no role, so must never be a candidate.
		{OrganisationNumber: unOrg, UserNumber: unTitleOnly, OrganisationUnitNumber: unKampala,
			AssignmentType: "manager", PositionTitle: "Finance Manager",
			Status: "ACTIVE", AddedByNumber: unRequester},
	}
	for _, placement := range placements {
		placement.Number = uuid.NewV4().String()
		if err := db.Create(&placement).Error; err != nil {
			t.Fatalf("seeding placement for %s: %v", placement.UserNumber, err)
		}
	}

	authorities := []authorityModel.DecisionAuthority{
		{DecisionType: unDecOrgWide, ApprovalLevel: 1, RoleNumber: unRoleFin},
		{DecisionType: unDecKampala, ApprovalLevel: 1, RoleNumber: unRoleFin, OrganisationUnitNumber: unKampala},
		{DecisionType: unDecJinja, ApprovalLevel: 1, RoleNumber: unRoleFin, OrganisationUnitNumber: unJinja},
		{DecisionType: unDecRegion, ApprovalLevel: 1, RoleNumber: unRoleFin, OrganisationUnitNumber: unRegion},
		// Same anchor, same level, DIFFERENT units: must stay two distinct anchors.
		{DecisionType: unDecCollapse, ApprovalLevel: 1, RoleNumber: unRoleFin, OrganisationUnitNumber: unKampala},
		{DecisionType: unDecCollapse, ApprovalLevel: 1, RoleNumber: unRoleFin, OrganisationUnitNumber: unRegion},
	}
	for _, authority := range authorities {
		authority.Number = uuid.NewV4().String()
		authority.OrganisationNumber = unOrg
		authority.Status = "ACTIVE"
		authority.AddedByNumber = unRequester
		if err := db.Create(&authority).Error; err != nil {
			t.Fatalf("seeding authority %s: %v", authority.DecisionType, err)
		}
	}

	tokens := map[string]string{}
	for _, userNumber := range users {
		value, err := authServices.JWTAuthService().GenerateToken(userNumber)
		if err != nil {
			t.Fatalf("generating token for %s: %v", userNumber, err)
		}
		tokens[userNumber] = value
	}

	return &unFixture{db: db, tokens: tokens}
}

// unHarness wraps the fixture with HTTP helpers.
type unHarness struct {
	*unFixture
	t      *testing.T
	server *httptest.Server
	client *http.Client
}

func newUnHarness(t *testing.T) *unHarness {
	fixture := setupUnitNarrowingFixture(t)
	server := httptest.NewServer(routes.SetupRoutes())
	t.Cleanup(server.Close)
	return &unHarness{unFixture: fixture, t: t, server: server, client: server.Client()}
}

func (h *unHarness) do(method, path, token string, payload any) (int, []byte) {
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

// createRequest creates an approval request for the requester and returns it.
func (h *unHarness) createRequest(decisionType, unitNumber string) eApprovalEnvelope {
	h.t.Helper()

	payload := map[string]any{
		"module": "finance", "resource_type": "purchase_order",
		"resource_number": uuid.NewV4().String(), "decision_type": decisionType,
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

// fetch reads a request back through the API.
func (h *unHarness) fetch(number string) eApprovalEnvelope {
	h.t.Helper()

	status, body := h.do(http.MethodGet, eApprovalsBase+"/"+number, h.token(unRequester), nil)
	if status != http.StatusOK {
		h.t.Fatalf("fetch %s = %d (%s)", number, status, string(body))
	}

	var envelope eApprovalEnvelope
	if err := json.Unmarshal(body, &envelope); err != nil {
		h.t.Fatalf("decoding fetched request: %v", err)
	}
	return envelope
}

func (h *unHarness) decide(token, number, action, comment string) (int, []byte) {
	h.t.Helper()
	return h.do(http.MethodPost, eApprovalsBase+"/"+number+"/"+action, token,
		map[string]any{"comment": comment})
}

// candidatesOf returns the sorted user numbers assigned to a level.
func candidatesOf(envelope eApprovalEnvelope, level int) []string {
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

// anchorKeysOf returns the sorted "type:value@unit" identity of a level's anchors.
func anchorKeysOf(envelope eApprovalEnvelope, level int) []string {
	step := stepAtLevel(envelope, level)
	if step == nil {
		return nil
	}
	result := make([]string, 0, len(step.Anchors))
	for _, anchor := range step.Anchors {
		result = append(result, anchor.Type+":"+anchor.Value+"@"+anchor.OrganisationUnitNumber)
	}
	sort.Strings(result)
	return result
}

func assertCandidates(t *testing.T, envelope eApprovalEnvelope, level int, want ...string) {
	t.Helper()
	sort.Strings(want)
	got := candidatesOf(envelope, level)
	if strings.Join(got, ",") != strings.Join(want, ",") {
		t.Fatalf("candidates = %v, want %v", got, want)
	}
}

// ============================================================================
// Authority scope and anchor identity
// ============================================================================

func TestUnitNarrowing(t *testing.T) {
	h := newUnHarness(t)

	t.Run("organisation_wide_authority_applies_no_placement_filter", func(t *testing.T) {
		envelope := h.createRequest(unDecOrgWide, "")
		// Every Finance Manager in the organisation qualifies, regardless of any
		// placement: narrowing is skipped entirely for an empty unit (D7).
		assertCandidates(t, envelope, 1,
			unAlice, unBob, unCarol, unDave, unErin, unExpired, unFuture, unInactive, unNoPlace)
		if keys := anchorKeysOf(envelope, 1); len(keys) != 1 || keys[0] != "role:"+unRoleFin+"@" {
			t.Fatalf("organisation-wide anchor must have an empty unit: %v", keys)
		}
	})

	t.Run("unit_specific_authority_narrows_to_placement", func(t *testing.T) {
		envelope := h.createRequest(unDecKampala, unKampala)
		assertCandidates(t, envelope, 1, unAlice, unCarol)
		if keys := anchorKeysOf(envelope, 1); len(keys) != 1 || keys[0] != "role:"+unRoleFin+"@"+unKampala {
			t.Fatalf("anchor must carry the Kampala unit: %v", keys)
		}
	})

	t.Run("jinja_request_resolves_only_jinja", func(t *testing.T) {
		envelope := h.createRequest(unDecJinja, unJinja)
		assertCandidates(t, envelope, 1, unBob)
	})

	t.Run("parent_unit_authority_applies_to_descendants", func(t *testing.T) {
		// A rule scoped to the region applies to a Kampala request (D2), and its
		// anchor is the REGION, so narrowing covers the whole region subtree.
		envelope := h.createRequest(unDecRegion, unKampala)
		assertCandidates(t, envelope, 1, unAlice, unBob, unCarol, unDave)
		if keys := anchorKeysOf(envelope, 1); len(keys) != 1 || keys[0] != "role:"+unRoleFin+"@"+unRegion {
			t.Fatalf("anchor must carry the region unit: %v", keys)
		}
	})

	t.Run("same_anchor_different_units_do_not_collapse", func(t *testing.T) {
		envelope := h.createRequest(unDecCollapse, unKampala)

		keys := anchorKeysOf(envelope, 1)
		wantKeys := []string{
			"role:" + unRoleFin + "@" + unKampala,
			"role:" + unRoleFin + "@" + unRegion,
		}
		sort.Strings(wantKeys)
		if strings.Join(keys, ",") != strings.Join(wantKeys, ",") {
			t.Fatalf("same role in two units must stay two distinct anchors: got %v want %v", keys, wantKeys)
		}

		// The candidate set is the union across both anchors.
		assertCandidates(t, envelope, 1, unAlice, unBob, unCarol, unDave)

		// Both anchors are persisted distinctly (D6).
		var rows []workflowModel.ApprovalStepAnchor
		if err := h.db.Where("step_number = ?", stepAtLevel(envelope, 1).Number).Find(&rows).Error; err != nil {
			t.Fatalf("reading step anchors: %v", err)
		}
		if len(rows) != 2 {
			t.Fatalf("expected 2 persisted anchors, got %d", len(rows))
		}
		units := map[string]bool{}
		for _, row := range rows {
			units[row.OrganisationUnitNumber] = true
		}
		if !units[unKampala] || !units[unRegion] {
			t.Fatalf("both units must be persisted: %+v", rows)
		}
	})

	t.Run("anchor_rows_persist_the_unit", func(t *testing.T) {
		kampala := h.createRequest(unDecKampala, unKampala)
		jinja := h.createRequest(unDecJinja, unJinja)

		type stored struct{ Value, Unit string }
		load := func(envelope eApprovalEnvelope) stored {
			var rows []workflowModel.ApprovalStepAnchor
			if err := h.db.Where("step_number = ?", stepAtLevel(envelope, 1).Number).Find(&rows).Error; err != nil {
				t.Fatalf("reading step anchors: %v", err)
			}
			if len(rows) != 1 {
				t.Fatalf("expected exactly 1 anchor row, got %d", len(rows))
			}
			return stored{Value: rows[0].AnchorValue, Unit: rows[0].OrganisationUnitNumber}
		}

		kampalaAnchor := load(kampala)
		jinjaAnchor := load(jinja)
		if kampalaAnchor.Value != unRoleFin || kampalaAnchor.Unit != unKampala {
			t.Fatalf("unexpected Kampala anchor row: %+v", kampalaAnchor)
		}
		if jinjaAnchor.Value != unRoleFin || jinjaAnchor.Unit != unJinja {
			t.Fatalf("unexpected Jinja anchor row: %+v", jinjaAnchor)
		}
	})

	t.Run("placement_alone_never_grants_authority", func(t *testing.T) {
		// unTitleOnly is placed at Kampala and even carries the job title, but holds
		// no role, so must never be a candidate.
		envelope := h.createRequest(unDecKampala, unKampala)
		for _, candidate := range candidatesOf(envelope, 1) {
			if candidate == unTitleOnly {
				t.Fatal("a placement (or a job title) must never grant approval authority")
			}
		}
	})
}

// ============================================================================
// Placement validity, snapshot immutability and action-time authorization
// ============================================================================

func TestUnitNarrowingLifecycle(t *testing.T) {
	h := newUnHarness(t)

	t.Run("snapshot_is_frozen_against_later_changes", func(t *testing.T) {
		envelope := h.createRequest(unDecKampala, unKampala)
		number := envelope.Data.Number
		assertCandidates(t, envelope, 1, unAlice, unCarol)

		// A new Kampala placement appears after creation …
		newcomer := assignmentModel.UserAssignment{
			Number: uuid.NewV4().String(), OrganisationNumber: unOrg, UserNumber: unNoPlace,
			OrganisationUnitNumber: unKampala, AssignmentType: "manager",
			Status: "ACTIVE", AddedByNumber: unRequester,
		}
		if err := h.db.Create(&newcomer).Error; err != nil {
			t.Fatalf("adding placement: %v", err)
		}
		// … and the authority rule is deactivated.
		if err := h.db.Model(&authorityModel.DecisionAuthority{}).
			Where("organisation_number = ? AND decision_type = ?", unOrg, unDecKampala).
			Update("status", "INACTIVE").Error; err != nil {
			t.Fatalf("deactivating authority: %v", err)
		}
		// Restore the shared fixture state for the following subtests.
		t.Cleanup(func() {
			h.db.Unscoped().Where("number = ?", newcomer.Number).Delete(&assignmentModel.UserAssignment{})
			h.db.Model(&authorityModel.DecisionAuthority{}).
				Where("organisation_number = ? AND decision_type = ?", unOrg, unDecKampala).
				Update("status", "ACTIVE")
		})

		// The snapshotted candidate list must be unchanged.
		assertCandidates(t, h.fetch(number), 1, unAlice, unCarol)
	})

	t.Run("eligible_candidate_can_act_and_ineligible_cannot", func(t *testing.T) {
		envelope := h.createRequest(unDecKampala, unKampala)
		number := envelope.Data.Number

		// Bob is a Finance Manager but placed in Jinja: never a candidate here.
		if status, _ := h.decide(h.token(unBob), number, "approve", "not mine"); status != http.StatusForbidden {
			t.Fatalf("a Jinja-placed approver on a Kampala request = %d, want %d", status, http.StatusForbidden)
		}

		// Alice is a snapshotted candidate and currently eligible.
		if status, body := h.decide(h.token(unAlice), number, "approve", "ok"); status != http.StatusOK {
			t.Fatalf("eligible candidate approve = %d (%s)", status, string(body))
		}
		if final := h.fetch(number); final.Data.Status != "APPROVED" {
			t.Fatalf("request status = %s, want APPROVED", final.Data.Status)
		}
	})

	t.Run("candidate_that_becomes_ineligible_is_rejected_and_not_replaced", func(t *testing.T) {
		envelope := h.createRequest(unDecKampala, unKampala)
		number := envelope.Data.Number
		assertCandidates(t, envelope, 1, unAlice, unCarol)

		// Alice's placement ends after the snapshot.
		past := time.Now().Add(-time.Hour)
		if err := h.db.Model(&assignmentModel.UserAssignment{}).
			Where("organisation_number = ? AND user_number = ?", unOrg, unAlice).
			Update("valid_to", past).Error; err != nil {
			t.Fatalf("expiring placement: %v", err)
		}
		// Restore the shared fixture state for the following subtests.
		t.Cleanup(func() {
			h.db.Model(&assignmentModel.UserAssignment{}).
				Where("organisation_number = ? AND user_number = ?", unOrg, unAlice).
				Update("valid_to", nil)
		})

		if status, _ := h.decide(h.token(unAlice), number, "approve", "should fail"); status != http.StatusForbidden {
			t.Fatalf("ineligible snapshotted candidate = %d, want %d", status, http.StatusForbidden)
		}

		// No automatic substitution: the candidate list is unchanged.
		assertCandidates(t, h.fetch(number), 1, unAlice, unCarol)
	})

	t.Run("all_candidates_ineligible_stalls_without_automatic_action", func(t *testing.T) {
		envelope := h.createRequest(unDecKampala, unKampala)
		number := envelope.Data.Number

		past := time.Now().Add(-time.Hour)
		for _, userNumber := range []string{unAlice, unCarol} {
			if err := h.db.Model(&assignmentModel.UserAssignment{}).
				Where("organisation_number = ? AND user_number = ?", unOrg, userNumber).
				Update("valid_to", past).Error; err != nil {
				t.Fatalf("expiring placement for %s: %v", userNumber, err)
			}
		}

		for _, userNumber := range []string{unAlice, unCarol} {
			if status, _ := h.decide(h.token(userNumber), number, "approve", "stalled"); status != http.StatusForbidden {
				t.Fatalf("%s acting on a stalled level = %d, want %d", userNumber, status, http.StatusForbidden)
			}
		}

		stalled := h.fetch(number)
		if stalled.Data.Status != "PENDING" {
			t.Fatalf("a stalled request must not be auto-approved/cancelled: status = %s", stalled.Data.Status)
		}
		assertCandidates(t, stalled, 1, unAlice, unCarol)
	})
}

// ============================================================================
// Backward compatibility (D9 / D11)
// ============================================================================

func TestUnitNarrowingBackwardCompatibility(t *testing.T) {
	h := newUnHarness(t)

	t.Run("legacy_step_without_anchor_rows_stays_organisation_wide", func(t *testing.T) {
		envelope := h.createRequest(unDecKampala, unKampala)
		number := envelope.Data.Number
		step := stepAtLevel(envelope, 1)
		assertCandidates(t, envelope, 1, unAlice, unCarol)

		// Emulate a pre-amendment snapshot: no normalised anchor rows exist.
		if err := h.db.Unscoped().
			Where("step_number = ?", step.Number).
			Delete(&workflowModel.ApprovalStepAnchor{}).Error; err != nil {
			t.Fatalf("removing anchor rows: %v", err)
		}

		legacy := h.fetch(number)
		if keys := anchorKeysOf(legacy, 1); len(keys) != 1 || keys[0] != "role:"+unRoleFin+"@" {
			t.Fatalf("a legacy step must read as organisation-wide: %v", keys)
		}
		// The snapshot itself is untouched and still actionable by its candidates.
		assertCandidates(t, legacy, 1, unAlice, unCarol)
		if status, body := h.decide(h.token(unAlice), number, "approve", "legacy ok"); status != http.StatusOK {
			t.Fatalf("legacy request approve = %d (%s)", status, string(body))
		}
	})

	t.Run("unit_less_authority_is_unchanged", func(t *testing.T) {
		envelope := h.createRequest(unDecOrgWide, "")
		if keys := anchorKeysOf(envelope, 1); keys[0] != "role:"+unRoleFin+"@" {
			t.Fatalf("unit-less authority must stay organisation-wide: %v", keys)
		}
	})
}

// ============================================================================
// Security and schema
// ============================================================================

func TestUnitNarrowingSecurity(t *testing.T) {
	h := newUnHarness(t)

	t.Run("cross_organisation_access_is_not_found", func(t *testing.T) {
		envelope := h.createRequest(unDecKampala, unKampala)
		number := envelope.Data.Number

		if status, _ := h.do(http.MethodGet, eApprovalsBase+"/"+number, h.token(unOtherOrg), nil); status != http.StatusNotFound {
			t.Fatalf("cross-organisation read = %d, want %d", status, http.StatusNotFound)
		}
		if status, _ := h.decide(h.token(unOtherOrg), number, "approve", "x"); status != http.StatusNotFound {
			t.Fatalf("cross-organisation approve = %d, want %d", status, http.StatusNotFound)
		}
	})
}

// TestUnitNarrowingSchema verifies the ratified database facts: the normalised
// anchor table exists and no stalled-status mechanism was introduced (D6, D11).
func TestUnitNarrowingSchema(t *testing.T) {
	h := newUnHarness(t)

	var anchorTables int64
	if err := h.db.Raw(
		"SELECT COUNT(*) FROM information_schema.TABLES WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'ems_approval_step_anchors'",
	).Scan(&anchorTables).Error; err != nil {
		t.Fatalf("checking anchor table: %v", err)
	}
	if anchorTables != 1 {
		t.Fatalf("ems_approval_step_anchors must exist (D6), found %d", anchorTables)
	}

	var stallColumns int64
	// Match stalled-status column names precisely: the earlier substring check
	// ('%stall%') also matched the legitimate max_installations column added by
	// Foundation F Increment 1. The guard still fails on any stalled-status
	// column (stalled, stall_status, is_stalled, ...).
	if err := h.db.Raw(
		"SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND (COLUMN_NAME LIKE 'stall%' OR COLUMN_NAME LIKE '%stall' OR COLUMN_NAME LIKE '%stalled%')",
	).Scan(&stallColumns).Error; err != nil {
		t.Fatalf("checking stalled columns: %v", err)
	}
	if stallColumns != 0 {
		t.Fatalf("no persisted stalled-status column may exist (D11/E3 D2), found %d", stallColumns)
	}
}
