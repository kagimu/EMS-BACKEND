package integration

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/routes"
	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	licensingModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/licensing/models"
	organisationModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/organisation/models"
	eosModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/eos/models"
)

func TestStage6A2ConcurrentSiteCapacity(t *testing.T) {
	fx := setupStage0BFixture(t)
	licenceNumber := fmt.Sprintf("s6a2c-%d", time.Now().UnixNano()%10000000000)
	seedStage1Licence(t, fx, s0bOrgA, licenceNumber)
	t.Setenv(commons.SoftwareOwnerUserNumbersEnv, s0bAdminA)
	now := time.Now()
	validFrom := now.Add(-time.Hour).Format(time.RFC3339)
	validTo := now.Add(48 * time.Hour).Format(time.RFC3339)

	server := httptest.NewServer(routes.SetupRoutes())
	t.Cleanup(server.Close)
	client := server.Client()
	req := func(method, path, token string, body any) (int, []byte) {
		return stage1Request(t, client, server.URL, method, path, token, body)
	}
	if status, body := req(http.MethodPut, "/main/ems/owner/licences/"+licenceNumber+"/flexible", fx.tokenAdminA, map[string]any{
		"scopes":       []map[string]any{{"scope_type": "SITE", "quantity": 1, "status": "ACTIVE", "valid_from": validFrom, "valid_to": validTo}},
		"entitlements": []any{}, "reason": "6A.2 concurrency",
	}); status != http.StatusOK {
		t.Fatalf("configure = %d (%s)", status, string(body))
	}
	var configuredScope licensingModel.LicenceScope
	fx.db.Where("licence_number = ? AND scope_type = ?", licenceNumber, licensingModel.ScopeSite).First(&configuredScope)
	t.Logf("configured site scope quantity=%d status=%s", configuredScope.Quantity, configuredScope.Status)
	t.Cleanup(func() {
		fx.db.Unscoped().Where("licence_number = ?", licenceNumber).Delete(&licensingModel.LicensedUnit{})
		fx.db.Unscoped().Where("licence_number = ?", licenceNumber).Delete(&licensingModel.LicenceScope{})
		fx.db.Unscoped().Where("licence_number = ?", licenceNumber).Delete(&licensingModel.LicenceEntitlement{})
		fx.db.Unscoped().Where("licence_number = ?", licenceNumber).Delete(&licensingModel.LicenceEvent{})
		fx.db.Unscoped().Where("organisation_number = ? AND code IN ?", s0bOrgA, []string{"s6a2-concurrent-a", "s6a2-concurrent-b"}).Delete(&organisationModel.OrganisationUnit{})
		fx.db.Unscoped().Where("number = ?", licenceNumber).Delete(&licensingModel.Licence{})
	})

	type result struct {
		status int
		body   string
	}
	results := make(chan result, 2)
	var wait sync.WaitGroup
	for _, site := range []struct{ name, code string }{
		{"Concurrent Site A", "s6a2-concurrent-a"},
		{"Concurrent Site B", "s6a2-concurrent-b"},
	} {
		wait.Add(1)
		go func(name, code string) {
			defer wait.Done()
			status, _ := req(http.MethodPost, "/main/ems/organisation-units", fx.tokenAdminA, map[string]any{"name": name, "unit_type": "site", "code": code})
			results <- result{status: status}
		}(site.name, site.code)
	}
	wait.Wait()
	close(results)
	success, conflict := 0, 0
	statuses := []int{}
	for item := range results {
		statuses = append(statuses, item.status)
		t.Logf("concurrent response status=%d body=%s", item.status, item.body)
		switch item.status {
		case http.StatusCreated:
			success++
		case http.StatusConflict:
			conflict++
		}
	}
	if success != 1 || conflict != 1 {
		t.Fatalf("concurrent results success=%d conflict=%d statuses=%v, want 1/1", success, conflict, statuses)
	}
	var activeUnits int64
	fx.db.Model(&licensingModel.LicensedUnit{}).Where("licence_number = ? AND unit_type = ? AND status = ?", licenceNumber, licensingModel.ScopeSite, licensingModel.UnitStatusActive).Count(&activeUnits)
	if activeUnits != 1 {
		t.Fatalf("active site units = %d, want 1", activeUnits)
	}
	var activeOrgUnits int64
	fx.db.Model(&organisationModel.OrganisationUnit{}).Where("organisation_number = ? AND code IN ?", s0bOrgA, []string{"s6a2-concurrent-a", "s6a2-concurrent-b"}).Count(&activeOrgUnits)
	if activeOrgUnits != 1 {
		t.Fatalf("active concurrent organisation units = %d, want 1", activeOrgUnits)
	}
	var events int64
	fx.db.Model(&licensingModel.LicenceEvent{}).Where("licence_number = ? AND reason = ?", licenceNumber, "SITE_REGISTERED").Count(&events)
	if events != 1 {
		t.Fatalf("site registration events = %d, want 1", events)
	}
}

func TestStage6A2ReplacementInsufficientCapacityRollsBack(t *testing.T) {
	fx := setupStage0BFixture(t)
	fx.db.Model(&licensingModel.Licence{}).Where("organisation_number = ? AND active_org_number IS NOT NULL", s0bOrgA).Updates(map[string]any{"status": licensingModel.StatusRevoked, "active_org_number": nil})
	licenceNumber := fmt.Sprintf("s6a2r-%d", time.Now().UnixNano()%10000000000)
	seedStage1Licence(t, fx, s0bOrgA, licenceNumber)
	t.Setenv(commons.SoftwareOwnerUserNumbersEnv, s0bAdminA)
	now := time.Now()
	validFrom := now.Add(-time.Hour).Format(time.RFC3339)
	validTo := now.Add(48 * time.Hour).Format(time.RFC3339)

	secondUser := authModel.User{Number: fmt.Sprintf("s6a2u-%d", time.Now().UnixNano()%10000000000), FirstName: "Second", SecondName: "Seat", LastName: "User", Gender: "M", Type: "Staff", Status: "Active"}
	if err := fx.db.Create(&secondUser).Error; err != nil {
		t.Fatalf("seed user: %v", err)
	}
	t.Cleanup(func() {
		fx.db.Unscoped().Where("user_number IN ?", []string{secondUser.Number, s0bUserOnlyB}).Where("entity_number = ?", s0bOrgA).Delete(&eosModel.EntityUser{})
		fx.db.Unscoped().Where("number = ?", secondUser.Number).Delete(&authModel.User{})
		fx.db.Unscoped().Where("licence_number = ?", licenceNumber).Delete(&licensingModel.LicensedUnit{})
		fx.db.Unscoped().Where("licence_number = ?", licenceNumber).Delete(&licensingModel.LicenceScope{})
		fx.db.Unscoped().Where("licence_number = ?", licenceNumber).Delete(&licensingModel.LicenceEntitlement{})
		fx.db.Unscoped().Where("licence_number = ?", licenceNumber).Delete(&licensingModel.LicenceEvent{})
	})

	server := httptest.NewServer(routes.SetupRoutes())
	t.Cleanup(server.Close)
	client := server.Client()
	req := func(method, path, token string, body any) (int, []byte) {
		return stage1Request(t, client, server.URL, method, path, token, body)
	}
	if status, body := req(http.MethodPut, "/main/ems/owner/licences/"+licenceNumber+"/flexible", fx.tokenAdminA, map[string]any{
		"scopes":       []map[string]any{{"scope_type": "USER_SEAT", "quantity": 2, "status": "ACTIVE", "valid_from": validFrom, "valid_to": validTo}},
		"entitlements": []any{}, "reason": "6A.2 setup",
	}); status != http.StatusOK {
		t.Fatalf("configure = %d (%s)", status, string(body))
	}
	for _, user := range []string{s0bUserOnlyB, secondUser.Number} {
		if status, body := req(http.MethodPost, "/main/eos/add-user-to-entity", fx.tokenAdminA, map[string]string{"entity_number": s0bOrgA, "user_number": user}); status != http.StatusOK {
			t.Fatalf("add member %s = %d (%s)", user, status, string(body))
		}
	}
	// Force an inconsistent purchase quantity to simulate an authorised reduction
	// that must be rejected by replacement, not silently applied.
	if err := fx.db.Model(&licensingModel.LicenceScope{}).Where("licence_number = ? AND scope_type = ?", licenceNumber, licensingModel.ScopeUserSeat).Update("quantity", 1).Error; err != nil {
		t.Fatalf("reduce scope: %v", err)
	}
	status, body := req(http.MethodPost, "/main/ems/licence-registry/licences/"+licenceNumber+"/replace", fx.tokenAdminA, map[string]any{
		"organisation_number": s0bOrgA,
		"valid_from":          now.Add(2 * time.Hour).Format(time.RFC3339),
		"valid_to":            now.Add(72 * time.Hour).Format(time.RFC3339),
		"max_installations":   2,
		"licensed_addons":     []string{"education"},
		"reason":              "6A.2 insufficient capacity",
	})
	if status != http.StatusBadRequest {
		t.Fatalf("replace insufficient = %d (%s), want 400", status, string(body))
	}
	var old licensingModel.Licence
	fx.db.Where("number = ?", licenceNumber).First(&old)
	if old.Status != licensingModel.StatusActive || old.ReplacedByNumber != nil {
		t.Fatalf("failed replacement mutated old licence: status=%s replaced_by=%v", old.Status, old.ReplacedByNumber)
	}
	var activeUnits int64
	fx.db.Model(&licensingModel.LicensedUnit{}).Where("licence_number = ? AND unit_type = ? AND status = ?", licenceNumber, licensingModel.ScopeUserSeat, licensingModel.UnitStatusActive).Count(&activeUnits)
	if activeUnits != 2 {
		t.Fatalf("active seats after rollback = %d, want 2", activeUnits)
	}
}
