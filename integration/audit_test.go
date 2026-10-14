package integration

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/routes"
	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	authServices "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/services"
	auditModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/audit/models"
	auditServices "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/audit/services"
	eosModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/eos/models"
	permissionsModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/permissions/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/internal/testutil"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/migrations"
	"gorm.io/gorm"
)

const (
	auOwner       = "ems-au-owner"
	auAdminNoPerm = "ems-au-admin-noperm"
	auClientA     = "ems-au-client-a"
	auClientB     = "ems-au-client-b"
	auNoPerm      = "ems-au-noperm"
	auNoCtx       = "ems-au-nocontext"

	auOrgA = "ems-au-org-a"
	auOrgB = "ems-au-org-b"

	auRoleRegistry = "ems-au-role-registry"
	auRoleView     = "ems-au-role-view"

	auAuditURL    = "/main/ems/audit"
	auRegistryURL = "/main/ems/audit-registry"
)

type auFixture struct {
	db          *gorm.DB
	tokenOwner  string
	tokenAdmin  string
	tokenA      string
	tokenB      string
	tokenNoPerm string
	tokenNoCtx  string
}

func setupAuditFixture(t *testing.T) *auFixture {
	t.Helper()

	db := testutil.RequireTestDB(t)
	if err := migrations.Init(); err != nil {
		t.Fatalf("running migrations: %v", err)
	}
	t.Setenv("JWT_SECRET", "ems-foundation-g-audit-test-secret")

	users := []string{auOwner, auAdminNoPerm, auClientA, auClientB, auNoPerm, auNoCtx}
	organisations := []string{auOrgA, auOrgB}
	roles := []string{auRoleRegistry, auRoleView}

	cleanup := func() {
		db.Unscoped().Where("organisation_number IN ?", organisations).Delete(&auditModel.AuditLog{})
		db.Unscoped().Where("role_number IN ?", roles).Delete(&permissionsModel.RolePermission{})
		db.Unscoped().Where("user_number IN ?", users).Delete(&permissionsModel.UserRole{})
		db.Unscoped().Where("number IN ?", roles).Delete(&permissionsModel.Role{})
		db.Unscoped().Where("user_number IN ?", users).Delete(&eosModel.ActiveUserEntity{})
		db.Unscoped().Where("number IN ?", organisations).Delete(&eosModel.Entity{})
		db.Unscoped().Where("number IN ?", users).Delete(&authModel.User{})
	}
	cleanup()
	t.Cleanup(cleanup)

	for _, spec := range []struct{ number, userType string }{
		{auOwner, "User"},
		{auAdminNoPerm, "Admin"},
		{auClientA, "User"},
		{auClientB, "User"},
		{auNoPerm, "User"},
		{auNoCtx, "User"},
	} {
		user := authModel.User{
			Number:     spec.number,
			FirstName:  "FoundationG",
			SecondName: "Audit",
			Gender:     "M",
			Type:       spec.userType,
			Status:     "ACTIVE",
		}
		if err := db.Create(&user).Error; err != nil {
			t.Fatalf("seeding user %s: %v", spec.number, err)
		}
	}

	for _, spec := range []struct{ number, name string }{
		{auOrgA, "Foundation G Audit Org A"},
		{auOrgB, "Foundation G Audit Org B"},
	} {
		entity := eosModel.Entity{Number: spec.number, Name: spec.name, AddedBy: auOwner}
		if err := db.Create(&entity).Error; err != nil {
			t.Fatalf("seeding entity %s: %v", spec.number, err)
		}
	}

	for _, spec := range []struct{ userNumber, entityNumber string }{
		{auOwner, auOrgA},
		{auAdminNoPerm, auOrgA},
		{auClientA, auOrgA},
		{auClientB, auOrgB},
		{auNoPerm, auOrgA},
		// auNoCtx intentionally has no active entity.
	} {
		active := eosModel.ActiveUserEntity{UserNumber: spec.userNumber, EntityNumber: spec.entityNumber}
		if err := db.Omit("Entity", "User").Create(&active).Error; err != nil {
			t.Fatalf("setting active entity for %s: %v", spec.userNumber, err)
		}
	}

	for _, number := range roles {
		role := permissionsModel.Role{
			Number: number, Name: number, EntityNumber: auOrgA, Status: "ACTIVE", UserNumber: auOwner,
		}
		if err := db.Omit("User", "Entity", "Permissions").Create(&role).Error; err != nil {
			t.Fatalf("seeding role %s: %v", number, err)
		}
	}

	for _, spec := range []struct{ roleNumber, permission string }{
		{auRoleRegistry, "manage_ems_audit"},
		{auRoleView, "view_ems_audit"},
	} {
		rolePermission := permissionsModel.RolePermission{
			RoleNumber: spec.roleNumber, Permission: spec.permission, AddedByNumber: auOwner,
		}
		if err := db.Omit("AddedBy", "Role").Create(&rolePermission).Error; err != nil {
			t.Fatalf("seeding permission %s: %v", spec.permission, err)
		}
	}

	expiry := time.Now().Add(365 * 24 * time.Hour)
	for _, spec := range []struct{ userNumber, roleNumber string }{
		{auOwner, auRoleRegistry},
		{auClientA, auRoleView},
		{auClientB, auRoleView},
		{auNoCtx, auRoleView},
	} {
		userRole := permissionsModel.UserRole{
			UserNumber: spec.userNumber, RoleNumber: spec.roleNumber, ExpiryDate: expiry, AddedBy: auOwner,
		}
		if err := db.Omit("User", "Role", "AddedByUser").Create(&userRole).Error; err != nil {
			t.Fatalf("seeding user role for %s: %v", spec.userNumber, err)
		}
	}

	token := func(userNumber string) string {
		value, err := authServices.JWTAuthService().GenerateToken(userNumber)
		if err != nil {
			t.Fatalf("generating token for %s: %v", userNumber, err)
		}
		return value
	}

	return &auFixture{
		db:          db,
		tokenOwner:  token(auOwner),
		tokenAdmin:  token(auAdminNoPerm),
		tokenA:      token(auClientA),
		tokenB:      token(auClientB),
		tokenNoPerm: token(auNoPerm),
		tokenNoCtx:  token(auNoCtx),
	}
}

type auAuditData struct {
	Number                 uint            `json:"number"`
	OrganisationNumber     string          `json:"organisation_number"`
	ActorUserNumber        string          `json:"actor_user_number"`
	Action                 string          `json:"action"`
	EntityType             string          `json:"entity_type"`
	EntityID               string          `json:"entity_id"`
	OldValues              json.RawMessage `json:"old_values"`
	NewValues              json.RawMessage `json:"new_values"`
	IP                     string          `json:"ip"`
	UserAgent              string          `json:"user_agent"`
	CorrelationID          string          `json:"correlation_id"`
	WorkflowInstanceNumber string          `json:"workflow_instance_number"`
	CreatedAt              time.Time       `json:"created_at"`
}

type auAuditEnvelope struct {
	Message string      `json:"message"`
	Data    auAuditData `json:"data"`
}

type auAuditListEnvelope struct {
	Message string        `json:"message"`
	Data    []auAuditData `json:"data"`
	Total   int64         `json:"total"`
}

// TestFoundationGAudit covers Foundation G Increment 1: the append-only audit
// fact store, organisation-scoped client reads, the explicit cross-organisation
// audit authority (no Admin bypass), immutability and organisation isolation.
func TestFoundationGAudit(t *testing.T) {
	fx := setupAuditFixture(t)

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
			t.Fatalf("building request %s %s: %v", method, path, err)
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

	var recordAID, recordBID uint

	t.Run("audit_facts_created_and_read_back", func(t *testing.T) {
		if err := auditServices.Record(auditServices.AuditRecord{
			OrganisationNumber:     auOrgA,
			ActorUserNumber:        auClientA,
			Action:                 "stakeholder.update",
			EntityType:             "ems_stakeholder",
			EntityID:               "stake-1",
			OldValues:              map[string]interface{}{"name": "old"},
			NewValues:              map[string]interface{}{"name": "new"},
			IP:                     "10.0.0.1",
			UserAgent:              "audit-test-agent",
			CorrelationID:          "corr-1",
			WorkflowInstanceNumber: "wf-1",
		}); err != nil {
			t.Fatalf("recording org A fact: %v", err)
		}
		if err := auditServices.Record(auditServices.AuditRecord{
			OrganisationNumber: auOrgB,
			ActorUserNumber:    auClientB,
			Action:             "assignment.change",
			EntityType:         "ems_user_assignment",
			EntityID:           "assign-1",
		}); err != nil {
			t.Fatalf("recording org B fact: %v", err)
		}

		status, body := request(http.MethodGet, auAuditURL, fx.tokenA, nil)
		if status != http.StatusOK {
			t.Fatalf("org A audit list = %d (%s)", status, string(body))
		}
		var list auAuditListEnvelope
		if err := json.Unmarshal(body, &list); err != nil {
			t.Fatalf("decoding audit list: %v", err)
		}
		if list.Total != 1 || len(list.Data) != 1 {
			t.Fatalf("org A should see exactly its own fact, got total=%d len=%d", list.Total, len(list.Data))
		}
		got := list.Data[0]
		recordAID = got.Number
		if got.OrganisationNumber != auOrgA || got.ActorUserNumber != auClientA ||
			got.Action != "stakeholder.update" || got.EntityType != "ems_stakeholder" || got.EntityID != "stake-1" ||
			got.IP != "10.0.0.1" || got.UserAgent != "audit-test-agent" ||
			got.CorrelationID != "corr-1" || got.WorkflowInstanceNumber != "wf-1" || got.CreatedAt.IsZero() {
			t.Fatalf("audit fact fields not preserved: %s", string(body))
		}
		if !bytes.Contains(got.OldValues, []byte(`"name":"old"`)) || !bytes.Contains(got.NewValues, []byte(`"name":"new"`)) {
			t.Fatalf("before/after values not preserved: %s", string(body))
		}
	})

	t.Run("organisation_scoped_read_and_isolation", func(t *testing.T) {
		// Org B sees only its own fact.
		status, body := request(http.MethodGet, auAuditURL, fx.tokenB, nil)
		if status != http.StatusOK {
			t.Fatalf("org B audit list = %d (%s)", status, string(body))
		}
		var list auAuditListEnvelope
		if err := json.Unmarshal(body, &list); err != nil {
			t.Fatalf("decoding audit list: %v", err)
		}
		if list.Total != 1 {
			t.Fatalf("org B should see exactly 1 fact, got %d", list.Total)
		}
		recordBID = list.Data[0].Number
		if bytes.Contains(body, []byte(auOrgA)) {
			t.Fatalf("org B must not see org A facts: %s", string(body))
		}

		// Org A cannot read org B's fact by id.
		if status, _ := request(http.MethodGet, auAuditURL+"/"+uintString(recordBID), fx.tokenA, nil); status != http.StatusNotFound {
			t.Fatalf("cross-organisation detail read = %d, want %d", status, http.StatusNotFound)
		}
	})

	t.Run("immutable_no_update_or_delete", func(t *testing.T) {
		// No PUT/PATCH/DELETE route is registered for audit.
		registered := make(map[string]bool)
		for _, route := range routes.SetupRoutes().Routes() {
			registered[route.Method+" "+route.Path] = true
		}
		for _, forbidden := range []string{
			http.MethodPut + " " + auAuditURL + "/:number",
			http.MethodPatch + " " + auAuditURL + "/:number",
			http.MethodDelete + " " + auAuditURL + "/:number",
		} {
			if registered[forbidden] {
				t.Fatalf("audit must be immutable: route registered %s", forbidden)
			}
		}
		// The mutable methods are not served.
		for _, method := range []string{http.MethodPut, http.MethodPatch, http.MethodDelete} {
			if status, _ := request(method, auAuditURL+"/"+uintString(recordAID), fx.tokenA, nil); status == http.StatusOK || status == http.StatusNoContent {
				t.Fatalf("%s audit record unexpectedly succeeded (%d)", method, status)
			}
		}
	})

	t.Run("client_permission_and_context", func(t *testing.T) {
		if status, _ := request(http.MethodGet, auAuditURL, "", nil); status != http.StatusUnauthorized {
			t.Fatalf("unauthenticated = %d, want %d", status, http.StatusUnauthorized)
		}
		if status, _ := request(http.MethodGet, auAuditURL, fx.tokenNoPerm, nil); status != http.StatusForbidden {
			t.Fatalf("user without view_ems_audit = %d, want %d", status, http.StatusForbidden)
		}
		if status, _ := request(http.MethodGet, auAuditURL, fx.tokenNoCtx, nil); status != http.StatusConflict {
			t.Fatalf("user without active organisation = %d, want %d", status, http.StatusConflict)
		}
	})

	t.Run("control_plane_explicit_authority_no_admin_bypass", func(t *testing.T) {
		// Admin without manage_ems_audit must NOT receive the control plane.
		if status, _ := request(http.MethodGet, auRegistryURL, fx.tokenAdmin, nil); status != http.StatusForbidden {
			t.Fatalf("admin without audit authority = %d, want %d", status, http.StatusForbidden)
		}
		// A client with only view_ems_audit must not reach the control plane.
		if status, _ := request(http.MethodGet, auRegistryURL, fx.tokenA, nil); status != http.StatusForbidden {
			t.Fatalf("client audit-registry access = %d, want %d", status, http.StatusForbidden)
		}
		// The explicit authority sees across organisations.
		status, body := request(http.MethodGet, auRegistryURL, fx.tokenOwner, nil)
		if status != http.StatusOK {
			t.Fatalf("owner audit-registry = %d (%s)", status, string(body))
		}
		if !bytes.Contains(body, []byte(auOrgA)) || !bytes.Contains(body, []byte(auOrgB)) {
			t.Fatalf("control plane should see both organisations: %s", string(body))
		}
	})

	t.Run("no_secret_fields_in_audit_response", func(t *testing.T) {
		status, body := request(http.MethodGet, auAuditURL, fx.tokenA, nil)
		if status != http.StatusOK {
			t.Fatalf("audit list = %d", status)
		}
		for _, forbidden := range []string{"licence_key", "password", "token", "secret"} {
			if bytes.Contains(bytes.ToLower(body), []byte(forbidden)) {
				t.Fatalf("audit response must not expose %q: %s", forbidden, string(body))
			}
		}
	})
}

func uintString(v uint) string {
	return strconv.FormatUint(uint64(v), 10)
}
