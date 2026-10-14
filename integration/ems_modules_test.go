package integration

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/routes"
	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	authServices "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/services"
	addonModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/addons/models"
	addonServices "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/addons/services"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/internal/testutil"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/migrations"
	"gorm.io/gorm"
)

const (
	emsAdminNumber = "ems-foundation-a-admin"
	emsUserNumber  = "ems-foundation-a-user"
)

// seedEmsUsers ensures a migrated schema and two users: one Admin and one
// non-admin with no permissions. Returns the database handle.
func seedEmsUsers(t *testing.T) *gorm.DB {
	t.Helper()

	db := testutil.RequireTestDB(t)

	if err := migrations.Init(); err != nil {
		t.Fatalf("running migrations: %v", err)
	}

	cleanup := func() {
		db.Unscoped().Where("number = ?", emsAdminNumber).Delete(&authModel.User{})
		db.Unscoped().Where("number = ?", emsUserNumber).Delete(&authModel.User{})
	}
	cleanup()
	t.Cleanup(cleanup)

	admin := authModel.User{
		Number:     emsAdminNumber,
		FirstName:  "Ems",
		SecondName: "Admin",
		Gender:     "M",
		Type:       "Admin",
		Status:     "ACTIVE",
	}
	if err := db.Create(&admin).Error; err != nil {
		t.Fatalf("seeding admin user: %v", err)
	}

	nonAdmin := authModel.User{
		Number:     emsUserNumber,
		FirstName:  "Ems",
		SecondName: "User",
		Gender:     "M",
		Type:       "User",
		Status:     "ACTIVE",
	}
	if err := db.Create(&nonAdmin).Error; err != nil {
		t.Fatalf("seeding non-admin user: %v", err)
	}

	return db
}

// TestEmsModuleRegistrySyncIsIdempotent proves repeated synchronisation creates
// no duplicates and never unexpectedly disables an addon.
func TestEmsModuleRegistrySyncIsIdempotent(t *testing.T) {
	db := seedEmsUsers(t)

	if err := addonServices.SyncCatalog(); err != nil {
		t.Fatalf("first sync: %v", err)
	}
	if err := addonServices.SyncCatalog(); err != nil {
		t.Fatalf("second sync: %v", err)
	}

	var educationCount int64
	if err := db.Model(&addonModel.Module{}).Where("code = ?", "education").Count(&educationCount).Error; err != nil {
		t.Fatalf("counting education rows: %v", err)
	}
	if educationCount != 1 {
		t.Fatalf("education rows = %d, want 1 after repeated sync", educationCount)
	}

	var education addonModel.Module
	if err := db.Where("code = ?", "education").First(&education).Error; err != nil {
		t.Fatalf("loading education module: %v", err)
	}
	if education.DependsOn != "core" {
		t.Errorf("education depends_on = %q, want %q", education.DependsOn, "core")
	}
	if !education.Enabled {
		t.Errorf("education should be enabled by default after first sync")
	}

	// Disabling must survive a re-sync.
	if _, err := addonServices.SetModuleEnabled("education", false); err != nil {
		t.Fatalf("disabling education: %v", err)
	}
	if err := addonServices.SyncCatalog(); err != nil {
		t.Fatalf("sync after disable: %v", err)
	}

	var afterSync addonModel.Module
	if err := db.Where("code = ?", "education").First(&afterSync).Error; err != nil {
		t.Fatalf("reloading education module: %v", err)
	}
	if afterSync.Enabled {
		t.Fatal("sync unexpectedly re-enabled a disabled addon")
	}

	// Restore default state for other tests.
	if _, err := addonServices.SetModuleEnabled("education", true); err != nil {
		t.Fatalf("re-enabling education: %v", err)
	}
}

// TestEmsModulesEndpointAuthorization covers 401 (no token), 403 (authenticated
// without permission) and 200 (authorised).
func TestEmsModulesEndpointAuthorization(t *testing.T) {
	seedEmsUsers(t)
	t.Setenv("JWT_SECRET", "ems-foundation-a-test-secret")

	server := httptest.NewServer(routes.SetupRoutes())
	t.Cleanup(server.Close)
	client := server.Client()

	endpoint := server.URL + "/main/ems/modules"

	call := func(authorization string) int {
		req, err := http.NewRequest(http.MethodGet, endpoint, nil)
		if err != nil {
			t.Fatalf("building request: %v", err)
		}
		if authorization != "" {
			req.Header.Set("Authorization", authorization)
		}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		defer resp.Body.Close()
		return resp.StatusCode
	}

	if status := call(""); status != http.StatusUnauthorized {
		t.Errorf("unauthenticated status = %d, want %d", status, http.StatusUnauthorized)
	}

	nonAdminToken, err := authServices.JWTAuthService().GenerateToken(emsUserNumber)
	if err != nil {
		t.Fatalf("generating non-admin token: %v", err)
	}
	if status := call("Bearer " + nonAdminToken); status != http.StatusForbidden {
		t.Errorf("unauthorised status = %d, want %d", status, http.StatusForbidden)
	}

	adminToken, err := authServices.JWTAuthService().GenerateToken(emsAdminNumber)
	if err != nil {
		t.Fatalf("generating admin token: %v", err)
	}
	if status := call("Bearer " + adminToken); status != http.StatusOK {
		t.Errorf("authorised status = %d, want %d", status, http.StatusOK)
	}
}
