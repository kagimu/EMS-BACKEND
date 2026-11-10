package integration

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/routes"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	licensingModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/licensing/models"
	organisationModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/organisation/models"
	eosModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/eos/models"
)

func TestStage6A1Integrity(t *testing.T) {
	fx := setupStage0BFixture(t)
	fx.db.Model(&licensingModel.Licence{}).Where("organisation_number = ? AND active_org_number IS NOT NULL", s0bOrgA).Updates(map[string]any{"status": licensingModel.StatusRevoked, "active_org_number": nil})
	licenceNumber := fmt.Sprintf("s6a1-%d", time.Now().UnixNano()%10000000000)
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
	configure := func() (int, []byte) {
		return req(http.MethodPut, "/main/ems/owner/licences/"+licenceNumber+"/flexible", fx.tokenAdminA, map[string]any{
			"scopes": []map[string]any{
				{"scope_type": "USER_SEAT", "quantity": 1, "status": "ACTIVE", "valid_from": validFrom, "valid_to": validTo},
				{"scope_type": "SITE", "quantity": 1, "status": "ACTIVE", "valid_from": validFrom, "valid_to": validTo},
			},
			"entitlements": []map[string]any{{"code": "education", "quantity": 1, "status": "ACTIVE", "valid_from": validFrom, "valid_to": validTo}},
			"reason":       "Stage 6A.1 test setup",
		})
	}
	if status, body := configure(); status != http.StatusOK {
		t.Fatalf("configure = %d (%s)", status, string(body))
	}

	t.Cleanup(func() {
		fx.db.Unscoped().Where("user_number = ? AND entity_number = ?", s0bUserOnlyB, s0bOrgA).Delete(&eosModel.EntityUser{})
		fx.db.Unscoped().Where("licence_number = ?", licenceNumber).Delete(&licensingModel.LicensedUnit{})
		fx.db.Unscoped().Where("licence_number = ?", licenceNumber).Delete(&licensingModel.LicenceScope{})
		fx.db.Unscoped().Where("licence_number = ?", licenceNumber).Delete(&licensingModel.LicenceEntitlement{})
	})

	if status, body := req(http.MethodPost, "/main/eos/add-user-to-entity", fx.tokenAdminA, map[string]string{"entity_number": s0bOrgA, "user_number": s0bUserOnlyB}); status != http.StatusOK {
		t.Fatalf("add member = %d (%s)", status, string(body))
	}

	var siteNumber string
	status, body := req(http.MethodPost, "/main/ems/organisation-units", fx.tokenAdminA, map[string]any{"name": "Stage6A1 Site", "unit_type": "site", "code": "s6a1-site"})
	if status != http.StatusCreated {
		t.Fatalf("create site = %d (%s)", status, string(body))
	}
	var site struct {
		Data struct {
			Number string `json:"number"`
		} `json:"data"`
	}
	_ = json.Unmarshal(body, &site)
	siteNumber = site.Data.Number
	t.Cleanup(func() {
		fx.db.Unscoped().Where("number = ?", siteNumber).Delete(&organisationModel.OrganisationUnit{})
		fx.db.Unscoped().Where("licence_number = ?", licenceNumber).Delete(&licensingModel.LicenceEvent{})
	})

	t.Run("seat_and_site_events_are_recorded", func(t *testing.T) {
		for _, reason := range []string{"USER_SEAT_REGISTERED", "SITE_REGISTERED"} {
			var count int64
			fx.db.Model(&licensingModel.LicenceEvent{}).Where("licence_number = ? AND reason = ?", licenceNumber, reason).Count(&count)
			if count != 1 {
				t.Fatalf("%s event count = %d, want 1", reason, count)
			}
		}
	})

	t.Run("site_release_and_reactivation_are_audited", func(t *testing.T) {
		if status, body := req(http.MethodPatch, "/main/ems/organisation-units/"+siteNumber+"/status", fx.tokenAdminA, map[string]string{"status": "INACTIVE"}); status != http.StatusOK {
			t.Fatalf("deactivate site = %d (%s)", status, string(body))
		}
		var released int64
		fx.db.Model(&licensingModel.LicenceEvent{}).Where("licence_number = ? AND reason = ?", licenceNumber, "SITE_RELEASED").Count(&released)
		if released != 1 {
			t.Fatalf("site release events = %d, want 1", released)
		}
		if status, body := req(http.MethodPatch, "/main/ems/organisation-units/"+siteNumber+"/status", fx.tokenAdminA, map[string]string{"status": "ACTIVE"}); status != http.StatusOK {
			t.Fatalf("reactivate site = %d (%s)", status, string(body))
		}
		var registered int64
		fx.db.Model(&licensingModel.LicenceEvent{}).Where("licence_number = ? AND reason = ?", licenceNumber, "SITE_REGISTERED").Count(&registered)
		if registered < 2 {
			t.Fatalf("site registration events = %d, want at least 2", registered)
		}
	})

	t.Run("membership_removal_is_idempotent_and_releases_once", func(t *testing.T) {
		if status, body := req(http.MethodPost, "/main/eos/remove-user-from-entity", fx.tokenAdminA, map[string]string{"entity_number": s0bOrgA, "user_number": s0bUserOnlyB}); status != http.StatusOK {
			t.Fatalf("remove member = %d (%s)", status, string(body))
		}
		if status, body := req(http.MethodPost, "/main/eos/remove-user-from-entity", fx.tokenAdminA, map[string]string{"entity_number": s0bOrgA, "user_number": s0bUserOnlyB}); status != http.StatusOK {
			t.Fatalf("repeat remove = %d (%s)", status, string(body))
		}
		var active int64
		fx.db.Model(&licensingModel.LicensedUnit{}).Where("licence_number = ? AND unit_type = ? AND status = ?", licenceNumber, licensingModel.ScopeUserSeat, licensingModel.UnitStatusActive).Count(&active)
		if active != 0 {
			t.Fatalf("active user seats = %d, want 0", active)
		}
		var releaseEvents int64
		fx.db.Model(&licensingModel.LicenceEvent{}).Where("licence_number = ? AND reason = ?", licenceNumber, "USER_SEAT_RELEASED").Count(&releaseEvents)
		if releaseEvents != 1 {
			t.Fatalf("user seat release events = %d, want 1", releaseEvents)
		}
	})

	t.Run("replacement_preserves_scopes_and_releases_old_units", func(t *testing.T) {
		// Re-add a seat and keep the site active so replacement has live usage.
		if status, body := req(http.MethodPost, "/main/eos/add-user-to-entity", fx.tokenAdminA, map[string]string{"entity_number": s0bOrgA, "user_number": s0bUserOnlyB}); status != http.StatusOK {
			t.Fatalf("re-add member = %d (%s)", status, string(body))
		}
		status, body := req(http.MethodPost, "/main/ems/licence-registry/licences/"+licenceNumber+"/replace", fx.tokenAdminA, map[string]any{
			"organisation_number": s0bOrgA,
			"valid_from":          now.Add(2 * time.Hour).Format(time.RFC3339),
			"valid_to":            now.Add(72 * time.Hour).Format(time.RFC3339),
			"max_installations":   3,
			"licensed_addons":     []string{"education"},
			"reason":              "Stage 6A.1 replacement",
		})
		if status != http.StatusCreated {
			t.Fatalf("replace = %d (%s)", status, string(body))
		}
		var issued struct {
			Data struct {
				Number     string `json:"number"`
				LicenceKey string `json:"licence_key"`
			} `json:"data"`
		}
		_ = json.Unmarshal(body, &issued)
		if issued.Data.Number == "" || issued.Data.LicenceKey == "" {
			t.Fatalf("replacement response missing licence: %s", string(body))
		}
		var old licensingModel.Licence
		fx.db.Where("number = ?", licenceNumber).First(&old)
		if old.ReplacedByNumber == nil || *old.ReplacedByNumber != issued.Data.Number {
			t.Fatalf("old licence replacement link missing")
		}
		var oldActiveUnits int64
		fx.db.Model(&licensingModel.LicensedUnit{}).Where("licence_number = ? AND status = ?", licenceNumber, licensingModel.UnitStatusActive).Count(&oldActiveUnits)
		if oldActiveUnits != 0 {
			t.Fatalf("old active units = %d, want 0", oldActiveUnits)
		}
		var newScopes []licensingModel.LicenceScope
		fx.db.Where("licence_number = ?", issued.Data.Number).Find(&newScopes)
		if len(newScopes) < 2 {
			t.Fatalf("replacement scopes = %d, want at least 2", len(newScopes))
		}
		var seatScope licensingModel.LicenceScope
		fx.db.Where("licence_number = ? AND scope_type = ?", issued.Data.Number, licensingModel.ScopeUserSeat).First(&seatScope)
		if seatScope.Quantity != 1 {
			t.Fatalf("replacement seat quantity = %d, want 1", seatScope.Quantity)
		}
		var newActiveUnits int64
		fx.db.Model(&licensingModel.LicensedUnit{}).Where("licence_number = ? AND status = ?", issued.Data.Number, licensingModel.UnitStatusActive).Count(&newActiveUnits)
		if newActiveUnits != 2 {
			t.Fatalf("replacement active units = %d, want 2 migrated seat/site units", newActiveUnits)
		}
	})
}
