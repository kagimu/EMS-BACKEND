package services

import (
	"testing"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/workflow/models"
)

// TestE4DerivedDeadlineConditions pins the ratified C4 reminder/overdue boundary
// semantics with a deterministic evaluation instant (no wall-clock dependency).
func TestE4DerivedDeadlineConditions(t *testing.T) {
	now := time.Date(2030, 1, 1, 12, 0, 0, 0, time.UTC)
	offset := func(d time.Duration) *time.Time {
		value := now.Add(d)
		return &value
	}

	cases := []struct {
		name         string
		status       string
		deadline     *time.Time
		createdAt    time.Time
		wantReminder bool
		wantOverdue  bool
	}{
		{"pending_more_than_24h", models.RequestStatusPending, offset(25 * time.Hour), now, false, false},
		{"pending_exactly_24h", models.RequestStatusPending, offset(24 * time.Hour), now, true, false},
		{"pending_less_than_24h", models.RequestStatusPending, offset(time.Hour), now, true, false},
		{"pending_exactly_at_deadline", models.RequestStatusPending, offset(0), now, true, false},
		{"pending_just_after_deadline", models.RequestStatusPending, offset(-time.Second), now, false, true},
		{"pending_significantly_overdue", models.RequestStatusPending, offset(-72 * time.Hour), now, false, true},
		{"non_pending_within_24h", models.RequestStatusApproved, offset(time.Hour), now, false, false},
		{"non_pending_overdue", models.RequestStatusRejected, offset(-72 * time.Hour), now, false, false},
		// No stored deadline: E4-1/E4-2 apply the configured default to CreatedAt,
		// so a request created just now has the full default window ahead of it.
		{"no_stored_deadline_uses_default", models.RequestStatusPending, nil, now, false, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			approval := models.ApprovalRequest{
				Status:     tc.status,
				DeadlineAt: tc.deadline,
			}
			approval.CreatedAt = tc.createdAt
			if got := requestIsApproachingDeadline(approval, now); got != tc.wantReminder {
				t.Fatalf("reminder = %v, want %v", got, tc.wantReminder)
			}
			if got := requestIsOverdue(approval, now); got != tc.wantOverdue {
				t.Fatalf("overdue = %v, want %v", got, tc.wantOverdue)
			}
		})
	}
}
