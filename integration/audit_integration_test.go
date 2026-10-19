package integration

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/routes"
	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	authServices "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/services"
	assignmentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/assignments/models"
	auditModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/audit/models"
	authorityModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/authority/models"
	organisationModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/organisation/models"
	stakeholderModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/stakeholders/models"
	workflowModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/workflow/models"
	eosModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/eos/models"
	permissionsModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/permissions/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/internal/testutil"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/migrations"
	"gorm.io/gorm"
)

const (
	g2Actor        = "ems-g2-actor"
	g2Person       = "ems-g2-person"
	g2RoleApprover = "ems-g2-role-approver"
	g2OrgA         = "ems-g2-org-a"
	g2OrgB         = "ems-g2-org-b"
)

type g2Fixture struct {
	db          *gorm.DB
	token       string
	personToken string
}

func setupAuditIntegrationFixture(t *testing.T) *g2Fixture {
	t.Helper()

	db := testutil.RequireTestDB(t)
	if err := migrations.Init(); err != nil {
		t.Fatalf("running migrations: %v", err)
	}
	t.Setenv("JWT_SECRET", "ems-foundation-g-integration-test-secret")

	users := []string{g2Actor, g2Person}
	organisations := []string{g2OrgA, g2OrgB}

	cleanup := func() {
		db.Unscoped().Where("organisation_number IN ?", organisations).Delete(&auditModel.AuditLog{})
		db.Unscoped().Where("organisation_number IN ?", organisations).Delete(&workflowModel.ApprovalStepAnchor{})
		db.Unscoped().Where("organisation_number IN ?", organisations).Delete(&workflowModel.ApprovalAssignment{})
		db.Unscoped().Where("organisation_number IN ?", organisations).Delete(&workflowModel.ApprovalStep{})
		db.Unscoped().Where("organisation_number IN ?", organisations).Delete(&workflowModel.ApprovalRequest{})
		db.Unscoped().Where("organisation_number IN ?", organisations).Delete(&authorityModel.DecisionAuthority{})
		db.Unscoped().Where("organisation_number IN ?", organisations).Delete(&assignmentModel.UserAssignment{})
		db.Unscoped().Where("organisation_number IN ?", organisations).Delete(&stakeholderModel.Stakeholder{})
		db.Unscoped().Where("organisation_number IN ?", organisations).Delete(&organisationModel.OrganisationUnit{})
		db.Unscoped().Where("role_number = ?", g2RoleApprover).Delete(&permissionsModel.RolePermission{})
		db.Unscoped().Where("role_number = ?", g2RoleApprover).Delete(&permissionsModel.UserRole{})
		db.Unscoped().Where("number = ?", g2RoleApprover).Delete(&permissionsModel.Role{})
		db.Unscoped().Where("user_number IN ?", users).Delete(&authModel.Contact{})
		db.Unscoped().Where("user_number IN ?", users).Delete(&authModel.UserPassword{})
		db.Unscoped().Where("user_number IN ?", users).Delete(&eosModel.ActiveUserEntity{})
		db.Unscoped().Where("number IN ?", organisations).Delete(&eosModel.Entity{})
		db.Unscoped().Where("number IN ?", users).Delete(&authModel.User{})
	}
	cleanup()
	t.Cleanup(cleanup)

	// An Admin actor bypasses A–E permissions; audit hooks do not rely on
	// permissions, only on the authenticated session actor and organisation.
	actor := authModel.User{Number: g2Actor, FirstName: "G2", SecondName: "Actor", Gender: "M", Type: "Admin", Status: "ACTIVE"}
	if err := db.Create(&actor).Error; err != nil {
		t.Fatalf("seeding actor: %v", err)
	}
	person := authModel.User{Number: g2Person, FirstName: "G2", SecondName: "Person", Gender: "F", Type: "User", Status: "Active"}
	if err := db.Create(&person).Error; err != nil {
		t.Fatalf("seeding person: %v", err)
	}
	// g2Person is an eligible approver: Foundation C2 requires a usable login.
	personContact := authModel.Contact{Contact: g2Person + "@example.com", ContactType: "Email", UserNumber: g2Person}
	if err := db.Omit("User").Create(&personContact).Error; err != nil {
		t.Fatalf("seeding person contact: %v", err)
	}
	personPassword := authModel.UserPassword{UserNumber: g2Person, UserPassword: "seeded-hash"}
	if err := db.Omit("User").Create(&personPassword).Error; err != nil {
		t.Fatalf("seeding person password: %v", err)
	}

	for _, number := range organisations {
		entity := eosModel.Entity{Number: number, Name: number, AddedBy: g2Actor}
		if err := db.Create(&entity).Error; err != nil {
			t.Fatalf("seeding entity %s: %v", number, err)
		}
	}
	for _, spec := range []struct{ userNumber, entityNumber string }{
		{g2Actor, g2OrgA},
		{g2Person, g2OrgA},
	} {
		active := eosModel.ActiveUserEntity{UserNumber: spec.userNumber, EntityNumber: spec.entityNumber}
		if err := db.Omit("Entity", "User").Create(&active).Error; err != nil {
			t.Fatalf("setting active entity for %s: %v", spec.userNumber, err)
		}
		if err := db.Omit("Entity", "User", "AddedBy").Create(&eosModel.EntityUser{
			EntityNumber:  spec.entityNumber,
			UserNumber:    spec.userNumber,
			AddedByNumber: spec.userNumber,
		}).Error; err != nil {
			t.Fatalf("seeding membership for %s: %v", spec.userNumber, err)
		}
	}

	// g2Person is an eligible approver, so needs the approval action permission.
	role := permissionsModel.Role{Number: g2RoleApprover, Name: g2RoleApprover, EntityNumber: g2OrgA, Status: "ACTIVE", UserNumber: g2Actor}
	if err := db.Omit("User", "Entity", "Permissions").Create(&role).Error; err != nil {
		t.Fatalf("seeding approver role: %v", err)
	}
	rolePermission := permissionsModel.RolePermission{RoleNumber: g2RoleApprover, Permission: "act_ems_approval", AddedByNumber: g2Actor}
	if err := db.Omit("AddedBy", "Role").Create(&rolePermission).Error; err != nil {
		t.Fatalf("seeding approver permission: %v", err)
	}
	userRole := permissionsModel.UserRole{UserNumber: g2Person, RoleNumber: g2RoleApprover, ExpiryDate: time.Now().Add(365 * 24 * time.Hour), AddedBy: g2Actor}
	if err := db.Omit("User", "Role", "AddedByUser").Create(&userRole).Error; err != nil {
		t.Fatalf("assigning approver role: %v", err)
	}

	token, err := authServices.JWTAuthService().GenerateToken(g2Actor)
	if err != nil {
		t.Fatalf("generating token: %v", err)
	}
	personToken, err := authServices.JWTAuthService().GenerateToken(g2Person)
	if err != nil {
		t.Fatalf("generating person token: %v", err)
	}
	return &g2Fixture{db: db, token: token, personToken: personToken}
}

type g2NumberEnvelope struct {
	Data struct {
		Number string `json:"number"`
	} `json:"data"`
}

// TestFoundationGAuditIntegration covers Foundation G Increment 2: additive,
// non-fatal audit hooks on the A–E business actions, with actor and
// organisation attribution and no business-semantic change.
func TestFoundationGAuditIntegration(t *testing.T) {
	fx := setupAuditIntegrationFixture(t)

	server := httptest.NewServer(routes.SetupRoutes())
	t.Cleanup(server.Close)
	client := server.Client()

	requestAs := func(token, method, path string, payload any) (int, []byte) {
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
		req.Header.Set("Authorization", "Bearer "+token)
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
	request := func(method, path string, payload any) (int, []byte) {
		return requestAs(fx.token, method, path, payload)
	}

	hasFact := func(action, entityType, entityID, actor, org string) bool {
		var count int64
		fx.db.Model(&auditModel.AuditLog{}).
			Where("action = ? AND entity_type = ? AND entity_id = ? AND actor_user_number = ? AND organisation_number = ?",
				action, entityType, entityID, actor, org).
			Count(&count)
		return count > 0
	}

	var unitNumber string

	t.Run("module_enable_disable", func(t *testing.T) {
		if status, body := request(http.MethodPatch, "/main/ems/modules/education/disable", nil); status != http.StatusOK {
			t.Fatalf("disable = %d (%s)", status, string(body))
		}
		if status, body := request(http.MethodPatch, "/main/ems/modules/education/enable", nil); status != http.StatusOK {
			t.Fatalf("enable = %d (%s)", status, string(body))
		}
		if !hasFact("module.disable", "ems_module", "education", g2Actor, g2OrgA) {
			t.Fatal("module.disable fact missing")
		}
		if !hasFact("module.enable", "ems_module", "education", g2Actor, g2OrgA) {
			t.Fatal("module.enable fact missing")
		}
	})

	t.Run("organisation_unit_lifecycle", func(t *testing.T) {
		status, body := request(http.MethodPost, "/main/ems/organisation-units",
			map[string]any{"name": "G2 Unit", "unit_type": "school", "code": "g2-u1"})
		if status != http.StatusCreated {
			t.Fatalf("create unit = %d (%s)", status, string(body))
		}
		var created g2NumberEnvelope
		if err := json.Unmarshal(body, &created); err != nil {
			t.Fatalf("decoding unit: %v", err)
		}
		unitNumber = created.Data.Number
		if !hasFact("organisation_unit.create", "ems_organisation_unit", unitNumber, g2Actor, g2OrgA) {
			t.Fatal("organisation_unit.create fact missing")
		}

		if status, body := request(http.MethodPut, "/main/ems/organisation-units/"+unitNumber,
			map[string]any{"name": "G2 Unit Renamed", "unit_type": "school", "code": "g2-u1"}); status != http.StatusOK {
			t.Fatalf("update unit = %d (%s)", status, string(body))
		}
		if !hasFact("organisation_unit.update", "ems_organisation_unit", unitNumber, g2Actor, g2OrgA) {
			t.Fatal("organisation_unit.update fact missing")
		}

		if status, body := request(http.MethodPatch, "/main/ems/organisation-units/"+unitNumber+"/status",
			map[string]any{"status": "INACTIVE"}); status != http.StatusOK {
			t.Fatalf("status unit inactive = %d (%s)", status, string(body))
		}
		if status, body := request(http.MethodPatch, "/main/ems/organisation-units/"+unitNumber+"/status",
			map[string]any{"status": "ACTIVE"}); status != http.StatusOK {
			t.Fatalf("status unit active = %d (%s)", status, string(body))
		}
		if !hasFact("organisation_unit.status", "ems_organisation_unit", unitNumber, g2Actor, g2OrgA) {
			t.Fatal("organisation_unit.status fact missing")
		}

		// A separate unit exercises the delete hook without removing unitNumber.
		status, body = request(http.MethodPost, "/main/ems/organisation-units",
			map[string]any{"name": "G2 Temp", "unit_type": "school", "code": "g2-u2"})
		if status != http.StatusCreated {
			t.Fatalf("create temp unit = %d (%s)", status, string(body))
		}
		var temp g2NumberEnvelope
		if err := json.Unmarshal(body, &temp); err != nil {
			t.Fatalf("decoding temp unit: %v", err)
		}
		if status, body := request(http.MethodDelete, "/main/ems/organisation-units/"+temp.Data.Number, nil); status != http.StatusOK {
			t.Fatalf("delete unit = %d (%s)", status, string(body))
		}
		if !hasFact("organisation_unit.delete", "ems_organisation_unit", temp.Data.Number, g2Actor, g2OrgA) {
			t.Fatal("organisation_unit.delete fact missing")
		}
	})

	var stakeholderNumber string

	t.Run("stakeholder_lifecycle", func(t *testing.T) {
		status, body := request(http.MethodPost, "/main/ems/stakeholders",
			map[string]any{"user_number": g2Person, "stakeholder_type": "director"})
		if status != http.StatusCreated {
			t.Fatalf("create stakeholder = %d (%s)", status, string(body))
		}
		var created g2NumberEnvelope
		if err := json.Unmarshal(body, &created); err != nil {
			t.Fatalf("decoding stakeholder: %v", err)
		}
		stakeholderNumber = created.Data.Number
		if !hasFact("stakeholder.create", "ems_stakeholder", stakeholderNumber, g2Actor, g2OrgA) {
			t.Fatal("stakeholder.create fact missing")
		}

		if status, body := request(http.MethodPut, "/main/ems/stakeholders/"+stakeholderNumber,
			map[string]any{"title": "Director"}); status != http.StatusOK {
			t.Fatalf("update stakeholder = %d (%s)", status, string(body))
		}
		if !hasFact("stakeholder.update", "ems_stakeholder", stakeholderNumber, g2Actor, g2OrgA) {
			t.Fatal("stakeholder.update fact missing")
		}

		if status, body := request(http.MethodPatch, "/main/ems/stakeholders/"+stakeholderNumber+"/status",
			map[string]any{"status": "INACTIVE"}); status != http.StatusOK {
			t.Fatalf("stakeholder inactive = %d (%s)", status, string(body))
		}
		if status, body := request(http.MethodPatch, "/main/ems/stakeholders/"+stakeholderNumber+"/status",
			map[string]any{"status": "ACTIVE"}); status != http.StatusOK {
			t.Fatalf("stakeholder active = %d (%s)", status, string(body))
		}
		if !hasFact("stakeholder.status", "ems_stakeholder", stakeholderNumber, g2Actor, g2OrgA) {
			t.Fatal("stakeholder.status fact missing")
		}
	})

	var authorityNumber string

	t.Run("authority_lifecycle", func(t *testing.T) {
		status, body := request(http.MethodPost, "/main/ems/decision-authorities",
			map[string]any{"decision_type": "purchase", "approval_level": 1, "min_amount": 0, "max_amount": 1000, "role_number": g2RoleApprover})
		if status != http.StatusCreated {
			t.Fatalf("create authority = %d (%s)", status, string(body))
		}
		var created g2NumberEnvelope
		if err := json.Unmarshal(body, &created); err != nil {
			t.Fatalf("decoding authority: %v", err)
		}
		authorityNumber = created.Data.Number
		if !hasFact("authority.create", "ems_decision_authority", authorityNumber, g2Actor, g2OrgA) {
			t.Fatal("authority.create fact missing")
		}

		if status, body := request(http.MethodPut, "/main/ems/decision-authorities/"+authorityNumber,
			map[string]any{"decision_type": "purchase", "approval_level": 1, "min_amount": 0, "max_amount": 2000, "role_number": g2RoleApprover}); status != http.StatusOK {
			t.Fatalf("update authority = %d (%s)", status, string(body))
		}
		if !hasFact("authority.update", "ems_decision_authority", authorityNumber, g2Actor, g2OrgA) {
			t.Fatal("authority.update fact missing")
		}
	})

	t.Run("assignment_lifecycle", func(t *testing.T) {
		status, body := request(http.MethodPost, "/main/ems/user-assignments",
			map[string]any{"user_number": g2Person, "organisation_unit_number": unitNumber, "assignment_type": "staff"})
		if status != http.StatusCreated {
			t.Fatalf("create assignment = %d (%s)", status, string(body))
		}
		var created g2NumberEnvelope
		if err := json.Unmarshal(body, &created); err != nil {
			t.Fatalf("decoding assignment: %v", err)
		}
		assignmentNumber := created.Data.Number
		if !hasFact("assignment.create", "ems_user_assignment", assignmentNumber, g2Actor, g2OrgA) {
			t.Fatal("assignment.create fact missing")
		}

		if status, body := request(http.MethodPut, "/main/ems/user-assignments/"+assignmentNumber,
			map[string]any{"organisation_unit_number": unitNumber, "assignment_type": "manager"}); status != http.StatusOK {
			t.Fatalf("update assignment = %d (%s)", status, string(body))
		}
		if !hasFact("assignment.update", "ems_user_assignment", assignmentNumber, g2Actor, g2OrgA) {
			t.Fatal("assignment.update fact missing")
		}

		if status, body := request(http.MethodPatch, "/main/ems/user-assignments/"+assignmentNumber+"/status",
			map[string]any{"status": "ENDED"}); status != http.StatusOK {
			t.Fatalf("assignment ended = %d (%s)", status, string(body))
		}
		if !hasFact("assignment.status", "ems_user_assignment", assignmentNumber, g2Actor, g2OrgA) {
			t.Fatal("assignment.status fact missing")
		}
	})

	t.Run("workflow_transition", func(t *testing.T) {
		status, body := request(http.MethodPost, "/main/ems/approvals",
			map[string]any{"module": "ems", "resource_type": "g2", "resource_number": "1", "decision_type": "purchase", "amount": 500})
		if status != http.StatusCreated {
			t.Fatalf("create approval = %d (%s)", status, string(body))
		}
		var created struct {
			Data struct {
				Number string `json:"number"`
				Status string `json:"status"`
			} `json:"data"`
		}
		if err := json.Unmarshal(body, &created); err != nil {
			t.Fatalf("decoding approval: %v", err)
		}
		approvalNumber := created.Data.Number
		if !hasFact("approval.create", "ems_approval_request", approvalNumber, g2Actor, g2OrgA) {
			t.Fatal("approval.create fact missing")
		}

		// The eligible approver (g2Person, holding the approval role) approves.
		if status, approveBody := requestAs(fx.personToken, http.MethodPost, "/main/ems/approvals/"+approvalNumber+"/approve",
			map[string]any{"comment": "ok"}); status != http.StatusOK {
			t.Fatalf("approve = %d (%s)", status, string(approveBody))
		}
		if !hasFact("approval.approve", "ems_approval_request", approvalNumber, g2Person, g2OrgA) {
			t.Fatal("approval.approve fact missing")
		}
	})
}
