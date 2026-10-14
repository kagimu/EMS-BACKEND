package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/routes"
	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	authServices "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/services"
	authorityModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/authority/models"
	organisationModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/organisation/models"
	stakeholderModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/stakeholders/models"
	workflowModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/workflow/models"
	eosModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/eos/models"
	permissionsModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/permissions/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/internal/testutil"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/migrations"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

const (
	eAdminA    = "ems-e-admin-a"
	eApprover1 = "ems-e-approver-1"
	eApprover2 = "ems-e-approver-2"
	eOwnerA    = "ems-e-owner-a"
	eSolo      = "ems-e-solo"
	eGhost     = "ems-e-ghost"
	eNoPerm    = "ems-e-noperm"
	eNoCtx     = "ems-e-nocontext"
	eAdminB    = "ems-e-admin-b"

	eOrgA = "ems-e-org-a"
	eOrgB = "ems-e-org-b"

	eRoleL1    = "ems-e-role-l1"
	eRoleSolo  = "ems-e-role-solo"
	eRoleGhost = "ems-e-role-ghost"
	eUnitA     = "ems-e-unit-a"

	eApprovalsBase = "/main/ems/approvals"
	eTasksURL      = "/main/ems/approval-tasks"
	eEduURL        = "/main/campus/all"
)

type eFixture struct {
	db          *gorm.DB
	tokenAdminA string
	tokenAppr1  string
	tokenAppr2  string
	tokenOwnerA string
	tokenSolo   string
	tokenNoPerm string
	tokenNoCtx  string
	tokenAdminB string
}

func setupEFixture(t *testing.T) *eFixture {
	t.Helper()

	db := testutil.RequireTestDB(t)
	if err := migrations.Init(); err != nil {
		t.Fatalf("running migrations: %v", err)
	}
	t.Setenv("JWT_SECRET", "ems-foundation-e-test-secret")

	users := []string{eAdminA, eApprover1, eApprover2, eOwnerA, eSolo, eGhost, eNoPerm, eNoCtx, eAdminB}
	organisations := []string{eOrgA, eOrgB}
	roles := []string{eRoleL1, eRoleSolo, eRoleGhost}

	cleanup := func() {
		db.Unscoped().Where("organisation_number IN ?", organisations).Delete(&workflowModel.ApprovalStepAnchor{})
		db.Unscoped().Where("organisation_number IN ?", organisations).Delete(&workflowModel.ApprovalAssignment{})
		db.Unscoped().Where("organisation_number IN ?", organisations).Delete(&workflowModel.ApprovalStep{})
		db.Unscoped().Where("organisation_number IN ?", organisations).Delete(&workflowModel.ApprovalRequest{})
		db.Unscoped().Where("organisation_number IN ?", organisations).Delete(&authorityModel.DecisionAuthority{})
		db.Unscoped().Where("organisation_number IN ?", organisations).Delete(&stakeholderModel.Stakeholder{})
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

	for _, spec := range []struct{ number, userType string }{
		{eAdminA, "Admin"}, {eApprover1, "Admin"}, {eApprover2, "Admin"}, {eOwnerA, "Admin"},
		{eSolo, "Admin"}, {eGhost, "User"}, {eNoPerm, "User"}, {eNoCtx, "Admin"}, {eAdminB, "Admin"},
	} {
		user := authModel.User{
			Number: spec.number, FirstName: "E", SecondName: spec.userType,
			Gender: "M", Type: spec.userType, Status: "Active",
		}
		if err := db.Create(&user).Error; err != nil {
			t.Fatalf("seeding user %s: %v", spec.number, err)
		}
	}

	for _, spec := range []struct{ number, name string }{
		{eOrgA, "E Test Org A"}, {eOrgB, "E Test Org B"},
	} {
		entity := eosModel.Entity{Number: spec.number, Name: spec.name, AddedBy: eAdminA}
		if err := db.Create(&entity).Error; err != nil {
			t.Fatalf("seeding entity %s: %v", spec.number, err)
		}
	}

	for _, spec := range []struct{ userNumber, entityNumber string }{
		{eAdminA, eOrgA}, {eApprover1, eOrgA}, {eApprover2, eOrgA}, {eOwnerA, eOrgA},
		{eSolo, eOrgA}, {eGhost, eOrgA}, {eNoPerm, eOrgA}, {eAdminB, eOrgB},
	} {
		active := eosModel.ActiveUserEntity{UserNumber: spec.userNumber, EntityNumber: spec.entityNumber}
		if err := db.Omit("Entity", "User").Create(&active).Error; err != nil {
			t.Fatalf("setting active entity for %s: %v", spec.userNumber, err)
		}
	}

	for _, number := range roles {
		role := permissionsModel.Role{
			Number: number, Name: number, EntityNumber: eOrgA, Status: "ACTIVE", UserNumber: eAdminA,
		}
		if err := db.Omit("User", "Entity", "Permissions").Create(&role).Error; err != nil {
			t.Fatalf("seeding role %s: %v", number, err)
		}
	}

	expiry := time.Now().Add(365 * 24 * time.Hour)
	for _, spec := range []struct{ userNumber, roleNumber string }{
		{eApprover1, eRoleL1}, {eApprover2, eRoleL1}, {eSolo, eRoleSolo}, {eGhost, eRoleGhost},
	} {
		userRole := permissionsModel.UserRole{
			UserNumber: spec.userNumber, RoleNumber: spec.roleNumber, ExpiryDate: expiry, AddedBy: eAdminA,
		}
		if err := db.Omit("User", "Role", "AddedByUser").Create(&userRole).Error; err != nil {
			t.Fatalf("seeding user role for %s: %v", spec.userNumber, err)
		}
	}

	// Foundation C2 eligibility: an approver must be able to sign in, which
	// requires an Email contact and a password row. eGhost deliberately has
	// neither and therefore can never be an approver.
	for _, spec := range []struct{ userNumber, contact string }{
		{eAdminA, "ems-e-admin-a@example.com"}, {eApprover1, "ems-e-approver-1@example.com"},
		{eApprover2, "ems-e-approver-2@example.com"}, {eOwnerA, "ems-e-owner-a@example.com"},
		{eSolo, "ems-e-solo@example.com"},
	} {
		contact := authModel.Contact{Contact: spec.contact, ContactType: "Email", UserNumber: spec.userNumber}
		if err := db.Omit("User").Create(&contact).Error; err != nil {
			t.Fatalf("seeding contact for %s: %v", spec.userNumber, err)
		}
		password := authModel.UserPassword{UserNumber: spec.userNumber, UserPassword: "seeded-hash"}
		if err := db.Omit("User").Create(&password).Error; err != nil {
			t.Fatalf("seeding password for %s: %v", spec.userNumber, err)
		}
	}

	stakeholder := stakeholderModel.Stakeholder{
		Number: uuid.NewV4().String(), OrganisationNumber: eOrgA, UserNumber: eOwnerA,
		StakeholderType: stakeholderModel.StakeholderTypeOwner, Status: "ACTIVE", AddedByNumber: eAdminA,
	}
	if err := db.Create(&stakeholder).Error; err != nil {
		t.Fatalf("seeding owner stakeholder: %v", err)
	}

	unit := organisationModel.OrganisationUnit{
		Number: eUnitA, OrganisationNumber: eOrgA, UnitType: "department",
		Name: "E Unit A", Status: "ACTIVE", AddedByNumber: eAdminA,
	}
	if err := db.Create(&unit).Error; err != nil {
		t.Fatalf("seeding unit: %v", err)
	}

	amount := func(v float64) *float64 { return &v }
	authorities := []authorityModel.DecisionAuthority{
		{Number: uuid.NewV4().String(), OrganisationNumber: eOrgA, DecisionType: "purchase_request",
			ApprovalLevel: 1, MaxAmount: amount(500000), RoleNumber: eRoleL1, Status: "ACTIVE", AddedByNumber: eAdminA},
		{Number: uuid.NewV4().String(), OrganisationNumber: eOrgA, DecisionType: "purchase_request",
			ApprovalLevel: 2, StakeholderType: stakeholderModel.StakeholderTypeOwner, Status: "ACTIVE", AddedByNumber: eAdminA},
		{Number: uuid.NewV4().String(), OrganisationNumber: eOrgA, DecisionType: "capped_request",
			ApprovalLevel: 1, MaxAmount: amount(1000), RoleNumber: eRoleL1, Status: "ACTIVE", AddedByNumber: eAdminA},
		{Number: uuid.NewV4().String(), OrganisationNumber: eOrgA, DecisionType: "unit_only_request",
			ApprovalLevel: 1, MaxAmount: amount(1000000), RoleNumber: eRoleL1,
			OrganisationUnitNumber: eUnitA, Status: "ACTIVE", AddedByNumber: eAdminA},
		{Number: uuid.NewV4().String(), OrganisationNumber: eOrgA, DecisionType: "ghost_request",
			ApprovalLevel: 1, RoleNumber: eRoleGhost, Status: "ACTIVE", AddedByNumber: eAdminA},
		{Number: uuid.NewV4().String(), OrganisationNumber: eOrgA, DecisionType: "solo_request",
			ApprovalLevel: 1, RoleNumber: eRoleSolo, Status: "ACTIVE", AddedByNumber: eAdminA},
	}
	for _, authority := range authorities {
		if err := db.Create(&authority).Error; err != nil {
			t.Fatalf("seeding authority %s: %v", authority.DecisionType, err)
		}
	}

	token := func(userNumber string) string {
		value, err := authServices.JWTAuthService().GenerateToken(userNumber)
		if err != nil {
			t.Fatalf("generating token for %s: %v", userNumber, err)
		}
		return value
	}

	return &eFixture{
		db:          db,
		tokenAdminA: token(eAdminA),
		tokenAppr1:  token(eApprover1),
		tokenAppr2:  token(eApprover2),
		tokenOwnerA: token(eOwnerA),
		tokenSolo:   token(eSolo),
		tokenNoPerm: token(eNoPerm),
		tokenNoCtx:  token(eNoCtx),
		tokenAdminB: token(eAdminB),
	}
}

type eAssignment struct {
	Number     string     `json:"number"`
	UserNumber string     `json:"user_number"`
	Status     string     `json:"status"`
	Comment    string     `json:"comment"`
	DecidedAt  *time.Time `json:"decided_at"`
}

type eStep struct {
	Number        string `json:"number"`
	ApprovalLevel int    `json:"approval_level"`
	Status        string `json:"status"`
	Actionable    bool   `json:"actionable"`
	Anchors       []struct {
		Type                   string `json:"type"`
		Value                  string `json:"value"`
		OrganisationUnitNumber string `json:"organisation_unit_number"`
	} `json:"anchors"`
	Assignments []eAssignment `json:"assignments"`
}

type eApprovalEnvelope struct {
	Message string `json:"message"`
	Data    struct {
		Number                 string   `json:"number"`
		OrganisationNumber     string   `json:"organisation_number"`
		OrganisationUnitNumber string   `json:"organisation_unit_number"`
		Module                 string   `json:"module"`
		ResourceType           string   `json:"resource_type"`
		ResourceNumber         string   `json:"resource_number"`
		DecisionType           string   `json:"decision_type"`
		Amount                 *float64 `json:"amount"`
		Status                 string   `json:"status"`
		CurrentLevel           int      `json:"current_level"`
		Stalled                bool     `json:"stalled"`
		RequestedByNumber      string   `json:"requested_by_number"`
		PreviousRequestNumber  string   `json:"previous_request_number"`
		Steps                  []eStep  `json:"steps"`
	} `json:"data"`
}

type eTaskEnvelope struct {
	Message string `json:"message"`
	Data    []struct {
		RequestNumber    string `json:"request_number"`
		StepNumber       string `json:"step_number"`
		AssignmentNumber string `json:"assignment_number"`
		ApprovalLevel    int    `json:"approval_level"`
		DecisionType     string `json:"decision_type"`
		RequestStatus    string `json:"request_status"`
	} `json:"data"`
}

// TestWorkflowApprovals covers EMS Foundation E1 end to end.
func TestWorkflowApprovals(t *testing.T) {
	fx := setupEFixture(t)

	server := httptest.NewServer(routes.SetupRoutes())
	t.Cleanup(server.Close)
	client := server.Client()

	request := func(method, path, token string, payload any) (int, []byte) {
		t.Helper()

		var reader io.Reader
		if payload != nil {
			raw, err := json.Marshal(payload)
			if err != nil {
				t.Fatalf("encoding payload: %v", err)
			}
			reader = bytes.NewReader(raw)
		}

		req, err := http.NewRequest(method, server.URL+path, reader)
		if err != nil {
			t.Fatalf("building request: %v", err)
		}
		if token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}
		if payload != nil {
			req.Header.Set("Content-Type", "application/json")
		}

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("request %s %s failed: %v", method, path, err)
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		return resp.StatusCode, body
	}

	create := func(t *testing.T, token, decisionType, resourceNumber string, amount *float64) (int, eApprovalEnvelope, []byte) {
		t.Helper()
		payload := map[string]any{
			"module": "finance", "resource_type": "purchase_order",
			"resource_number": resourceNumber, "decision_type": decisionType,
		}
		if amount != nil {
			payload["amount"] = *amount
		}
		status, body := request(http.MethodPost, eApprovalsBase, token, payload)
		var envelope eApprovalEnvelope
		if status == http.StatusCreated {
			if err := json.Unmarshal(body, &envelope); err != nil {
				t.Fatalf("decoding created approval %s: %v", string(body), err)
			}
		}
		return status, envelope, body
	}

	decide := func(t *testing.T, token, number, action, comment string) (int, eApprovalEnvelope, []byte) {
		t.Helper()
		status, body := request(http.MethodPost, eApprovalsBase+"/"+number+"/"+action, token,
			map[string]any{"comment": comment})
		var envelope eApprovalEnvelope
		if status == http.StatusOK {
			if err := json.Unmarshal(body, &envelope); err != nil {
				t.Fatalf("decoding decision %s: %v", string(body), err)
			}
		}
		return status, envelope, body
	}

	stepAt := func(envelope eApprovalEnvelope, level int) *eStep {
		for i := range envelope.Data.Steps {
			if envelope.Data.Steps[i].ApprovalLevel == level {
				return &envelope.Data.Steps[i]
			}
		}
		return nil
	}

	assignmentFor := func(step *eStep, userNumber string) *eAssignment {
		for i := range step.Assignments {
			if step.Assignments[i].UserNumber == userNumber {
				return &step.Assignments[i]
			}
		}
		return nil
	}

	f := func(v float64) *float64 { return &v }

	t.Run("authentication_boundaries", func(t *testing.T) {
		if status, _ := request(http.MethodGet, eApprovalsBase, "", nil); status != http.StatusUnauthorized {
			t.Fatalf("unauthenticated list = %d, want %d", status, http.StatusUnauthorized)
		}
		if status, _ := request(http.MethodGet, eTasksURL, "", nil); status != http.StatusUnauthorized {
			t.Fatalf("unauthenticated tasks = %d, want %d", status, http.StatusUnauthorized)
		}
	})

	t.Run("permission_and_context_boundaries", func(t *testing.T) {
		if status, _, _ := create(t, fx.tokenNoPerm, "purchase_request", "perm-1", f(10)); status != http.StatusForbidden {
			t.Fatalf("no-permission create = %d, want %d", status, http.StatusForbidden)
		}
		if status, _, _ := create(t, fx.tokenNoCtx, "purchase_request", "ctx-1", f(10)); status != http.StatusConflict {
			t.Fatalf("no-active-organisation create = %d, want %d", status, http.StatusConflict)
		}
	})

	t.Run("create_snapshots_path_and_candidates", func(t *testing.T) {
		status, envelope, body := create(t, fx.tokenAdminA, "purchase_request", "po-single", f(100000))
		if status != http.StatusCreated {
			t.Fatalf("create = %d (%s)", status, string(body))
		}
		if envelope.Data.Status != "PENDING" || envelope.Data.CurrentLevel != 1 {
			t.Fatalf("unexpected initial state: %s", string(body))
		}
		if envelope.Data.OrganisationNumber != eOrgA || envelope.Data.RequestedByNumber != eAdminA {
			t.Fatalf("organisation and requester must come from the authenticated context: %s", string(body))
		}
		if len(envelope.Data.Steps) != 1 {
			t.Fatalf("expected one snapshotted level, got %s", string(body))
		}

		step := stepAt(envelope, 1)
		if step == nil || !step.Actionable || step.Status != "PENDING" {
			t.Fatalf("level 1 must be actionable: %s", string(body))
		}
		if len(step.Anchors) != 1 || step.Anchors[0].Type != "role" || step.Anchors[0].Value != eRoleL1 {
			t.Fatalf("level 1 must snapshot the configured role anchor: %s", string(body))
		}

		// BD-E5: every eligible candidate is represented; none is "primary".
		if assignmentFor(step, eApprover1) == nil || assignmentFor(step, eApprover2) == nil {
			t.Fatalf("both role holders must be candidates: %s", string(body))
		}
		// BD-E6: the requester is never a candidate.
		if assignmentFor(step, eAdminA) != nil {
			t.Fatalf("the requester must not be an approver of their own request: %s", string(body))
		}
		for _, assignment := range step.Assignments {
			if assignment.Status != "UNDECIDED" {
				t.Fatalf("candidates must start undecided: %s", string(body))
			}
		}
	})

	t.Run("task_inbox_is_per_candidate", func(t *testing.T) {
		status, created, body := create(t, fx.tokenAdminA, "purchase_request", "po-inbox", f(100000))
		if status != http.StatusCreated {
			t.Fatalf("create = %d (%s)", status, string(body))
		}
		number := created.Data.Number

		for _, token := range []string{fx.tokenAppr1, fx.tokenAppr2} {
			taskStatus, taskBody := request(http.MethodGet, eTasksURL, token, nil)
			if taskStatus != http.StatusOK {
				t.Fatalf("tasks = %d (%s)", taskStatus, string(taskBody))
			}
			var tasks eTaskEnvelope
			if err := json.Unmarshal(taskBody, &tasks); err != nil {
				t.Fatalf("decoding tasks: %v", err)
			}
			found := false
			for _, task := range tasks.Data {
				if task.RequestNumber == number {
					found = true
					if task.ApprovalLevel != 1 || task.RequestStatus != "PENDING" || task.DecisionType != "purchase_request" {
						t.Fatalf("unexpected task payload: %s", string(taskBody))
					}
				}
			}
			if !found {
				t.Fatalf("an assigned approver must see their task: %s", string(taskBody))
			}
		}

		// A user with no assignment never sees the work.
		_, otherBody := request(http.MethodGet, eTasksURL, fx.tokenAdminA, nil)
		if bytes.Contains(otherBody, []byte(number)) {
			t.Fatalf("a non-candidate must not see another person's task: %s", string(otherBody))
		}
	})

	t.Run("any_candidate_completes_the_level", func(t *testing.T) {
		status, created, body := create(t, fx.tokenAdminA, "purchase_request", "po-any", f(100000))
		if status != http.StatusCreated {
			t.Fatalf("create = %d (%s)", status, string(body))
		}
		number := created.Data.Number

		status, approved, body := decide(t, fx.tokenAppr2, number, "approve", "looks good")
		if status != http.StatusOK {
			t.Fatalf("approve = %d (%s)", status, string(body))
		}
		if approved.Data.Status != "APPROVED" {
			t.Fatalf("a single valid approval must complete the request: %s", string(body))
		}
		step := stepAt(approved, 1)
		if step.Status != "APPROVED" || step.Actionable {
			t.Fatalf("the completed step must not stay actionable: %s", string(body))
		}
		decided := assignmentFor(step, eApprover2)
		if decided == nil || decided.Status != "APPROVED" || decided.DecidedAt == nil || decided.Comment != "looks good" {
			t.Fatalf("the decision must be recorded on the acting candidate: %s", string(body))
		}
		// BD-E1: the other candidates of a completed level are blocked, and their
		// rows stay UNDECIDED forever.
		other := assignmentFor(step, eApprover1)
		if other == nil || other.Status != "UNDECIDED" || other.DecidedAt != nil {
			t.Fatalf("a second candidate must not record anything on a completed level: %s", string(body))
		}

		if status, _, _ := decide(t, fx.tokenAppr1, number, "approve", ""); status != http.StatusConflict {
			t.Fatalf("a second approval on a completed request = %d, want %d", status, http.StatusConflict)
		}
	})

	t.Run("unassigned_admin_cannot_approve", func(t *testing.T) {
		status, created, body := create(t, fx.tokenAdminA, "purchase_request", "po-unassigned", f(100000))
		if status != http.StatusCreated {
			t.Fatalf("create = %d (%s)", status, string(body))
		}

		// eAdminA holds the management permission (Admin) but holds no assignment,
		// so permission alone must not allow approving.
		if status, _, _ := decide(t, fx.tokenAdminA, created.Data.Number, "approve", ""); status != http.StatusForbidden {
			t.Fatalf("unassigned admin approve = %d, want %d", status, http.StatusForbidden)
		}
	})

	t.Run("sequential_levels", func(t *testing.T) {
		status, created, body := create(t, fx.tokenAdminA, "purchase_request", "po-multi", f(600000))
		if status != http.StatusCreated {
			t.Fatalf("create = %d (%s)", status, string(body))
		}
		number := created.Data.Number

		if len(created.Data.Steps) != 2 {
			t.Fatalf("600000 must resolve a two-level path: %s", string(body))
		}
		level1, level2 := stepAt(created, 1), stepAt(created, 2)
		if !level1.Actionable || level2.Actionable {
			t.Fatalf("only level 1 may be actionable initially (BD-E2): %s", string(body))
		}
		if level2.Status != "PENDING" {
			t.Fatalf("a later level stays pending but not actionable: %s", string(body))
		}
		if len(level2.Anchors) != 1 || level2.Anchors[0].Type != "stakeholder_type" || level2.Anchors[0].Value != "owner" {
			t.Fatalf("level 2 must snapshot the stakeholder anchor: %s", string(body))
		}
		if assignmentFor(level2, eOwnerA) == nil {
			t.Fatalf("the active owner must be a candidate at level 2: %s", string(body))
		}

		// The level-2 approver cannot jump the queue.
		if status, _, _ := decide(t, fx.tokenOwnerA, number, "approve", ""); status != http.StatusForbidden {
			t.Fatalf("level-2 approver acting on level 1 = %d, want %d", status, http.StatusForbidden)
		}

		status, afterLevel1, body := decide(t, fx.tokenAppr1, number, "approve", "")
		if status != http.StatusOK {
			t.Fatalf("level 1 approve = %d (%s)", status, string(body))
		}
		if afterLevel1.Data.Status != "PENDING" || afterLevel1.Data.CurrentLevel != 2 {
			t.Fatalf("the request must advance to level 2, not complete: %s", string(body))
		}
		if !stepAt(afterLevel1, 2).Actionable || stepAt(afterLevel1, 1).Actionable {
			t.Fatalf("level 2 must become the only actionable step: %s", string(body))
		}

		// A level-1 candidate has no standing on level 2.
		if status, _, _ := decide(t, fx.tokenAppr2, number, "approve", ""); status != http.StatusForbidden {
			t.Fatalf("level-1 candidate acting on level 2 = %d, want %d", status, http.StatusForbidden)
		}

		status, final, body := decide(t, fx.tokenOwnerA, number, "approve", "owner sign-off")
		if status != http.StatusOK {
			t.Fatalf("level 2 approve = %d (%s)", status, string(body))
		}
		if final.Data.Status != "APPROVED" || final.Data.CurrentLevel != 2 {
			t.Fatalf("the final level must complete the request: %s", string(body))
		}
	})

	t.Run("self_approval_blocks_the_request", func(t *testing.T) {
		// eSolo is the only holder of the anchor role, so removing the requester
		// leaves no eligible approver.
		status, created, body := create(t, fx.tokenSolo, "solo_request", "solo-1", f(10))
		if status != http.StatusCreated {
			t.Fatalf("create = %d (%s)", status, string(body))
		}
		if created.Data.Status != "BLOCKED" {
			t.Fatalf("a request with no eligible approver must be BLOCKED (BD-E6/BD-E7): %s", string(body))
		}
		step := stepAt(created, 1)
		if step.Status != "BLOCKED" || step.Actionable || len(step.Assignments) != 0 {
			t.Fatalf("the blocked level must be BLOCKED with no assignments and not actionable: %s", string(body))
		}

		if status, _, _ := decide(t, fx.tokenSolo, created.Data.Number, "approve", ""); status != http.StatusConflict {
			t.Fatalf("approving a blocked request = %d, want %d", status, http.StatusConflict)
		}
		// A BLOCKED request has exactly one exit: the requester cancels it, which
		// also closes the blocked step (PENDING/BLOCKED -> CANCELLED).
		status, cancelled, body := decide(t, fx.tokenSolo, created.Data.Number, "cancel", "")
		if status != http.StatusOK {
			t.Fatalf("cancelling a blocked request = %d (%s)", status, string(body))
		}
		if cancelled.Data.Status != "CANCELLED" {
			t.Fatalf("blocked request cancel status = %s", cancelled.Data.Status)
		}
		if blocked := stepAt(cancelled, 1); blocked.Status != "CANCELLED" || blocked.Actionable {
			t.Fatalf("cancelling a request must also close a blocked step: %s", string(body))
		}
	})

	t.Run("non_authenticatable_candidate_is_not_an_approver", func(t *testing.T) {
		// eGhost holds the anchor role but has no login, so it is excluded by the
		// re-used Foundation C2 derivation.
		status, created, body := create(t, fx.tokenAdminA, "ghost_request", "ghost-1", f(10))
		if status != http.StatusCreated {
			t.Fatalf("create = %d (%s)", status, string(body))
		}
		if created.Data.Status != "BLOCKED" {
			t.Fatalf("a role holder without a usable login is not an approver: %s", string(body))
		}
		step := stepAt(created, 1)
		if len(step.Assignments) != 0 {
			t.Fatalf("no assignment may be created for a person who cannot act: %s", string(body))
		}
	})

	t.Run("rejection_is_required_reasoned_and_terminal", func(t *testing.T) {
		status, created, body := create(t, fx.tokenAdminA, "purchase_request", "po-reject", f(100000))
		if status != http.StatusCreated {
			t.Fatalf("create = %d (%s)", status, string(body))
		}
		number := created.Data.Number

		if status, _, _ := decide(t, fx.tokenAppr1, number, "reject", "   "); status != http.StatusBadRequest {
			t.Fatalf("reject without a reason = %d, want %d", status, http.StatusBadRequest)
		}

		status, rejected, body := decide(t, fx.tokenAppr1, number, "reject", "budget exhausted")
		if status != http.StatusOK {
			t.Fatalf("reject = %d (%s)", status, string(body))
		}
		if rejected.Data.Status != "REJECTED" {
			t.Fatalf("rejection must be terminal: %s", string(body))
		}
		step := stepAt(rejected, 1)
		if step.Status != "REJECTED" || step.Actionable {
			t.Fatalf("the rejected step must be closed: %s", string(body))
		}
		recorded := assignmentFor(step, eApprover1)
		if recorded == nil || recorded.Status != "REJECTED" || recorded.Comment != "budget exhausted" || recorded.DecidedAt == nil {
			t.Fatalf("the rejection must be recorded immutably with its reason: %s", string(body))
		}

		// The other candidate cannot complete or override the rejection.
		if status, _, _ := decide(t, fx.tokenAppr2, number, "approve", ""); status != http.StatusConflict {
			t.Fatalf("approving after rejection = %d, want %d", status, http.StatusConflict)
		}
		if status, _, _ := decide(t, fx.tokenAppr1, number, "reject", "again"); status != http.StatusConflict {
			t.Fatalf("re-rejecting = %d, want %d", status, http.StatusConflict)
		}
	})

	t.Run("cancellation_is_requester_only", func(t *testing.T) {
		status, created, body := create(t, fx.tokenAdminA, "purchase_request", "po-cancel", f(100000))
		if status != http.StatusCreated {
			t.Fatalf("create = %d (%s)", status, string(body))
		}
		number := created.Data.Number

		// An admin with the management permission still cannot cancel someone
		// else's request.
		if status, _, _ := decide(t, fx.tokenAdminB, number, "cancel", ""); status == http.StatusOK {
			t.Fatal("another organisation's admin must never cancel this request")
		}
		if status, _, _ := decide(t, fx.tokenAppr1, number, "cancel", ""); status != http.StatusForbidden {
			t.Fatalf("non-requester cancel = %d, want %d", status, http.StatusForbidden)
		}

		status, cancelled, body := decide(t, fx.tokenAdminA, number, "cancel", "")
		if status != http.StatusOK {
			t.Fatalf("requester cancel = %d (%s)", status, string(body))
		}
		if cancelled.Data.Status != "CANCELLED" {
			t.Fatalf("cancel status = %s", cancelled.Data.Status)
		}
		if step := stepAt(cancelled, 1); step.Status != "CANCELLED" || step.Actionable {
			t.Fatalf("pending steps must be cancelled with the request: %s", string(body))
		}
		// After cancellation the approver can no longer act.
		if status, _, _ := decide(t, fx.tokenAppr1, number, "approve", ""); status != http.StatusConflict {
			t.Fatalf("approving a cancelled request = %d, want %d", status, http.StatusConflict)
		}
	})

	t.Run("action_time_eligibility_is_rechecked", func(t *testing.T) {
		status, created, body := create(t, fx.tokenAdminA, "purchase_request", "po-eligibility", f(100000))
		if status != http.StatusCreated {
			t.Fatalf("create = %d (%s)", status, string(body))
		}
		number := created.Data.Number

		// Foundation D was already called at creation; removing the role now must
		// not rebuild the path, but must stop the now-ineligible actor (BD-E8).
		if err := fx.db.Unscoped().Where("user_number = ? AND role_number = ?", eApprover1, eRoleL1).
			Delete(&permissionsModel.UserRole{}).Error; err != nil {
			t.Fatalf("removing role: %v", err)
		}

		if status, _, _ := decide(t, fx.tokenAppr1, number, "approve", ""); status != http.StatusForbidden {
			t.Fatalf("a user whose role was removed = %d, want %d", status, http.StatusForbidden)
		}
		// The snapshot is unchanged: the other candidate can still complete it.
		status, approved, body := decide(t, fx.tokenAppr2, number, "approve", "")
		if status != http.StatusOK {
			t.Fatalf("remaining candidate approve = %d (%s)", status, string(body))
		}
		if approved.Data.Status != "APPROVED" {
			t.Fatalf("the unchanged snapshot must still complete: %s", string(body))
		}
		// The ineligible actor keeps their (undecided) assignment row as history.
		if assignmentFor(stepAt(approved, 1), eApprover1) == nil {
			t.Fatal("the snapshotted assignment must survive an eligibility change")
		}
	})

	t.Run("foundation_d_outcomes_create_no_request", func(t *testing.T) {
		// Not configured at all.
		if status, _, _ := create(t, fx.tokenAdminA, "never_configured_request", "d-1", f(10)); status != http.StatusConflict {
			t.Fatalf("unconfigured decision type = %d, want %d", status, http.StatusConflict)
		}
		// Configured, but nothing applies to this scope (unit-less request vs a
		// unit-scoped authority).
		if status, _, _ := create(t, fx.tokenAdminA, "unit_only_request", "d-2", f(10)); status != http.StatusConflict {
			t.Fatalf("no applicable authority = %d, want %d", status, http.StatusConflict)
		}
		// Above every declared limit: E must not pretend authority exists.
		if status, _, _ := create(t, fx.tokenAdminA, "capped_request", "d-3", f(5000)); status != http.StatusConflict {
			t.Fatalf("exceeding declared authority = %d, want %d", status, http.StatusConflict)
		}

		var count int64
		if err := fx.db.Model(&workflowModel.ApprovalRequest{}).
			Where("organisation_number = ? AND resource_number IN ?", eOrgA, []string{"d-1", "d-2", "d-3"}).
			Count(&count).Error; err != nil {
			t.Fatalf("counting requests: %v", err)
		}
		if count != 0 {
			t.Fatalf("no approval row may be written when D yields no path (got %d)", count)
		}
	})

	t.Run("validation_and_unknown_unit", func(t *testing.T) {
		if status, _ := request(http.MethodPost, eApprovalsBase, fx.tokenAdminA, map[string]any{
			"resource_type": "purchase_order", "resource_number": "v-1", "decision_type": "purchase_request",
		}); status != http.StatusBadRequest {
			t.Fatalf("missing module = %d, want %d", status, http.StatusBadRequest)
		}
		if status, _ := request(http.MethodPost, eApprovalsBase, fx.tokenAdminA, map[string]any{
			"module": "finance", "resource_type": "purchase_order", "resource_number": "v-2", "decision_type": "",
		}); status != http.StatusBadRequest {
			t.Fatalf("missing decision type = %d, want %d", status, http.StatusBadRequest)
		}
		if status, _ := request(http.MethodPost, eApprovalsBase, fx.tokenAdminA, map[string]any{
			"module": "finance", "resource_type": "purchase_order", "resource_number": "v-3",
			"decision_type": "purchase_request", "amount": -1,
		}); status != http.StatusBadRequest {
			t.Fatalf("negative amount = %d, want %d", status, http.StatusBadRequest)
		}
		if status, _ := request(http.MethodPost, eApprovalsBase, fx.tokenAdminA, map[string]any{
			"module": "finance", "resource_type": "purchase_order", "resource_number": "v-4",
			"decision_type": "purchase_request", "organisation_unit_number": "does-not-exist",
		}); status != http.StatusBadRequest {
			t.Fatalf("unknown unit = %d, want %d", status, http.StatusBadRequest)
		}
	})

	t.Run("organisation_isolation_and_idor", func(t *testing.T) {
		status, created, body := create(t, fx.tokenAdminA, "purchase_request", "po-isolation", f(100000))
		if status != http.StatusCreated {
			t.Fatalf("create = %d (%s)", status, string(body))
		}
		number := created.Data.Number

		if status, _ := request(http.MethodGet, eApprovalsBase+"/"+number, fx.tokenAdminB, nil); status != http.StatusNotFound {
			t.Fatalf("cross-organisation get = %d, want %d", status, http.StatusNotFound)
		}
		if status, _, _ := decide(t, fx.tokenAdminB, number, "cancel", ""); status != http.StatusNotFound {
			t.Fatalf("cross-organisation cancel = %d, want %d", status, http.StatusNotFound)
		}
		_, listBody := request(http.MethodGet, eApprovalsBase, fx.tokenAdminB, nil)
		if bytes.Contains(listBody, []byte(number)) {
			t.Fatalf("another organisation must not see this request: %s", string(listBody))
		}
	})

	t.Run("list_is_bounded_and_pageable", func(t *testing.T) {
		// The default page is bounded even when the caller supplies nothing.
		status, body := request(http.MethodGet, eApprovalsBase, fx.tokenAdminA, nil)
		if status != http.StatusOK {
			t.Fatalf("list = %d (%s)", status, string(body))
		}
		var page struct {
			Data []struct {
				Number string `json:"number"`
			} `json:"data"`
		}
		if err := json.Unmarshal(body, &page); err != nil {
			t.Fatalf("decoding list: %v", err)
		}
		if len(page.Data) > 200 {
			t.Fatalf("the default page must stay bounded, got %d rows", len(page.Data))
		}

		// limit=1 returns at most one row, and offset moves the window.
		_, firstBody := request(http.MethodGet, eApprovalsBase+"?limit=1&offset=0", fx.tokenAdminA, nil)
		_, secondBody := request(http.MethodGet, eApprovalsBase+"?limit=1&offset=1", fx.tokenAdminA, nil)
		for _, raw := range []string{string(firstBody), string(secondBody)} {
			var limited struct {
				Data []struct {
					Number string `json:"number"`
				} `json:"data"`
			}
			if err := json.Unmarshal([]byte(raw), &limited); err != nil {
				t.Fatalf("decoding paged list: %v", err)
			}
			if len(limited.Data) > 1 {
				t.Fatalf("limit=1 must return at most one row: %s", raw)
			}
		}

		// Malformed paging is a client error, not a silent default.
		if status, _ := request(http.MethodGet, eApprovalsBase+"?limit=abc", fx.tokenAdminA, nil); status != http.StatusBadRequest {
			t.Fatalf("non-numeric limit = %d, want %d", status, http.StatusBadRequest)
		}
		if status, _ := request(http.MethodGet, eApprovalsBase+"?offset=-1", fx.tokenAdminA, nil); status != http.StatusBadRequest {
			t.Fatalf("negative offset = %d, want %d", status, http.StatusBadRequest)
		}
	})

	t.Run("education_scope_isolation", func(t *testing.T) {
		statusBefore, bodyBefore := request(http.MethodGet, eEduURL, fx.tokenAdminA, nil)

		status, created, body := create(t, fx.tokenAdminA, "purchase_request", "po-scope", f(100000))
		if status != http.StatusCreated {
			t.Fatalf("scope-check create = %d (%s)", status, string(body))
		}
		decide(t, fx.tokenAdminA, created.Data.Number, "cancel", "")

		statusAfter, bodyAfter := request(http.MethodGet, eEduURL, fx.tokenAdminA, nil)
		if statusBefore != statusAfter || !bytes.Equal(bodyBefore, bodyAfter) {
			t.Fatalf("Education scoped result changed after EMS approval activity:\nbefore: %s\nafter:  %s",
				string(bodyBefore), string(bodyAfter))
		}
	})
}

type eDecisionSpec struct {
	token   string
	action  string
	comment string
}

// stepAtLevel returns the step at the given approval level, or nil.
func stepAtLevel(envelope eApprovalEnvelope, level int) *eStep {
	for i := range envelope.Data.Steps {
		if envelope.Data.Steps[i].ApprovalLevel == level {
			return &envelope.Data.Steps[i]
		}
	}
	return nil
}

// TestWorkflowApprovalsConcurrency verifies the single-writer guard against a
// real database: two eligible approvers firing at the same ANY step, and an
// approve racing a reject, must never both succeed and must never advance the
// request twice or activate the next step more than once.
func TestWorkflowApprovalsConcurrency(t *testing.T) {
	fx := setupEFixture(t)

	server := httptest.NewServer(routes.SetupRoutes())
	t.Cleanup(server.Close)
	client := server.Client()

	// post is goroutine-safe: it never calls a testing helper.
	post := func(token, path string, payload any) (int, []byte, error) {
		raw, err := json.Marshal(payload)
		if err != nil {
			return 0, nil, err
		}
		req, err := http.NewRequest(http.MethodPost, server.URL+path, bytes.NewReader(raw))
		if err != nil {
			return 0, nil, err
		}
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			return 0, nil, err
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		return resp.StatusCode, body, nil
	}

	createRequest := func(t *testing.T, resourceNumber string, amount float64) string {
		t.Helper()
		status, body, err := post(fx.tokenAdminA, eApprovalsBase, map[string]any{
			"module": "finance", "resource_type": "purchase_order",
			"resource_number": resourceNumber, "decision_type": "purchase_request", "amount": amount,
		})
		if err != nil {
			t.Fatalf("creating request: %v", err)
		}
		if status != http.StatusCreated {
			t.Fatalf("create = %d (%s)", status, string(body))
		}
		var envelope eApprovalEnvelope
		if err := json.Unmarshal(body, &envelope); err != nil {
			t.Fatalf("decoding created request: %v", err)
		}
		return envelope.Data.Number
	}

	fetch := func(t *testing.T, number string) eApprovalEnvelope {
		t.Helper()
		req, err := http.NewRequest(http.MethodGet, server.URL+eApprovalsBase+"/"+number, nil)
		if err != nil {
			t.Fatalf("building fetch: %v", err)
		}
		req.Header.Set("Authorization", "Bearer "+fx.tokenAdminA)
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("fetching request: %v", err)
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("fetch = %d (%s)", resp.StatusCode, string(body))
		}
		var envelope eApprovalEnvelope
		if err := json.Unmarshal(body, &envelope); err != nil {
			t.Fatalf("decoding fetched request: %v", err)
		}
		return envelope
	}

	// race releases both decisions at the same instant and returns their statuses.
	race := func(t *testing.T, number string, first, second eDecisionSpec) (int, int) {
		t.Helper()

		start := make(chan struct{})
		statuses := make([]int, 2)
		errs := make([]error, 2)

		var wait sync.WaitGroup
		fire := func(index int, spec eDecisionSpec) {
			defer wait.Done()
			<-start
			statuses[index], _, errs[index] = post(spec.token,
				eApprovalsBase+"/"+number+"/"+spec.action, map[string]any{"comment": spec.comment})
		}

		wait.Add(2)
		go fire(0, first)
		go fire(1, second)
		close(start)
		wait.Wait()

		for index, err := range errs {
			if err != nil {
				t.Fatalf("concurrent decision %d failed: %v", index, err)
			}
		}
		return statuses[0], statuses[1]
	}

	countDecisions := func(envelope eApprovalEnvelope, status string) int {
		count := 0
		for _, step := range envelope.Data.Steps {
			for _, assignment := range step.Assignments {
				if assignment.Status == status {
					count++
				}
			}
		}
		return count
	}

	t.Run("simultaneous_approvals_single_level", func(t *testing.T) {
		// Both approvers hold the level-1 role, so both are equally actionable
		// (BD-E5 ANY). Exactly one may win.
		for attempt := 0; attempt < 5; attempt++ {
			number := createRequest(t, fmt.Sprintf("po-race-single-%d", attempt), 100000)

			first, second := race(t, number,
				eDecisionSpec{token: fx.tokenAppr1, action: "approve", comment: "one"},
				eDecisionSpec{token: fx.tokenAppr2, action: "approve", comment: "two"})

			switch {
			case first == http.StatusOK && second == http.StatusConflict:
			case second == http.StatusOK && first == http.StatusConflict:
			default:
				t.Fatalf("attempt %d: expected exactly one 200 and one 409, got %d and %d", attempt, first, second)
			}

			state := fetch(t, number)
			t.Logf("single-level race attempt %d: statuses = %d/%d, final request status = %s", attempt, first, second, state.Data.Status)
			if state.Data.Status != "APPROVED" {
				t.Fatalf("attempt %d: request status = %s, want APPROVED", attempt, state.Data.Status)
			}
			if decided := countDecisions(state, "APPROVED"); decided != 1 {
				t.Fatalf("attempt %d: %d successful decisions on one ANY step, want exactly 1", attempt, decided)
			}
			if undecided := countDecisions(state, "UNDECIDED"); undecided != 1 {
				t.Fatalf("attempt %d: %d undecided candidates remain, want 1", attempt, undecided)
			}
			if step := stepAtLevel(state, 1); step.Status != "APPROVED" || step.Actionable {
				t.Fatalf("attempt %d: the step must be completed exactly once and not actionable: %s", attempt, step.Status)
			}
		}
	})

	t.Run("simultaneous_approvals_do_not_advance_twice", func(t *testing.T) {
		// A two-level path: the race happens at level 1, so a double advance would
		// move the request past level 2 or skip activation of the next step.
		for attempt := 0; attempt < 3; attempt++ {
			number := createRequest(t, fmt.Sprintf("po-race-multi-%d", attempt), 600000)

			first, second := race(t, number,
				eDecisionSpec{token: fx.tokenAppr1, action: "approve", comment: "one"},
				eDecisionSpec{token: fx.tokenAppr2, action: "approve", comment: "two"})

			successes := 0
			for _, status := range []int{first, second} {
				if status == http.StatusOK {
					successes++
				} else if status != http.StatusConflict && status != http.StatusForbidden {
					t.Fatalf("attempt %d: loser status = %d, want 409 (stale step) or 403 (not an approver of the advanced step)", attempt, status)
				}
			}
			if successes != 1 {
				t.Fatalf("attempt %d: %d simultaneous approvals succeeded, want exactly 1", attempt, successes)
			}

			state := fetch(t, number)
			t.Logf("multi-level race attempt %d: statuses = %d/%d, current_level = %d", attempt, first, second, state.Data.CurrentLevel)
			if state.Data.CurrentLevel != 2 {
				t.Fatalf("attempt %d: current_level = %d, want 2 (advanced exactly once)", attempt, state.Data.CurrentLevel)
			}
			if state.Data.Status != "PENDING" {
				t.Fatalf("attempt %d: request status = %s, want PENDING at level 2", attempt, state.Data.Status)
			}
			if level1 := stepAtLevel(state, 1); level1.Status != "APPROVED" {
				t.Fatalf("attempt %d: level 1 status = %s, want APPROVED", attempt, level1.Status)
			}
			if level2 := stepAtLevel(state, 2); level2.Status != "PENDING" || !level2.Actionable {
				t.Fatalf("attempt %d: level 2 must be activated exactly once: status=%s actionable=%v", attempt, level2.Status, level2.Actionable)
			}
			if decided := countDecisions(state, "APPROVED"); decided != 1 {
				t.Fatalf("attempt %d: %d successful decisions, want exactly 1", attempt, decided)
			}
		}
	})

	t.Run("simultaneous_approve_and_reject", func(t *testing.T) {
		for attempt := 0; attempt < 3; attempt++ {
			number := createRequest(t, fmt.Sprintf("po-race-mixed-%d", attempt), 100000)

			approveStatus, rejectStatus := race(t, number,
				eDecisionSpec{token: fx.tokenAppr1, action: "approve", comment: "yes"},
				eDecisionSpec{token: fx.tokenAppr2, action: "reject", comment: "no"})

			successes := 0
			for _, status := range []int{approveStatus, rejectStatus} {
				if status == http.StatusOK {
					successes++
				} else if status != http.StatusConflict {
					t.Fatalf("attempt %d: loser status = %d, want 409", attempt, status)
				}
			}
			if successes != 1 {
				t.Fatalf("attempt %d: %d decisions succeeded, want exactly 1", attempt, successes)
			}

			state := fetch(t, number)
			t.Logf("approve-vs-reject attempt %d: approve=%d reject=%d, final status = %s", attempt, approveStatus, rejectStatus, state.Data.Status)
			if approveStatus == http.StatusOK && state.Data.Status != "APPROVED" {
				t.Fatalf("attempt %d: approve won but status = %s", attempt, state.Data.Status)
			}
			if rejectStatus == http.StatusOK && state.Data.Status != "REJECTED" {
				t.Fatalf("attempt %d: reject won but status = %s", attempt, state.Data.Status)
			}
			decided := countDecisions(state, "APPROVED") + countDecisions(state, "REJECTED")
			if decided != 1 {
				t.Fatalf("attempt %d: %d decisions recorded, want exactly 1", attempt, decided)
			}
		}
	})
}
