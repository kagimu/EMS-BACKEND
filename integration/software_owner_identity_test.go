package integration

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/routes"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
)

func TestSoftwareOwnerIdentityFlag(t *testing.T) {
	fx := setupStage0BFixture(t)
	t.Setenv(commons.SoftwareOwnerUserNumbersEnv, s0bAdminA)
	server := httptest.NewServer(routes.SetupRoutes())
	t.Cleanup(server.Close)
	client := server.Client()

	status, body := stage1Request(t, client, server.URL, http.MethodGet, "/main/profile", fx.tokenAdminA, nil)
	if status != http.StatusOK || !strings.Contains(string(body), `"is_software_owner":true`) {
		t.Fatalf("owner profile flag = %d (%s)", status, string(body))
	}
	status, body = stage1Request(t, client, server.URL, http.MethodGet, "/main/profile", fx.tokenBoth, nil)
	if status != http.StatusOK || !strings.Contains(string(body), `"is_software_owner":false`) {
		t.Fatalf("customer profile flag = %d (%s)", status, string(body))
	}
}
