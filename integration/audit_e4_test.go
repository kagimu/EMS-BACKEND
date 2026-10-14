package integration

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	assignmentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/assignments/models"
	auditModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/audit/models"
)

type e4AuditView struct {
	Number                string `json:"number"`
	Status                string `json:"status"`
	PreviousRequestNumber string `json:"previous_request_number"`
}

// TestFoundationE4AuditIntegration covers Foundation E4 Increment 4: the
// authoritative workflow actions are integrated with the existing G audit
// foundation (append-only ems_audit_logs), while derived observations
// (reminder/overdue/stalled) never create audit facts.
func TestFoundationE4AuditIntegration(t *testing.T) {
	fx := setupEFixture(t)
	organisations := []string{eOrgA, eOrgB}
	fx.db.Unscoped().Where("organisation_number IN ?", organisations).Delete(&auditModel.AuditLog{})
	fx.db.Unscoped().Where("organisation_number IN ?", organisations).Delete(&assignmentModel.UserAssignment{})
	g4CleanNotifications(t, fx.db, organisations)
	t.Cleanup(func() {
		fx.db.Unscoped().Where("organisation_number IN ?", organisations).Delete(&auditModel.AuditLog{})
		fx.db.Unscoped().Where("organisation_number IN ?", organisations).Delete(&assignmentModel.UserAssignment{})
		g4CleanNotifications(t, fx.db, organisations)
	})

	h := newG4Harness(t)

	hasFact := func(action, entityID, actor, org string) bool {
		var count int64
		fx.db.Model(&auditModel.AuditLog{}).
			Where("action = ? AND entity_type = ? AND entity_id = ? AND actor_user_number = ? AND organisation_number = ?",
				action, "ems_approval_request", entityID, actor, org).
			Count(&count)
		return count > 0
	}
	auditForEntity := func(entityID string) int64 {
		var count int64
		fx.db.Model(&auditModel.AuditLog{}).Where("entity_id = ?", entityID).Count(&count)
		return count
	}

	create := func(payload map[string]any) (int, []byte) {
		base := map[string]any{
			"module": "ems", "resource_type": "e4a", "resource_number": "1",
			"decision_type": "purchase_request", "amount": 100.0,
		}
		for key, value := range payload {
			base[key] = value
		}
		return h.do(fx.tokenAdminA, http.MethodPost, eApprovalsBase, base)
	}
	mustCreate := func(payload map[string]any) e4AuditView {
		t.Helper()
		status, body := create(payload)
		if status != http.StatusCreated {
			t.Fatalf("create = %d (%s)", status, string(body))
		}
		var envelope struct {
			Data e4AuditView `json:"data"`
		}
		if err := json.Unmarshal(body, &envelope); err != nil {
			t.Fatalf("decoding create: %v", err)
		}
		return envelope.Data
	}

	t.Run("creation_records_audit_fact", func(t *testing.T) {
		created := mustCreate(nil)
		if !hasFact("approval.create", created.Number, eAdminA, eOrgA) {
			t.Fatal("approval.create audit fact missing")
		}
	})

	t.Run("approval_records_actual_actor", func(t *testing.T) {
		created := mustCreate(nil)
		if status, body := h.do(fx.tokenAppr1, http.MethodPost, eApprovalsBase+"/"+created.Number+"/approve", map[string]any{"comment": "ok"}); status != http.StatusOK {
			t.Fatalf("approve = %d (%s)", status, string(body))
		}
		if !hasFact("approval.approve", created.Number, eApprover1, eOrgA) {
			t.Fatal("approval.approve fact must record the actual approver")
		}
		if hasFact("approval.approve", created.Number, eAdminA, eOrgA) {
			t.Fatal("approval.approve must not record the requester as actor")
		}
	})

	t.Run("advancement_records_actor", func(t *testing.T) {
		created := mustCreate(map[string]any{"amount": 600000.0}) // two-level path
		if status, body := h.do(fx.tokenAppr1, http.MethodPost, eApprovalsBase+"/"+created.Number+"/approve", map[string]any{"comment": "ok"}); status != http.StatusOK {
			t.Fatalf("level-1 approve = %d (%s)", status, string(body))
		}
		if !hasFact("approval.approve", created.Number, eApprover1, eOrgA) {
			t.Fatal("advancement must be captured by the approval audit fact")
		}
	})

	t.Run("rejection_records_actor", func(t *testing.T) {
		created := mustCreate(nil)
		if status, body := h.do(fx.tokenAppr1, http.MethodPost, eApprovalsBase+"/"+created.Number+"/reject", map[string]any{"comment": "no"}); status != http.StatusOK {
			t.Fatalf("reject = %d (%s)", status, string(body))
		}
		if !hasFact("approval.reject", created.Number, eApprover1, eOrgA) {
			t.Fatal("approval.reject fact missing")
		}
	})

	t.Run("return_records_actor", func(t *testing.T) {
		created := mustCreate(nil)
		if status, body := h.do(fx.tokenAppr1, http.MethodPost, eApprovalsBase+"/"+created.Number+"/return", map[string]any{"comment": "fix"}); status != http.StatusOK {
			t.Fatalf("return = %d (%s)", status, string(body))
		}
		if !hasFact("approval.return", created.Number, eApprover1, eOrgA) {
			t.Fatal("approval.return fact missing")
		}
	})

	t.Run("cancellation_records_requester", func(t *testing.T) {
		created := mustCreate(nil)
		if status, body := h.do(fx.tokenAdminA, http.MethodPost, eApprovalsBase+"/"+created.Number+"/cancel", map[string]any{}); status != http.StatusOK {
			t.Fatalf("cancel = %d (%s)", status, string(body))
		}
		if !hasFact("approval.cancel", created.Number, eAdminA, eOrgA) {
			t.Fatal("approval.cancel fact missing")
		}
	})

	t.Run("resubmission_records_create", func(t *testing.T) {
		original := mustCreate(nil)
		if status, body := h.do(fx.tokenAppr1, http.MethodPost, eApprovalsBase+"/"+original.Number+"/return", map[string]any{"comment": "fix"}); status != http.StatusOK {
			t.Fatalf("return = %d (%s)", status, string(body))
		}
		resubmitted := mustCreate(map[string]any{"previous_request_number": original.Number, "resource_number": "2"})
		if resubmitted.PreviousRequestNumber != original.Number {
			t.Fatalf("resubmission linkage not recorded: %q", resubmitted.PreviousRequestNumber)
		}
		if !hasFact("approval.create", resubmitted.Number, eAdminA, eOrgA) {
			t.Fatal("resubmission must be audited (approval.create)")
		}
	})

	t.Run("deadline_override_records_deadline_action", func(t *testing.T) {
		override := mustCreate(map[string]any{"deadline_at": time.Now().UTC().Add(48 * time.Hour)})
		if !hasFact("approval.deadline_set", override.Number, eAdminA, eOrgA) {
			t.Fatal("an explicit per-request deadline override must record approval.deadline_set")
		}
		// The configured default is not a user action and must not be recorded.
		defaulted := mustCreate(nil)
		if hasFact("approval.deadline_set", defaulted.Number, eAdminA, eOrgA) {
			t.Fatal("the configured default deadline must not record a deadline action")
		}
	})

	t.Run("failed_action_records_no_success_fact", func(t *testing.T) {
		created := mustCreate(nil)
		// Return without a reason fails validation.
		if status, body := h.do(fx.tokenAppr1, http.MethodPost, eApprovalsBase+"/"+created.Number+"/return", map[string]any{}); status != http.StatusBadRequest {
			t.Fatalf("return without reason = %d (%s)", status, string(body))
		}
		if hasFact("approval.return", created.Number, eApprover1, eOrgA) {
			t.Fatal("a failed return must not record a success audit fact")
		}
		// An ineligible actor cannot approve.
		if status, _ := h.do(fx.tokenSolo, http.MethodPost, eApprovalsBase+"/"+created.Number+"/approve", map[string]any{"comment": "x"}); status != http.StatusForbidden {
			t.Fatalf("ineligible approve = %d, want 403", status)
		}
		if hasFact("approval.approve", created.Number, eSolo, eOrgA) {
			t.Fatal("a failed approval must not record a success audit fact")
		}
	})

	t.Run("derived_observations_create_no_audit", func(t *testing.T) {
		created := mustCreate(nil)
		fx.db.Table("ems_approval_requests").Where("number = ?", created.Number).
			Update("deadline_at", time.Now().UTC().Add(-time.Hour))

		before := auditForEntity(created.Number)
		// Overdue detail read, list read, and inbox read.
		if status, body := h.do(fx.tokenAdminA, http.MethodGet, eApprovalsBase+"/"+created.Number, nil); status != http.StatusOK {
			t.Fatalf("detail = %d (%s)", status, string(body))
		}
		h.do(fx.tokenAdminA, http.MethodGet, eApprovalsBase, nil)
		h.do(fx.tokenAppr1, http.MethodGet, "/main/ems/approval-tasks", nil)
		if after := auditForEntity(created.Number); after != before {
			t.Fatalf("reads must not create audit facts for derived observations (%d -> %d)", before, after)
		}

		var derived int64
		fx.db.Model(&auditModel.AuditLog{}).
			Where("action LIKE '%reminder%' OR action LIKE '%overdue%' OR action LIKE '%stalled%'").Count(&derived)
		if derived != 0 {
			t.Fatalf("derived conditions must not be audited, found %d", derived)
		}
	})

	t.Run("organisation_isolation", func(t *testing.T) {
		created := mustCreate(nil)
		if !hasFact("approval.create", created.Number, eAdminA, eOrgA) {
			t.Fatal("approval.create fact missing")
		}
		var cross int64
		fx.db.Model(&auditModel.AuditLog{}).
			Where("entity_id = ? AND organisation_number = ?", created.Number, eOrgB).Count(&cross)
		if cross != 0 {
			t.Fatal("workflow audit must not be attributed to another organisation")
		}
		// Another organisation's client audit read cannot see it.
		status, body := h.do(fx.tokenAdminB, http.MethodGet, "/main/ems/audit?entity_id="+created.Number, nil)
		if status != http.StatusOK {
			t.Fatalf("org B audit read = %d (%s)", status, string(body))
		}
		var envelope struct {
			Total int64 `json:"total"`
		}
		json.Unmarshal(body, &envelope)
		if envelope.Total != 0 {
			t.Fatalf("org B saw %d of org A's audit facts", envelope.Total)
		}
	})

	t.Run("audit_is_append_only", func(t *testing.T) {
		// No mutation routes exist for audit records.
		if status, _ := h.do(fx.tokenAdminA, http.MethodPut, "/main/ems/audit/1", map[string]any{}); status == http.StatusOK {
			t.Fatal("audit records must not be mutable")
		}
		if status, _ := h.do(fx.tokenAdminA, http.MethodDelete, "/main/ems/audit/1", nil); status == http.StatusOK {
			t.Fatal("audit records must not be deletable")
		}
		if status, _ := h.do(fx.tokenAdminA, http.MethodDelete, "/main/ems/audit-registry", nil); status == http.StatusOK {
			t.Fatal("audit records must not be deletable")
		}
	})
}
