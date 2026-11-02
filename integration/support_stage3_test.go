package integration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/routes"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	supportModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/support/models"
)

func TestStage3SupportCentre(t *testing.T) {
	fx := setupStage0BFixture(t)
	seedStage1Licence(t, fx, s0bOrgA, "stage3-licence-a")
	t.Setenv(commons.SoftwareOwnerUserNumbersEnv, s0bAdminA)

	server := httptest.NewServer(routes.SetupRoutes())
	t.Cleanup(server.Close)
	client := server.Client()

	req := func(method, path, token string, body any) (int, []byte) {
		return stage1Request(t, client, server.URL, method, path, token, body)
	}

	ticketsToCleanup := []string{}
	t.Cleanup(func() {
		for _, number := range ticketsToCleanup {
			fx.db.Unscoped().Where("ticket_number = ?", number).Delete(&supportModels.Message{})
			fx.db.Unscoped().Where("ticket_number = ?", number).Delete(&supportModels.Audit{})
			fx.db.Unscoped().Where("number = ?", number).Delete(&supportModels.Ticket{})
		}
	})

	t.Run("unauthenticated_owner_support_is_rejected", func(t *testing.T) {
		if status, _ := req(http.MethodGet, "/main/ems/owner/support/overview", "", nil); status != http.StatusUnauthorized {
			t.Fatalf("status = %d, want %d", status, http.StatusUnauthorized)
		}
	})

	t.Run("client_without_permission_is_rejected", func(t *testing.T) {
		if status, _ := req(http.MethodGet, "/main/ems/support/tickets", fx.tokenInactive, nil); status != http.StatusForbidden {
			t.Fatalf("status = %d, want %d", status, http.StatusForbidden)
		}
	})

	var ownerTicketA, ownerTicketB, clientTicket string

	t.Run("owner_creates_tickets_for_two_organisations", func(t *testing.T) {
		status, body := req(http.MethodPost, "/main/ems/owner/support/tickets", fx.tokenAdminA, map[string]any{
			"organisation_number": s0bOrgA, "subject": "Stage3 owner ticket A", "description": "Licence activation issue",
			"category": "LICENCE_ACTIVATION", "priority": "HIGH",
		})
		if status != http.StatusCreated {
			t.Fatalf("create A = %d (%s)", status, string(body))
		}
		var decoded struct {
			Data struct {
				Ticket struct {
					Number string `json:"number"`
				} `json:"ticket"`
			} `json:"data"`
		}
		_ = json.Unmarshal(body, &decoded)
		ownerTicketA = decoded.Data.Ticket.Number
		if ownerTicketA == "" {
			t.Fatalf("missing ticket number: %s", string(body))
		}
		ticketsToCleanup = append(ticketsToCleanup, ownerTicketA)

		status, body = req(http.MethodPost, "/main/ems/owner/support/tickets", fx.tokenAdminA, map[string]any{
			"organisation_number": s0bOrgB, "subject": "Stage3 owner ticket B", "description": "Other organisation",
			"category": "GENERAL_ENQUIRY", "priority": "MEDIUM",
		})
		if status != http.StatusCreated {
			t.Fatalf("create B = %d (%s)", status, string(body))
		}
		_ = json.Unmarshal(body, &decoded)
		ownerTicketB = decoded.Data.Ticket.Number
		ticketsToCleanup = append(ticketsToCleanup, ownerTicketB)
	})

	t.Run("client_creates_own_ticket", func(t *testing.T) {
		status, body := req(http.MethodPost, "/main/ems/support/tickets", fx.tokenBoth, map[string]any{
			"subject": "Stage3 client ticket", "description": "Cannot see a module",
			"category": "MODULE_ACCESS", "priority": "LOW",
		})
		if status != http.StatusCreated {
			t.Fatalf("create client = %d (%s)", status, string(body))
		}
		var decoded struct {
			Data struct {
				Ticket struct {
					Number string `json:"number"`
				} `json:"ticket"`
			} `json:"data"`
		}
		_ = json.Unmarshal(body, &decoded)
		clientTicket = decoded.Data.Ticket.Number
		if clientTicket == "" {
			t.Fatalf("missing client ticket number")
		}
		ticketsToCleanup = append(ticketsToCleanup, clientTicket)
	})

	t.Run("client_listing_is_organisation_scoped", func(t *testing.T) {
		status, body := req(http.MethodGet, "/main/ems/support/tickets?per_page=100", fx.tokenBoth, nil)
		if status != http.StatusOK {
			t.Fatalf("status = %d (%s)", status, string(body))
		}
		if !strings.Contains(string(body), clientTicket) || strings.Contains(string(body), ownerTicketB) {
			t.Fatalf("client listing leaked or missed tickets: %s", string(body))
		}
		if status, _ := req(http.MethodGet, "/main/ems/support/tickets/"+ownerTicketB, fx.tokenBoth, nil); status != http.StatusNotFound {
			t.Fatalf("cross-organisation detail = %d, want %d", status, http.StatusNotFound)
		}
	})

	t.Run("internal_notes_are_never_exposed_to_customers", func(t *testing.T) {
		internalBody := "INTERNAL-ONLY-DIAGNOSTIC"
		if status, body := req(http.MethodPost, "/main/ems/owner/support/tickets/"+ownerTicketA+"/messages", fx.tokenAdminA, map[string]any{
			"body": internalBody, "visibility": "INTERNAL",
		}); status != http.StatusCreated {
			t.Fatalf("internal note = %d (%s)", status, string(body))
		}

		if status, body := req(http.MethodPost, "/main/ems/owner/support/tickets/"+ownerTicketA+"/messages", fx.tokenAdminA, map[string]any{
			"body": "Customer-visible update", "visibility": "CUSTOMER",
		}); status != http.StatusCreated {
			t.Fatalf("customer reply = %d (%s)", status, string(body))
		}

		status, body := req(http.MethodGet, "/main/ems/support/tickets/"+ownerTicketA, fx.tokenBoth, nil)
		if status != http.StatusOK {
			t.Fatalf("client detail = %d (%s)", status, string(body))
		}
		if strings.Contains(string(body), internalBody) || !strings.Contains(string(body), "Customer-visible update") {
			t.Fatalf("visibility filter failed: %s", string(body))
		}
	})

	t.Run("client_reply_and_owner_conversation", func(t *testing.T) {
		if status, body := req(http.MethodPost, "/main/ems/support/tickets/"+ownerTicketA+"/messages", fx.tokenBoth, map[string]any{"body": "Customer additional information"}); status != http.StatusCreated {
			t.Fatalf("client reply = %d (%s)", status, string(body))
		}
		if status, body := req(http.MethodGet, "/main/ems/owner/support/tickets/"+ownerTicketA, fx.tokenAdminA, nil); status != http.StatusOK || !strings.Contains(string(body), "Customer additional information") {
			t.Fatalf("owner conversation = %d (%s)", status, string(body))
		}
	})

	t.Run("owner_updates_assignment_and_status", func(t *testing.T) {
		status, body := req(http.MethodPatch, "/main/ems/owner/support/tickets/"+ownerTicketA, fx.tokenAdminA, map[string]any{
			"assigned_to_number": s0bAdminA, "priority": "CRITICAL", "status": "IN_PROGRESS",
		})
		if status != http.StatusOK {
			t.Fatalf("update = %d (%s)", status, string(body))
		}

		if status, _ := req(http.MethodPatch, "/main/ems/owner/support/tickets/"+ownerTicketA, fx.tokenAdminA, map[string]any{"status": "RESOLVED"}); status != http.StatusBadRequest {
			t.Fatalf("resolve without summary = %d, want %d", status, http.StatusBadRequest)
		}
		if status, body := req(http.MethodPatch, "/main/ems/owner/support/tickets/"+ownerTicketA, fx.tokenAdminA, map[string]any{
			"status": "RESOLVED", "resolution_summary": "Restored activation",
		}); status != http.StatusOK {
			t.Fatalf("resolve = %d (%s)", status, string(body))
		}

		if status, body := req(http.MethodPatch, "/main/ems/owner/support/tickets/"+ownerTicketA, fx.tokenAdminA, map[string]any{"status": "OPEN"}); status != http.StatusOK {
			t.Fatalf("reopen = %d (%s)", status, string(body))
		}
		req(http.MethodPatch, "/main/ems/owner/support/tickets/"+ownerTicketA, fx.tokenAdminA, map[string]any{"status": "CLOSED"})
		if status, _ := req(http.MethodPatch, "/main/ems/owner/support/tickets/"+ownerTicketA, fx.tokenAdminA, map[string]any{"status": "IN_PROGRESS"}); status != http.StatusConflict {
			t.Fatalf("closed ticket transition = %d, want %d", status, http.StatusConflict)
		}
	})

	t.Run("support_actions_are_audited_and_overview_counts", func(t *testing.T) {
		var count int64
		if err := fx.db.Model(&supportModels.Audit{}).Where("ticket_number = ?", ownerTicketA).Count(&count).Error; err != nil {
			t.Fatalf("counting support audit: %v", err)
		}
		if count < 5 {
			t.Fatalf("support audit count = %d, want at least 5", count)
		}
		if status, body := req(http.MethodGet, "/main/ems/owner/support/overview", fx.tokenAdminA, nil); status != http.StatusOK || !strings.Contains(string(body), "total") {
			t.Fatalf("overview = %d (%s)", status, string(body))
		}
	})
}
