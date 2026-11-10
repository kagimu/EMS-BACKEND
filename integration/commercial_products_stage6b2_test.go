package integration

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/routes"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	commercialModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/commercial/models"
	licensingModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/licensing/models"
	eosModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/eos/models"
)

func TestStage6B2CommercialProducts(t *testing.T) {
	fx := setupStage0BFixture(t)
	seedStage1Licence(t, fx, s0bOrgA, "stage6b2-licence-a")
	t.Setenv(commons.SoftwareOwnerUserNumbersEnv, s0bAdminA)
	prefix := fmt.Sprintf("S6B2-%d", time.Now().UnixNano()%10000000000)
	now := time.Now()
	scope := licensingModel.LicenceScope{
		Number: "stage6b2-seat-scope", LicenceNumber: "stage6b2-licence-a", ScopeType: licensingModel.ScopeUserSeat,
		Quantity: 5, Status: licensingModel.FlexStatusActive, ValidFrom: now.Add(-time.Hour), ValidTo: now.Add(48 * time.Hour),
	}
	if err := fx.db.Create(&scope).Error; err != nil {
		t.Fatalf("seed scope: %v", err)
	}
	unit := licensingModel.LicensedUnit{Number: "stage6b2-seat-unit", LicenceNumber: "stage6b2-licence-a", ScopeNumber: scope.Number, UnitType: licensingModel.ScopeUserSeat, UnitIdentifier: "stage6b2-user", OrganisationNumber: s0bOrgA, Status: licensingModel.UnitStatusActive, RegisteredAt: now}
	if err := fx.db.Create(&unit).Error; err != nil {
		t.Fatalf("seed unit: %v", err)
	}
	siteScope := licensingModel.LicenceScope{
		Number: "stage6b2-site-scope", LicenceNumber: "stage6b2-licence-a", ScopeType: licensingModel.ScopeSite,
		Quantity: 5, Status: licensingModel.FlexStatusActive, ValidFrom: now.Add(-time.Hour), ValidTo: now.Add(48 * time.Hour),
	}
	if err := fx.db.Create(&siteScope).Error; err != nil {
		t.Fatalf("seed site scope: %v", err)
	}
	siteUnit := licensingModel.LicensedUnit{
		Number: "stage6b2-site-unit", LicenceNumber: "stage6b2-licence-a", ScopeNumber: siteScope.Number,
		UnitType: licensingModel.ScopeSite, UnitIdentifier: "stage6b2-site", SiteNumber: "stage6b2-site",
		OrganisationNumber: s0bOrgA, Status: licensingModel.UnitStatusActive, RegisteredAt: now,
	}
	if err := fx.db.Create(&siteUnit).Error; err != nil {
		t.Fatalf("seed site unit: %v", err)
	}
	var membershipBefore, installationBefore int64
	fx.db.Model(&eosModel.EntityUser{}).Where("entity_number = ?", s0bOrgA).Count(&membershipBefore)
	fx.db.Model(&licensingModel.Installation{}).Where("organisation_number = ?", s0bOrgA).Count(&installationBefore)
	t.Cleanup(func() {
		fx.db.Unscoped().Where("product_number IN (SELECT number FROM ems_commercial_products WHERE product_code LIKE ?)", prefix+"%").Delete(&commercialModels.CommercialProductComponent{})
		fx.db.Unscoped().Where("product_code LIKE ?", prefix+"%").Delete(&commercialModels.CommercialProduct{})
		fx.db.Unscoped().Where("number = ?", unit.Number).Delete(&licensingModel.LicensedUnit{})
		fx.db.Unscoped().Where("number = ?", unit.Number).Delete(&licensingModel.LicensedUnit{})
		fx.db.Unscoped().Where("number = ?", siteUnit.Number).Delete(&licensingModel.LicensedUnit{})
		fx.db.Unscoped().Where("number = ?", scope.Number).Delete(&licensingModel.LicenceScope{})
		fx.db.Unscoped().Where("number = ?", siteScope.Number).Delete(&licensingModel.LicenceScope{})
	})

	server := httptest.NewServer(routes.SetupRoutes())
	t.Cleanup(server.Close)
	client := server.Client()
	req := func(method, path, token string, body any) (int, []byte) {
		return stage1Request(t, client, server.URL, method, path, token, body)
	}
	productPayload := func(code, name string) map[string]any {
		return map[string]any{"product_code": code, "name": name, "description": "Commercial product " + name, "product_type": "ADDON", "display_order": 1}
	}

	t.Run("owner_security_boundary", func(t *testing.T) {
		if status, _ := req(http.MethodGet, "/main/ems/owner/products", "", nil); status != http.StatusUnauthorized {
			t.Fatalf("unauth = %d", status)
		}
		if status, _ := req(http.MethodGet, "/main/ems/owner/products", fx.tokenBoth, nil); status != http.StatusForbidden {
			t.Fatalf("non-owner = %d", status)
		}
	})

	var productNumber, concurrentProductNumber string
	t.Run("product_crud_and_duplicate_validation", func(t *testing.T) {
		status, body := req(http.MethodPost, "/main/ems/owner/products", fx.tokenAdminA, productPayload(prefix+"-CORE", "EMS Core Product"))
		if status != http.StatusCreated {
			t.Fatalf("create product = %d (%s)", status, string(body))
		}
		var created struct {
			Data struct {
				Product struct {
					Number string `json:"number"`
					Status string `json:"status"`
				} `json:"product"`
			} `json:"data"`
		}
		_ = unmarshalStage6B2(body, &created)
		productNumber = created.Data.Product.Number
		if productNumber == "" || created.Data.Product.Status != "DRAFT" {
			t.Fatalf("product not created as DRAFT: %s", string(body))
		}
		if status, _ := req(http.MethodPost, "/main/ems/owner/products", fx.tokenAdminA, productPayload(prefix+"-CORE", "Duplicate")); status != http.StatusConflict {
			t.Fatalf("duplicate product = %d, want 409", status)
		}
		if status, _ := req(http.MethodPost, "/main/ems/owner/products", fx.tokenAdminA, productPayload("bad code!", "Bad")); status != http.StatusBadRequest {
			t.Fatalf("invalid code = %d, want 400", status)
		}
		if status, _ := req(http.MethodGet, "/main/ems/owner/products?search="+prefix, fx.tokenAdminA, nil); status != http.StatusOK {
			t.Fatalf("list = %d", status)
		}
		if status, body := req(http.MethodPut, "/main/ems/owner/products/"+productNumber, fx.tokenAdminA, map[string]any{"name": "EMS Core Product Updated", "description": "updated", "product_type": "CORE", "display_order": 2}); status != http.StatusOK {
			t.Fatalf("update product = %d (%s)", status, string(body))
		}
		if status, _ := req(http.MethodPatch, "/main/ems/owner/products/"+productNumber+"/status", fx.tokenAdminA, map[string]string{"status": "BOGUS"}); status != http.StatusBadRequest {
			t.Fatalf("invalid status = %d", status)
		}
		if status, body := req(http.MethodPatch, "/main/ems/owner/products/"+productNumber+"/status", fx.tokenAdminA, map[string]string{"status": "ACTIVE"}); status != http.StatusOK {
			t.Fatalf("activate product = %d (%s)", status, string(body))
		}
	})

	t.Run("component_validation_and_lifecycle", func(t *testing.T) {
		if status, body := req(http.MethodPost, "/main/ems/owner/products/"+productNumber+"/components", fx.tokenAdminA, map[string]any{"component_code": "core", "is_included": true, "is_optional": false, "display_order": 1}); status != http.StatusCreated {
			t.Fatalf("add core = %d (%s)", status, string(body))
		}
		if status, _ := req(http.MethodPost, "/main/ems/owner/products/"+productNumber+"/components", fx.tokenAdminA, map[string]any{"component_code": "core", "is_included": true, "is_optional": false}); status != http.StatusConflict {
			t.Fatalf("duplicate component = %d, want 409", status)
		}
		if status, _ := req(http.MethodPost, "/main/ems/owner/products/"+productNumber+"/components", fx.tokenAdminA, map[string]any{"component_code": "not-an-addon", "is_included": true, "is_optional": false}); status != http.StatusBadRequest {
			t.Fatalf("bad addon = %d, want 400", status)
		}
		if status, _ := req(http.MethodPost, "/main/ems/owner/products/"+productNumber+"/components", fx.tokenAdminA, map[string]any{"component_code": "education", "is_included": true, "is_optional": true}); status != http.StatusBadRequest {
			t.Fatalf("both flags = %d, want 400", status)
		}
		status, body := req(http.MethodPost, "/main/ems/owner/products/"+productNumber+"/components", fx.tokenAdminA, map[string]any{"component_code": "education", "is_included": false, "is_optional": true, "display_order": 2})
		if status != http.StatusCreated {
			t.Fatalf("add education = %d (%s)", status, string(body))
		}
		var component struct {
			Data struct {
				Number string `json:"number"`
			} `json:"data"`
		}
		_ = unmarshalStage6B2(body, &component)
		if status, body := req(http.MethodPut, "/main/ems/owner/products/"+productNumber+"/components/"+component.Data.Number, fx.tokenAdminA, map[string]any{"component_code": "education", "is_included": true, "is_optional": false, "display_order": 3}); status != http.StatusOK {
			t.Fatalf("update component = %d (%s)", status, string(body))
		}
		if status, body := req(http.MethodPatch, "/main/ems/owner/products/"+productNumber+"/components/"+component.Data.Number+"/status", fx.tokenAdminA, map[string]string{"status": "INACTIVE"}); status != http.StatusOK {
			t.Fatalf("deactivate component = %d (%s)", status, string(body))
		}
	})

	t.Run("component_belongs_to_another_product_is_rejected", func(t *testing.T) {
		status, body := req(http.MethodPost, "/main/ems/owner/products", fx.tokenAdminA, productPayload(prefix+"-OTHER", "Other Product"))
		if status != http.StatusCreated {
			t.Fatalf("create other product = %d (%s)", status, string(body))
		}
		var other struct {
			Data struct {
				Product struct {
					Number string `json:"number"`
				} `json:"product"`
			} `json:"data"`
		}
		_ = unmarshalStage6B2(body, &other)
		status, body = req(http.MethodPost, "/main/ems/owner/products/"+other.Data.Product.Number+"/components", fx.tokenAdminA, map[string]any{"component_code": "core", "is_included": true, "is_optional": false})
		if status != http.StatusCreated {
			t.Fatalf("add other component = %d (%s)", status, string(body))
		}
		var component struct {
			Data struct {
				Number string `json:"number"`
			} `json:"data"`
		}
		_ = unmarshalStage6B2(body, &component)
		if status, _ := req(http.MethodPut, "/main/ems/owner/products/"+productNumber+"/components/"+component.Data.Number, fx.tokenAdminA, map[string]any{"component_code": "core", "is_included": false, "is_optional": true}); status != http.StatusNotFound {
			t.Fatalf("cross-product component = %d, want 404", status)
		}
	})

	t.Run("concurrent_duplicate_component", func(t *testing.T) {
		status, body := req(http.MethodPost, "/main/ems/owner/products", fx.tokenAdminA, productPayload(prefix+"-CONCURRENT", "Concurrent Product"))
		if status != http.StatusCreated {
			t.Fatalf("create concurrent product = %d (%s)", status, string(body))
		}
		var created struct {
			Data struct {
				Product struct {
					Number string `json:"number"`
				} `json:"product"`
			} `json:"data"`
		}
		_ = unmarshalStage6B2(body, &created)
		concurrentProductNumber = created.Data.Product.Number
		payload := map[string]any{"component_code": "education", "is_included": true, "is_optional": false}
		results := make(chan int, 2)
		var wait sync.WaitGroup
		for i := 0; i < 2; i++ {
			wait.Add(1)
			go func() {
				defer wait.Done()
				status, _ := req(http.MethodPost, "/main/ems/owner/products/"+concurrentProductNumber+"/components", fx.tokenAdminA, payload)
				results <- status
			}()
		}
		wait.Wait()
		close(results)
		success, conflict := 0, 0
		for status := range results {
			if status == http.StatusCreated {
				success++
			}
			if status == http.StatusConflict {
				conflict++
			}
		}
		if success != 1 || conflict != 1 {
			t.Fatalf("concurrent component success=%d conflict=%d, want 1/1", success, conflict)
		}
		var active int64
		fx.db.Model(&commercialModels.CommercialProductComponent{}).Where("product_number = ? AND component_code = ?", concurrentProductNumber, "education").Count(&active)
		if active != 1 {
			t.Fatalf("component rows = %d, want 1", active)
		}
	})

	t.Run("audit_and_licensing_boundaries", func(t *testing.T) {
		var events int64
		fx.db.Model(&commercialModels.CommercialEvent{}).Where("entity_number = ? AND actor_user_number = ?", productNumber, s0bAdminA).Count(&events)
		if events < 3 {
			t.Fatalf("product audit events = %d, want at least 3", events)
		}
		var componentEvents int64
		fx.db.Model(&commercialModels.CommercialEvent{}).Where("entity_type = ? AND after_state LIKE ? AND actor_user_number = ?", "COMMERCIAL_PRODUCT_COMPONENT", "%"+productNumber+"%", s0bAdminA).Count(&componentEvents)
		if componentEvents < 4 {
			t.Fatalf("component audit events = %d, want at least 4", componentEvents)
		}
		var wronglyScoped int64
		fx.db.Model(&commercialModels.CommercialEvent{}).Where("entity_number = ? AND organisation_number <> ''", productNumber).Count(&wronglyScoped)
		if wronglyScoped != 0 {
			t.Fatalf("global product events wrongly attached to organisations: %d", wronglyScoped)
		}
		var componentScoped int64
		fx.db.Model(&commercialModels.CommercialEvent{}).Where("entity_type = ? AND after_state LIKE ? AND organisation_number <> ''", "COMMERCIAL_PRODUCT_COMPONENT", "%"+productNumber+"%").Count(&componentScoped)
		if componentScoped != 0 {
			t.Fatalf("global component events wrongly attached to organisations: %d", componentScoped)
		}
		var membershipAfter, installationAfter int64
		fx.db.Model(&eosModel.EntityUser{}).Where("entity_number = ?", s0bOrgA).Count(&membershipAfter)
		fx.db.Model(&licensingModel.Installation{}).Where("organisation_number = ?", s0bOrgA).Count(&installationAfter)
		var activeUnits int64
		fx.db.Model(&licensingModel.LicensedUnit{}).Where("licence_number = ? AND unit_type = ? AND status = ?", "stage6b2-licence-a", licensingModel.ScopeUserSeat, licensingModel.UnitStatusActive).Count(&activeUnits)
		var activeSiteUnits int64
		fx.db.Model(&licensingModel.LicensedUnit{}).Where("licence_number = ? AND unit_type = ? AND status = ?", "stage6b2-licence-a", licensingModel.ScopeSite, licensingModel.UnitStatusActive).Count(&activeSiteUnits)
		if membershipAfter != membershipBefore || installationAfter != installationBefore || activeUnits != 1 || activeSiteUnits != 1 {
			t.Fatalf("commercial products changed licensing/membership: memberships %d->%d installs %d->%d units=%d", membershipBefore, membershipAfter, installationBefore, installationAfter, activeUnits)
		}
		var licence licensingModel.Licence
		fx.db.Where("number = ?", "stage6b2-licence-a").First(&licence)
		if licence.Status != licensingModel.StatusActive {
			t.Fatalf("licence status changed to %s", licence.Status)
		}
	})

	t.Run("inactive_product_remains_identifiable", func(t *testing.T) {
		if status, body := req(http.MethodPatch, "/main/ems/owner/products/"+productNumber+"/status", fx.tokenAdminA, map[string]string{"status": "INACTIVE"}); status != http.StatusOK {
			t.Fatalf("deactivate product = %d (%s)", status, string(body))
		}
		status, body := req(http.MethodGet, "/main/ems/owner/products/"+productNumber, fx.tokenAdminA, nil)
		if status != http.StatusOK {
			t.Fatalf("inactive detail = %d", status)
		}
		if !containsStage6(string(body), `"status":"INACTIVE"`) {
			t.Fatalf("inactive product not identifiable: %s", string(body))
		}
		if status, _ := req(http.MethodPatch, "/main/ems/owner/products/"+productNumber+"/status", fx.tokenAdminA, map[string]string{"status": "ACTIVE"}); status != http.StatusOK {
			t.Fatalf("reactivate product = %d", status)
		}
	})
}

func unmarshalStage6B2(body []byte, target any) error {
	return jsonUnmarshalStage6B2(body, target)
}
func jsonUnmarshalStage6B2(body []byte, target any) error { return json.Unmarshal(body, target) }
