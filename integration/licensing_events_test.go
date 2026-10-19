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
	sevOwner       = "ems-sev-owner"
	sevAdminNoPerm = "ems-sev-admin-noperm"
	sevClient      = "ems-sev-client"
	sevRateUser    = "ems-sev-rate-user"

	sevOrgA = "ems-sev-org-a"
	sevOrgB = "ems-sev-org-b"

	sevRoleRegistry = "ems-sev-role-registry"
	sevRoleClient   = "ems-sev-role-client"

	sevEventsURL = "/main/ems/licence-registry/events"
)

type sevFixture struct {
	db          *gorm.DB
	tokenOwner  string
	tokenAdmin  string
	tokenClient string
	tokenRate   string
}

func setupFEventsFixture(t *testing.T) *sevFixture {
	t.Helper()

	db := testutil.RequireTestDB(t)
	if err := migrations.Init(); err != nil {
		t.Fatalf("running migrations: %v", err)
	}
	t.Setenv("JWT_SECRET", "ems-foundation-f-events-test-secret")

	users := []string{sevOwner, sevAdminNoPerm, sevClient, sevRateUser}
	organisations := []string{sevOrgA, sevOrgB}
	clientRoleB := sevRoleClient + "-b"
	roles := []string{sevRoleRegistry, sevRoleClient, clientRoleB}

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
		{sevOwner, "User"},
		{sevAdminNoPerm, "Admin"},
		{sevClient, "User"},
		{sevRateUser, "User"},
	} {
		user := authModel.User{
			Number:     spec.number,
			FirstName:  "FoundationF",
			SecondName: "Events",
			Gender:     "M",
			Type:       spec.userType,
			Status:     "ACTIVE",
		}
		if err := db.Create(&user).Error; err != nil {
			t.Fatalf("seeding user %s: %v", spec.number, err)
		}
	}

	for _, spec := range []struct{ number, name string }{
		{sevOrgA, "Foundation F Events Org A"},
		{sevOrgB, "Foundation F Events Org B"},
	} {
		entity := eosModel.Entity{Number: spec.number, Name: spec.name, AddedBy: sevOwner}
		if err := db.Create(&entity).Error; err != nil {
			t.Fatalf("seeding entity %s: %v", spec.number, err)
		}
	}

	for _, spec := range []struct{ userNumber, entityNumber string }{
		{sevOwner, sevOrgA},
		{sevAdminNoPerm, sevOrgA},
		{sevClient, sevOrgA},
		{sevRateUser, sevOrgB},
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

	for _, number := range []string{sevRoleRegistry, sevRoleClient} {
		role := permissionsModel.Role{
			Number: number, Name: number, EntityNumber: sevOrgA, Status: "ACTIVE", UserNumber: sevOwner,
		}
		if err := db.Omit("User", "Entity", "Permissions").Create(&role).Error; err != nil {
			t.Fatalf("seeding role %s: %v", number, err)
		}
	}

	for _, spec := range []struct{ roleNumber, permission string }{
		{sevRoleRegistry, "manage_licence_registry"},
		{sevRoleClient, "view_licence_status"},
		{sevRoleClient, "activate_licence_installation"},
		{sevRoleClient, "deactivate_own_installation"},
		{sevRoleClient, "request_licence_renewal"},
	} {
		rolePermission := permissionsModel.RolePermission{
			RoleNumber: spec.roleNumber, Permission: spec.permission, AddedByNumber: sevOwner,
		}
		if err := db.Omit("AddedBy", "Role").Create(&rolePermission).Error; err != nil {
			t.Fatalf("seeding permission %s: %v", spec.permission, err)
		}
	}

	expiry := time.Now().Add(365 * 24 * time.Hour)
	cloneRoleForOrganisation(t, db, sevRoleClient, clientRoleB, sevOrgB, sevOwner)
	for _, spec := range []struct{ userNumber, roleNumber string }{
		{sevOwner, sevRoleRegistry},
		{sevClient, sevRoleClient},
		{sevRateUser, clientRoleB},
	} {
		userRole := permissionsModel.UserRole{
			UserNumber: spec.userNumber, RoleNumber: spec.roleNumber, ExpiryDate: expiry, AddedBy: sevOwner,
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

	return &sevFixture{
		db:          db,
		tokenOwner:  token(sevOwner),
		tokenAdmin:  token(sevAdminNoPerm),
		tokenClient: token(sevClient),
		tokenRate:   token(sevRateUser),
	}
}

// TestFoundationFEventsAndRateLimit covers Foundation F Increment 5: durable
// append-only licensing facts for Foundation G, the control-plane facts read
// endpoint, and the in-process activation rate limit. Notifications, payment,
// the remote service and the owner UI are boundaries only and are not built.
func TestFoundationFEventsAndRateLimit(t *testing.T) {
	fx := setupFEventsFixture(t)

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

	issueBody := func(org string) map[string]any {
		return map[string]any{
			"organisation_number": org,
			"valid_from":          time.Now().Add(-time.Hour).Format(time.RFC3339),
			"valid_to":            time.Now().Add(30 * 24 * time.Hour).Format(time.RFC3339),
			"max_installations":   3,
			"licensed_addons":     []string{"education"},
		}
	}

	issue := func(t *testing.T, org string) string {
		t.Helper()
		status, body := request(http.MethodPost, flRegistryURL, fx.tokenOwner, issueBody(org))
		if status != http.StatusCreated {
			t.Fatalf("issuing licence for %s = %d (%s)", org, status, string(body))
		}
		var envelope flLicenceEnvelope
		if err := json.Unmarshal(body, &envelope); err != nil {
			t.Fatalf("decoding licence: %v", err)
		}
		return envelope.Data.Number
	}

	hasEvent := func(t *testing.T, licenceNumber, eventType, actor string) bool {
		t.Helper()
		var count int64
		fx.db.Model(&licensingModel.LicenceEvent{}).
			Where("licence_number = ? AND event_type = ? AND actor_user_number = ?",
				licenceNumber, eventType, actor).
			Count(&count)
		return count > 0
	}

	var licenceA string
	var aKey string
	var activeANumber string

	t.Run("licence_lifecycle_facts_recorded", func(t *testing.T) {
		status, body := request(http.MethodPost, flRegistryURL, fx.tokenOwner, issueBody(sevOrgA))
		if status != http.StatusCreated {
			t.Fatalf("issue = %d (%s)", status, string(body))
		}
		var envelope flLicenceEnvelope
		if err := json.Unmarshal(body, &envelope); err != nil {
			t.Fatalf("decoding licence: %v", err)
		}
		licenceA = envelope.Data.Number
		aKey = envelope.Data.LicenceKey

		if !hasEvent(t, licenceA, string(licensingModel.LicenceEventIssued), sevOwner) {
			t.Fatal("LICENCE_ISSUED fact not recorded with the issuing actor")
		}

		if status, body := request(http.MethodPut, flRegistryURL+"/"+licenceA, fx.tokenOwner, map[string]any{
			"max_installations": 4,
		}); status != http.StatusOK {
			t.Fatalf("configure = %d (%s)", status, string(body))
		}
		if !hasEvent(t, licenceA, string(licensingModel.LicenceEventConfigured), sevOwner) {
			t.Fatal("LICENCE_CONFIGURED fact not recorded with actor")
		}

		if status, body := request(http.MethodPost, flRegistryURL+"/"+licenceA+"/renew", fx.tokenOwner, map[string]any{
			"valid_to": time.Now().Add(60 * 24 * time.Hour).Format(time.RFC3339),
		}); status != http.StatusOK {
			t.Fatalf("renew = %d (%s)", status, string(body))
		}
		if !hasEvent(t, licenceA, string(licensingModel.LicenceEventRenewed), sevOwner) {
			t.Fatal("LICENCE_RENEWED fact not recorded with actor")
		}

		if status, body := request(http.MethodPost, flRegistryURL+"/"+licenceA+"/suspend", fx.tokenOwner, nil); status != http.StatusOK {
			t.Fatalf("suspend = %d (%s)", status, string(body))
		}
		if !hasEvent(t, licenceA, string(licensingModel.LicenceEventSuspended), sevOwner) {
			t.Fatal("LICENCE_SUSPENDED fact not recorded with actor")
		}

		if status, body := request(http.MethodPost, flRegistryURL+"/"+licenceA+"/revoke", fx.tokenOwner, nil); status != http.StatusOK {
			t.Fatalf("revoke = %d (%s)", status, string(body))
		}
		if !hasEvent(t, licenceA, string(licensingModel.LicenceEventRevoked), sevOwner) {
			t.Fatal("LICENCE_REVOKED fact not recorded with actor")
		}
	})

	t.Run("installation_facts_recorded", func(t *testing.T) {
		number := issue(t, sevOrgA)
		activeANumber = number
		identity := uuid.NewV4().String()

		// Read the freshly issued key through the DB (control-plane only).
		var stored licensingModel.Licence
		if err := fx.db.Where("number = ?", number).First(&stored).Error; err != nil {
			t.Fatalf("loading licence: %v", err)
		}
		freshKey := stored.LicenceKey

		status, body := request(http.MethodPost, faActivateURL, fx.tokenClient, map[string]any{
			"licence_key":           freshKey,
			"installation_identity": identity,
		})
		if status != http.StatusCreated {
			t.Fatalf("activation = %d (%s)", status, string(body))
		}
		if !hasEvent(t, number, string(licensingModel.InstallationEventActivated), sevClient) {
			t.Fatal("INSTALLATION_ACTIVATED fact not recorded")
		}

		var installation licensingModel.Installation
		if err := fx.db.Where("licence_number = ?", number).First(&installation).Error; err != nil {
			t.Fatalf("loading installation: %v", err)
		}

		if status, body := request(http.MethodPost, faClientInst+"/"+installation.Number+"/deactivate", fx.tokenClient, nil); status != http.StatusOK {
			t.Fatalf("deactivate = %d (%s)", status, string(body))
		}
		if !hasEvent(t, number, string(licensingModel.InstallationEventDeactivated), sevClient) {
			t.Fatal("INSTALLATION_DEACTIVATED fact not recorded")
		}

		if status, body := request(http.MethodPost, faActivateURL, fx.tokenClient, map[string]any{
			"licence_key":           freshKey,
			"installation_identity": identity,
		}); status != http.StatusCreated {
			t.Fatalf("reactivate = %d (%s)", status, string(body))
		}
		if !hasEvent(t, number, string(licensingModel.InstallationEventReactivated), sevClient) {
			t.Fatal("INSTALLATION_REACTIVATED fact not recorded")
		}

		if status, body := request(http.MethodPost, faRegistryIns+"/"+installation.Number+"/release", fx.tokenOwner, map[string]any{
			"reason": "hardware retired",
		}); status != http.StatusOK {
			t.Fatalf("release = %d (%s)", status, string(body))
		}
		if !hasEvent(t, number, string(licensingModel.InstallationEventReleased), sevOwner) {
			t.Fatal("INSTALLATION_RELEASED fact not recorded")
		}
	})

	t.Run("replacement_fact_recorded", func(t *testing.T) {
		// Organisation A already holds an ACTIVE licence (from the installation
		// subtest); supersede it by replacement.
		number := activeANumber
		status, body := request(http.MethodPost, flRegistryURL+"/"+number+"/replace", fx.tokenOwner, map[string]any{
			"organisation_number": sevOrgA,
			"valid_from":          time.Now().Add(-time.Hour).Format(time.RFC3339),
			"valid_to":            time.Now().Add(30 * 24 * time.Hour).Format(time.RFC3339),
			"max_installations":   2,
			"licensed_addons":     []string{"education"},
			"reason":              "package change",
		})
		if status != http.StatusCreated {
			t.Fatalf("replace = %d (%s)", status, string(body))
		}
		var envelope flLicenceEnvelope
		if err := json.Unmarshal(body, &envelope); err != nil {
			t.Fatalf("decoding replacement: %v", err)
		}
		activeANumber = envelope.Data.Number
		if !hasEvent(t, number, string(licensingModel.LicenceEventReplaced), sevOwner) {
			t.Fatal("LICENCE_REPLACED fact not recorded")
		}
	})

	t.Run("client_renewal_request_fact_recorded", func(t *testing.T) {
		if status, body := request(http.MethodPost, "/main/ems/licence/renewal-request", fx.tokenClient, nil); status != http.StatusOK {
			t.Fatalf("renewal request = %d (%s)", status, string(body))
		}
		if !hasEvent(t, activeANumber, string(licensingModel.ClientEventRenewalRequested), sevClient) {
			t.Fatal("CLIENT_RENEWAL_REQUESTED fact not recorded")
		}
	})

	t.Run("events_endpoint_control_plane_only_and_no_key", func(t *testing.T) {
		if status, _ := request(http.MethodGet, sevEventsURL, fx.tokenClient, nil); status != http.StatusForbidden {
			t.Fatalf("client events access = %d, want %d", status, http.StatusForbidden)
		}
		if status, _ := request(http.MethodGet, sevEventsURL, fx.tokenAdmin, nil); status != http.StatusForbidden {
			t.Fatalf("admin without permission = %d, want %d", status, http.StatusForbidden)
		}
		status, body := request(http.MethodGet, sevEventsURL, fx.tokenOwner, nil)
		if status != http.StatusOK {
			t.Fatalf("owner events access = %d, want %d (%s)", status, http.StatusOK, string(body))
		}
		if aKey != "" && bytes.Contains(body, []byte(aKey)) {
			t.Fatalf("events must never expose the licence key: %s", string(body))
		}
		if !bytes.Contains(body, []byte("LICENCE_ISSUED")) {
			t.Fatalf("events response should include recorded facts: %s", string(body))
		}
	})

	t.Run("activation_rate_limited_in_process", func(t *testing.T) {
		// The in-process limiter allows 30 attempts per user per minute.
		var lastStatus int
		for i := 0; i < 31; i++ {
			status, _ := request(http.MethodPost, faActivateURL, fx.tokenRate, map[string]any{
				"licence_key":           "bogus-key",
				"installation_identity": uuid.NewV4().String(),
			})
			lastStatus = status
		}
		if lastStatus != http.StatusTooManyRequests {
			t.Fatalf("31st activation attempt = %d, want %d", lastStatus, http.StatusTooManyRequests)
		}
		// Another user is unaffected by the first user's burst.
		if status, _ := request(http.MethodPost, faActivateURL, fx.tokenClient, map[string]any{
			"licence_key":           "bogus-key-2",
			"installation_identity": uuid.NewV4().String(),
		}); status == http.StatusTooManyRequests {
			t.Fatalf("a different user must not inherit another user's rate limit")
		}
	})
}
