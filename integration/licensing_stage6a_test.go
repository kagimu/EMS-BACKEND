package integration

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/routes"
	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	licensingModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/licensing/models"
	organisationModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/organisation/models"
	eosModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/eos/models"
)

func stage6Request(t *testing.T, client *http.Client, url, method, path, token string, body any) (int, []byte) {
	t.Helper()
	return stage1Request(t, client, url, method, path, token, body)
}

func TestStage6ASeatAndSiteEnforcement(t *testing.T) {
	fx := setupStage0BFixture(t)
	licenceNumber := fmt.Sprintf("s6a-lic-%d", time.Now().UnixNano()%10000000000)
	seedStage1Licence(t, fx, s0bOrgA, licenceNumber)
	t.Setenv(commons.SoftwareOwnerUserNumbersEnv, s0bAdminA)

	now := time.Now()
	validFrom := now.Add(-time.Hour).Format(time.RFC3339)
	validTo := now.Add(24 * time.Hour).Format(time.RFC3339)
	secondUser := authModel.User{Number: fmt.Sprintf("s6a-%d", time.Now().UnixNano()%10000000000), FirstName: "Seat", SecondName: "Two", LastName: "User", Gender: "M", Type: "Staff", Status: "Active"}
	if err := fx.db.Create(&secondUser).Error; err != nil {
		t.Fatalf("seeding second user: %v", err)
	}
	t.Cleanup(func() {
		fx.db.Unscoped().Where("user_number = ? AND entity_number = ?", secondUser.Number, s0bOrgA).Delete(&eosModel.EntityUser{})
		fx.db.Unscoped().Where("user_number = ?", secondUser.Number).Delete(&authModel.User{})
		fx.db.Unscoped().Where("unit_identifier = ? AND organisation_number = ?", secondUser.Number, s0bOrgA).Delete(&licensingModel.LicensedUnit{})
		fx.db.Unscoped().Where("user_number = ? AND entity_number = ?", s0bUserOnlyB, s0bOrgA).Delete(&eosModel.EntityUser{})
		fx.db.Unscoped().Where("unit_identifier = ? AND organisation_number = ?", s0bUserOnlyB, s0bOrgA).Delete(&licensingModel.LicensedUnit{})
	})

	server := httptest.NewServer(routes.SetupRoutes())
	t.Cleanup(server.Close)
	client := server.Client()

	req := func(method, path, token string, body any) (int, []byte) {
		return stage6Request(t, client, server.URL, method, path, token, body)
	}

	configure := func(scopes []map[string]any) (int, []byte) {
		return req(http.MethodPut, "/main/ems/owner/licences/"+licenceNumber+"/flexible", fx.tokenAdminA, map[string]any{
			"scopes": scopes, "entitlements": []any{}, "reason": "Stage 6A test",
		})
	}

	if status, body := configure([]map[string]any{{
		"scope_type": "USER_SEAT", "quantity": 1, "status": "ACTIVE", "valid_from": validFrom, "valid_to": validTo,
	}}); status != http.StatusOK {
		t.Fatalf("configure seat scope = %d (%s)", status, string(body))
	}

	t.Run("seat_membership_capacity_and_release", func(t *testing.T) {
		if status, body := req(http.MethodPost, "/main/eos/add-user-to-entity", fx.tokenAdminA, map[string]string{"entity_number": s0bOrgA, "user_number": s0bUserOnlyB}); status != http.StatusOK {
			t.Fatalf("add first member = %d (%s)", status, string(body))
		}
		if status, body := req(http.MethodPost, "/main/eos/add-user-to-entity", fx.tokenAdminA, map[string]string{"entity_number": s0bOrgA, "user_number": s0bUserOnlyB}); status != http.StatusOK {
			t.Fatalf("idempotent add = %d (%s)", status, string(body))
		}
		if status, _ := req(http.MethodPost, "/main/eos/add-user-to-entity", fx.tokenAdminA, map[string]string{"entity_number": s0bOrgA, "user_number": secondUser.Number}); status != http.StatusConflict {
			t.Fatalf("over-capacity membership = %d, want %d", status, http.StatusConflict)
		}
		if status, body := req(http.MethodPost, "/main/eos/remove-user-from-entity", fx.tokenAdminA, map[string]string{"entity_number": s0bOrgA, "user_number": s0bUserOnlyB}); status != http.StatusOK {
			t.Fatalf("remove member = %d (%s)", status, string(body))
		}
		if status, body := req(http.MethodPost, "/main/eos/add-user-to-entity", fx.tokenAdminA, map[string]string{"entity_number": s0bOrgA, "user_number": secondUser.Number}); status != http.StatusOK {
			t.Fatalf("add after release = %d (%s)", status, string(body))
		}
	})

	t.Run("seat_quantity_reduction_below_usage_is_rejected", func(t *testing.T) {
		if status, _ := configure([]map[string]any{{
			"scope_type": "USER_SEAT", "quantity": 0, "status": "ACTIVE", "valid_from": validFrom, "valid_to": validTo,
		}}); status != http.StatusBadRequest {
			t.Fatalf("unsafe reduction = %d, want %d", status, http.StatusBadRequest)
		}
	})

	if status, body := configure([]map[string]any{
		{"scope_type": "USER_SEAT", "quantity": 1, "status": "ACTIVE", "valid_from": validFrom, "valid_to": validTo},
		{"scope_type": "SITE", "quantity": 1, "status": "ACTIVE", "valid_from": validFrom, "valid_to": validTo},
	}); status != http.StatusOK {
		t.Fatalf("configure site scope = %d (%s)", status, string(body))
	}

	var firstSite string
	t.Cleanup(func() {
		fx.db.Unscoped().Where("number = ?", firstSite).Delete(&organisationModel.OrganisationUnit{})
	})

	t.Run("site_capacity_and_reactivation", func(t *testing.T) {
		status, body := req(http.MethodPost, "/main/ems/organisation-units", fx.tokenAdminA, map[string]any{
			"name": "Stage6 Site One", "unit_type": "site", "code": "stage6-site-one",
		})
		if status != http.StatusCreated {
			t.Fatalf("create site one = %d (%s)", status, string(body))
		}
		var created struct {
			Data struct {
				Number string `json:"number"`
			} `json:"data"`
		}
		_ = json.Unmarshal(body, &created)
		firstSite = created.Data.Number
		if firstSite == "" {
			t.Fatalf("site number missing: %s", string(body))
		}

		if status, _ := req(http.MethodPost, "/main/ems/organisation-units", fx.tokenAdminA, map[string]any{
			"name": "Stage6 Site Two", "unit_type": "site", "code": "stage6-site-two",
		}); status != http.StatusConflict {
			t.Fatalf("site over-capacity = %d, want %d", status, http.StatusConflict)
		}

		if status, body := req(http.MethodPatch, "/main/ems/organisation-units/"+firstSite+"/status", fx.tokenAdminA, map[string]string{"status": "INACTIVE"}); status != http.StatusOK {
			t.Fatalf("deactivate site = %d (%s)", status, string(body))
		}

		status, body = req(http.MethodPost, "/main/ems/organisation-units", fx.tokenAdminA, map[string]any{
			"name": "Stage6 Site Two", "unit_type": "site", "code": "stage6-site-two",
		})
		if status != http.StatusCreated {
			t.Fatalf("create site two after release = %d (%s)", status, string(body))
		}
		var second struct {
			Data struct {
				Number string `json:"number"`
			} `json:"data"`
		}
		_ = json.Unmarshal(body, &second)
		t.Cleanup(func() {
			fx.db.Unscoped().Where("number = ?", second.Data.Number).Delete(&organisationModel.OrganisationUnit{})
		})

		if status, _ := req(http.MethodPatch, "/main/ems/organisation-units/"+firstSite+"/status", fx.tokenAdminA, map[string]string{"status": "ACTIVE"}); status != http.StatusConflict {
			t.Fatalf("site reactivation over capacity = %d, want %d", status, http.StatusConflict)
		}
	})

	t.Run("usage_is_exposed_to_owner_and_client", func(t *testing.T) {
		if status, body := req(http.MethodGet, "/main/ems/owner/licences/"+licenceNumber+"/flexible", fx.tokenAdminA, nil); status != http.StatusOK || !strings.Contains(string(body), `"scope_type":"USER_SEAT"`) {
			t.Fatalf("owner usage = %d (%s)", status, string(body))
		}
		if status, body := req(http.MethodGet, "/main/ems/licence", fx.tokenBoth, nil); status != http.StatusOK || !strings.Contains(string(body), `"scopes"`) {
			t.Fatalf("client usage = %d (%s)", status, string(body))
		}
	})
}
