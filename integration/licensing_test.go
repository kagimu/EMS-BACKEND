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
	licensingServices "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/licensing/services"
	eosModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/eos/models"
	permissionsModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/permissions/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/internal/testutil"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/migrations"
	"gorm.io/gorm"
)

const (
	flOwner       = "ems-f-owner"
	flClientA     = "ems-f-client-a"
	flClientB     = "ems-f-client-b"
	flClientNoCtx = "ems-f-client-nocontext"
	flAdminNoPerm = "ems-f-admin-noperm"
	flUserNoPerm  = "ems-f-noperm"

	flOrgA = "ems-f-org-a"
	flOrgB = "ems-f-org-b"

	flRoleRegistry = "ems-f-role-registry"
	flRoleView     = "ems-f-role-view"

	flRegistryURL = "/main/ems/licence-registry/licences"
	flClientURL   = "/main/ems/licence"
)

type flFixture struct {
	db               *gorm.DB
	tokenOwner       string
	tokenClientA     string
	tokenClientB     string
	tokenNoCtx       string
	tokenAdminNoPerm string
	tokenNoPerm      string
}

func setupFLicensingFixture(t *testing.T) *flFixture {
	t.Helper()

	db := testutil.RequireTestDB(t)
	if err := migrations.Init(); err != nil {
		t.Fatalf("running migrations: %v", err)
	}
	t.Setenv("JWT_SECRET", "ems-foundation-f-test-secret")

	users := []string{flOwner, flClientA, flClientB, flClientNoCtx, flAdminNoPerm, flUserNoPerm}
	organisations := []string{flOrgA, flOrgB}
	roles := []string{flRoleRegistry, flRoleView}

	cleanup := func() {
		db.Unscoped().Where("organisation_number IN ?", organisations).Delete(&licensingModel.LicenceEvent{})
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
		{flOwner, "User"},
		{flClientA, "User"},
		{flClientB, "User"},
		{flClientNoCtx, "User"},
		{flAdminNoPerm, "Admin"},
		{flUserNoPerm, "User"},
	} {
		user := authModel.User{
			Number:     spec.number,
			FirstName:  "FoundationF",
			SecondName: spec.userType,
			Gender:     "M",
			Type:       spec.userType,
			Status:     "ACTIVE",
		}
		if err := db.Create(&user).Error; err != nil {
			t.Fatalf("seeding user %s: %v", spec.number, err)
		}
	}

	for _, spec := range []struct{ number, name string }{
		{flOrgA, "Foundation F Test Org A"},
		{flOrgB, "Foundation F Test Org B"},
	} {
		entity := eosModel.Entity{Number: spec.number, Name: spec.name, AddedBy: flOwner}
		if err := db.Create(&entity).Error; err != nil {
			t.Fatalf("seeding entity %s: %v", spec.number, err)
		}
	}

	for _, spec := range []struct{ userNumber, entityNumber string }{
		{flOwner, flOrgA},
		{flClientA, flOrgA},
		{flClientB, flOrgB},
		{flAdminNoPerm, flOrgA},
		{flUserNoPerm, flOrgA},
		// flClientNoCtx intentionally has no active entity (409 case).
	} {
		active := eosModel.ActiveUserEntity{UserNumber: spec.userNumber, EntityNumber: spec.entityNumber}
		if err := db.Omit("Entity", "User").Create(&active).Error; err != nil {
			t.Fatalf("setting active entity for %s: %v", spec.userNumber, err)
		}
	}

	for _, number := range roles {
		role := permissionsModel.Role{
			Number: number, Name: number, EntityNumber: flOrgA, Status: "ACTIVE", UserNumber: flOwner,
		}
		if err := db.Omit("User", "Entity", "Permissions").Create(&role).Error; err != nil {
			t.Fatalf("seeding role %s: %v", number, err)
		}
	}

	for _, spec := range []struct{ roleNumber, permission string }{
		{flRoleRegistry, "manage_licence_registry"},
		{flRoleView, "view_licence_status"},
	} {
		rolePermission := permissionsModel.RolePermission{
			RoleNumber: spec.roleNumber, Permission: spec.permission, AddedByNumber: flOwner,
		}
		if err := db.Omit("AddedBy", "Role").Create(&rolePermission).Error; err != nil {
			t.Fatalf("seeding permission %s: %v", spec.permission, err)
		}
	}

	expiry := time.Now().Add(365 * 24 * time.Hour)
	for _, spec := range []struct{ userNumber, roleNumber string }{
		{flOwner, flRoleRegistry},
		{flClientA, flRoleView},
		{flClientB, flRoleView},
		{flClientNoCtx, flRoleView},
	} {
		userRole := permissionsModel.UserRole{
			UserNumber: spec.userNumber, RoleNumber: spec.roleNumber, ExpiryDate: expiry, AddedBy: flOwner,
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

	return &flFixture{
		db:               db,
		tokenOwner:       token(flOwner),
		tokenClientA:     token(flClientA),
		tokenClientB:     token(flClientB),
		tokenNoCtx:       token(flClientNoCtx),
		tokenAdminNoPerm: token(flAdminNoPerm),
		tokenNoPerm:      token(flUserNoPerm),
	}
}

type flLicenceData struct {
	Number             string    `json:"number"`
	OrganisationNumber string    `json:"organisation_number"`
	LicenceKey         string    `json:"licence_key"`
	ValidFrom          time.Time `json:"valid_from"`
	ValidTo            time.Time `json:"valid_to"`
	GraceDays          int       `json:"grace_days"`
	MaxInstallations   int       `json:"max_installations"`
	LicensedAddons     []string  `json:"licensed_addons"`
	Status             string    `json:"status"`
	EffectiveStatus    string    `json:"effective_status"`
}

type flLicenceEnvelope struct {
	Message string        `json:"message"`
	Data    flLicenceData `json:"data"`
}

type flLicenceListEnvelope struct {
	Message string          `json:"message"`
	Data    []flLicenceData `json:"data"`
}

// TestFoundationFLicensing covers Foundation F Increment 1 end to end with one
// shared fixture (slow migration and seeding run once).
func TestFoundationFLicensing(t *testing.T) {
	fx := setupFLicensingFixture(t)

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

	issuePayload := func(org string, from, to time.Time, max int, addons []string) map[string]any {
		return map[string]any{
			"organisation_number": org,
			"valid_from":          from.Format(time.RFC3339),
			"valid_to":            to.Format(time.RFC3339),
			"max_installations":   max,
			"licensed_addons":     addons,
		}
	}

	decodeLicence := func(t *testing.T, body []byte) flLicenceEnvelope {
		t.Helper()
		var envelope flLicenceEnvelope
		if err := json.Unmarshal(body, &envelope); err != nil {
			t.Fatalf("decoding response %s: %v", string(body), err)
		}
		return envelope
	}

	from := time.Now().Add(-time.Hour)
	to := time.Now().Add(30 * 24 * time.Hour)

	var aNumber, aKey string
	var bNumber string
	var newANumber string

	t.Run("unauthenticated_returns_401", func(t *testing.T) {
		if status, _ := request(http.MethodGet, flRegistryURL, "", nil); status != http.StatusUnauthorized {
			t.Fatalf("registry status = %d, want %d", status, http.StatusUnauthorized)
		}
		if status, _ := request(http.MethodGet, flClientURL, "", nil); status != http.StatusUnauthorized {
			t.Fatalf("client status = %d, want %d", status, http.StatusUnauthorized)
		}
	})

	t.Run("control_plane_requires_explicit_authority", func(t *testing.T) {
		// An Admin user-type must NOT be granted the control plane (the
		// dedicated guard deliberately does not honour the Admin bypass).
		if status, _ := request(http.MethodGet, flRegistryURL, fx.tokenAdminNoPerm, nil); status != http.StatusForbidden {
			t.Fatalf("admin without permission = %d, want %d", status, http.StatusForbidden)
		}
		if status, _ := request(http.MethodGet, flRegistryURL, fx.tokenNoPerm, nil); status != http.StatusForbidden {
			t.Fatalf("user without permission = %d, want %d", status, http.StatusForbidden)
		}
		if status, _ := request(http.MethodGet, flRegistryURL, fx.tokenOwner, nil); status != http.StatusOK {
			t.Fatalf("explicit registry actor = %d, want %d", status, http.StatusOK)
		}
	})

	t.Run("issue_returns_key_once", func(t *testing.T) {
		status, body := request(http.MethodPost, flRegistryURL, fx.tokenOwner, issuePayload(flOrgA, from, to, 3, []string{"education"}))
		if status != http.StatusCreated {
			t.Fatalf("issue status = %d, want %d (%s)", status, http.StatusCreated, string(body))
		}
		issued := decodeLicence(t, body)
		if issued.Data.Number == "" || issued.Data.LicenceKey == "" {
			t.Fatalf("issued licence must carry number and key: %s", string(body))
		}
		if issued.Data.GraceDays != 14 || issued.Data.EffectiveStatus != "ACTIVE" {
			t.Fatalf("unexpected issued licence: %s", string(body))
		}
		aNumber, aKey = issued.Data.Number, issued.Data.LicenceKey

		// The key is returned at issuance only; ordinary reads must mask it.
		if status, body = request(http.MethodGet, flRegistryURL+"/"+aNumber, fx.tokenOwner, nil); status != http.StatusOK {
			t.Fatalf("detail status = %d, want %d", status, http.StatusOK)
		}
		if bytes.Contains(body, []byte(aKey)) {
			t.Fatalf("licence detail must not expose the licence key: %s", string(body))
		}
		if status, body = request(http.MethodGet, flRegistryURL, fx.tokenOwner, nil); status != http.StatusOK {
			t.Fatalf("list status = %d, want %d", status, http.StatusOK)
		}
		if bytes.Contains(body, []byte(aKey)) {
			t.Fatalf("licence list must not expose the licence key: %s", string(body))
		}
	})

	t.Run("invalid_issuance_requests", func(t *testing.T) {
		if status, _ := request(http.MethodPost, flRegistryURL, fx.tokenOwner, issuePayload("ems-f-unknown-org", from, to, 1, []string{"education"})); status != http.StatusBadRequest {
			t.Fatalf("unknown organisation = %d, want %d", status, http.StatusBadRequest)
		}
		if status, _ := request(http.MethodPost, flRegistryURL, fx.tokenOwner, issuePayload(flOrgB, to, from, 1, []string{"education"})); status != http.StatusBadRequest {
			t.Fatalf("inverted validity = %d, want %d", status, http.StatusBadRequest)
		}
		if status, _ := request(http.MethodPost, flRegistryURL, fx.tokenOwner, issuePayload(flOrgB, from, to, 0, []string{"education"})); status != http.StatusBadRequest {
			t.Fatalf("zero installations = %d, want %d", status, http.StatusBadRequest)
		}
		if status, _ := request(http.MethodPost, flRegistryURL, fx.tokenOwner, issuePayload(flOrgB, from, to, 1, []string{"finance"})); status != http.StatusBadRequest {
			t.Fatalf("unknown add-on = %d, want %d", status, http.StatusBadRequest)
		}
	})

	t.Run("duplicate_active_conflict", func(t *testing.T) {
		if status, _ := request(http.MethodPost, flRegistryURL, fx.tokenOwner, issuePayload(flOrgA, from, to, 1, []string{"education"})); status != http.StatusConflict {
			t.Fatalf("duplicate active licence = %d, want %d", status, http.StatusConflict)
		}
	})

	t.Run("configure_and_renew", func(t *testing.T) {
		extended := time.Now().Add(60 * 24 * time.Hour)
		status, body := request(http.MethodPut, flRegistryURL+"/"+aNumber, fx.tokenOwner, map[string]any{
			"valid_to":          extended.Format(time.RFC3339),
			"max_installations": 5,
			"licensed_addons":   []string{"education"},
		})
		if status != http.StatusOK {
			t.Fatalf("configure status = %d, want %d (%s)", status, http.StatusOK, string(body))
		}
		if updated := decodeLicence(t, body); updated.Data.MaxInstallations != 5 {
			t.Fatalf("configure did not apply: %s", string(body))
		}

		renewedTo := time.Now().Add(90 * 24 * time.Hour)
		status, body = request(http.MethodPost, flRegistryURL+"/"+aNumber+"/renew", fx.tokenOwner, map[string]any{
			"valid_to": renewedTo.Format(time.RFC3339),
		})
		if status != http.StatusOK {
			t.Fatalf("renew status = %d, want %d (%s)", status, http.StatusOK, string(body))
		}
		renewed := decodeLicence(t, body)
		if renewed.Data.Number != aNumber {
			t.Fatalf("renewal must preserve the licence number: %s", string(body))
		}

		var stored licensingModel.Licence
		if err := fx.db.Where("number = ?", aNumber).First(&stored).Error; err != nil {
			t.Fatalf("loading stored licence: %v", err)
		}
		if stored.LicenceKey != aKey {
			t.Fatal("renewal must preserve the licence key")
		}
		if stored.OrganisationNumber != flOrgA || stored.MaxInstallations != 5 {
			t.Fatalf("renewal must preserve organisation and entitlements/capacity: %+v", stored)
		}
		var count int64
		fx.db.Model(&licensingModel.Licence{}).Where("organisation_number = ?", flOrgA).Count(&count)
		if count != 1 {
			t.Fatalf("renewal must not create a replacement licence; count = %d", count)
		}
	})

	t.Run("renewal_requires_extension", func(t *testing.T) {
		shorter := time.Now().Add(10 * 24 * time.Hour)
		if status, _ := request(http.MethodPost, flRegistryURL+"/"+aNumber+"/renew", fx.tokenOwner, map[string]any{
			"valid_to": shorter.Format(time.RFC3339),
		}); status != http.StatusBadRequest {
			t.Fatalf("non-extending renewal = %d, want %d", status, http.StatusBadRequest)
		}
	})

	t.Run("entitlements", func(t *testing.T) {
		licensed, err := licensingServices.IsAddonLicensed(flOrgA, "education")
		if err != nil || !licensed {
			t.Fatalf("education must be licensed for org A (licensed=%v, err=%v)", licensed, err)
		}
		unlicensed, err := licensingServices.IsAddonLicensed(flOrgA, "core")
		if err != nil || unlicensed {
			t.Fatalf("core must not be licensed unless listed (licensed=%v, err=%v)", unlicensed, err)
		}
	})

	t.Run("client_permission_context_and_missing_licence", func(t *testing.T) {
		if status, _ := request(http.MethodGet, flClientURL, fx.tokenNoPerm, nil); status != http.StatusForbidden {
			t.Fatalf("client without permission = %d, want %d", status, http.StatusForbidden)
		}
		if status, _ := request(http.MethodGet, flClientURL, fx.tokenNoCtx, nil); status != http.StatusConflict {
			t.Fatalf("client without active organisation = %d, want %d", status, http.StatusConflict)
		}
		if status, _ := request(http.MethodGet, flClientURL, fx.tokenClientB, nil); status != http.StatusNotFound {
			t.Fatalf("client without a licence = %d, want %d", status, http.StatusNotFound)
		}
	})

	t.Run("client_status_and_isolation", func(t *testing.T) {
		status, body := request(http.MethodGet, flClientURL, fx.tokenClientA, nil)
		if status != http.StatusOK {
			t.Fatalf("client A status = %d, want %d (%s)", status, http.StatusOK, string(body))
		}
		licence := decodeLicence(t, body)
		if licence.Data.Number != aNumber {
			t.Fatalf("client A must see its own licence: %s", string(body))
		}
		if licence.Data.LicenceKey != "" || bytes.Contains(body, []byte(aKey)) {
			t.Fatalf("client status must not expose the licence key: %s", string(body))
		}
		if bytes.Contains(body, []byte(flOrgB)) {
			t.Fatalf("client A must not see another organisation: %s", string(body))
		}

		if status, _ := request(http.MethodGet, flRegistryURL+"/"+aNumber, fx.tokenClientA, nil); status != http.StatusForbidden {
			t.Fatalf("client registry access = %d, want %d", status, http.StatusForbidden)
		}
	})

	t.Run("concurrent_issuance_one_active", func(t *testing.T) {
		type result struct {
			status int
			err    error
		}
		results := make(chan result, 2)
		for i := 0; i < 2; i++ {
			go func() {
				raw, _ := json.Marshal(issuePayload(flOrgB, from, to, 3, []string{"education"}))
				req, err := http.NewRequest(http.MethodPost, server.URL+flRegistryURL, bytes.NewReader(raw))
				if err != nil {
					results <- result{err: err}
					return
				}
				req.Header.Set("Authorization", "Bearer "+fx.tokenOwner)
				req.Header.Set("Content-Type", "application/json")
				resp, err := client.Do(req)
				if err != nil {
					results <- result{err: err}
					return
				}
				defer resp.Body.Close()
				_, _ = io.Copy(io.Discard, resp.Body)
				results <- result{status: resp.StatusCode}
			}()
		}

		created, conflicts := 0, 0
		for i := 0; i < 2; i++ {
			outcome := <-results
			if outcome.err != nil {
				t.Fatalf("concurrent issuance request failed: %v", outcome.err)
			}
			switch outcome.status {
			case http.StatusCreated:
				created++
			case http.StatusConflict:
				conflicts++
			}
		}
		if created != 1 || conflicts != 1 {
			t.Fatalf("concurrent issuance = %d created / %d conflicts, want 1/1", created, conflicts)
		}

		var activeCount int64
		fx.db.Model(&licensingModel.Licence{}).Where("active_org_number = ?", flOrgB).Count(&activeCount)
		if activeCount != 1 {
			t.Fatalf("active licence count for org B = %d, want 1", activeCount)
		}

		var stored licensingModel.Licence
		if err := fx.db.Where("active_org_number = ?", flOrgB).First(&stored).Error; err != nil {
			t.Fatalf("loading org B licence: %v", err)
		}
		bNumber = stored.Number

		licensed, err := licensingServices.IsAddonLicensed(flOrgB, "education")
		if err != nil || !licensed {
			t.Fatalf("education must be licensed for org B (licensed=%v, err=%v)", licensed, err)
		}

		if status, body := request(http.MethodGet, flClientURL, fx.tokenClientB, nil); status != http.StatusOK || !bytes.Contains(body, []byte(bNumber)) {
			t.Fatalf("client B must see its own licence (status=%d): %s", status, string(body))
		}
	})

	t.Run("suspend_revoke_lifecycle", func(t *testing.T) {
		// Suspend the original org A licence: the client keeps visibility of
		// its status, and the active guard is cleared so one new ACTIVE
		// licence may be issued (F-Q7).
		status, body := request(http.MethodPost, flRegistryURL+"/"+aNumber+"/suspend", fx.tokenOwner, nil)
		if status != http.StatusOK {
			t.Fatalf("suspend status = %d, want %d (%s)", status, http.StatusOK, string(body))
		}
		if suspended := decodeLicence(t, body); suspended.Data.EffectiveStatus != "SUSPENDED" || suspended.Data.Status != "SUSPENDED" {
			t.Fatalf("suspend not applied: %s", string(body))
		}
		if status, body = request(http.MethodGet, flClientURL, fx.tokenClientA, nil); status != http.StatusOK {
			t.Fatalf("client status after suspend = %d, want %d", status, http.StatusOK)
		} else if suspended := decodeLicence(t, body); suspended.Data.EffectiveStatus != "SUSPENDED" {
			t.Fatalf("client must see SUSPENDED: %s", string(body))
		}

		if status, _ = request(http.MethodPost, flRegistryURL+"/"+aNumber+"/suspend", fx.tokenOwner, nil); status != http.StatusConflict {
			t.Fatalf("double suspend = %d, want %d", status, http.StatusConflict)
		}
		if status, _ = request(http.MethodPut, flRegistryURL+"/"+aNumber, fx.tokenOwner, map[string]any{"max_installations": 9}); status != http.StatusConflict {
			t.Fatalf("configure suspended licence = %d, want %d", status, http.StatusConflict)
		}

		// A new ACTIVE licence for the same organisation is now allowed.
		status, body = request(http.MethodPost, flRegistryURL, fx.tokenOwner, issuePayload(flOrgA, from, to, 2, []string{"education"}))
		if status != http.StatusCreated {
			t.Fatalf("re-issue after suspend = %d, want %d (%s)", status, http.StatusCreated, string(body))
		}
		newANumber = decodeLicence(t, body).Data.Number

		if status, _ = request(http.MethodPost, flRegistryURL, fx.tokenOwner, issuePayload(flOrgA, from, to, 2, []string{"education"})); status != http.StatusConflict {
			t.Fatalf("second active licence = %d, want %d", status, http.StatusConflict)
		}

		// Revoke the new licence (terminal) and confirm the guards cleared.
		status, body = request(http.MethodPost, flRegistryURL+"/"+newANumber+"/revoke", fx.tokenOwner, nil)
		if status != http.StatusOK {
			t.Fatalf("revoke status = %d, want %d (%s)", status, http.StatusOK, string(body))
		}
		if revoked := decodeLicence(t, body); revoked.Data.EffectiveStatus != "REVOKED" {
			t.Fatalf("revoke not applied: %s", string(body))
		}
		if status, _ = request(http.MethodPost, flRegistryURL+"/"+newANumber+"/revoke", fx.tokenOwner, nil); status != http.StatusConflict {
			t.Fatalf("double revoke = %d, want %d", status, http.StatusConflict)
		}
		if status, _ = request(http.MethodPost, flRegistryURL+"/"+newANumber+"/suspend", fx.tokenOwner, nil); status != http.StatusConflict {
			t.Fatalf("suspend revoked licence = %d, want %d", status, http.StatusConflict)
		}

		if status, body = request(http.MethodGet, flClientURL, fx.tokenClientA, nil); status != http.StatusOK {
			t.Fatalf("client status after revoke = %d, want %d", status, http.StatusOK)
		} else if revoked := decodeLicence(t, body); revoked.Data.EffectiveStatus != "REVOKED" {
			t.Fatalf("client must see REVOKED: %s", string(body))
		}

		for _, number := range []string{aNumber, newANumber} {
			var stored licensingModel.Licence
			if err := fx.db.Where("number = ?", number).First(&stored).Error; err != nil {
				t.Fatalf("loading licence %s: %v", number, err)
			}
			if stored.ActiveOrgNumber != nil {
				t.Fatalf("non-ACTIVE licence %s must clear the active guard", number)
			}
		}
	})

	t.Run("invalid_status_filter", func(t *testing.T) {
		if status, _ := request(http.MethodGet, flRegistryURL+"?status=IN_GRACE", fx.tokenOwner, nil); status != http.StatusBadRequest {
			t.Fatalf("derived status filter = %d, want %d", status, http.StatusBadRequest)
		}
	})
}
