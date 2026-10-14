package integration

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/routes"
	assignmentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/assignments/models"
	notificationModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/notifications/models"
	"gorm.io/gorm"
)

// g4CleanNotifications removes notifications for the given organisations so G4
// tests do not observe rows left by other tests.
func g4CleanNotifications(t *testing.T, db *gorm.DB, organisations []string) {
	t.Helper()
	if err := db.Unscoped().Where("organisation_number IN ?", organisations).Delete(&notificationModel.Notification{}).Error; err != nil {
		t.Fatalf("cleaning notifications: %v", err)
	}
}

func g4NotificationsFor(db *gorm.DB, userNumber string) []notificationModel.Notification {
	var rows []notificationModel.Notification
	db.Where("recipient_user_number = ?", userNumber).Order("id asc").Find(&rows)
	return rows
}

// g4HasNotification reports whether the recipient has a notification of the
// given type with the given action-required classification.
func g4HasNotification(db *gorm.DB, userNumber, eventType string, actionRequired bool) bool {
	var rows []notificationModel.Notification
	db.Where("recipient_user_number = ? AND type = ?", userNumber, eventType).Find(&rows)
	for _, row := range rows {
		var payload struct {
			ActionRequired bool `json:"action_required"`
		}
		if json.Unmarshal([]byte(row.Payload), &payload) == nil && payload.ActionRequired == actionRequired {
			return true
		}
	}
	return false
}

func g4CountType(db *gorm.DB, eventType string) int64 {
	var count int64
	db.Model(&notificationModel.Notification{}).Where("type = ?", eventType).Count(&count)
	return count
}

// g4Harness provides a server and a request helper for G4 HTTP-level tests.
type g4Harness struct {
	server *httptest.Server
	client *http.Client
	t      *testing.T
}

func newG4Harness(t *testing.T) *g4Harness {
	t.Helper()
	server := httptest.NewServer(routes.SetupRoutes())
	t.Cleanup(server.Close)
	return &g4Harness{server: server, client: server.Client(), t: t}
}

func (h *g4Harness) do(token, method, path string, payload any) (int, []byte) {
	h.t.Helper()
	var reader io.Reader
	if payload != nil {
		raw, err := json.Marshal(payload)
		if err != nil {
			h.t.Fatalf("encoding payload: %v", err)
		}
		reader = bytes.NewReader(raw)
	}
	req, err := http.NewRequest(method, h.server.URL+path, reader)
	if err != nil {
		h.t.Fatalf("building request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := h.client.Do(req)
	if err != nil {
		h.t.Fatalf("request %s %s failed: %v", method, path, err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return resp.StatusCode, body
}

// TestFoundationG4ApprovalNotifications covers the ratified approval map
// (D-GQ1-1, D-GQ1-8/9/10): approver action notifications on create/advance,
// requester outcome notifications, and no notification on cancellation
// (D-GQ1-7).
func TestFoundationG4ApprovalNotifications(t *testing.T) {
	fx := setupEFixture(t)
	organisations := []string{eOrgA, eOrgB}
	g4CleanNotifications(t, fx.db, organisations)
	t.Cleanup(func() { g4CleanNotifications(t, fx.db, organisations) })

	h := newG4Harness(t)

	create := func(resource string, amount float64) string {
		status, body := h.do(fx.tokenAdminA, http.MethodPost, eApprovalsBase, map[string]any{
			"module": "ems", "resource_type": "g4", "resource_number": resource,
			"decision_type": "purchase_request", "amount": amount,
		})
		if status != http.StatusCreated {
			t.Fatalf("create approval = %d (%s)", status, string(body))
		}
		var envelope struct {
			Data struct {
				Number string `json:"number"`
			} `json:"data"`
		}
		if err := json.Unmarshal(body, &envelope); err != nil {
			t.Fatalf("decoding approval: %v", err)
		}
		return envelope.Data.Number
	}

	t.Run("create_notifies_current_approvers", func(t *testing.T) {
		create("g4-a1", 100)
		if !g4HasNotification(fx.db, eApprover1, "approval.create", true) {
			t.Fatal("approver 1 was not notified of the pending approval")
		}
		if !g4HasNotification(fx.db, eApprover2, "approval.create", true) {
			t.Fatal("approver 2 was not notified of the pending approval")
		}
		// The level-2 approver is not current yet, and the requester is not the
		// recipient of a create notification.
		if len(g4NotificationsFor(fx.db, eOwnerA)) != 0 {
			t.Fatal("level-2 approver was notified before being current")
		}
		if len(g4NotificationsFor(fx.db, eAdminA)) != 0 {
			t.Fatal("requester was notified on creation")
		}
	})

	t.Run("advance_and_terminal_outcome", func(t *testing.T) {
		// Amount above the level-1 ceiling so the resolver returns a two-level
		// path; approving level 1 then advances to the level-2 owner.
		number := create("g4-a2", 600000)
		if status, body := h.do(fx.tokenAppr1, http.MethodPost, eApprovalsBase+"/"+number+"/approve", map[string]any{"comment": "ok"}); status != http.StatusOK {
			t.Fatalf("level 1 approve = %d (%s)", status, string(body))
		}
		if !g4HasNotification(fx.db, eOwnerA, "approval.advance", true) {
			t.Fatal("new current approver was not notified on advance")
		}
		if status, body := h.do(fx.tokenOwnerA, http.MethodPost, eApprovalsBase+"/"+number+"/approve", map[string]any{"comment": "ok"}); status != http.StatusOK {
			t.Fatalf("level 2 approve = %d (%s)", status, string(body))
		}
		if !g4HasNotification(fx.db, eAdminA, "approval.approve", false) {
			t.Fatal("requester was not informed of full approval")
		}
	})

	t.Run("rejection_informs_requester", func(t *testing.T) {
		number := create("g4-a3", 100)
		if status, body := h.do(fx.tokenAppr1, http.MethodPost, eApprovalsBase+"/"+number+"/reject", map[string]any{"comment": "no"}); status != http.StatusOK {
			t.Fatalf("reject = %d (%s)", status, string(body))
		}
		if !g4HasNotification(fx.db, eAdminA, "approval.reject", false) {
			t.Fatal("requester was not informed of rejection")
		}
	})

	t.Run("return_requires_requester_action", func(t *testing.T) {
		number := create("g4-a4", 100)
		if status, body := h.do(fx.tokenAppr1, http.MethodPost, eApprovalsBase+"/"+number+"/return", map[string]any{"comment": "fix"}); status != http.StatusOK {
			t.Fatalf("return = %d (%s)", status, string(body))
		}
		if !g4HasNotification(fx.db, eAdminA, "approval.return", true) {
			t.Fatal("requester was not told to act on a return")
		}
	})

	t.Run("cancellation_is_audit_only", func(t *testing.T) {
		number := create("g4-a5", 100)
		if status, body := h.do(fx.tokenAdminA, http.MethodPost, eApprovalsBase+"/"+number+"/cancel", map[string]any{}); status != http.StatusOK {
			t.Fatalf("cancel = %d (%s)", status, string(body))
		}
		if count := g4CountType(fx.db, "approval.cancel"); count != 0 {
			t.Fatalf("cancellation produced %d notifications, want 0", count)
		}
	})
}

// TestFoundationG4StakeholderNotifications covers D-GQ1-2/11: only the affected
// stakeholder is notified, informationally.
func TestFoundationG4StakeholderNotifications(t *testing.T) {
	fx := setupEFixture(t)
	organisations := []string{eOrgA, eOrgB}
	g4CleanNotifications(t, fx.db, organisations)
	t.Cleanup(func() { g4CleanNotifications(t, fx.db, organisations) })

	h := newG4Harness(t)

	status, body := h.do(fx.tokenAdminA, http.MethodPost, "/main/ems/stakeholders", map[string]any{
		"user_number": eApprover1, "stakeholder_type": "staff",
	})
	if status != http.StatusCreated {
		t.Fatalf("create stakeholder = %d (%s)", status, string(body))
	}
	var envelope struct {
		Data struct {
			Number string `json:"number"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		t.Fatalf("decoding stakeholder: %v", err)
	}

	if !g4HasNotification(fx.db, eApprover1, "stakeholder.changed", false) {
		t.Fatal("affected stakeholder was not notified")
	}
	if len(g4NotificationsFor(fx.db, eApprover2)) != 0 {
		t.Fatal("an unrelated user received a stakeholder notification")
	}

	if status, body := h.do(fx.tokenAdminA, http.MethodPut, "/main/ems/stakeholders/"+envelope.Data.Number,
		map[string]any{"title": "Lead"}); status != http.StatusOK {
		t.Fatalf("update stakeholder = %d (%s)", status, string(body))
	}
	if status, body := h.do(fx.tokenAdminA, http.MethodPatch, "/main/ems/stakeholders/"+envelope.Data.Number+"/status",
		map[string]any{"status": "INACTIVE"}); status != http.StatusOK {
		t.Fatalf("stakeholder status = %d (%s)", status, string(body))
	}
	if len(g4NotificationsFor(fx.db, eApprover1)) < 3 {
		t.Fatalf("expected create/update/status notifications, got %d", len(g4NotificationsFor(fx.db, eApprover1)))
	}
}

// TestFoundationG4AssignmentNotifications covers D-GQ1-3/12/13: a new
// assignment is action-required; updates are informational; only the assignee
// is notified.
func TestFoundationG4AssignmentNotifications(t *testing.T) {
	fx := setupEFixture(t)
	organisations := []string{eOrgA, eOrgB}
	g4CleanNotifications(t, fx.db, organisations)
	fx.db.Unscoped().Where("organisation_number IN ?", organisations).Delete(&assignmentModel.UserAssignment{})
	t.Cleanup(func() {
		g4CleanNotifications(t, fx.db, organisations)
		fx.db.Unscoped().Where("organisation_number IN ?", organisations).Delete(&assignmentModel.UserAssignment{})
	})

	h := newG4Harness(t)

	status, body := h.do(fx.tokenAdminA, http.MethodPost, "/main/ems/user-assignments", map[string]any{
		"user_number": eApprover2, "organisation_unit_number": eUnitA, "assignment_type": "member",
	})
	if status != http.StatusCreated {
		t.Fatalf("create assignment = %d (%s)", status, string(body))
	}
	var envelope struct {
		Data struct {
			Number string `json:"number"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		t.Fatalf("decoding assignment: %v", err)
	}

	if !g4HasNotification(fx.db, eApprover2, "assignment.create", true) {
		t.Fatal("assignee was not notified of a new action-required assignment")
	}
	if len(g4NotificationsFor(fx.db, eApprover1)) != 0 {
		t.Fatal("an unrelated user received an assignment notification")
	}

	if status, body := h.do(fx.tokenAdminA, http.MethodPut, "/main/ems/user-assignments/"+envelope.Data.Number,
		map[string]any{"organisation_unit_number": eUnitA, "assignment_type": "member", "position_title": "Updated"}); status != http.StatusOK {
		t.Fatalf("update assignment = %d (%s)", status, string(body))
	}
	if !g4HasNotification(fx.db, eApprover2, "assignment.update", false) {
		t.Fatal("assignee was not informed of the assignment update")
	}
}
