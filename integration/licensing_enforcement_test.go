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
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/migrations"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

const (
	feOwner       = "ems-fe-owner"
	feAdminNoPerm = "ems-fe-admin-noperm"
	feClientA     = "ems-fe-client-a"
	feClientB     = "ems-fe-client-b"
	feClientC     = "ems-fe-client-c"

	feOrgA = "ems-fe-org-a"
	feOrgB = "ems-fe-org-b"
	feOrgC = "ems-fe-org-c"

	feRoleRegistry = "ems-fe-role-registry"
	feRoleClient   = "ems-fe-role-client"

	// Education route used to observe the enforcement boundary. It carries its
	// own permission, which the fixture client role holds, so the licence gate
	// is the only difference between the permitted and denied cases.
	feEducationRoute = "/main/curriculum/all"
)

type feFixture struct {
	db         *gorm.DB
	tokenOwner string
	tokenAdmin string
	tokenA     string
	tokenB     string
	tokenC     string
}

func setupFEnforcementFixture(t *testing.T) *feFixture {
	t.Helper()

	db := testutil.RequireTestDB(t)
	if err := migrations.Init(); err != nil {
		t.Fatalf("running migrations: %v", err)
	}
	t.Setenv("JWT_SECRET", "ems-foundation-f-enforcement-test-secret")

	users := []string{feOwner, feAdminNoPerm, feClientA, feClientB, feClientC}
	organisations := []string{feOrgA, feOrgB, feOrgC}
	roles := []string{feRoleRegistry, feRoleClient}

	cleanup := func() {
		db.Unscoped().Where("organisation_number IN ?", organisations).Delete(&licensingModel.LicenceEvent{})
		db.Unscoped().Where("organisation_number IN ?", organisations).Delete(&licensingModel.Installation{})
		db.Unscoped().Where("organisation_number IN ?", organisations).Delete(&licensingModel.Licence{})
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
		{feOwner, "User"},
		{feAdminNoPerm, "Admin"},
		{feClientA, "User"},
		{feClientB, "User"},
		{feClientC, "User"},
	} {
		user := authModel.User{
			Number:     spec.number,
			FirstName:  "FoundationF",
			SecondName: "Enforcement",
			Gender:     "M",
			Type:       spec.userType,
			Status:     "ACTIVE",
		}
		if err := db.Create(&user).Error; err != nil {
			t.Fatalf("seeding user %s: %v", spec.number, err)
		}
	}

	for _, spec := range []struct{ number, name string }{
		{feOrgA, "Foundation F Enforcement Org A"},
		{feOrgB, "Foundation F Enforcement Org B"},
		{feOrgC, "Foundation F Enforcement Org C"},
	} {
		entity := eosModel.Entity{Number: spec.number, Name: spec.name, AddedBy: feOwner}
		if err := db.Create(&entity).Error; err != nil {
			t.Fatalf("seeding entity %s: %v", spec.number, err)
		}
	}

	for _, spec := range []struct{ userNumber, entityNumber string }{
		{feOwner, feOrgA},
		{feAdminNoPerm, feOrgA},
		{feClientA, feOrgA},
		{feClientB, feOrgB},
		{feClientC, feOrgC},
	} {
		active := eosModel.ActiveUserEntity{UserNumber: spec.userNumber, EntityNumber: spec.entityNumber}
		if err := db.Omit("Entity", "User").Create(&active).Error; err != nil {
			t.Fatalf("setting active entity for %s: %v", spec.userNumber, err)
		}
	}

	for _, number := range roles {
		role := permissionsModel.Role{
			Number: number, Name: number, EntityNumber: feOrgA, Status: "ACTIVE", UserNumber: feOwner,
		}
		if err := db.Omit("User", "Entity", "Permissions").Create(&role).Error; err != nil {
			t.Fatalf("seeding role %s: %v", number, err)
		}
	}

	for _, spec := range []struct{ roleNumber, permission string }{
		{feRoleRegistry, "manage_licence_registry"},
		{feRoleRegistry, "view_ems_modules"},
		{feRoleClient, "view_licence_status"},
		{feRoleClient, "activate_licence_installation"},
		{feRoleClient, "deactivate_own_installation"},
		{feRoleClient, "view_curriculum"},
	} {
		rolePermission := permissionsModel.RolePermission{
			RoleNumber: spec.roleNumber, Permission: spec.permission, AddedByNumber: feOwner,
		}
		if err := db.Omit("AddedBy", "Role").Create(&rolePermission).Error; err != nil {
			t.Fatalf("seeding permission %s: %v", spec.permission, err)
		}
	}

	expiry := time.Now().Add(365 * 24 * time.Hour)
	for _, spec := range []struct{ userNumber, roleNumber string }{
		{feOwner, feRoleRegistry},
		{feClientA, feRoleClient},
		{feClientB, feRoleClient},
		{feClientC, feRoleClient},
	} {
		userRole := permissionsModel.UserRole{
			UserNumber: spec.userNumber, RoleNumber: spec.roleNumber, ExpiryDate: expiry, AddedBy: feOwner,
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

	return &feFixture{
		db:         db,
		tokenOwner: token(feOwner),
		tokenAdmin: token(feAdminNoPerm),
		tokenA:     token(feClientA),
		tokenB:     token(feClientB),
		tokenC:     token(feClientC),
	}
}

// TestFoundationFEntitlementEnforcement covers Foundation F Increment 3: the
// additive licence entitlement boundary plus the licensing exemptions.
func TestFoundationFEntitlementEnforcement(t *testing.T) {
	fx := setupFEnforcementFixture(t)

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

	issue := func(t *testing.T, org string, from, to time.Time, max int, addons []string) flLicenceData {
		t.Helper()
		status, body := request(http.MethodPost, flRegistryURL, fx.tokenOwner, map[string]any{
			"organisation_number": org,
			"valid_from":          from.Format(time.RFC3339),
			"valid_to":            to.Format(time.RFC3339),
			"max_installations":   max,
			"licensed_addons":     addons,
		})
		if status != http.StatusCreated {
			t.Fatalf("issuing licence for %s = %d, want %d (%s)", org, status, http.StatusCreated, string(body))
		}
		var envelope flLicenceEnvelope
		if err := json.Unmarshal(body, &envelope); err != nil {
			t.Fatalf("decoding licence: %v", err)
		}
		return envelope.Data
	}

	from := time.Now().Add(-time.Hour)
	to := time.Now().Add(30 * 24 * time.Hour)
	var aLicence flLicenceData

	t.Run("licensed_access_passes_through", func(t *testing.T) {
		aLicence = issue(t, feOrgA, from, to, 3, []string{"education"})

		if status, body := request(http.MethodGet, feEducationRoute, fx.tokenA, nil); status != http.StatusOK {
			t.Fatalf("education request with entitled licence = %d, want %d (%s)", status, http.StatusOK, string(body))
		}
		// Profile remains accessible (exempt).
		if status, _ := request(http.MethodGet, "/main/profile", fx.tokenA, nil); status != http.StatusOK {
			t.Fatalf("profile = %d, want %d", status, http.StatusOK)
		}
	})

	t.Run("grace_access_passes_through", func(t *testing.T) {
		issue(t, feOrgB, time.Now().Add(-48*time.Hour), time.Now().Add(-time.Hour), 1, []string{"education"})
		if status, body := request(http.MethodGet, feEducationRoute, fx.tokenB, nil); status != http.StatusOK {
			t.Fatalf("education request in grace = %d, want %d (%s)", status, http.StatusOK, string(body))
		}
	})

	t.Run("expired_denies_module_but_licensing_accessible", func(t *testing.T) {
		// Move organisation A's licence beyond its 14-day grace.
		status, body := request(http.MethodPut, flRegistryURL+"/"+aLicence.Number, fx.tokenOwner, map[string]any{
			"valid_from": time.Now().Add(-40 * 24 * time.Hour).Format(time.RFC3339),
			"valid_to":   time.Now().Add(-30 * 24 * time.Hour).Format(time.RFC3339),
		})
		if status != http.StatusOK {
			t.Fatalf("configuring expired licence = %d (%s)", status, string(body))
		}

		if status, _ := request(http.MethodGet, feEducationRoute, fx.tokenA, nil); status != http.StatusConflict {
			t.Fatalf("expired education request = %d, want %d", status, http.StatusConflict)
		}

		// Licence status stays reachable and reports EXPIRED (derived).
		status, body = request(http.MethodGet, "/main/ems/licence", fx.tokenA, nil)
		if status != http.StatusOK || !bytes.Contains(body, []byte("EXPIRED")) {
			t.Fatalf("licence status after expiry = %d (%s)", status, string(body))
		}

		// Activation stays reachable: it reaches the service and fails on the
		// expired window (not on the enforcement gate).
		status, body = request(http.MethodPost, faActivateURL, fx.tokenA, map[string]any{
			"licence_key":           aLicence.LicenceKey,
			"installation_identity": uuid.NewV4().String(),
		})
		if status != http.StatusConflict || !bytes.Contains(body, []byte("expired")) {
			t.Fatalf("activation after expiry = %d (%s), want reachable service error", status, string(body))
		}

		// Renewal stays reachable: extending validity restores access.
		status, body = request(http.MethodPost, flRegistryURL+"/"+aLicence.Number+"/renew", fx.tokenOwner, map[string]any{
			"valid_to": time.Now().Add(60 * 24 * time.Hour).Format(time.RFC3339),
		})
		if status != http.StatusOK {
			t.Fatalf("renewal after expiry = %d (%s)", status, string(body))
		}
		if status, body := request(http.MethodGet, feEducationRoute, fx.tokenA, nil); status != http.StatusOK {
			t.Fatalf("education request after renewal = %d, want %d (%s)", status, http.StatusOK, string(body))
		}
	})

	t.Run("suspended_and_revoked_deny_module", func(t *testing.T) {
		if status, body := request(http.MethodPost, flRegistryURL+"/"+aLicence.Number+"/suspend", fx.tokenOwner, nil); status != http.StatusOK {
			t.Fatalf("suspend = %d (%s)", status, string(body))
		}
		if status, _ := request(http.MethodGet, feEducationRoute, fx.tokenA, nil); status != http.StatusConflict {
			t.Fatalf("suspended education request = %d, want %d", status, http.StatusConflict)
		}
		if status, body := request(http.MethodGet, "/main/ems/licence", fx.tokenA, nil); status != http.StatusOK || !bytes.Contains(body, []byte("SUSPENDED")) {
			t.Fatalf("status while suspended = %d (%s)", status, string(body))
		}

		if status, body := request(http.MethodPost, flRegistryURL+"/"+aLicence.Number+"/revoke", fx.tokenOwner, nil); status != http.StatusOK {
			t.Fatalf("revoke = %d (%s)", status, string(body))
		}
		if status, _ := request(http.MethodGet, feEducationRoute, fx.tokenA, nil); status != http.StatusConflict {
			t.Fatalf("revoked education request = %d, want %d", status, http.StatusConflict)
		}
		if status, body := request(http.MethodGet, "/main/ems/licence", fx.tokenA, nil); status != http.StatusOK || !bytes.Contains(body, []byte("REVOKED")) {
			t.Fatalf("status while revoked = %d (%s)", status, string(body))
		}
	})

	t.Run("entitlement_absent_denies_module", func(t *testing.T) {
		// A new licence for organisation A that does NOT include education.
		issue(t, feOrgA, from, to, 1, []string{"core"})

		if status, _ := request(http.MethodGet, feEducationRoute, fx.tokenA, nil); status != http.StatusConflict {
			t.Fatalf("unentitled education request = %d, want %d", status, http.StatusConflict)
		}
		if status, body := request(http.MethodGet, "/main/ems/licence", fx.tokenA, nil); status != http.StatusOK || !bytes.Contains(body, []byte("core")) {
			t.Fatalf("status for unentitled licence = %d (%s)", status, string(body))
		}

		// EMS Core is not a licensed addon: it remains reachable.
		if status, body := request(http.MethodGet, "/main/ems/modules", fx.tokenOwner, nil); status != http.StatusOK {
			t.Fatalf("EMS Core modules route = %d, want %d (%s)", status, http.StatusOK, string(body))
		}
	})

	t.Run("organisation_isolation_cannot_borrow_another_licence", func(t *testing.T) {
		// Organisation C has an active licence without education.
		issue(t, feOrgC, from, to, 1, []string{"core"})

		if status, _ := request(http.MethodGet, feEducationRoute, fx.tokenC, nil); status != http.StatusConflict {
			t.Fatalf("unentitled organisation C = %d, want %d", status, http.StatusConflict)
		}
		// A request parameter must not override the authenticated
		// organisation's licence decision.
		if status, _ := request(http.MethodGet, feEducationRoute+"?organisation_number="+feOrgB, fx.tokenC, nil); status != http.StatusConflict {
			t.Fatalf("organisation override attempt = %d, want %d", status, http.StatusConflict)
		}
		if status, _ := request(http.MethodGet, feEducationRoute+"?organisation="+feOrgA, fx.tokenC, nil); status != http.StatusConflict {
			t.Fatalf("organisation override attempt (alt param) = %d, want %d", status, http.StatusConflict)
		}
		// The entitled organisation still passes.
		if status, body := request(http.MethodGet, feEducationRoute, fx.tokenB, nil); status != http.StatusOK {
			t.Fatalf("entitled organisation B = %d, want %d (%s)", status, http.StatusOK, string(body))
		}
	})

	t.Run("registry_authority_unchanged", func(t *testing.T) {
		if status, _ := request(http.MethodGet, faRegistryIns, fx.tokenAdmin, nil); status != http.StatusForbidden {
			t.Fatalf("Admin without registry permission = %d, want %d", status, http.StatusForbidden)
		}
		if status, _ := request(http.MethodGet, faRegistryIns, fx.tokenA, nil); status != http.StatusForbidden {
			t.Fatalf("client registry access = %d, want %d", status, http.StatusForbidden)
		}
		if status, _ := request(http.MethodGet, faRegistryIns, fx.tokenOwner, nil); status != http.StatusOK {
			t.Fatalf("registry actor access = %d, want %d", status, http.StatusOK)
		}
	})
}
