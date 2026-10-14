package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	assignmentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/assignments/models"
	auditModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/audit/models"
	notificationModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/notifications/models"
	workflowServices "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/workflow/services"
	"gorm.io/gorm/logger"
)

type e4RequestView struct {
	Number     string     `json:"number"`
	Status     string     `json:"status"`
	Stalled    bool       `json:"stalled"`
	Reminder   bool       `json:"reminder"`
	Overdue    bool       `json:"overdue"`
	DeadlineAt *time.Time `json:"deadline_at"`
}

type e4Envelope struct {
	Data e4RequestView `json:"data"`
}

type e4Row struct {
	Status       string
	CurrentLevel int
	DeadlineAt   *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// e4CountingLogger counts SQL statements issued through GORM so a test can prove
// the deadline read path adds no per-row query (C5 / no N+1).
type e4CountingLogger struct{ queries int64 }

func (l *e4CountingLogger) LogMode(logger.LogLevel) logger.Interface { return l }
func (l *e4CountingLogger) Info(context.Context, string, ...interface{}) {
}
func (l *e4CountingLogger) Warn(context.Context, string, ...interface{}) {
}
func (l *e4CountingLogger) Error(context.Context, string, ...interface{}) {
}
func (l *e4CountingLogger) Trace(context.Context, time.Time, func() (string, int64), error) {
	atomic.AddInt64(&l.queries, 1)
}

// TestFoundationE4AdvisoryDeadline covers Foundation E4 under the ratified
// Option B / C4 architecture: advisory per-request deadlines and the read-time
// DERIVED reminder/overdue conditions, surfaced on the approval/inbox reads.
// Deadline evaluation performs no durable write and creates no notification.
func TestFoundationE4AdvisoryDeadline(t *testing.T) {
	fx := setupEFixture(t)
	organisations := []string{eOrgA, eOrgB}
	g4CleanNotifications(t, fx.db, organisations)
	fx.db.Unscoped().Where("organisation_number IN ?", organisations).Delete(&assignmentModel.UserAssignment{})
	t.Cleanup(func() {
		g4CleanNotifications(t, fx.db, organisations)
		fx.db.Unscoped().Where("organisation_number IN ?", organisations).Delete(&assignmentModel.UserAssignment{})
	})

	h := newG4Harness(t)

	create := func(deadline any) (int, []byte) {
		payload := map[string]any{
			"module": "ems", "resource_type": "e4", "resource_number": "d",
			"decision_type": "purchase_request", "amount": 100,
		}
		if deadline != nil {
			payload["deadline_at"] = deadline
		}
		return h.do(fx.tokenAdminA, http.MethodPost, eApprovalsBase, payload)
	}
	fetch := func(token, number string) e4RequestView {
		t.Helper()
		status, body := h.do(token, http.MethodGet, eApprovalsBase+"/"+number, nil)
		if status != http.StatusOK {
			t.Fatalf("fetch %s = %d (%s)", number, status, string(body))
		}
		var env e4Envelope
		if err := json.Unmarshal(body, &env); err != nil {
			t.Fatalf("decoding request view: %v", err)
		}
		return env.Data
	}
	loadRow := func(number string) e4Row {
		t.Helper()
		var row e4Row
		fx.db.Table("ems_approval_requests").
			Select("status, current_level, deadline_at, created_at, updated_at").
			Where("number = ?", number).Scan(&row)
		return row
	}
	mustCreate := func(deadline any) e4RequestView {
		t.Helper()
		status, body := create(deadline)
		if status != http.StatusCreated {
			t.Fatalf("create = %d (%s)", status, string(body))
		}
		var env e4Envelope
		if err := json.Unmarshal(body, &env); err != nil {
			t.Fatalf("decoding create: %v", err)
		}
		return env.Data
	}
	setDeadline := func(number string, deadline time.Time) {
		t.Helper()
		if err := fx.db.Table("ems_approval_requests").Where("number = ?", number).
			Update("deadline_at", deadline).Error; err != nil {
			t.Fatalf("setting deadline: %v", err)
		}
	}
	notificationCount := func() int64 {
		var count int64
		fx.db.Model(&notificationModel.Notification{}).Where("organisation_number = ?", eOrgA).Count(&count)
		return count
	}

	t.Run("default_deadline_when_no_override", func(t *testing.T) {
		created := mustCreate(nil)
		view := fetch(fx.tokenAdminA, created.Number)
		if view.DeadlineAt == nil {
			t.Fatal("effective deadline must be exposed even without an override")
		}
		if view.Overdue || view.Reminder {
			t.Fatalf("a request created just now must be neither reminder nor overdue (reminder=%v overdue=%v)", view.Reminder, view.Overdue)
		}
		row := loadRow(created.Number)
		// Foundation E4-2: the resolved default deadline is persisted at creation.
		if row.DeadlineAt == nil {
			t.Fatal("the resolved default deadline must be persisted (E4-2 historical invariant)")
		}
		if !row.DeadlineAt.UTC().Equal(view.DeadlineAt.UTC()) {
			t.Fatal("the response deadline must equal the persisted deadline")
		}
		delta := view.DeadlineAt.UTC().Sub(row.CreatedAt.UTC())
		if delta < 167*time.Hour || delta > 169*time.Hour {
			t.Fatalf("default deadline should be ~7d after creation, got %s", delta)
		}
	})

	t.Run("override_is_honoured_and_stored_utc", func(t *testing.T) {
		target := time.Date(2031, 3, 4, 5, 6, 7, 0, time.FixedZone("EAT", 3*3600))
		created := mustCreate(target)
		view := fetch(fx.tokenAdminA, created.Number)
		if view.DeadlineAt == nil || !view.DeadlineAt.UTC().Equal(target.UTC()) {
			t.Fatalf("override not honoured: got %v want %v", view.DeadlineAt, target.UTC())
		}
		if view.DeadlineAt.Location() != time.UTC {
			t.Fatalf("deadline must be UTC, got location %v", view.DeadlineAt.Location())
		}
		if view.Overdue || view.Reminder {
			t.Fatal("a far-future deadline must be neither reminder nor overdue")
		}
		row := loadRow(created.Number)
		if row.DeadlineAt == nil || !row.DeadlineAt.UTC().Equal(target.UTC()) {
			t.Fatalf("override not persisted correctly: %v", row.DeadlineAt)
		}
	})

	t.Run("reminder_window_surfacing", func(t *testing.T) {
		// More than 24h away: neither condition.
		far := mustCreate(time.Now().UTC().Add(48 * time.Hour))
		if v := fetch(fx.tokenAdminA, far.Number); v.Reminder || v.Overdue {
			t.Fatalf("48h remaining => reminder=false overdue=false, got reminder=%v overdue=%v", v.Reminder, v.Overdue)
		}
		// Within one hour (<= 24h and still in the future): approaching deadline.
		near := mustCreate(time.Now().UTC().Add(time.Hour))
		if v := fetch(fx.tokenAdminA, near.Number); !v.Reminder || v.Overdue {
			t.Fatalf("1h remaining => reminder=true overdue=false, got reminder=%v overdue=%v", v.Reminder, v.Overdue)
		}
		// Exactly at the deadline boundary (set to now): reminder true, overdue false.
		at := mustCreate(time.Now().UTC().Add(2 * time.Hour))
		setDeadline(at.Number, time.Now().UTC().Add(30*time.Second))
		if v := fetch(fx.tokenAdminA, at.Number); !v.Reminder || v.Overdue {
			t.Fatalf("~at deadline => reminder=true overdue=false, got reminder=%v overdue=%v", v.Reminder, v.Overdue)
		}
	})

	t.Run("overdue_is_mutually_exclusive_with_reminder", func(t *testing.T) {
		created := mustCreate(time.Now().UTC().Add(72 * time.Hour))
		setDeadline(created.Number, time.Now().UTC().Add(-time.Hour))
		view := fetch(fx.tokenAdminA, created.Number)
		if !view.Overdue {
			t.Fatal("a past deadline on a PENDING request must read as overdue")
		}
		if view.Reminder {
			t.Fatal("a reminder must not be active once overdue")
		}
		if view.Status != "PENDING" {
			t.Fatalf("overdue must not change the status: got %s", view.Status)
		}
	})

	t.Run("non_pending_is_neither", func(t *testing.T) {
		// Cancel a request with a near deadline; a non-PENDING request is neither.
		created := mustCreate(time.Now().UTC().Add(time.Hour))
		if status, body := h.do(fx.tokenAdminA, http.MethodPost, eApprovalsBase+"/"+created.Number+"/cancel", map[string]any{}); status != http.StatusOK {
			t.Fatalf("cancel = %d (%s)", status, string(body))
		}
		view := fetch(fx.tokenAdminA, created.Number)
		if view.Reminder || view.Overdue {
			t.Fatalf("non-PENDING => reminder=false overdue=false, got reminder=%v overdue=%v status=%s", view.Reminder, view.Overdue, view.Status)
		}
	})

	t.Run("reads_create_no_notification", func(t *testing.T) {
		created := mustCreate(time.Now().UTC().Add(72 * time.Hour))
		setDeadline(created.Number, time.Now().UTC().Add(-2*time.Hour)) // overdue

		before := loadRow(created.Number)
		notifBefore := notificationCount()
		var auditBefore int64
		fx.db.Model(&auditModel.AuditLog{}).Where("entity_id = ?", created.Number).Count(&auditBefore)

		// Detail read, list read, and inbox read — none may write.
		fetch(fx.tokenAdminA, created.Number)
		if status, body := h.do(fx.tokenAdminA, http.MethodGet, eApprovalsBase, nil); status != http.StatusOK {
			t.Fatalf("list = %d (%s)", status, string(body))
		}
		if status, body := h.do(fx.tokenAppr1, http.MethodGet, "/main/ems/approval-tasks", nil); status != http.StatusOK {
			t.Fatalf("tasks = %d (%s)", status, string(body))
		}

		after := loadRow(created.Number)
		if !after.UpdatedAt.Equal(before.UpdatedAt) || after.Status != before.Status || after.CurrentLevel != before.CurrentLevel {
			t.Fatal("computing overdue must not modify the request row")
		}
		if notificationCount() != notifBefore {
			t.Fatal("reads must not create durable notifications (Option B)")
		}
		var auditAfter int64
		fx.db.Model(&auditModel.AuditLog{}).Where("entity_id = ?", created.Number).Count(&auditAfter)
		if auditAfter != auditBefore {
			t.Fatal("reading an overdue request must not create an audit fact")
		}
	})

	t.Run("no_e4_notification_types", func(t *testing.T) {
		var count int64
		if err := fx.db.Model(&notificationModel.Notification{}).
			Where("type LIKE '%deadline%' OR type LIKE '%overdue%' OR type LIKE '%reminder%'").
			Count(&count).Error; err != nil {
			t.Fatalf("counting E4 notifications: %v", err)
		}
		if count != 0 {
			t.Fatalf("E4 reminder/overdue must not create durable notifications, found %d", count)
		}
	})

	t.Run("inbox_surfaces_overdue", func(t *testing.T) {
		created := mustCreate(time.Now().UTC().Add(72 * time.Hour))
		setDeadline(created.Number, time.Now().UTC().Add(-time.Hour))

		status, body := h.do(fx.tokenAppr1, http.MethodGet, "/main/ems/approval-tasks", nil)
		if status != http.StatusOK {
			t.Fatalf("tasks = %d (%s)", status, string(body))
		}
		var tasks struct {
			Data []struct {
				RequestNumber string `json:"request_number"`
				Reminder      bool   `json:"reminder"`
				Overdue       bool   `json:"overdue"`
			} `json:"data"`
		}
		if err := json.Unmarshal(body, &tasks); err != nil {
			t.Fatalf("decoding tasks: %v", err)
		}
		found := false
		for _, task := range tasks.Data {
			if task.RequestNumber == created.Number {
				found = true
				if !task.Overdue || task.Reminder {
					t.Fatalf("inbox task should be overdue and not reminder: overdue=%v reminder=%v", task.Overdue, task.Reminder)
				}
			}
		}
		if !found {
			t.Fatal("the eligible approver's inbox must still contain the overdue task")
		}
	})

	t.Run("no_persisted_overdue_column", func(t *testing.T) {
		var columns int64
		if err := fx.db.Raw(
			"SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND COLUMN_NAME LIKE '%overdue%'",
		).Scan(&columns).Error; err != nil {
			t.Fatalf("checking overdue columns: %v", err)
		}
		if columns != 0 {
			t.Fatalf("E4 is read-time only: no overdue column may exist, found %d", columns)
		}
	})

	t.Run("overdue_does_not_change_eligibility_or_grant_authority", func(t *testing.T) {
		created := mustCreate(time.Now().UTC().Add(72 * time.Hour))
		setDeadline(created.Number, time.Now().UTC().Add(-time.Hour))
		if view := fetch(fx.tokenAdminA, created.Number); !view.Overdue {
			t.Fatal("expected overdue")
		}

		// A user outside the snapshot still cannot act: overdue grants nothing.
		if status, _ := h.do(fx.tokenSolo, http.MethodPost, eApprovalsBase+"/"+created.Number+"/approve", map[string]any{"comment": "x"}); status != http.StatusForbidden {
			t.Fatalf("non-approver on an overdue request = %d, want 403", status)
		}
	})

	t.Run("organisation_isolation", func(t *testing.T) {
		created := mustCreate(time.Now().UTC().Add(72 * time.Hour))
		if status, _ := h.do(fx.tokenAdminB, http.MethodGet, eApprovalsBase+"/"+created.Number, nil); status != http.StatusNotFound {
			t.Fatalf("cross-organisation deadline read must be 404, got %d", status)
		}
	})

	t.Run("invalid_override_rejected", func(t *testing.T) {
		if status, body := create(time.Now().UTC().Add(-time.Hour)); status != http.StatusBadRequest {
			t.Fatalf("past deadline = %d, want 400 (%s)", status, string(body))
		}
		if status, body := h.do(fx.tokenAdminA, http.MethodPost, eApprovalsBase, map[string]any{
			"module": "ems", "resource_type": "e4", "resource_number": "bad",
			"decision_type": "purchase_request", "amount": 100, "deadline_at": "not-a-timestamp",
		}); status != http.StatusBadRequest {
			t.Fatalf("malformed deadline = %d, want 400 (%s)", status, string(body))
		}
	})

	t.Run("list_deadline_evaluation_adds_no_query_per_row", func(t *testing.T) {
		// Count SQL statements for the list read, then add more pending requests
		// and count again. The deadline/reminder/overdue computation must be
		// in-memory: the query count must not grow with the number of approvals.
		counter := &e4CountingLogger{}
		original := config.DB.Logger
		config.DB.Logger = counter
		defer func() { config.DB.Logger = original }()

		list := func() int64 {
			atomic.StoreInt64(&counter.queries, 0)
			if _, err := workflowServices.ListApprovalRequests(eOrgA, "", "", "", "", 50, 0); err != nil {
				t.Fatalf("list: %v", err)
			}
			return atomic.LoadInt64(&counter.queries)
		}

		first := list()
		for i := 0; i < 3; i++ {
			mustCreate(time.Now().UTC().Add(72 * time.Hour))
		}
		second := list()
		if second != first {
			t.Fatalf("list query count grew with the number of approvals (%d -> %d): deadline evaluation must not be N+1", first, second)
		}
	})
}
