package integration

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/routes"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	licensingModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/licensing/models"
)

func TestStage2OwnerConsole(t *testing.T) {
	fx := setupStage0BFixture(t)
	seedStage1Licence(t, fx, s0bOrgA, "stage2-licence-a")
	t.Setenv(commons.SoftwareOwnerUserNumbersEnv, s0bAdminA)

	installation := licensingModel.Installation{
		Number:               "stage2-installation-a",
		InstallationIdentity: "stage2-identity-a",
		OrganisationNumber:   s0bOrgA,
		LicenceNumber:        "stage2-licence-a",
		ActivatedAt:          time.Now().Add(-time.Hour),
		Status:               licensingModel.InstallationStatusActive,
	}
	if err := fx.db.Create(&installation).Error; err != nil {
		t.Fatalf("seeding installation: %v", err)
	}
	t.Cleanup(func() {
		fx.db.Unscoped().Where("number = ?", installation.Number).Delete(&licensingModel.Installation{})
		fx.db.Unscoped().Where("installation_number = ?", installation.Number).Delete(&licensingModel.LicenceEvent{})
	})

	server := httptest.NewServer(routes.SetupRoutes())
	t.Cleanup(server.Close)
	client := server.Client()

	req := func(method, path, token string, body any) (int, []byte) {
		return stage1Request(t, client, server.URL, method, path, token, body)
	}

	t.Run("unauthenticated_owner_api_is_rejected", func(t *testing.T) {
		if status, _ := req(http.MethodGet, "/main/ems/owner/overview", "", nil); status != http.StatusUnauthorized {
			t.Fatalf("status = %d, want %d", status, http.StatusUnauthorized)
		}
	})

	t.Run("client_admin_is_rejected_from_owner_api", func(t *testing.T) {
		if status, _ := req(http.MethodGet, "/main/ems/owner/overview", fx.tokenBoth, nil); status == http.StatusOK {
			t.Fatal("ordinary user must not reach owner overview")
		}
		if status, _ := req(http.MethodPost, "/main/ems/licence-registry/installations/"+installation.Number+"/release", fx.tokenBoth, map[string]any{"reason": "test"}); status != http.StatusForbidden {
			t.Fatalf("client release status = %d, want %d", status, http.StatusForbidden)
		}
	})

	t.Run("owner_can_read_owner_console", func(t *testing.T) {
		for _, path := range []string{
			"/main/ems/owner/access",
			"/main/ems/owner/overview",
			"/main/ems/owner/organisations?per_page=100",
			"/main/ems/owner/organisations/" + s0bOrgA,
			"/main/ems/owner/licences?per_page=100",
			"/main/ems/owner/licences/stage2-licence-a",
			"/main/ems/owner/installations?per_page=100",
			"/main/ems/owner/installations/" + installation.Number,
			"/main/ems/owner/addons",
		} {
			status, body := req(http.MethodGet, path, fx.tokenAdminA, nil)
			if status != http.StatusOK {
				t.Fatalf("%s = %d (%s)", path, status, string(body))
			}
			if strings.Contains(string(body), `"licence_key"`) {
				t.Fatalf("%s exposed a licence key", path)
			}
		}
	})

	t.Run("owner_can_release_and_action_is_audited", func(t *testing.T) {
		status, body := req(http.MethodPost, "/main/ems/licence-registry/installations/"+installation.Number+"/release", fx.tokenAdminA, map[string]any{"reason": "Stage 2 owner release"})
		if status != http.StatusOK {
			t.Fatalf("release status = %d (%s)", status, string(body))
		}
		var updated licensingModel.Installation
		if err := fx.db.Where("number = ?", installation.Number).First(&updated).Error; err != nil {
			t.Fatalf("loading released installation: %v", err)
		}
		if updated.Status != licensingModel.InstallationStatusReleased {
			t.Fatalf("installation status = %s, want %s", updated.Status, licensingModel.InstallationStatusReleased)
		}
		if updated.OrganisationNumber != s0bOrgA || updated.LicenceNumber != "stage2-licence-a" {
			t.Fatalf("installation association changed: %+v", updated)
		}
		var eventCount int64
		if err := fx.db.Model(&licensingModel.LicenceEvent{}).
			Where("installation_number = ? AND event_type = ?", installation.Number, licensingModel.InstallationEventReleased).
			Count(&eventCount).Error; err != nil {
			t.Fatalf("counting release event: %v", err)
		}
		if eventCount != 1 {
			t.Fatalf("release event count = %d, want 1", eventCount)
		}
	})
}
