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
	licensingModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/licensing/models"
	eosModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/eos/models"
	permissionsModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/permissions/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/internal/testutil"
	uuid "github.com/satori/go.uuid"
)

const (
	// Unique prefixes to avoid collision with other test fixtures.
	bUserNoOrg      = "ems-bootstrap-no-org"
	bUserWithOrg    = "ems-bootstrap-with-org"
	bOrgWithLicence = "ems-bootstrap-org-licence"
	bRegistryUser   = "ems-bootstrap-registry"
	bRegistryRole   = "ems-bootstrap-role-registry"
)

// TestBootstrapLicenceGateExemption verifies the narrowly scoped exemption
// for first-time organisation bootstrap endpoints (POST entity/add,
// GET my-entities, POST entity/current). It exercises A–F from the
// diagnostic report.
func TestBootstrapLicenceGateExemption(t *testing.T) {
	// Skip if no test DB configured.
	db := testutil.RequireTestDB(t)
	t.Setenv("JWT_SECRET", "ems-bootstrap-test-secret")

	// --- Cleanup helpers ---------------------------------------------------
	// FK-safe teardown: dependent rows (role_permissions, user_roles, roles,
	// licensing facts) are removed before their parents so the test is
	// repeatable and never leaves fixture residue behind.
	clientRoleNumber := bUserWithOrg + "-client"
	cleanup := func() {
		db.Unscoped().Where("role_number IN ?", []string{bRegistryRole, clientRoleNumber}).Delete(&permissionsModel.RolePermission{})
		db.Unscoped().Where("user_number IN ?", []string{bUserNoOrg, bUserWithOrg, bRegistryUser}).Delete(&permissionsModel.UserRole{})
		db.Unscoped().Where("number IN ?", []string{bRegistryRole, clientRoleNumber}).Delete(&permissionsModel.Role{})
		db.Unscoped().Where("organisation_number = ?", bOrgWithLicence).Delete(&licensingModel.LicenceEvent{})
		db.Unscoped().Where("organisation_number = ?", bOrgWithLicence).Delete(&licensingModel.Licence{})
		db.Unscoped().Where("entity_number = ?", bOrgWithLicence).Delete(&eosModel.EntityUser{})
		db.Unscoped().Where("entity_number = ?", bOrgWithLicence).Delete(&eosModel.ActiveUserEntity{})
		db.Unscoped().Where("number = ?", bOrgWithLicence).Delete(&eosModel.Entity{})
		db.Unscoped().Where("user_number IN ?", []string{bUserNoOrg, bUserWithOrg, bRegistryUser}).Delete(&eosModel.ActiveUserEntity{})
		db.Unscoped().Where("number IN ?", []string{bUserNoOrg, bUserWithOrg, bRegistryUser}).Delete(&authModel.User{})
	}
	cleanup()
	t.Cleanup(cleanup)

	// --- Seed users ---------------------------------------------------------
	users := []authModel.User{
		{Number: bUserNoOrg, FirstName: "Bootstrap", SecondName: "NoOrg", Gender: "M", Type: "User", Status: "ACTIVE"},
		{Number: bUserWithOrg, FirstName: "Bootstrap", SecondName: "WithOrg", Gender: "M", Type: "User", Status: "ACTIVE"},
		{Number: bRegistryUser, FirstName: "Bootstrap", SecondName: "Registry", Gender: "M", Type: "User", Status: "ACTIVE"},
	}
	for _, u := range users {
		if err := db.Create(&u).Error; err != nil {
			t.Fatalf("seeding user %s: %v", u.Number, err)
		}
	}

	// --- Create the organisation entity FIRST (needed for registry role FK) ---
	orgWithOrg := eosModel.Entity{
		Number: bOrgWithLicence, Name: "Bootstrap Org With User", AddedBy: bRegistryUser,
	}
	if err := db.Create(&orgWithOrg).Error; err != nil {
		t.Fatalf("seeding org: %v", err)
	}
	entityUser := eosModel.EntityUser{EntityNumber: bOrgWithLicence, UserNumber: bUserWithOrg, AddedByNumber: bRegistryUser}
	if err := db.Omit("Entity", "User").Create(&entityUser).Error; err != nil {
		t.Fatalf("seeding entity_user: %v", err)
	}
	activeUE := eosModel.ActiveUserEntity{UserNumber: bUserWithOrg, EntityNumber: bOrgWithLicence}
	if err := db.Omit("Entity", "User").Create(&activeUE).Error; err != nil {
		t.Fatalf("seeding active_user_entity: %v", err)
	}

	// --- Registry user with manage_licence_registry permission --------------
	registryRole := &permissionsModel.Role{
		Number: bRegistryRole, Name: bRegistryRole, EntityNumber: bOrgWithLicence, Status: "ACTIVE", UserNumber: bRegistryUser,
	}
	if err := db.Omit("User", "Entity", "Permissions").Create(registryRole).Error; err != nil {
		t.Fatalf("seeding registry role: %v", err)
	}
	regPerm := &permissionsModel.RolePermission{
		RoleNumber: bRegistryRole, Permission: "manage_licence_registry", AddedByNumber: bRegistryUser,
	}
	if err := db.Omit("AddedBy", "Role").Create(regPerm).Error; err != nil {
		t.Fatalf("seeding registry permission: %v", err)
	}
	expiry := time.Now().Add(365 * 24 * time.Hour)
	regUserRole := &permissionsModel.UserRole{
		UserNumber: bRegistryUser, RoleNumber: bRegistryRole, ExpiryDate: expiry, AddedBy: bRegistryUser,
	}
	if err := db.Omit("User", "Role", "AddedByUser").Create(regUserRole).Error; err != nil {
		t.Fatalf("seeding registry user role: %v", err)
	}

	// --- User WITH org needs view_curriculum permission for test E ---
	clientRole := &permissionsModel.Role{
		Number: bUserWithOrg + "-client", Name: "Client Role", EntityNumber: bOrgWithLicence, Status: "ACTIVE", UserNumber: bRegistryUser,
	}
	if err := db.Omit("User", "Entity", "Permissions").Create(clientRole).Error; err != nil {
		t.Fatalf("seeding client role: %v", err)
	}
	clientPerm := &permissionsModel.RolePermission{
		RoleNumber: clientRole.Number, Permission: "view_curriculum", AddedByNumber: bRegistryUser,
	}
	if err := db.Omit("AddedBy", "Role").Create(clientPerm).Error; err != nil {
		t.Fatalf("seeding client permission: %v", err)
	}
	clientUserRole := &permissionsModel.UserRole{
		UserNumber: bUserWithOrg, RoleNumber: clientRole.Number, ExpiryDate: expiry, AddedBy: bRegistryUser,
	}
	if err := db.Omit("User", "Role", "AddedByUser").Create(clientUserRole).Error; err != nil {
		t.Fatalf("seeding client user role: %v", err)
	}

	// --- Tokens -------------------------------------------------------------
	token := func(userNumber string) string {
		val, err := authServices.JWTAuthService().GenerateToken(userNumber)
		if err != nil {
			t.Fatalf("generating token for %s: %v", userNumber, err)
		}
		return val
	}
	tokenNoOrg := token(bUserNoOrg)
	tokenWithOrg := token(bUserWithOrg)
	tokenRegistry := token(bRegistryUser)

	// --- Server -------------------------------------------------------------
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

	licenceGateMessage := "Your EMS licence does not permit access to this module"

	// -------------------------------------------------------------------------
	// A. Authenticated user WITHOUT an organisation can reach the bootstrap
	//    endpoints without being rejected by the licence gate.
	// -------------------------------------------------------------------------
	t.Run("A_no_org_user_reaches_bootstrap_endpoints", func(t *testing.T) {
		// POST /main/eos/entity/add with INVALID payload (missing required name)
		// → must reach handler, fail validation (400), NOT licence 409.
		status, body := request(http.MethodPost, "/main/eos/entity/add", tokenNoOrg, map[string]any{
			"phone": "123", "email": "a@b.c", "address": "addr",
		})
		if status == http.StatusConflict && bytes.Contains(body, []byte(licenceGateMessage)) {
			t.Fatalf("POST entity/add blocked by licence gate: %d %s", status, string(body))
		}
		if status != http.StatusBadRequest {
			t.Fatalf("POST entity/add with invalid payload = %d, want 400 (%s)", status, string(body))
		}

		// GET /main/eos/my-entities → handler reached, returns 200 (empty list)
		status, body = request(http.MethodGet, "/main/eos/my-entities", tokenNoOrg, nil)
		if status == http.StatusConflict && bytes.Contains(body, []byte(licenceGateMessage)) {
			t.Fatalf("GET my-entities blocked by licence gate: %d %s", status, string(body))
		}
		if status != http.StatusOK {
			t.Fatalf("GET my-entities = %d, want 200 (%s)", status, string(body))
		}

		// POST /main/eos/entity/current with non-existent entity_number
		// → membership/validation error (400), NOT licence 409.
		status, body = request(http.MethodPost, "/main/eos/entity/current?entity_number="+uuid.NewV4().String(), tokenNoOrg, nil)
		if status == http.StatusConflict && bytes.Contains(body, []byte(licenceGateMessage)) {
			t.Fatalf("POST entity/current blocked by licence gate: %d %s", status, string(body))
		}
		if status != http.StatusBadRequest && status != http.StatusNotFound {
			t.Fatalf("POST entity/current non-existent = %d, want 400/404 (%s)", status, string(body))
		}
	})

	// -------------------------------------------------------------------------
	// B. Unauthenticated requests remain rejected (JWT still enforced).
	// -------------------------------------------------------------------------
	t.Run("B_unauthenticated_rejected", func(t *testing.T) {
		status, body := request(http.MethodGet, "/main/eos/my-entities", "", nil)
		if status != http.StatusUnauthorized {
			t.Fatalf("unauth GET my-entities = %d, want 401 (%s)", status, string(body))
		}
		status, body = request(http.MethodPost, "/main/eos/entity/add", "", map[string]any{
			"name": "x", "phone": "y", "email": "z", "address": "a",
		})
		if status != http.StatusUnauthorized {
			t.Fatalf("unauth POST entity/add = %d, want 401 (%s)", status, string(body))
		}
	})

	// -------------------------------------------------------------------------
	// C. Existing membership and validation checks remain effective.
	// -------------------------------------------------------------------------
	t.Run("C_validation_and_membership_enforced", func(t *testing.T) {
		// Validation: POST entity/add with missing name → 400
		status, body := request(http.MethodPost, "/main/eos/entity/add", tokenNoOrg, map[string]any{
			"phone": "123", "email": "a@b.c", "address": "addr",
		})
		if status != http.StatusBadRequest || !bytes.Contains(body, []byte("provide entity name")) {
			t.Fatalf("validation missing name = %d, want 400 with 'provide entity name' (%s)", status, string(body))
		}

		// Membership: POST entity/current with non-existent entity → 400/404
		// (handler checks FindEntityUser which requires membership)
		status, body = request(http.MethodPost, "/main/eos/entity/current?entity_number="+uuid.NewV4().String(), tokenNoOrg, nil)
		if status != http.StatusBadRequest && status != http.StatusNotFound {
			t.Fatalf("membership check = %d, want 400/404 (%s)", status, string(body))
		}
	})

	// -------------------------------------------------------------------------
	// D. Unrelated licence-protected module blocked when entitlement absent.
	//    User has an active organisation but NO licence → curriculum blocked.
	// -------------------------------------------------------------------------
	t.Run("D_unrelated_module_blocked_when_no_licence", func(t *testing.T) {
		// bUserWithOrg has active org bOrgWithLicence but no licence issued.
		status, body := request(http.MethodGet, "/main/curriculum/all", tokenWithOrg, nil)
		if status != http.StatusConflict {
			t.Fatalf("education route with no licence = %d, want 409 (%s)", status, string(body))
		}
		if !bytes.Contains(body, []byte(licenceGateMessage)) {
			t.Fatalf("409 body missing licence message: %s", string(body))
		}
	})

	// -------------------------------------------------------------------------
	// E. Existing licensed-user behaviour intact.
	//    Issue an education licence for bOrgWithLicence via registry API,
	//    then education route passes for bUserWithOrg.
	// -------------------------------------------------------------------------
	t.Run("E_licensed_user_access_restored", func(t *testing.T) {
		from := time.Now().Add(-time.Hour)
		to := time.Now().Add(30 * 24 * time.Hour)
		// Issue licence via registry (control-plane) using registry actor.
		status, body := request(http.MethodPost, flRegistryURL, tokenRegistry, map[string]any{
			"organisation_number": bOrgWithLicence,
			"valid_from":          from.Format(time.RFC3339),
			"valid_to":            to.Format(time.RFC3339),
			"max_installations":   1,
			"licensed_addons":     []string{"education"},
		})
		if status != http.StatusCreated {
			t.Fatalf("issuing licence = %d, want 201 (%s)", status, string(body))
		}

		// Now education route should be accessible for the user in that org.
		status, body = request(http.MethodGet, "/main/curriculum/all", tokenWithOrg, nil)
		if status != http.StatusOK {
			t.Fatalf("education route with entitled licence = %d, want 200 (%s)", status, string(body))
		}

		// Profile and EMS Core remain reachable (exempt).
		if status, _ := request(http.MethodGet, "/main/profile", tokenWithOrg, nil); status != http.StatusOK {
			t.Fatalf("profile = %d, want 200", status)
		}
	})

	// -------------------------------------------------------------------------
	// F. No unrelated route gains an exemption.
	//    - POST /main/eos/entity/update → 409 for no-org user
	//    - GET /main/eos/entity/current → 409 for no-org user (not exempted per scope)
	//    - GET /main/curriculum/all → 409 for no-org user (no licence, no org)
	// -------------------------------------------------------------------------
	t.Run("F_unrelated_routes_not_exempt", func(t *testing.T) {
		// entity/update requires permission but licence gate runs first.
		status, body := request(http.MethodPost, "/main/eos/entity/update", tokenNoOrg, map[string]any{
			"entity_number": "x", "name": "y", "phone": "z", "email": "a", "address": "b",
		})
		if status != http.StatusConflict || !bytes.Contains(body, []byte(licenceGateMessage)) {
			t.Fatalf("entity/update not blocked by licence gate = %d (%s)", status, string(body))
		}

		// GET entity/current (only POST is exempted)
		status, body = request(http.MethodGet, "/main/eos/entity/current", tokenNoOrg, nil)
		if status != http.StatusConflict || !bytes.Contains(body, []byte(licenceGateMessage)) {
			t.Fatalf("GET entity/current not blocked by licence gate = %d (%s)", status, string(body))
		}

		// curriculum/all (education module) for no-org user
		status, body = request(http.MethodGet, "/main/curriculum/all", tokenNoOrg, nil)
		if status != http.StatusConflict || !bytes.Contains(body, []byte(licenceGateMessage)) {
			t.Fatalf("curriculum/all not blocked by licence gate = %d (%s)", status, string(body))
		}
	})
}
