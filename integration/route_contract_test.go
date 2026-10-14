package integration

import (
	"net/http"
	"testing"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/routes"
)

// TestCriticalRoutesAreRegistered is a database-independent contract test. It
// proves that the routes the eight Education regression flows depend on are
// still registered after the Release 1 security changes.
func TestCriticalRoutesAreRegistered(t *testing.T) {
	engine := routes.SetupRoutes()

	registered := make(map[string]bool)
	for _, route := range engine.Routes() {
		registered[route.Method+" "+route.Path] = true
	}

	required := []string{
		http.MethodPost + " /auth/login",
		http.MethodPost + " /auth/register",
		http.MethodPost + " /auth/set-password",
		http.MethodPost + " /auth/reset-password",
		http.MethodGet + " /main/profile",
		http.MethodGet + " /main/campus/all",
		http.MethodGet + " /main/curriculum/all",
		http.MethodGet + " /main/invoices",
		http.MethodGet + " /main/applications/all",
		http.MethodGet + " /main/assessment-areas",
		http.MethodGet + " /main/catering/meal-calendar",
		http.MethodGet + " /main/lesson-plan/export-template/:scheme_sub_topic_number",
		// EMS Core Foundation E1: workflow and approvals.
		http.MethodGet + " /main/ems/approvals",
		http.MethodPost + " /main/ems/approvals",
		http.MethodGet + " /main/ems/approvals/:number",
		http.MethodPost + " /main/ems/approvals/:number/approve",
		http.MethodPost + " /main/ems/approvals/:number/reject",
		http.MethodPost + " /main/ems/approvals/:number/return",
		http.MethodPost + " /main/ems/approvals/:number/cancel",
		http.MethodGet + " /main/ems/approval-tasks",
		// EMS Core Foundation F (Increment 1): licence registry and client
		// licence status.
		http.MethodPost + " /main/ems/licence-registry/licences",
		http.MethodGet + " /main/ems/licence-registry/licences",
		http.MethodGet + " /main/ems/licence-registry/licences/:number",
		http.MethodPut + " /main/ems/licence-registry/licences/:number",
		http.MethodPost + " /main/ems/licence-registry/licences/:number/renew",
		http.MethodPost + " /main/ems/licence-registry/licences/:number/suspend",
		http.MethodPost + " /main/ems/licence-registry/licences/:number/revoke",
		http.MethodPost + " /main/ems/licence-registry/licences/:number/replace",
		http.MethodGet + " /main/ems/licence",
		// EMS Core Foundation F (Increment 2): installation activation and
		// control-plane release.
		http.MethodPost + " /main/ems/licence/activate",
		http.MethodGet + " /main/ems/licence/installations",
		http.MethodPost + " /main/ems/licence/installations/:number/deactivate",
		http.MethodPost + " /main/ems/licence/renewal-request",
		http.MethodGet + " /main/ems/licence-registry/installations",
		http.MethodGet + " /main/ems/licence-registry/installations/:number",
		http.MethodPost + " /main/ems/licence-registry/installations/:number/release",
		// EMS Core Foundation F (Increment 5): append-only licensing facts.
		http.MethodGet + " /main/ems/licence-registry/events",
		// EMS Core Foundation G (Increment 1): cross-cutting audit facts
		// (read-only; no update/delete route exists).
		http.MethodGet + " /main/ems/audit",
		http.MethodGet + " /main/ems/audit/:number",
		http.MethodGet + " /main/ems/audit-registry",
	}

	for _, want := range required {
		if !registered[want] {
			t.Errorf("critical route not registered: %s", want)
		}
	}
}
