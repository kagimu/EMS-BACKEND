package integration

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/routes"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	supportModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/support/models"
)

func multipartStage4(t *testing.T, client *http.Client, url, token, fieldName, filename string, content []byte) (int, []byte) {
	t.Helper()
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile(fieldName, filename)
	if err != nil {
		t.Fatalf("creating form file: %v", err)
	}
	if _, err := part.Write(content); err != nil {
		t.Fatalf("writing form file: %v", err)
	}
	_ = writer.WriteField("visibility", "INTERNAL")
	_ = writer.Close()
	req, err := http.NewRequest(http.MethodPost, url, &body)
	if err != nil {
		t.Fatalf("building multipart request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("multipart request: %v", err)
	}
	defer resp.Body.Close()
	responseBody, _ := ioReadAll(resp)
	return resp.StatusCode, responseBody
}

func ioReadAll(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(resp.Body)
	return buf.Bytes(), err
}

func TestStage4SupportOperations(t *testing.T) {
	fx := setupStage0BFixture(t)
	seedStage1Licence(t, fx, s0bOrgA, "stage4-licence-a")
	t.Setenv(commons.SoftwareOwnerUserNumbersEnv, s0bAdminA)
	attachmentDir := t.TempDir()
	t.Setenv("EMS_SUPPORT_ATTACHMENT_DIR", attachmentDir)

	server := httptest.NewServer(routes.SetupRoutes())
	t.Cleanup(server.Close)
	client := server.Client()
	req := func(method, path, token string, body any) (int, []byte) {
		return stage1Request(t, client, server.URL, method, path, token, body)
	}

	if status, body := req(http.MethodPost, "/main/ems/owner/support/staff", fx.tokenAdminA, map[string]any{"user_number": s0bUserBoth, "role": "AGENT"}); status != http.StatusOK {
		t.Fatalf("add agent = %d (%s)", status, string(body))
	}
	if status, body := req(http.MethodPost, "/main/ems/owner/support/staff", fx.tokenAdminA, map[string]any{"user_number": s0bUserNoOrg, "role": "MANAGER"}); status != http.StatusOK {
		t.Fatalf("add manager = %d (%s)", status, string(body))
	}

	var ticketNumber string
	status, body := req(http.MethodPost, "/main/ems/owner/support/tickets", fx.tokenAdminA, map[string]any{
		"organisation_number": s0bOrgA, "subject": "Stage4 critical ticket", "description": "Critical issue",
		"category": "TECHNICAL_ERROR", "priority": "CRITICAL",
	})
	if status != http.StatusCreated {
		t.Fatalf("create ticket = %d (%s)", status, string(body))
	}
	var created struct {
		Data struct {
			Ticket struct {
				Number           string `json:"number"`
				SLAResolutionDue string `json:"sla_resolution_due"`
			} `json:"ticket"`
		} `json:"data"`
	}
	_ = json.Unmarshal(body, &created)
	ticketNumber = created.Data.Ticket.Number
	if ticketNumber == "" || created.Data.Ticket.SLAResolutionDue == "" {
		t.Fatalf("ticket/SLA missing: %s", string(body))
	}
	t.Cleanup(func() {
		fx.db.Unscoped().Where("ticket_number = ?", ticketNumber).Delete(&supportModels.Attachment{})
		fx.db.Unscoped().Where("ticket_number = ?", ticketNumber).Delete(&supportModels.Notification{})
		fx.db.Unscoped().Where("ticket_number = ?", ticketNumber).Delete(&supportModels.Message{})
		fx.db.Unscoped().Where("ticket_number = ?", ticketNumber).Delete(&supportModels.Audit{})
		fx.db.Unscoped().Where("number = ?", ticketNumber).Delete(&supportModels.Ticket{})
	})

	t.Run("manager_and_agent_scope", func(t *testing.T) {
		if status, body := req(http.MethodGet, "/main/ems/support/admin/tickets", fx.tokenNoOrg, nil); status != http.StatusOK || !strings.Contains(string(body), ticketNumber) {
			t.Fatalf("manager list = %d (%s)", status, string(body))
		}
		if status, body := req(http.MethodGet, "/main/ems/support/admin/tickets", fx.tokenBoth, nil); status != http.StatusOK || strings.Contains(string(body), ticketNumber) {
			t.Fatalf("agent unassigned list = %d (%s)", status, string(body))
		}
	})

	t.Run("assignment_and_inactive_staff", func(t *testing.T) {
		if status, body := req(http.MethodPatch, "/main/ems/owner/support/tickets/"+ticketNumber, fx.tokenAdminA, map[string]any{"assigned_to_number": s0bUserBoth, "assignment_reason": "Route to agent"}); status != http.StatusOK {
			t.Fatalf("assign = %d (%s)", status, string(body))
		}
		if status, _ := req(http.MethodPatch, "/main/ems/support/admin/tickets/"+ticketNumber, fx.tokenBoth, map[string]any{"assigned_to_number": s0bUserNoOrg}); status != http.StatusForbidden {
			t.Fatalf("agent assignment = %d, want %d", status, http.StatusForbidden)
		}
		if status, body := req(http.MethodPatch, "/main/ems/support/admin/tickets/"+ticketNumber, fx.tokenBoth, map[string]any{"status": "IN_PROGRESS", "priority": "HIGH"}); status != http.StatusOK {
			t.Fatalf("agent update = %d (%s)", status, string(body))
		}
		if status, body := req(http.MethodPost, "/main/ems/support/admin/tickets/"+ticketNumber+"/escalate", fx.tokenBoth, map[string]any{"reason": "Needs manager"}); status != http.StatusOK {
			t.Fatalf("escalate = %d (%s)", status, string(body))
		}
		if status, body := req(http.MethodPatch, "/main/ems/owner/support/staff/"+s0bUserBoth+"/status", fx.tokenAdminA, map[string]any{"status": "INACTIVE"}); status != http.StatusOK {
			t.Fatalf("deactivate = %d (%s)", status, string(body))
		}
		if status, _ := req(http.MethodPatch, "/main/ems/owner/support/tickets/"+ticketNumber, fx.tokenAdminA, map[string]any{"assigned_to_number": s0bUserBoth}); status != http.StatusBadRequest {
			t.Fatalf("inactive assignment = %d, want %d", status, http.StatusBadRequest)
		}
	})

	t.Run("notifications_and_attachment_visibility", func(t *testing.T) {
		var notificationCount int64
		if err := fx.db.Model(&supportModels.Notification{}).Where("ticket_number = ?", ticketNumber).Count(&notificationCount).Error; err != nil {
			t.Fatalf("counting notifications: %v", err)
		}
		if notificationCount == 0 {
			t.Fatal("expected support notification record")
		}
		var notification supportModels.Notification
		fx.db.Where("ticket_number = ?", ticketNumber).First(&notification)
		if notification.DeliveryStatus != supportModels.DeliveryInApp {
			t.Fatalf("delivery status = %s", notification.DeliveryStatus)
		}

		status, body := multipartStage4(t, client, server.URL+"/main/ems/support/admin/tickets/"+ticketNumber+"/attachments", fx.tokenNoOrg, "file", "note.txt", []byte("internal diagnostic"))
		if status != http.StatusCreated {
			t.Fatalf("attachment upload = %d (%s)", status, string(body))
		}
		var attachment struct {
			Data struct {
				Number string `json:"number"`
			} `json:"data"`
		}
		_ = json.Unmarshal(body, &attachment)
		if attachment.Data.Number == "" {
			t.Fatalf("missing attachment number")
		}

		status, body = req(http.MethodGet, "/main/ems/support/tickets/"+ticketNumber, fx.tokenBoth, nil)
		if status != http.StatusOK {
			t.Fatalf("client detail = %d (%s)", status, string(body))
		}
		if strings.Contains(string(body), attachment.Data.Number) {
			t.Fatal("internal attachment leaked to customer")
		}
	})

	t.Run("sla_overdue_and_report", func(t *testing.T) {
		fx.db.Model(&supportModels.Ticket{}).Where("number = ?", ticketNumber).UpdateColumns(map[string]any{"sla_resolution_due": "2020-01-01 00:00:00"})
		status, body := req(http.MethodGet, "/main/ems/owner/support/tickets/"+ticketNumber, fx.tokenAdminA, nil)
		if status != http.StatusOK || !strings.Contains(string(body), `"resolution_overdue":true`) {
			t.Fatalf("overdue detail = %d (%s)", status, string(body))
		}
		if status, body := req(http.MethodGet, "/main/ems/owner/support/reports", fx.tokenAdminA, nil); status != http.StatusOK || !strings.Contains(string(body), "tickets_received") {
			t.Fatalf("report = %d (%s)", status, string(body))
		}
		if status, _ := req(http.MethodGet, "/main/ems/owner/support/sla", fx.tokenAdminA, nil); status != http.StatusOK {
			t.Fatalf("sla targets = %d", status)
		}
	})

	_ = os.Getenv("EMS_SUPPORT_ATTACHMENT_DIR")
}
