package services

import (
	"errors"
	"testing"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/audit/models"
	"gorm.io/gorm"
)

// TestRecordRejectsIncompleteFacts proves an incomplete fact is rejected by the
// strict recording path.
func TestRecordRejectsIncompleteFacts(t *testing.T) {
	if err := Record(AuditRecord{}); !errors.Is(err, ErrInvalidAuditRecord) {
		t.Fatalf("Record(empty) = %v, want ErrInvalidAuditRecord", err)
	}
}

// TestRecordFactIsNonFatal proves that a failed audit write is swallowed: the
// non-fatal contract means a business action is never changed by an audit
// failure (Foundation G §3.5/§7).
func TestRecordFactIsNonFatal(t *testing.T) {
	original := appendAuditLog
	appendAuditLog = func(*models.AuditLog, ...*gorm.DB) error {
		return errors.New("simulated audit store failure")
	}
	defer func() { appendAuditLog = original }()

	// Must not panic and must not surface the failure.
	RecordFact(AuditRecord{
		OrganisationNumber: "org-1",
		ActorUserNumber:    "user-1",
		Action:             "stakeholder.update",
		EntityType:         "ems_stakeholder",
		EntityID:           "sta-1",
	})

	// The strict path still reports the error to callers that want it.
	if err := Record(AuditRecord{
		Action:     "stakeholder.update",
		EntityType: "ems_stakeholder",
	}); err == nil {
		t.Fatal("Record must report persistence errors")
	}
}
