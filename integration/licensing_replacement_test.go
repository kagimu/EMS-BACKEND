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
	frOwner       = "ems-fr-owner"
	frAdminNoPerm = "ems-fr-admin-noperm"
	frClientA     = "ems-fr-client-a"
	frClientB     = "ems-fr-client-b"
	frClientC     = "ems-fr-client-c"

	frOrgA = "ems-fr-org-a"
	frOrgB = "ems-fr-org-b"
	frOrgC = "ems-fr-org-c"

	frRoleRegistry = "ems-fr-role-registry"
	frRoleClient   = "ems-fr-role-client"
)

type frFixture struct {
	db         *gorm.DB
	tokenOwner string
	tokenAdmin string
	tokenA     string
	tokenB     string
	tokenC     string
}

func setupFReplacementFixture(t *testing.T) *frFixture {
	t.Helper()

	db := testutil.RequireTestDB(t)
	if err := migrations.Init(); err != nil {
		t.Fatalf("running migrations: %v", err)
	}
	t.Setenv("JWT_SECRET", "ems-foundation-f-replacement-test-secret")

	users := []string{frOwner, frAdminNoPerm, frClientA, frClientB, frClientC}
	organisations := []string{frOrgA, frOrgB, frOrgC}
	roles := []string{frRoleRegistry, frRoleClient}

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
		{frOwner, "User"},
		{frAdminNoPerm, "Admin"},
		{frClientA, "User"},
		{frClientB, "User"},
		{frClientC, "User"},
	} {
		user := authModel.User{
			Number:     spec.number,
			FirstName:  "FoundationF",
			SecondName: "Replacement",
			Gender:     "M",
			Type:       spec.userType,
			Status:     "ACTIVE",
		}
		if err := db.Create(&user).Error; err != nil {
			t.Fatalf("seeding user %s: %v", spec.number, err)
		}
	}

	for _, spec := range []struct{ number, name string }{
		{frOrgA, "Foundation F Replacement Org A"},
		{frOrgB, "Foundation F Replacement Org B"},
		{frOrgC, "Foundation F Replacement Org C"},
	} {
		entity := eosModel.Entity{Number: spec.number, Name: spec.name, AddedBy: frOwner}
		if err := db.Create(&entity).Error; err != nil {
			t.Fatalf("seeding entity %s: %v", spec.number, err)
		}
	}

	for _, spec := range []struct{ userNumber, entityNumber string }{
		{frOwner, frOrgA},
		{frAdminNoPerm, frOrgA},
		{frClientA, frOrgA},
		{frClientB, frOrgB},
		{frClientC, frOrgC},
	} {
		active := eosModel.ActiveUserEntity{UserNumber: spec.userNumber, EntityNumber: spec.entityNumber}
		if err := db.Omit("Entity", "User").Create(&active).Error; err != nil {
			t.Fatalf("setting active entity for %s: %v", spec.userNumber, err)
		}
	}

	for _, number := range roles {
		role := permissionsModel.Role{
			Number: number, Name: number, EntityNumber: frOrgA, Status: "ACTIVE", UserNumber: frOwner,
		}
		if err := db.Omit("User", "Entity", "Permissions").Create(&role).Error; err != nil {
			t.Fatalf("seeding role %s: %v", number, err)
		}
	}

	for _, spec := range []struct{ roleNumber, permission string }{
		{frRoleRegistry, "manage_licence_registry"},
		{frRoleClient, "view_licence_status"},
		{frRoleClient, "activate_licence_installation"},
		{frRoleClient, "deactivate_own_installation"},
	} {
		rolePermission := permissionsModel.RolePermission{
			RoleNumber: spec.roleNumber, Permission: spec.permission, AddedByNumber: frOwner,
		}
		if err := db.Omit("AddedBy", "Role").Create(&rolePermission).Error; err != nil {
			t.Fatalf("seeding permission %s: %v", spec.permission, err)
		}
	}

	expiry := time.Now().Add(365 * 24 * time.Hour)
	for _, spec := range []struct{ userNumber, roleNumber string }{
		{frOwner, frRoleRegistry},
		{frClientA, frRoleClient},
		{frClientB, frRoleClient},
		{frClientC, frRoleClient},
	} {
		userRole := permissionsModel.UserRole{
			UserNumber: spec.userNumber, RoleNumber: spec.roleNumber, ExpiryDate: expiry, AddedBy: frOwner,
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

	return &frFixture{
		db:         db,
		tokenOwner: token(frOwner),
		tokenAdmin: token(frAdminNoPerm),
		tokenA:     token(frClientA),
		tokenB:     token(frClientB),
		tokenC:     token(frClientC),
	}
}

// TestFoundationFReplacement covers Foundation F Increment 4: atomic
// new-licence issuance that supersedes an ACTIVE licence, explicitly releases
// the old licence's capacity-holding installations, and never transfers
// installations or capacity.
func TestFoundationFReplacement(t *testing.T) {
	fx := setupFReplacementFixture(t)

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

	replaceBody := func(org string, from, to time.Time, max int, addons []string, reason string) map[string]any {
		return map[string]any{
			"organisation_number": org,
			"valid_from":          from.Format(time.RFC3339),
			"valid_to":            to.Format(time.RFC3339),
			"max_installations":   max,
			"licensed_addons":     addons,
			"reason":              reason,
		}
	}
	replace := func(t *testing.T, number string, payload map[string]any) (int, []byte) {
		t.Helper()
		return request(http.MethodPost, flRegistryURL+"/"+number+"/replace", fx.tokenOwner, payload)
	}

	activate := func(t *testing.T, token, key, identity string) (int, []byte) {
		t.Helper()
		return request(http.MethodPost, faActivateURL, token, map[string]any{
			"licence_key":           key,
			"installation_identity": identity,
		})
	}

	activeLicences := func(t *testing.T, org string) int64 {
		t.Helper()
		var count int64
		fx.db.Model(&licensingModel.Licence{}).Where("active_org_number = ?", org).Count(&count)
		return count
	}
	activeInstallations := func(t *testing.T, licenceNumber string) int64 {
		t.Helper()
		var count int64
		fx.db.Model(&licensingModel.Installation{}).
			Where("licence_number = ? AND status = ?", licenceNumber, licensingModel.InstallationStatusActive).
			Count(&count)
		return count
	}

	from := time.Now().Add(-time.Hour)
	to := time.Now().Add(30 * 24 * time.Hour)

	var l1, l2 flLicenceData
	var i1, i2 faActivationData

	t.Run("replacement_is_atomic_and_preserves_history", func(t *testing.T) {
		l1 = issue(t, frOrgA, from, to, 3, []string{"education"})

		status, body := activate(t, fx.tokenA, l1.LicenceKey, uuid.NewV4().String())
		if status != http.StatusCreated {
			t.Fatalf("activating first installation = %d (%s)", status, string(body))
		}
		var envelope faActivationEnvelope
		if err := json.Unmarshal(body, &envelope); err != nil {
			t.Fatalf("decoding activation: %v", err)
		}
		i1 = envelope.Data

		status, body = activate(t, fx.tokenA, l1.LicenceKey, uuid.NewV4().String())
		if status != http.StatusCreated {
			t.Fatalf("activating second installation = %d (%s)", status, string(body))
		}
		if err := json.Unmarshal(body, &envelope); err != nil {
			t.Fatalf("decoding activation: %v", err)
		}
		i2 = envelope.Data

		// Client-deactivated installations still hold capacity until release.
		if status, body := request(http.MethodPost, faClientInst+"/"+i2.InstallationNumber+"/deactivate", fx.tokenA, nil); status != http.StatusOK {
			t.Fatalf("client deactivation = %d (%s)", status, string(body))
		}

		status, body = replace(t, l1.Number, replaceBody(frOrgA, from, time.Now().Add(60*24*time.Hour), 2, []string{"education"}, "package change"))
		if status != http.StatusCreated {
			t.Fatalf("replacement = %d, want %d (%s)", status, http.StatusCreated, string(body))
		}
		var issued flLicenceEnvelope
		if err := json.Unmarshal(body, &issued); err != nil {
			t.Fatalf("decoding replacement: %v", err)
		}
		l2 = issued.Data
		if l2.Number == "" || l2.Number == l1.Number || l2.LicenceKey == "" || l2.LicenceKey == l1.LicenceKey {
			t.Fatalf("replacement must be a new licence with a new key: %s", string(body))
		}

		var oldStored, newStored licensingModel.Licence
		if err := fx.db.Where("number = ?", l1.Number).First(&oldStored).Error; err != nil {
			t.Fatalf("loading old licence: %v", err)
		}
		if err := fx.db.Where("number = ?", l2.Number).First(&newStored).Error; err != nil {
			t.Fatalf("loading new licence: %v", err)
		}
		if oldStored.Status != licensingModel.StatusRevoked || oldStored.ActiveOrgNumber != nil {
			t.Fatalf("old licence must be revoked with a cleared guard: %+v", oldStored)
		}
		if oldStored.ReplacedByNumber == nil || *oldStored.ReplacedByNumber != l2.Number {
			t.Fatalf("old licence must link to the replacement licence")
		}
		if newStored.Status != licensingModel.StatusActive || newStored.ActiveOrgNumber == nil || *newStored.ActiveOrgNumber != frOrgA {
			t.Fatalf("new licence must be the active licence: %+v", newStored)
		}
		if newStored.MaxInstallations != 2 {
			t.Fatalf("new licence max_installations = %d, want 2", newStored.MaxInstallations)
		}
		if count := activeLicences(t, frOrgA); count != 1 {
			t.Fatalf("active licences for org A = %d, want exactly 1", count)
		}

		// Old installations remain bound to the old licence, released with
		// actor, timestamp and reason; none is moved to the new licence.
		for _, number := range []string{i1.InstallationNumber, i2.InstallationNumber} {
			var stored licensingModel.Installation
			if err := fx.db.Where("number = ?", number).First(&stored).Error; err != nil {
				t.Fatalf("loading installation %s: %v", number, err)
			}
			if stored.LicenceNumber != l1.Number {
				t.Fatalf("installation %s moved to licence %s; history must stay bound", number, stored.LicenceNumber)
			}
			if stored.Status != licensingModel.InstallationStatusReleased {
				t.Fatalf("installation %s status = %s, want RELEASED", number, stored.Status)
			}
			if stored.ReleasedBy == nil || *stored.ReleasedBy != frOwner || stored.ReleasedAt == nil ||
				stored.ReleaseReason == nil || *stored.ReleaseReason != "package change" {
				t.Fatalf("installation %s release audit incomplete: %+v", number, stored)
			}
		}
		var transferred int64
		fx.db.Model(&licensingModel.Installation{}).Where("licence_number = ?", l2.Number).Count(&transferred)
		if transferred != 0 {
			t.Fatalf("new licence must start with 0 installations, found %d", transferred)
		}

		// The new key is returned once and never by ordinary reads.
		status, body = request(http.MethodGet, flRegistryURL+"/"+l2.Number, fx.tokenOwner, nil)
		if status != http.StatusOK || bytes.Contains(body, []byte(l2.LicenceKey)) {
			t.Fatalf("ordinary read must not expose the key (status %d): %s", status, string(body))
		}
	})

	t.Run("old_key_dead_and_fresh_activation_uses_new_capacity", func(t *testing.T) {
		if status, _ := activate(t, fx.tokenA, l1.LicenceKey, uuid.NewV4().String()); status != http.StatusConflict {
			t.Fatalf("activation with superseded key = %d, want %d", status, http.StatusConflict)
		}

		// The old computer must activate fresh, under the new key, creating a
		// new installation record; its old record stays with the old licence.
		status, body := activate(t, fx.tokenA, l2.LicenceKey, i1.InstallationIdentity)
		if status != http.StatusCreated {
			t.Fatalf("fresh activation under new licence = %d (%s)", status, string(body))
		}
		var envelope faActivationEnvelope
		if err := json.Unmarshal(body, &envelope); err != nil {
			t.Fatalf("decoding activation: %v", err)
		}
		if envelope.Data.InstallationNumber == i1.InstallationNumber {
			t.Fatal("fresh activation must create a new installation record")
		}
		if envelope.Data.CapacityUsed != 1 {
			t.Fatalf("new licence capacity used = %d, want 1 (old capacity must not transfer)", envelope.Data.CapacityUsed)
		}

		var oldRecord licensingModel.Installation
		fx.db.Where("number = ?", i1.InstallationNumber).First(&oldRecord)
		if oldRecord.LicenceNumber != l1.Number || oldRecord.Status != licensingModel.InstallationStatusReleased {
			t.Fatal("old installation history must remain unchanged")
		}

		// New capacity is its own: max 2 allows exactly one more activation.
		if status, _ := activate(t, fx.tokenA, l2.LicenceKey, uuid.NewV4().String()); status != http.StatusCreated {
			t.Fatalf("second fresh activation = %d, want %d", status, http.StatusCreated)
		}
		if status, _ := activate(t, fx.tokenA, l2.LicenceKey, uuid.NewV4().String()); status != http.StatusConflict {
			t.Fatalf("activation beyond new capacity = %d, want %d", status, http.StatusConflict)
		}
		if count := activeInstallations(t, l2.Number); count != 2 {
			t.Fatalf("active installations for new licence = %d, want 2", count)
		}
	})

	t.Run("security_and_validation", func(t *testing.T) {
		if status, _ := request(http.MethodPost, flRegistryURL+"/"+l2.Number+"/replace", fx.tokenA, replaceBody(frOrgA, from, to, 1, []string{"education"}, "client attempt")); status != http.StatusForbidden {
			t.Fatalf("client replacement = %d, want %d", status, http.StatusForbidden)
		}
		if status, _ := request(http.MethodPost, flRegistryURL+"/"+l2.Number+"/replace", fx.tokenAdmin, replaceBody(frOrgA, from, to, 1, []string{"education"}, "admin attempt")); status != http.StatusForbidden {
			t.Fatalf("admin without permission = %d, want %d", status, http.StatusForbidden)
		}
		if status, _ := replace(t, "does-not-exist", replaceBody(frOrgA, from, to, 1, []string{"education"}, "missing")); status != http.StatusNotFound {
			t.Fatalf("replacement of unknown licence = %d, want %d", status, http.StatusNotFound)
		}
		if status, _ := replace(t, l1.Number, replaceBody(frOrgA, from, to, 1, []string{"education"}, "already revoked")); status != http.StatusConflict {
			t.Fatalf("replacement of revoked licence = %d, want %d", status, http.StatusConflict)
		}

		// A missing reason fails validation and changes nothing.
		if status, _ := replace(t, l2.Number, replaceBody(frOrgA, from, to, 1, []string{"education"}, "   ")); status != http.StatusBadRequest {
			t.Fatalf("replacement without reason = %d, want %d", status, http.StatusBadRequest)
		}
		var stored licensingModel.Licence
		fx.db.Where("number = ?", l2.Number).First(&stored)
		if stored.Status != licensingModel.StatusActive {
			t.Fatal("failed replacement must leave the active licence unchanged")
		}
	})

	t.Run("renewal_regression", func(t *testing.T) {
		status, body := request(http.MethodPost, flRegistryURL+"/"+l2.Number+"/renew", fx.tokenOwner, map[string]any{
			"valid_to": time.Now().Add(120 * 24 * time.Hour).Format(time.RFC3339),
		})
		if status != http.StatusOK {
			t.Fatalf("renewal = %d (%s)", status, string(body))
		}

		var stored licensingModel.Licence
		if err := fx.db.Where("number = ?", l2.Number).First(&stored).Error; err != nil {
			t.Fatalf("loading renewed licence: %v", err)
		}
		if stored.Status != licensingModel.StatusActive || stored.LicenceKey != l2.LicenceKey || stored.MaxInstallations != 2 {
			t.Fatalf("renewal must preserve identity, key and capacity: %+v", stored)
		}
		if count := activeInstallations(t, l2.Number); count != 2 {
			t.Fatalf("renewal must preserve installations: %d active", count)
		}
	})

	t.Run("organisation_change", func(t *testing.T) {
		status, body := replace(t, l2.Number, replaceBody(frOrgB, from, time.Now().Add(90*24*time.Hour), 1, []string{"education"}, "organisation change"))
		if status != http.StatusCreated {
			t.Fatalf("organisation-change replacement = %d (%s)", status, string(body))
		}
		var issued flLicenceEnvelope
		if err := json.Unmarshal(body, &issued); err != nil {
			t.Fatalf("decoding replacement: %v", err)
		}
		l3 := issued.Data

		var stored licensingModel.Licence
		fx.db.Where("number = ?", l3.Number).First(&stored)
		if stored.OrganisationNumber != frOrgB || stored.Status != licensingModel.StatusActive {
			t.Fatalf("new licence must belong to organisation B: %+v", stored)
		}
		if count := activeLicences(t, frOrgB); count != 1 {
			t.Fatalf("active licences for org B = %d, want 1", count)
		}
		if count := activeInstallations(t, l2.Number); count != 0 {
			t.Fatalf("superseded licence must have all installations released: %d active", count)
		}

		if status, body := request(http.MethodGet, "/main/ems/licence", fx.tokenB, nil); status != http.StatusOK || !bytes.Contains(body, []byte(l3.Number)) {
			t.Fatalf("organisation B client status = %d (%s)", status, string(body))
		}
		if status, body := request(http.MethodGet, "/main/ems/licence", fx.tokenA, nil); status != http.StatusOK || !bytes.Contains(body, []byte("REVOKED")) {
			t.Fatalf("organisation A client status = %d (%s)", status, string(body))
		}
	})

	t.Run("concurrent_replacement_one_wins", func(t *testing.T) {
		cLicence := issue(t, frOrgC, from, to, 2, []string{"education"})
		for i := 0; i < 2; i++ {
			if status, body := activate(t, fx.tokenC, cLicence.LicenceKey, uuid.NewV4().String()); status != http.StatusCreated {
				t.Fatalf("seeding org C installation %d = %d (%s)", i, status, string(body))
			}
		}

		type result struct {
			status int
			err    error
		}
		results := make(chan result, 2)
		for i := 0; i < 2; i++ {
			go func(validTo time.Time, reason string) {
				raw, _ := json.Marshal(replaceBody(frOrgC, from, validTo, 2, []string{"education"}, reason))
				req, err := http.NewRequest(http.MethodPost, server.URL+flRegistryURL+"/"+cLicence.Number+"/replace", bytes.NewReader(raw))
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
			}(time.Now().Add(time.Duration(30+i*10)*24*time.Hour), "concurrent replacement")
		}

		created, conflicts := 0, 0
		for i := 0; i < 2; i++ {
			outcome := <-results
			if outcome.err != nil {
				t.Fatalf("concurrent replacement request failed: %v", outcome.err)
			}
			switch outcome.status {
			case http.StatusCreated:
				created++
			case http.StatusConflict:
				conflicts++
			}
		}
		if created != 1 || conflicts != 1 {
			t.Fatalf("concurrent replacements = %d created / %d conflicts, want 1/1", created, conflicts)
		}
		if count := activeLicences(t, frOrgC); count != 1 {
			t.Fatalf("active licences for org C = %d, want exactly 1", count)
		}
		if count := activeInstallations(t, cLicence.Number); count != 0 {
			t.Fatalf("superseded licence installations must all be released: %d active", count)
		}
		var suppressed licensingModel.Licence
		fx.db.Where("number = ?", cLicence.Number).First(&suppressed)
		if suppressed.Status != licensingModel.StatusRevoked {
			t.Fatalf("superseded licence status = %s, want REVOKED", suppressed.Status)
		}
		var winner licensingModel.Licence
		fx.db.Where("active_org_number = ?", frOrgC).First(&winner)
		var carried int64
		fx.db.Model(&licensingModel.Installation{}).Where("licence_number = ?", winner.Number).Count(&carried)
		if carried != 0 {
			t.Fatalf("winner must start with 0 installations, found %d", carried)
		}
	})
}
