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
	faOwner       = "ems-fa-owner"
	faAdminNoPerm = "ems-fa-admin-noperm"
	faClientA     = "ems-fa-client-a"
	faClientB     = "ems-fa-client-b"
	faClientC     = "ems-fa-client-c"
	faClientD     = "ems-fa-client-d"

	faOrgA = "ems-fa-org-a"
	faOrgB = "ems-fa-org-b"
	faOrgC = "ems-fa-org-c"
	faOrgD = "ems-fa-org-d"

	faRoleRegistry = "ems-fa-role-registry"
	faRoleClient   = "ems-fa-role-client"

	faActivateURL = "/main/ems/licence/activate"
	faClientInst  = "/main/ems/licence/installations"
	faRegistryIns = "/main/ems/licence-registry/installations"
)

type faFixture struct {
	db         *gorm.DB
	tokenOwner string
	tokenAdmin string
	tokenA     string
	tokenB     string
	tokenC     string
	tokenD     string
}

func setupFActivationFixture(t *testing.T) *faFixture {
	t.Helper()

	db := testutil.RequireTestDB(t)
	if err := migrations.Init(); err != nil {
		t.Fatalf("running migrations: %v", err)
	}
	t.Setenv("JWT_SECRET", "ems-foundation-f-activation-test-secret")

	users := []string{faOwner, faAdminNoPerm, faClientA, faClientB, faClientC, faClientD}
	organisations := []string{faOrgA, faOrgB, faOrgC, faOrgD}
	roles := []string{faRoleRegistry, faRoleClient}

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
		{faOwner, "User"},
		{faAdminNoPerm, "Admin"},
		{faClientA, "User"},
		{faClientB, "User"},
		{faClientC, "User"},
		{faClientD, "User"},
	} {
		user := authModel.User{
			Number:     spec.number,
			FirstName:  "FoundationF",
			SecondName: "Activation",
			Gender:     "M",
			Type:       spec.userType,
			Status:     "ACTIVE",
		}
		if err := db.Create(&user).Error; err != nil {
			t.Fatalf("seeding user %s: %v", spec.number, err)
		}
	}

	for _, spec := range []struct{ number, name string }{
		{faOrgA, "Foundation F Activation Org A"},
		{faOrgB, "Foundation F Activation Org B"},
		{faOrgC, "Foundation F Activation Org C"},
		{faOrgD, "Foundation F Activation Org D"},
	} {
		entity := eosModel.Entity{Number: spec.number, Name: spec.name, AddedBy: faOwner}
		if err := db.Create(&entity).Error; err != nil {
			t.Fatalf("seeding entity %s: %v", spec.number, err)
		}
	}

	for _, spec := range []struct{ userNumber, entityNumber string }{
		{faOwner, faOrgA},
		{faAdminNoPerm, faOrgA},
		{faClientA, faOrgA},
		{faClientB, faOrgB},
		{faClientC, faOrgC},
		{faClientD, faOrgD},
	} {
		active := eosModel.ActiveUserEntity{UserNumber: spec.userNumber, EntityNumber: spec.entityNumber}
		if err := db.Omit("Entity", "User").Create(&active).Error; err != nil {
			t.Fatalf("setting active entity for %s: %v", spec.userNumber, err)
		}
	}

	for _, number := range roles {
		role := permissionsModel.Role{
			Number: number, Name: number, EntityNumber: faOrgA, Status: "ACTIVE", UserNumber: faOwner,
		}
		if err := db.Omit("User", "Entity", "Permissions").Create(&role).Error; err != nil {
			t.Fatalf("seeding role %s: %v", number, err)
		}
	}

	for _, spec := range []struct{ roleNumber, permission string }{
		{faRoleRegistry, "manage_licence_registry"},
		{faRoleClient, "view_licence_status"},
		{faRoleClient, "activate_licence_installation"},
		{faRoleClient, "deactivate_own_installation"},
	} {
		rolePermission := permissionsModel.RolePermission{
			RoleNumber: spec.roleNumber, Permission: spec.permission, AddedByNumber: faOwner,
		}
		if err := db.Omit("AddedBy", "Role").Create(&rolePermission).Error; err != nil {
			t.Fatalf("seeding permission %s: %v", spec.permission, err)
		}
	}

	expiry := time.Now().Add(365 * 24 * time.Hour)
	for _, spec := range []struct{ userNumber, roleNumber string }{
		{faOwner, faRoleRegistry},
		{faClientA, faRoleClient},
		{faClientB, faRoleClient},
		{faClientC, faRoleClient},
		{faClientD, faRoleClient},
	} {
		userRole := permissionsModel.UserRole{
			UserNumber: spec.userNumber, RoleNumber: spec.roleNumber, ExpiryDate: expiry, AddedBy: faOwner,
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

	return &faFixture{
		db:         db,
		tokenOwner: token(faOwner),
		tokenAdmin: token(faAdminNoPerm),
		tokenA:     token(faClientA),
		tokenB:     token(faClientB),
		tokenC:     token(faClientC),
		tokenD:     token(faClientD),
	}
}

type faActivationData struct {
	InstallationNumber   string   `json:"installation_number"`
	InstallationIdentity string   `json:"installation_identity"`
	OrganisationNumber   string   `json:"organisation_number"`
	LicenceNumber        string   `json:"licence_number"`
	Status               string   `json:"status"`
	LicenceStatus        string   `json:"licence_status"`
	EffectiveStatus      string   `json:"effective_status"`
	MaxInstallations     int      `json:"max_installations"`
	CapacityUsed         int      `json:"capacity_used"`
	LicensedAddons       []string `json:"licensed_addons"`
}

type faInstallationData struct {
	Number               string     `json:"number"`
	InstallationIdentity string     `json:"installation_identity"`
	OrganisationNumber   string     `json:"organisation_number"`
	LicenceNumber        string     `json:"licence_number"`
	Status               string     `json:"status"`
	ActivatedAt          time.Time  `json:"activated_at"`
	LastSeenAt           *time.Time `json:"last_seen_at"`
	ClientDeactivatedAt  *time.Time `json:"client_deactivated_at"`
	ClientDeactivatedBy  string     `json:"client_deactivated_by"`
	ReleasedBy           string     `json:"released_by"`
	ReleasedAt           *time.Time `json:"released_at"`
	ReleaseReason        string     `json:"release_reason"`
	ReplacedByNumber     string     `json:"replaced_by_number"`
}

type faActivationEnvelope struct {
	Message string           `json:"message"`
	Data    faActivationData `json:"data"`
}

type faInstallationEnvelope struct {
	Message string             `json:"message"`
	Data    faInstallationData `json:"data"`
}

type faInstallationListEnvelope struct {
	Message string               `json:"message"`
	Data    []faInstallationData `json:"data"`
}

// TestFoundationFInstallationActivation covers Foundation F Increment 2 end to
// end with one shared fixture.
func TestFoundationFInstallationActivation(t *testing.T) {
	fx := setupFActivationFixture(t)

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

	decodeLicenceData := func(t *testing.T, body []byte) flLicenceData {
		t.Helper()
		var envelope flLicenceEnvelope
		if err := json.Unmarshal(body, &envelope); err != nil {
			t.Fatalf("decoding licence response %s: %v", string(body), err)
		}
		return envelope.Data
	}

	decodeActivation := func(t *testing.T, body []byte) faActivationData {
		t.Helper()
		var envelope faActivationEnvelope
		if err := json.Unmarshal(body, &envelope); err != nil {
			t.Fatalf("decoding activation response %s: %v", string(body), err)
		}
		return envelope.Data
	}

	decodeInstallation := func(t *testing.T, body []byte) faInstallationData {
		t.Helper()
		var envelope faInstallationEnvelope
		if err := json.Unmarshal(body, &envelope); err != nil {
			t.Fatalf("decoding installation response %s: %v", string(body), err)
		}
		return envelope.Data
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
		return decodeLicenceData(t, body)
	}

	activate := func(t *testing.T, token, key, identity, addon, replaces string) (int, []byte) {
		t.Helper()
		payload := map[string]any{
			"licence_key":           key,
			"installation_identity": identity,
		}
		if addon != "" {
			payload["addon_code"] = addon
		}
		if replaces != "" {
			payload["replaces_installation_number"] = replaces
		}
		return request(http.MethodPost, faActivateURL, token, payload)
	}

	release := func(t *testing.T, installationNumber, reason string) (int, []byte) {
		t.Helper()
		return request(http.MethodPost, faRegistryIns+"/"+installationNumber+"/release", fx.tokenOwner, map[string]any{
			"reason": reason,
		})
	}

	activeCount := func(t *testing.T, licenceNumber string) int64 {
		t.Helper()
		var count int64
		fx.db.Model(&licensingModel.Installation{}).
			Where("licence_number = ? AND status = ?", licenceNumber, licensingModel.InstallationStatusActive).
			Count(&count)
		return count
	}

	from := time.Now().Add(-time.Hour)
	to := time.Now().Add(30 * 24 * time.Hour)

	var aLicence flLicenceData
	var a1, a2, a3, a4, a5, addonInstallation faActivationData
	a1Identity := uuid.NewV4().String()
	a2Identity := uuid.NewV4().String()
	a3Identity := uuid.NewV4().String()
	a4Identity := uuid.NewV4().String()
	addonIdentity := uuid.NewV4().String()

	t.Run("unauthenticated_returns_401", func(t *testing.T) {
		if status, _ := request(http.MethodPost, faActivateURL, "", nil); status != http.StatusUnauthorized {
			t.Fatalf("activate status = %d, want %d", status, http.StatusUnauthorized)
		}
		if status, _ := request(http.MethodGet, faClientInst, "", nil); status != http.StatusUnauthorized {
			t.Fatalf("installations status = %d, want %d", status, http.StatusUnauthorized)
		}
	})

	t.Run("activation_succeeds_and_masks_key", func(t *testing.T) {
		aLicence = issue(t, faOrgA, from, to, 3, []string{"education"})

		status, body := activate(t, fx.tokenA, aLicence.LicenceKey, a1Identity, "", "")
		if status != http.StatusCreated {
			t.Fatalf("activation = %d, want %d (%s)", status, http.StatusCreated, string(body))
		}
		if bytes.Contains(body, []byte(aLicence.LicenceKey)) {
			t.Fatalf("activation response must not expose the licence key: %s", string(body))
		}
		a1 = decodeActivation(t, body)
		if a1.InstallationNumber == "" || a1.Status != "ACTIVE" || a1.CapacityUsed != 1 {
			t.Fatalf("unexpected activation result: %s", string(body))
		}

		var stored licensingModel.Installation
		if err := fx.db.Where("number = ?", a1.InstallationNumber).First(&stored).Error; err != nil {
			t.Fatalf("loading installation: %v", err)
		}
		if stored.InstallationIdentity != a1Identity || stored.OrganisationNumber != faOrgA ||
			stored.LicenceNumber != aLicence.Number || stored.Status != licensingModel.InstallationStatusActive {
			t.Fatalf("unexpected stored installation: %+v", stored)
		}
	})

	t.Run("invalid_key_rejected", func(t *testing.T) {
		if status, _ := activate(t, fx.tokenA, "not-a-real-key", uuid.NewV4().String(), "", ""); status != http.StatusBadRequest {
			t.Fatalf("invalid key = %d, want %d", status, http.StatusBadRequest)
		}
	})

	t.Run("wrong_organisation_rejected", func(t *testing.T) {
		// The key exists but belongs to organisation A: organisation B must
		// receive the same non-disclosing rejection as an unknown key.
		before := activeCount(t, aLicence.Number)
		status, body := activate(t, fx.tokenB, aLicence.LicenceKey, uuid.NewV4().String(), "", "")
		if status != http.StatusBadRequest {
			t.Fatalf("wrong organisation = %d, want %d (%s)", status, http.StatusBadRequest, string(body))
		}
		if after := activeCount(t, aLicence.Number); after != before {
			t.Fatalf("wrong-organisation activation changed capacity: %d -> %d", before, after)
		}
	})

	t.Run("idempotent_same_identity", func(t *testing.T) {
		status, body := activate(t, fx.tokenA, aLicence.LicenceKey, a1Identity, "", "")
		if status != http.StatusCreated {
			t.Fatalf("idempotent activation = %d, want %d (%s)", status, http.StatusCreated, string(body))
		}
		if again := decodeActivation(t, body); again.InstallationNumber != a1.InstallationNumber || again.CapacityUsed != 1 {
			t.Fatalf("same identity must be capacity-neutral: %s", string(body))
		}
		var rows int64
		fx.db.Model(&licensingModel.Installation{}).
			Where("licence_number = ? AND installation_identity = ?", aLicence.Number, a1Identity).
			Count(&rows)
		if rows != 1 {
			t.Fatalf("same identity created duplicate records: %d", rows)
		}
	})

	t.Run("unlicensed_addon_rejected", func(t *testing.T) {
		if status, _ := activate(t, fx.tokenA, aLicence.LicenceKey, uuid.NewV4().String(), "core", ""); status != http.StatusForbidden {
			t.Fatalf("unlicensed addon = %d, want %d", status, http.StatusForbidden)
		}

		status, body := activate(t, fx.tokenA, aLicence.LicenceKey, addonIdentity, "education", "")
		if status != http.StatusCreated {
			t.Fatalf("licensed addon activation = %d, want %d (%s)", status, http.StatusCreated, string(body))
		}
		addonInstallation = decodeActivation(t, body)
		if status, body := release(t, addonInstallation.InstallationNumber, "test cleanup"); status != http.StatusOK {
			t.Fatalf("releasing addon installation = %d, want %d (%s)", status, http.StatusOK, string(body))
		}
	})

	t.Run("capacity_rules_and_owner_release", func(t *testing.T) {
		status, body := activate(t, fx.tokenA, aLicence.LicenceKey, a2Identity, "", "")
		if status != http.StatusCreated {
			t.Fatalf("second activation = %d (%s)", status, string(body))
		}
		a2 = decodeActivation(t, body)

		status, body = activate(t, fx.tokenA, aLicence.LicenceKey, a3Identity, "", "")
		if status != http.StatusCreated {
			t.Fatalf("third activation = %d (%s)", status, string(body))
		}
		a3 = decodeActivation(t, body)
		if a3.CapacityUsed != 3 {
			t.Fatalf("capacity used = %d, want 3", a3.CapacityUsed)
		}

		status, body = activate(t, fx.tokenA, aLicence.LicenceKey, a4Identity, "", "")
		if status != http.StatusConflict {
			t.Fatalf("fourth activation = %d, want %d (%s)", status, http.StatusConflict, string(body))
		}

		// Client self-deactivation must NOT free capacity.
		status, body = request(http.MethodPost, faClientInst+"/"+a2.InstallationNumber+"/deactivate", fx.tokenA, nil)
		if status != http.StatusOK {
			t.Fatalf("self-deactivation = %d, want %d (%s)", status, http.StatusOK, string(body))
		}
		deactivated := decodeInstallation(t, body)
		if deactivated.ClientDeactivatedAt == nil || deactivated.ClientDeactivatedBy != faClientA {
			t.Fatalf("self-deactivation must record actor/time: %s", string(body))
		}
		if status, _ = activate(t, fx.tokenA, aLicence.LicenceKey, a4Identity, "", ""); status != http.StatusConflict {
			t.Fatalf("activation after self-deactivation = %d, want %d (capacity is not released)", status, http.StatusConflict)
		}
		if count := activeCount(t, aLicence.Number); count != 3 {
			t.Fatalf("client-deactivated installation must keep consuming capacity: %d", count)
		}

		// Owner release DOES free capacity.
		status, body = release(t, a2.InstallationNumber, "computer replaced")
		if status != http.StatusOK {
			t.Fatalf("owner release = %d, want %d (%s)", status, http.StatusOK, string(body))
		}
		released := decodeInstallation(t, body)
		if released.Status != "RELEASED" || released.ReleasedBy != faOwner || released.ReleasedAt == nil || released.ReleaseReason != "computer replaced" {
			t.Fatalf("release must record actor/time/reason: %s", string(body))
		}
		if count := activeCount(t, aLicence.Number); count != 2 {
			t.Fatalf("owner release must free capacity: %d", count)
		}

		status, body = activate(t, fx.tokenA, aLicence.LicenceKey, a4Identity, "", "")
		if status != http.StatusCreated {
			t.Fatalf("activation after release = %d, want %d (%s)", status, http.StatusCreated, string(body))
		}
		a4 = decodeActivation(t, body)
	})

	t.Run("self_deactivation_reactivation", func(t *testing.T) {
		status, body := request(http.MethodPost, faClientInst+"/"+a3.InstallationNumber+"/deactivate", fx.tokenA, nil)
		if status != http.StatusOK {
			t.Fatalf("deactivation = %d (%s)", status, string(body))
		}
		capacityDuringDeactivation := activeCount(t, aLicence.Number)

		// Reactivating the same identity is capacity-neutral.
		status, body = activate(t, fx.tokenA, aLicence.LicenceKey, a3Identity, "", "")
		if status != http.StatusCreated {
			t.Fatalf("reactivation = %d, want %d (%s)", status, http.StatusCreated, string(body))
		}
		reactivated := decodeActivation(t, body)
		if reactivated.InstallationNumber != a3.InstallationNumber || reactivated.CapacityUsed != int(capacityDuringDeactivation) {
			t.Fatalf("reactivation must be capacity-neutral: %s", string(body))
		}

		var stored licensingModel.Installation
		if err := fx.db.Where("number = ?", a3.InstallationNumber).First(&stored).Error; err != nil {
			t.Fatalf("loading reactivated installation: %v", err)
		}
		if stored.ClientDeactivatedAt != nil || stored.ClientDeactivatedBy != nil {
			t.Fatal("reactivation must clear the client-deactivation marker")
		}
		var rows int64
		fx.db.Model(&licensingModel.Installation{}).
			Where("licence_number = ? AND installation_identity = ?", aLicence.Number, a3Identity).
			Count(&rows)
		if rows != 1 {
			t.Fatalf("reactivation created duplicate records: %d", rows)
		}
	})

	t.Run("replacement_link_history", func(t *testing.T) {
		// Owner releases a computer; the replacement links to it explicitly.
		if status, body := release(t, a4.InstallationNumber, "machine retired"); status != http.StatusOK {
			t.Fatalf("release before replacement = %d (%s)", status, string(body))
		}

		// Invalid replacement links fail without creating an installation.
		if status, _ := activate(t, fx.tokenA, aLicence.LicenceKey, uuid.NewV4().String(), "", "does-not-exist"); status != http.StatusBadRequest {
			t.Fatalf("invalid replacement link = %d, want %d", status, http.StatusBadRequest)
		}

		status, body := activate(t, fx.tokenA, aLicence.LicenceKey, uuid.NewV4().String(), "", a4.InstallationNumber)
		if status != http.StatusCreated {
			t.Fatalf("replacement activation = %d, want %d (%s)", status, http.StatusCreated, string(body))
		}
		a5 = decodeActivation(t, body)

		var old licensingModel.Installation
		if err := fx.db.Where("number = ?", a4.InstallationNumber).First(&old).Error; err != nil {
			t.Fatalf("loading released installation: %v", err)
		}
		if old.Status != licensingModel.InstallationStatusReleased || old.ReplacedByNumber == nil || *old.ReplacedByNumber != a5.InstallationNumber {
			t.Fatalf("replacement history link not recorded: %+v", old)
		}
		if old.LicenceNumber != aLicence.Number {
			t.Fatal("released installation must remain bound to its original licence")
		}

		// A released identity is not silently re-permitted.
		if status, _ := activate(t, fx.tokenA, aLicence.LicenceKey, a2Identity, "", ""); status != http.StatusConflict {
			t.Fatalf("released identity reactivation = %d, want %d", status, http.StatusConflict)
		}
	})

	t.Run("client_cannot_release_or_touch_other_org", func(t *testing.T) {
		// Clients cannot release capacity.
		if status, _ := request(http.MethodPost, faRegistryIns+"/"+a1.InstallationNumber+"/release", fx.tokenA, map[string]any{"reason": "self release"}); status != http.StatusForbidden {
			t.Fatalf("client release = %d, want %d", status, http.StatusForbidden)
		}
		// An Admin user-type without the explicit licensing permission cannot
		// release either.
		if status, _ := request(http.MethodPost, faRegistryIns+"/"+a1.InstallationNumber+"/release", fx.tokenAdmin, map[string]any{"reason": "admin release"}); status != http.StatusForbidden {
			t.Fatalf("admin without permission release = %d, want %d", status, http.StatusForbidden)
		}
		// Cross-organisation deactivation is invisible (404).
		if status, _ := request(http.MethodPost, faClientInst+"/"+a1.InstallationNumber+"/deactivate", fx.tokenB, nil); status != http.StatusNotFound {
			t.Fatalf("cross-organisation deactivation = %d, want %d", status, http.StatusNotFound)
		}
	})

	t.Run("suspended_revoked_expired_and_grace", func(t *testing.T) {
		bLicence := issue(t, faOrgB, from, to, 1, []string{"education"})
		bIdentity := uuid.NewV4().String()
		if status, body := activate(t, fx.tokenB, bLicence.LicenceKey, bIdentity, "", ""); status != http.StatusCreated {
			t.Fatalf("org B activation = %d (%s)", status, string(body))
		}

		if status, body := request(http.MethodPost, flRegistryURL+"/"+bLicence.Number+"/suspend", fx.tokenOwner, nil); status != http.StatusOK {
			t.Fatalf("suspending org B licence = %d (%s)", status, string(body))
		}
		if status, _ := activate(t, fx.tokenB, bLicence.LicenceKey, uuid.NewV4().String(), "", ""); status != http.StatusConflict {
			t.Fatalf("activation on suspended licence = %d, want %d", status, http.StatusConflict)
		}

		if status, body := request(http.MethodPost, flRegistryURL+"/"+bLicence.Number+"/revoke", fx.tokenOwner, nil); status != http.StatusOK {
			t.Fatalf("revoking org B licence = %d (%s)", status, string(body))
		}
		if status, _ := activate(t, fx.tokenB, bLicence.LicenceKey, uuid.NewV4().String(), "", ""); status != http.StatusConflict {
			t.Fatalf("activation on revoked licence = %d, want %d", status, http.StatusConflict)
		}

		// A licence inside the 14-day grace window still activates (F-Q2).
		graceLicence := issue(t, faOrgB, time.Now().Add(-48*time.Hour), time.Now().Add(-time.Hour), 1, []string{"education"})
		if status, body := activate(t, fx.tokenB, graceLicence.LicenceKey, uuid.NewV4().String(), "", ""); status != http.StatusCreated {
			t.Fatalf("activation in grace = %d, want %d (%s)", status, http.StatusCreated, string(body))
		}

		// Expired beyond grace fails closed.
		expiredLicence := issue(t, faOrgD, time.Now().Add(-40*24*time.Hour), time.Now().Add(-30*24*time.Hour), 1, []string{"education"})
		if status, _ := activate(t, fx.tokenD, expiredLicence.LicenceKey, uuid.NewV4().String(), "", ""); status != http.StatusConflict {
			t.Fatalf("activation on expired licence = %d, want %d", status, http.StatusConflict)
		}
	})

	t.Run("concurrent_activation_cannot_exceed_capacity", func(t *testing.T) {
		cLicence := issue(t, faOrgC, from, to, 3, []string{"education"})

		type result struct {
			status int
			err    error
		}
		results := make(chan result, 4)
		for i := 0; i < 4; i++ {
			go func() {
				identity := uuid.NewV4().String()
				raw, _ := json.Marshal(map[string]any{
					"licence_key":           cLicence.LicenceKey,
					"installation_identity": identity,
				})
				req, err := http.NewRequest(http.MethodPost, server.URL+faActivateURL, bytes.NewReader(raw))
				if err != nil {
					results <- result{err: err}
					return
				}
				req.Header.Set("Authorization", "Bearer "+fx.tokenC)
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
		for i := 0; i < 4; i++ {
			outcome := <-results
			if outcome.err != nil {
				t.Fatalf("concurrent activation request failed: %v", outcome.err)
			}
			switch outcome.status {
			case http.StatusCreated:
				created++
			case http.StatusConflict:
				conflicts++
			}
		}
		if created != 3 || conflicts != 1 {
			t.Fatalf("concurrent activations = %d created / %d conflicts, want 3/1", created, conflicts)
		}
		if count := activeCount(t, cLicence.Number); count != 3 {
			t.Fatalf("active installations = %d, want 3", count)
		}
	})

	t.Run("isolation_release_reason_and_key_masking", func(t *testing.T) {
		// Client B sees only its own installations.
		status, body := request(http.MethodGet, faClientInst, fx.tokenB, nil)
		if status != http.StatusOK {
			t.Fatalf("client B installations = %d (%s)", status, string(body))
		}
		if bytes.Contains(body, []byte(a1.InstallationNumber)) || bytes.Contains(body, []byte(aLicence.Number)) {
			t.Fatalf("client B must not see organisation A installations: %s", string(body))
		}

		// Client reads never expose the licence key.
		status, body = request(http.MethodGet, faClientInst, fx.tokenA, nil)
		if status != http.StatusOK {
			t.Fatalf("client A installations = %d (%s)", status, string(body))
		}
		if bytes.Contains(body, []byte(aLicence.LicenceKey)) {
			t.Fatalf("client installation read must not expose the licence key: %s", string(body))
		}

		// Control-plane reads never expose the licence key either.
		status, body = request(http.MethodGet, faRegistryIns, fx.tokenOwner, nil)
		if status != http.StatusOK {
			t.Fatalf("registry installations = %d (%s)", status, string(body))
		}
		if bytes.Contains(body, []byte(aLicence.LicenceKey)) {
			t.Fatalf("control-plane installation read must not expose the licence key: %s", string(body))
		}

		status, body = request(http.MethodGet, faRegistryIns+"/"+a1.InstallationNumber, fx.tokenOwner, nil)
		if status != http.StatusOK {
			t.Fatalf("registry installation detail = %d (%s)", status, string(body))
		}
		if bytes.Contains(body, []byte(aLicence.LicenceKey)) {
			t.Fatalf("control-plane installation detail must not expose the licence key: %s", string(body))
		}

		// Release requires an explicit reason.
		if status, _ := release(t, a1.InstallationNumber, "   "); status != http.StatusBadRequest {
			t.Fatalf("release without reason = %d, want %d", status, http.StatusBadRequest)
		}

		// Cross-organisation control-plane access is impossible for clients.
		if status, _ := request(http.MethodGet, faRegistryIns, fx.tokenA, nil); status != http.StatusForbidden {
			t.Fatalf("client registry read = %d, want %d", status, http.StatusForbidden)
		}
	})
}
