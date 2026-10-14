package integration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/routes"
)

// TestSetupStatusEndpoint verifies the public, read-only installation-status
// contract used by the frontend root:
//   - it is reachable without authentication,
//   - it returns exactly one field ("state") whose value is a known routing
//     state that is consistent with the database,
//   - it leaks no identities, counts, tokens or licence data.
//
// It is read-only and skips unless the connected database is the dedicated
// disposable onboarding database (TEST11 convention).
func TestSetupStatusEndpoint(t *testing.T) {
	db := requireOnboardingDB(t)

	server := httptest.NewServer(routes.SetupRoutes())
	t.Cleanup(server.Close)

	resp, err := server.Client().Get(server.URL + "/setup/status")
	if err != nil {
		t.Fatalf("GET setup/status: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("GET setup/status = %d, want 200", resp.StatusCode)
	}

	var body map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode body: %v", err)
	}

	// Minimum disclosure: exactly one field.
	if len(body) != 1 {
		t.Fatalf("setup/status disclosed extra fields: %#v", body)
	}
	state, ok := body["state"].(string)
	if !ok {
		t.Fatalf("setup/status state missing or not a string: %#v", body)
	}
	if state != "initialized" && state != "uninitialized" {
		t.Fatalf("setup/status state = %q, want initialized|uninitialized", state)
	}

	// Must be consistent with the authoritative source (active owner exists).
	var activeOwners int64
	if err := db.Raw("SELECT COUNT(*) FROM ems_stakeholders WHERE stakeholder_type = 'owner' AND status = 'ACTIVE'").Scan(&activeOwners).Error; err != nil {
		t.Fatalf("count active owners: %v", err)
	}
	want := "uninitialized"
	if activeOwners > 0 {
		want = "initialized"
	}
	if state != want {
		t.Fatalf("setup/status = %q, want %q for %d active owner(s)", state, want, activeOwners)
	}
}
