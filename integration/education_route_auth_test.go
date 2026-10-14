package integration

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/routes"
	authServices "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/services"
	licensingModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/licensing/models"
	eosModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/eos/models"
)

// TestEducationPermissionEnforcedBeforeHandler is the regression test for the
// middleware-ordering hotfix.
//
// A non-admin user without the required permission must receive 403 before the
// handler runs. If the handler ran first it would write 200 and Gin would log
// "Headers were already written. Wanted to override status code 200 with 403",
// so a clean 403 proves the permission middleware now precedes the handler.
//
// Foundation F Increment 3 adds the licence entitlement boundary upstream of
// the permission middleware, so the fixture provides the user with an active,
// entitled licence: the request passes the licence gate and still reaches the
// permission middleware, keeping this test's assertion (403 before the
// handler) meaningful. No production security behaviour is changed.
func TestEducationPermissionEnforcedBeforeHandler(t *testing.T) {
	db := seedEmsUsers(t)
	t.Setenv("JWT_SECRET", "ems-foundation-a-test-secret")

	const (
		orgNumber     = "ems-foundation-a-org"
		licenceNumber = "ems-foundation-a-licence"
	)
	now := time.Now()
	activeOrg := orgNumber
	entity := eosModel.Entity{Number: orgNumber, Name: "EMS Foundation A Org", AddedBy: emsAdminNumber}
	licence := licensingModel.Licence{
		Number:             licenceNumber,
		LicenceKey:         "ems-foundation-a-licence-key",
		OrganisationNumber: orgNumber,
		ValidFrom:          now.Add(-time.Hour),
		ValidTo:            now.Add(24 * time.Hour),
		GraceDays:          licensingModel.DefaultGraceDays,
		MaxInstallations:   1,
		LicensedAddons:     `["education"]`,
		Status:             licensingModel.StatusActive,
		ActiveOrgNumber:    &activeOrg,
		IssuedBy:           emsAdminNumber,
	}
	active := eosModel.ActiveUserEntity{UserNumber: emsUserNumber, EntityNumber: orgNumber}

	// Remove any residue from an earlier interrupted run before seeding.
	db.Unscoped().Where("number = ?", licenceNumber).Delete(&licensingModel.Licence{})
	db.Unscoped().Where("user_number = ?", emsUserNumber).Delete(&eosModel.ActiveUserEntity{})
	db.Unscoped().Where("number = ?", orgNumber).Delete(&eosModel.Entity{})

	if err := db.Create(&entity).Error; err != nil {
		t.Fatalf("seeding entity: %v", err)
	}
	if err := db.Omit("Entity", "User").Create(&active).Error; err != nil {
		t.Fatalf("setting active entity: %v", err)
	}
	if err := db.Create(&licence).Error; err != nil {
		t.Fatalf("seeding licence: %v", err)
	}
	t.Cleanup(func() {
		db.Unscoped().Where("number = ?", licenceNumber).Delete(&licensingModel.Licence{})
		db.Unscoped().Where("user_number = ?", emsUserNumber).Delete(&eosModel.ActiveUserEntity{})
		db.Unscoped().Where("number = ?", orgNumber).Delete(&eosModel.Entity{})
	})

	server := httptest.NewServer(routes.SetupRoutes())
	t.Cleanup(server.Close)
	client := server.Client()

	token, err := authServices.JWTAuthService().GenerateToken(emsUserNumber)
	if err != nil {
		t.Fatalf("generating non-admin token: %v", err)
	}

	req, err := http.NewRequest(http.MethodGet, server.URL+"/main/curriculum/all", nil)
	if err != nil {
		t.Fatalf("building request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusForbidden {
		t.Fatalf("GET /main/curriculum/all = %d, want %d; the permission middleware must run before the handler",
			resp.StatusCode, http.StatusForbidden)
	}
}
