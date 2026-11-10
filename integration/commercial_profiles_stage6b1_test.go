package integration

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/routes"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	commercialModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/commercial/models"
	licensingModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/licensing/models"
)

func TestStage6B1CustomerCommercialProfiles(t *testing.T) {
	fx := setupStage0BFixture(t)
	seedStage1Licence(t, fx, s0bOrgA, "stage6b1-licence-a")
	t.Setenv(commons.SoftwareOwnerUserNumbersEnv, s0bAdminA)
	now := time.Now()
	scope := licensingModel.LicenceScope{
		Number: "stage6b1-seat-scope", LicenceNumber: "stage6b1-licence-a", ScopeType: licensingModel.ScopeUserSeat,
		Quantity: 5, Status: licensingModel.FlexStatusActive, ValidFrom: now.Add(-time.Hour), ValidTo: now.Add(48 * time.Hour),
	}
	if err := fx.db.Create(&scope).Error; err != nil {
		t.Fatalf("seed scope: %v", err)
	}
	unit := licensingModel.LicensedUnit{
		Number: "stage6b1-seat-unit", LicenceNumber: "stage6b1-licence-a", ScopeNumber: scope.Number,
		UnitType: licensingModel.ScopeUserSeat, UnitIdentifier: "stage6b1-user", OrganisationNumber: s0bOrgA,
		Status: licensingModel.UnitStatusActive, RegisteredAt: now,
	}
	if err := fx.db.Create(&unit).Error; err != nil {
		t.Fatalf("seed unit: %v", err)
	}
	t.Cleanup(func() {
		fx.db.Unscoped().Where("organisation_number = ?", s0bOrgA).Delete(&commercialModels.CustomerProfile{})
		fx.db.Unscoped().Where("organisation_number = ?", s0bOrgA).Delete(&commercialModels.CommercialEvent{})
		fx.db.Unscoped().Where("number = ?", unit.Number).Delete(&licensingModel.LicensedUnit{})
		fx.db.Unscoped().Where("number = ?", scope.Number).Delete(&licensingModel.LicenceScope{})
	})

	server := httptest.NewServer(routes.SetupRoutes())
	t.Cleanup(server.Close)
	client := server.Client()
	req := func(method, path, token string, body any) (int, []byte) {
		return stage1Request(t, client, server.URL, method, path, token, body)
	}
	profilePath := "/main/ems/owner/customers/" + s0bOrgA + "/profile"
	payload := func(status, category, email, manager string) map[string]any {
		return map[string]any{
			"commercial_status": status, "customer_category": category, "billing_email": email,
			"registered_business_name": "Stage 6B.1 Ltd", "billing_contact_name": "Billing Contact",
			"billing_phone": "0700000000", "physical_address": "Kampala", "country": "Uganda",
			"relationship_manager_number": manager, "internal_notes": "Owner-only note",
		}
	}

	t.Run("unauthenticated_and_non_owner_rejected", func(t *testing.T) {
		if status, _ := req(http.MethodGet, profilePath, "", nil); status != http.StatusUnauthorized {
			t.Fatalf("unauth = %d, want 401", status)
		}
		if status, _ := req(http.MethodPost, profilePath, fx.tokenBoth, payload("ACTIVE", "BUSINESS", "billing@example.com", "")); status != http.StatusForbidden {
			t.Fatalf("non-owner = %d, want 403", status)
		}
	})

	t.Run("create_duplicate_and_invalid_inputs", func(t *testing.T) {
		status, body := req(http.MethodPost, profilePath, fx.tokenAdminA, payload("ACTIVE", "BUSINESS", "billing@example.com", ""))
		if status != http.StatusCreated {
			t.Fatalf("create = %d (%s)", status, string(body))
		}
		if status, _ := req(http.MethodPost, profilePath, fx.tokenAdminA, payload("ACTIVE", "BUSINESS", "billing@example.com", "")); status != http.StatusConflict {
			t.Fatalf("duplicate = %d, want 409", status)
		}
		if status, _ := req(http.MethodPost, "/main/ems/owner/customers/does-not-exist/profile", fx.tokenAdminA, payload("ACTIVE", "BUSINESS", "billing@example.com", "")); status != http.StatusNotFound {
			t.Fatalf("invalid organisation = %d, want 404", status)
		}
		if status, _ := req(http.MethodPut, profilePath, fx.tokenAdminA, payload("NOT_A_STATUS", "BUSINESS", "billing@example.com", "")); status != http.StatusBadRequest {
			t.Fatalf("invalid status = %d, want 400", status)
		}
		if status, _ := req(http.MethodPut, profilePath, fx.tokenAdminA, payload("ACTIVE", "NOT_A_CATEGORY", "billing@example.com", "")); status != http.StatusBadRequest {
			t.Fatalf("invalid category = %d, want 400", status)
		}
		if status, _ := req(http.MethodPut, profilePath, fx.tokenAdminA, payload("ACTIVE", "BUSINESS", "billing@example.com", "missing-manager")); status != http.StatusBadRequest {
			t.Fatalf("invalid manager = %d, want 400", status)
		}
	})

	t.Run("update_profile_audits_actor_and_preserves_licensing", func(t *testing.T) {
		var licenceBefore licensingModel.Licence
		fx.db.Where("number = ?", "stage6b1-licence-a").First(&licenceBefore)
		var unitsBefore int64
		fx.db.Model(&licensingModel.LicensedUnit{}).Where("licence_number = ? AND status = ?", "stage6b1-licence-a", licensingModel.UnitStatusActive).Count(&unitsBefore)

		status, body := req(http.MethodPut, profilePath, fx.tokenAdminA, payload("CHURNED", "GOVERNMENT", "billing@example.com", ""))
		if status != http.StatusOK {
			t.Fatalf("update = %d (%s)", status, string(body))
		}

		var profile commercialModels.CustomerProfile
		fx.db.Where("organisation_number = ?", s0bOrgA).First(&profile)
		if profile.CommercialStatus != commercialModels.CommercialStatusChurned || profile.UpdatedByNumber != s0bAdminA {
			t.Fatalf("profile not updated correctly: %+v", profile)
		}
		var event commercialModels.CommercialEvent
		fx.db.Where("entity_number = ? AND action = ?", profile.Number, "CUSTOMER_STATUS_CHANGED").First(&event)
		if event.ActorUserNumber != s0bAdminA || event.OrganisationNumber != s0bOrgA {
			t.Fatalf("audit actor/organisation incorrect: %+v", event)
		}

		var licenceAfter licensingModel.Licence
		fx.db.Where("number = ?", "stage6b1-licence-a").First(&licenceAfter)
		var unitsAfter int64
		fx.db.Model(&licensingModel.LicensedUnit{}).Where("licence_number = ? AND status = ?", "stage6b1-licence-a", licensingModel.UnitStatusActive).Count(&unitsAfter)
		if licenceAfter.Status != licenceBefore.Status || unitsAfter != unitsBefore || unitsBefore != 1 {
			t.Fatalf("licensing changed by commercial status: licence %s->%s units %d->%d", licenceBefore.Status, licenceAfter.Status, unitsBefore, unitsAfter)
		}
	})

	t.Run("directory_and_detail_include_profile", func(t *testing.T) {
		status, body := req(http.MethodGet, "/main/ems/owner/organisations?per_page=100", fx.tokenAdminA, nil)
		if status != http.StatusOK {
			t.Fatalf("directory = %d", status)
		}
		if !containsStage6(value(body), `"has_commercial_profile":true`) || !containsStage6(value(body), `"commercial_status":"CHURNED"`) {
			t.Fatalf("directory profile fields missing: %s", value(body))
		}
		status, body = req(http.MethodGet, "/main/ems/owner/organisations/"+s0bOrgA, fx.tokenAdminA, nil)
		if status != http.StatusOK || !containsStage6(value(body), `"commercial_profile"`) {
			t.Fatalf("detail profile missing: %d %s", status, value(body))
		}
	})
}

func value(body []byte) string { return fmt.Sprintf("%s", body) }
func containsStage6(value, target string) bool {
	return len(value) >= len(target) && (value == target || indexOfStage6(value, target) >= 0)
}
func indexOfStage6(value, target string) int {
	for index := 0; index+len(target) <= len(value); index++ {
		if value[index:index+len(target)] == target {
			return index
		}
	}
	return -1
}
