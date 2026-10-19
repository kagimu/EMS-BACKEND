package integration

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	authServices "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/services"
	assignmentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/assignments/models"
	auditModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/audit/models"
	notificationModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/notifications/models"
	workflowModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/workflow/models"
	eosModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/eos/models"
	permissionsModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/permissions/models"
)

const (
	e4WfManager = "ems-e4-wf-manager"
	e4WfRole    = "ems-e4-wf-role"

	e4DeadlineConfigURL = "/main/ems/workflow/deadline-config"
)

type e4Cfg struct {
	OrganisationNumber   string `json:"organisation_number"`
	DefaultDeadlineHours int    `json:"default_deadline_hours"`
}

type e4CfgEnv struct {
	Data e4Cfg `json:"data"`
}

// TestFoundationE4DeadlineConfig covers Foundation E4 Increment 2: the
// organisation-scoped configurable default deadline, its authority, validation,
// isolation, historical safety and audit.
func TestFoundationE4DeadlineConfig(t *testing.T) {
	fx := setupEFixture(t)
	organisations := []string{eOrgA, eOrgB}
	g4CleanNotifications(t, fx.db, organisations)
	fx.db.Unscoped().Where("organisation_number IN ?", organisations).Delete(&assignmentModel.UserAssignment{})
	fx.db.Unscoped().Where("organisation_number IN ?", organisations).Delete(&workflowModel.WorkflowDeadlineConfig{})

	// A non-admin workflow manager proves permission-based (not Admin-bypass) access.
	wfUser := authModel.User{Number: e4WfManager, FirstName: "E4", SecondName: "Manager", Gender: "M", Type: "User", Status: "ACTIVE"}
	if err := fx.db.Create(&wfUser).Error; err != nil {
		t.Fatalf("seeding workflow manager: %v", err)
	}
	if err := fx.db.Omit("Entity", "User").Create(&eosModel.ActiveUserEntity{UserNumber: e4WfManager, EntityNumber: eOrgA}).Error; err != nil {
		t.Fatalf("activating workflow manager entity: %v", err)
	}
	if err := fx.db.Omit("Entity", "User", "AddedBy").Create(&eosModel.EntityUser{
		EntityNumber:  eOrgA,
		UserNumber:    e4WfManager,
		AddedByNumber: e4WfManager,
	}).Error; err != nil {
		t.Fatalf("seeding workflow manager membership: %v", err)
	}
	role := permissionsModel.Role{Number: e4WfRole, Name: e4WfRole, EntityNumber: eOrgA, Status: "ACTIVE", UserNumber: e4WfManager}
	if err := fx.db.Omit("User", "Entity", "Permissions").Create(&role).Error; err != nil {
		t.Fatalf("seeding workflow role: %v", err)
	}
	for _, permission := range []string{"view_ems_workflow", "manage_ems_workflow"} {
		rp := permissionsModel.RolePermission{RoleNumber: e4WfRole, Permission: permission, AddedByNumber: e4WfManager}
		if err := fx.db.Omit("AddedBy", "Role").Create(&rp).Error; err != nil {
			t.Fatalf("seeding workflow permission: %v", err)
		}
	}
	expiry := time.Now().Add(365 * 24 * time.Hour)
	if err := fx.db.Omit("User", "Role", "AddedByUser").Create(&permissionsModel.UserRole{UserNumber: e4WfManager, RoleNumber: e4WfRole, ExpiryDate: expiry, AddedBy: e4WfManager}).Error; err != nil {
		t.Fatalf("assigning workflow role: %v", err)
	}

	t.Cleanup(func() {
		g4CleanNotifications(t, fx.db, organisations)
		fx.db.Unscoped().Where("organisation_number IN ?", organisations).Delete(&assignmentModel.UserAssignment{})
		fx.db.Unscoped().Where("organisation_number IN ?", organisations).Delete(&workflowModel.WorkflowDeadlineConfig{})
		fx.db.Unscoped().Where("role_number = ?", e4WfRole).Delete(&permissionsModel.UserRole{})
		fx.db.Unscoped().Where("role_number = ?", e4WfRole).Delete(&permissionsModel.RolePermission{})
		fx.db.Unscoped().Where("number = ?", e4WfRole).Delete(&permissionsModel.Role{})
		fx.db.Unscoped().Where("user_number = ?", e4WfManager).Delete(&eosModel.ActiveUserEntity{})
		fx.db.Unscoped().Where("number = ?", e4WfManager).Delete(&authModel.User{})
	})

	managerToken, err := authServices.JWTAuthService().GenerateToken(e4WfManager)
	if err != nil {
		t.Fatalf("manager token: %v", err)
	}

	h := newG4Harness(t)

	putCfg := func(token string, hours int) (int, []byte) {
		return h.do(token, http.MethodPut, e4DeadlineConfigURL, map[string]any{"default_deadline_hours": hours})
	}
	getCfg := func(token string) (int, e4Cfg) {
		t.Helper()
		status, body := h.do(token, http.MethodGet, e4DeadlineConfigURL, nil)
		var env e4CfgEnv
		json.Unmarshal(body, &env)
		return status, env.Data
	}
	createReq := func(deadline any) e4RequestView {
		t.Helper()
		payload := map[string]any{
			"module": "ems", "resource_type": "e4", "resource_number": "c",
			"decision_type": "purchase_request", "amount": 100,
		}
		if deadline != nil {
			payload["deadline_at"] = deadline
		}
		status, body := h.do(fx.tokenAdminA, http.MethodPost, eApprovalsBase, payload)
		if status != http.StatusCreated {
			t.Fatalf("create request = %d (%s)", status, string(body))
		}
		var env e4Envelope
		json.Unmarshal(body, &env)
		return env.Data
	}
	fetchReq := func(token, number string) e4RequestView {
		t.Helper()
		status, body := h.do(token, http.MethodGet, eApprovalsBase+"/"+number, nil)
		if status != http.StatusOK {
			t.Fatalf("fetch request = %d (%s)", status, string(body))
		}
		var env e4Envelope
		json.Unmarshal(body, &env)
		return env.Data
	}
	within := func(got time.Time, hours float64) bool {
		delta := time.Until(got.UTC()).Hours()
		return delta > hours-1 && delta < hours+1
	}

	t.Run("default_is_seven_days_when_unset", func(t *testing.T) {
		if status, cfg := getCfg(managerToken); status != http.StatusOK || cfg.DefaultDeadlineHours != 168 {
			t.Fatalf("unset default = %d %+v, want 200/168", status, cfg)
		}
	})

	t.Run("authorised_manager_can_change", func(t *testing.T) {
		if status, body := putCfg(managerToken, 24); status != http.StatusOK {
			t.Fatalf("manager update = %d (%s)", status, string(body))
		}
		if _, cfg := getCfg(managerToken); cfg.DefaultDeadlineHours != 24 {
			t.Fatalf("configured default = %d, want 24", cfg.DefaultDeadlineHours)
		}
	})

	t.Run("unauthorised_forbidden", func(t *testing.T) {
		if status, _ := putCfg(fx.tokenNoPerm, 24); status != http.StatusForbidden {
			t.Fatalf("no-permission update = %d, want 403", status)
		}
		if status, _ := getCfg(fx.tokenNoPerm); status != http.StatusForbidden {
			t.Fatalf("no-permission read = %d, want 403", status)
		}
	})

	t.Run("invalid_durations_rejected", func(t *testing.T) {
		for _, hours := range []int{0, -5, 8761, 999999} {
			if status, _ := putCfg(managerToken, hours); status != http.StatusBadRequest {
				t.Fatalf("hours=%d = %d, want 400", hours, status)
			}
		}
		if status, _ := h.do(managerToken, http.MethodPut, e4DeadlineConfigURL, map[string]any{"default_deadline_hours": "x"}); status != http.StatusBadRequest {
			t.Fatalf("malformed body = %d, want 400", status)
		}
	})

	t.Run("organisation_isolation", func(t *testing.T) {
		if status, body := putCfg(fx.tokenAdminB, 12); status != http.StatusOK {
			t.Fatalf("org B update = %d (%s)", status, string(body))
		}
		if _, cfg := getCfg(managerToken); cfg.DefaultDeadlineHours != 24 {
			t.Fatalf("org A config changed by org B action: %d", cfg.DefaultDeadlineHours)
		}
		if _, cfg := getCfg(fx.tokenAdminB); cfg.DefaultDeadlineHours != 12 {
			t.Fatalf("org B config = %d, want 12", cfg.DefaultDeadlineHours)
		}
	})

	t.Run("new_request_uses_configured_default", func(t *testing.T) {
		if status, _ := putCfg(managerToken, 24); status != http.StatusOK {
			t.Fatalf("set 24 failed")
		}
		created := createReq(nil)
		view := fetchReq(fx.tokenAdminA, created.Number)
		if view.DeadlineAt == nil || !within(*view.DeadlineAt, 24) {
			t.Fatalf("new request deadline %v not ~24h from now", view.DeadlineAt)
		}
	})

	t.Run("override_beats_config", func(t *testing.T) {
		override := time.Now().UTC().Add(100 * time.Hour)
		created := createReq(override)
		view := fetchReq(fx.tokenAdminA, created.Number)
		if view.DeadlineAt == nil || !view.DeadlineAt.UTC().Truncate(time.Second).Equal(override.Truncate(time.Second)) {
			t.Fatalf("override not honoured over configured default: got %v want %v", view.DeadlineAt, override)
		}
	})

	t.Run("existing_request_retains_deadline", func(t *testing.T) {
		if status, _ := putCfg(managerToken, 168); status != http.StatusOK {
			t.Fatalf("set 168 failed")
		}
		existing := createReq(nil)
		original := fetchReq(fx.tokenAdminA, existing.Number).DeadlineAt
		if original == nil || !within(*original, 168) {
			t.Fatalf("existing deadline not ~168h: %v", original)
		}

		if status, _ := putCfg(managerToken, 1); status != http.StatusOK {
			t.Fatalf("set 1 failed")
		}
		after := fetchReq(fx.tokenAdminA, existing.Number)
		if after.DeadlineAt == nil || !after.DeadlineAt.UTC().Equal(original.UTC()) {
			t.Fatalf("configuration change rewrote an existing deadline: %v -> %v", original, after.DeadlineAt)
		}
		if after.Status != "PENDING" {
			t.Fatalf("configuration change altered request status: %s", after.Status)
		}
		// A new request now uses the new default.
		fresh := createReq(nil)
		if freshView := fetchReq(fx.tokenAdminA, fresh.Number); freshView.DeadlineAt == nil || !within(*freshView.DeadlineAt, 1) {
			t.Fatalf("new request did not use the new default: %v", freshView.DeadlineAt)
		}
	})

	t.Run("overdue_still_derived", func(t *testing.T) {
		created := createReq(nil)
		past := time.Now().UTC().Add(-time.Hour)
		fx.db.Table("ems_approval_requests").Where("number = ?", created.Number).Update("deadline_at", past)
		view := fetchReq(fx.tokenAdminA, created.Number)
		if !view.Overdue || view.Status != "PENDING" {
			t.Fatalf("overdue must remain derived and non-terminal: overdue=%v status=%s", view.Overdue, view.Status)
		}
	})

	t.Run("config_change_creates_audit_fact", func(t *testing.T) {
		var before int64
		fx.db.Model(&auditModel.AuditLog{}).
			Where("action = ? AND organisation_number = ?", "workflow.deadline_config.update", eOrgA).Count(&before)
		if status, _ := putCfg(managerToken, 36); status != http.StatusOK {
			t.Fatalf("set 36 failed")
		}
		var after int64
		fx.db.Model(&auditModel.AuditLog{}).
			Where("action = ? AND organisation_number = ?", "workflow.deadline_config.update", eOrgA).Count(&after)
		if after != before+1 {
			t.Fatalf("expected one new audit fact, got %d -> %d", before, after)
		}
		var row auditModel.AuditLog
		fx.db.Where("action = ? AND organisation_number = ?", "workflow.deadline_config.update", eOrgA).
			Order("id desc").First(&row)
		if row.ActorUserNumber != e4WfManager || row.EntityType != "ems_workflow_deadline_config" || row.EntityID != eOrgA {
			t.Fatalf("audit fact missing actor/entity: %+v", row)
		}
		if !jsonHasHours(row.OldValues, 1) || !jsonHasHours(row.NewValues, 36) {
			t.Fatalf("audit fact missing old/new values: old=%s new=%s", row.OldValues, row.NewValues)
		}
	})

	t.Run("reads_create_no_deadline_audit", func(t *testing.T) {
		var before int64
		fx.db.Model(&auditModel.AuditLog{}).
			Where("action = ? AND organisation_number = ?", "workflow.deadline_config.update", eOrgA).Count(&before)
		getCfg(managerToken)
		created := createReq(nil)
		fetchReq(fx.tokenAdminA, created.Number)
		var after int64
		fx.db.Model(&auditModel.AuditLog{}).
			Where("action = ? AND organisation_number = ?", "workflow.deadline_config.update", eOrgA).Count(&after)
		if after != before {
			t.Fatalf("reads must not create deadline audit facts: %d -> %d", before, after)
		}
	})

	t.Run("config_change_creates_no_notification", func(t *testing.T) {
		var before int64
		fx.db.Model(&notificationModel.Notification{}).Where("organisation_number = ?", eOrgA).Count(&before)
		if status, _ := putCfg(managerToken, 48); status != http.StatusOK {
			t.Fatalf("set 48 failed")
		}
		var after int64
		fx.db.Model(&notificationModel.Notification{}).Where("organisation_number = ?", eOrgA).Count(&after)
		if after != before {
			t.Fatalf("E4-2 must not create notifications: %d -> %d", before, after)
		}
	})
}

// jsonHasHours reports whether a JSON document carries default_deadline_hours
// equal to hours.
func jsonHasHours(raw string, hours int) bool {
	var doc map[string]any
	if err := json.Unmarshal([]byte(raw), &doc); err != nil {
		return false
	}
	value, ok := doc["default_deadline_hours"].(float64)
	return ok && int(value) == hours
}
