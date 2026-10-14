package integration

import (
	"encoding/json"
	"net/http"
	"testing"

	auditModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/audit/models"
	licensingModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/licensing/models"
	licensingServices "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/licensing/services"
)

// TestFoundationG5LicensingNotDuplicatedIntoAudit verifies the Foundation G
// safety requirement that `ems_licence_events` remains the sole licensing fact
// source: recording a licensing fact must not create a row in `ems_audit_logs`
// (no second audit engine, no duplicated licensing history). It also confirms
// the fact still flows through the G4 decision layer to a delivered
// notification.
func TestFoundationG5LicensingNotDuplicatedIntoAudit(t *testing.T) {
	db := setupG4LicensingFixture(t)

	var auditBefore int64
	db.Model(&auditModel.AuditLog{}).Where("organisation_number = ?", g4LicOrg).Count(&auditBefore)

	if err := licensingServices.RecordLicenceIssued(g4Licence(g4LicenceNumber, g4LicOrg), g4LicControl); err != nil {
		t.Fatalf("recording licensing fact: %v", err)
	}

	var events int64
	db.Model(&licensingModel.LicenceEvent{}).Where("organisation_number = ?", g4LicOrg).Count(&events)
	if events == 0 {
		t.Fatal("the authoritative licensing fact was not recorded in ems_licence_events")
	}

	var auditAfter int64
	db.Model(&auditModel.AuditLog{}).Where("organisation_number = ?", g4LicOrg).Count(&auditAfter)
	if auditAfter != auditBefore {
		t.Fatalf("a licensing fact was duplicated into ems_audit_logs (%d -> %d)", auditBefore, auditAfter)
	}

	// The fact must still reach the ratified G4 recipient (decision + delivery).
	if !g4HasNotification(db, g4LicOwner, "LICENCE_ISSUED", false) {
		t.Fatal("the licensing fact did not flow through the G4 decision layer")
	}
}

// TestFoundationG5FactDecisionDeliveryLayers verifies that the three Foundation
// G concerns remain distinct for one business action: the authoritative audit
// fact is stored in `ems_audit_logs` (FACT) and the notification is stored in
// `ems_notifications` (DECISION + DELIVERY). The layers are not merged.
func TestFoundationG5FactDecisionDeliveryLayers(t *testing.T) {
	fx := setupEFixture(t)
	organisations := []string{eOrgA, eOrgB}
	g4CleanNotifications(t, fx.db, organisations)
	t.Cleanup(func() { g4CleanNotifications(t, fx.db, organisations) })

	h := newG4Harness(t)

	status, body := h.do(fx.tokenAdminA, http.MethodPost, eApprovalsBase, map[string]any{
		"module": "ems", "resource_type": "g5", "resource_number": "layer-1",
		"decision_type": "purchase_request", "amount": 100,
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

	// FACT layer: the G2 audit hook wrote the authoritative fact.
	var auditCount int64
	fx.db.Model(&auditModel.AuditLog{}).
		Where("action = ? AND organisation_number = ? AND entity_id = ?", "approval.create", eOrgA, envelope.Data.Number).
		Count(&auditCount)
	if auditCount == 0 {
		t.Fatal("the authoritative audit fact was not recorded (FACT layer)")
	}

	// DECISION + DELIVERY layer: the G4 decision produced a recipient-scoped
	// notification in the separate notifications table.
	if !g4HasNotification(fx.db, eApprover1, "approval.create", true) {
		t.Fatal("the notification was not delivered (DECISION/DELIVERY layer)")
	}
}
