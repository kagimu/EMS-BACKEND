package services

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/licensing/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/licensing/models"
	licenceRepositories "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/licensing/repositories"
	notificationDecision "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/notifications/decision"
)

// Foundation F records durable licensing facts for Foundation G to reference
// later (§17/§18). These helpers are append-only and synchronous; they never
// start background work and never deliver notifications. A recording failure
// is returned to the caller but is treated as non-fatal by callers so that the
// primary licensing action is never rolled back because of the fact log.

// EventOption customises a licensing fact before it is persisted.
type EventOption func(*models.LicenceEvent)

// WithInstallation attaches installation context to a fact.
func WithInstallation(number, identity string) EventOption {
	return func(e *models.LicenceEvent) {
		n := number
		i := identity
		e.InstallationNumber = &n
		e.InstallationIdentity = &i
	}
}

// WithReason attaches the required reason for the fact.
func WithReason(reason string) EventOption {
	return func(e *models.LicenceEvent) {
		r := reason
		e.Reason = &r
	}
}

// WithBeforeState attaches a JSON snapshot of the prior state.
func WithBeforeState(state interface{}) EventOption {
	return func(e *models.LicenceEvent) {
		e.BeforeState = mustMarshal(state)
	}
}

// WithAfterState attaches a JSON snapshot of the resulting state.
func WithAfterState(state interface{}) EventOption {
	return func(e *models.LicenceEvent) {
		e.AfterState = mustMarshal(state)
	}
}

// WithClientContext marks a fact as client-initiated and records client context.
func WithClientContext(clientIP, userAgent string) EventOption {
	return func(e *models.LicenceEvent) {
		e.ClientInitiated = true
		e.ClientIP = clientIP
		e.UserAgent = userAgent
	}
}

// recordLicenceFact appends a licence-level fact.
func recordLicenceFact(eventType models.LicenceEventType, licence *models.Licence, actorUserNumber string, opts ...EventOption) error {
	if licence == nil {
		return nil
	}
	event := models.LicenceEvent{
		EventType:          eventType,
		LicenceNumber:      licence.Number,
		OrganisationNumber: licence.OrganisationNumber,
		ActorUserNumber:    actorUserNumber,
	}
	for _, opt := range opts {
		opt(&event)
	}
	if err := licenceRepositories.CreateLicenceEvent(&event); err != nil {
		return err
	}

	// Foundation G (G4): best-effort notification decision, outside the
	// licensing flow. Never changes the licensing result (§19); failures are
	// swallowed inside the decision layer.
	notificationDecision.NotifyLicenceFact(string(event.EventType), event.OrganisationNumber, event.LicenceNumber, strconv.FormatUint(uint64(event.ID), 10))
	return nil
}

// recordInstallationFact appends an installation-level fact.
func recordInstallationFact(eventType models.LicenceEventType, installation *models.Installation, actorUserNumber string, opts ...EventOption) error {
	if installation == nil {
		return nil
	}
	event := models.LicenceEvent{
		EventType:          eventType,
		LicenceNumber:      installation.LicenceNumber,
		OrganisationNumber: installation.OrganisationNumber,
		ActorUserNumber:    actorUserNumber,
	}
	if installation.Number != "" {
		number := installation.Number
		identity := installation.InstallationIdentity
		event.InstallationNumber = &number
		event.InstallationIdentity = &identity
	}
	for _, opt := range opts {
		opt(&event)
	}
	if err := licenceRepositories.CreateLicenceEvent(&event); err != nil {
		return err
	}

	// Foundation G (G4): best-effort notification decision, outside the
	// licensing flow. Never changes the licensing result (§19).
	notificationDecision.NotifyLicenceFact(string(event.EventType), event.OrganisationNumber, event.LicenceNumber, strconv.FormatUint(uint64(event.ID), 10))
	return nil
}

// RecordLicenceIssued records licence issuance.
func RecordLicenceIssued(licence *models.Licence, actorUserNumber string) error {
	return recordLicenceFact(models.LicenceEventIssued, licence, actorUserNumber,
		WithAfterState(map[string]interface{}{
			"number":              licence.Number,
			"organisation_number": licence.OrganisationNumber,
			"valid_from":          licence.ValidFrom,
			"valid_to":            licence.ValidTo,
			"max_installations":   licence.MaxInstallations,
		}))
}

// RecordLicenceConfigured records a configuration change.
func RecordLicenceConfigured(licence *models.Licence, actorUserNumber string, before, after map[string]interface{}) error {
	return recordLicenceFact(models.LicenceEventConfigured, licence, actorUserNumber,
		WithBeforeState(before), WithAfterState(after))
}

// RecordLicenceRenewed records an ordinary renewal that preserved identity.
func RecordLicenceRenewed(licence *models.Licence, actorUserNumber string, previousValidTo time.Time) error {
	return recordLicenceFact(models.LicenceEventRenewed, licence, actorUserNumber,
		WithBeforeState(map[string]interface{}{"valid_to": previousValidTo}),
		WithAfterState(map[string]interface{}{"valid_to": licence.ValidTo}))
}

// RecordLicenceTransition records a status transition (suspend/revoke).
func RecordLicenceTransition(eventType models.LicenceEventType, licence *models.Licence, actorUserNumber string, previousStatus string) error {
	return recordLicenceFact(eventType, licence, actorUserNumber,
		WithBeforeState(map[string]interface{}{"status": previousStatus}),
		WithAfterState(map[string]interface{}{"status": licence.Status}))
}

// RecordLicenceReplaced records the supersession of one licence by another.
func RecordLicenceReplaced(oldLicence, newLicence *models.Licence, actorUserNumber, reason string) error {
	if oldLicence == nil || newLicence == nil {
		return nil
	}
	return recordLicenceFact(models.LicenceEventReplaced, oldLicence, actorUserNumber,
		WithReason(reason),
		WithBeforeState(map[string]interface{}{
			"number": oldLicence.Number,
			"status": oldLicence.Status,
		}),
		WithAfterState(map[string]interface{}{
			"number": newLicence.Number,
			"status": newLicence.Status,
		}))
}

// RecordInstallationActivated records a fresh installation activation.
func RecordInstallationActivated(installation *models.Installation, actorUserNumber string) error {
	return recordInstallationFact(models.InstallationEventActivated, installation, actorUserNumber,
		WithAfterState(map[string]interface{}{
			"number":   installation.Number,
			"licence":  installation.LicenceNumber,
			"status":   installation.Status,
			"identity": installation.InstallationIdentity,
		}))
}

// RecordInstallationReactivated records a capacity-neutral reactivation.
func RecordInstallationReactivated(installation *models.Installation, actorUserNumber string) error {
	return recordInstallationFact(models.InstallationEventReactivated, installation, actorUserNumber,
		WithAfterState(map[string]interface{}{
			"number":  installation.Number,
			"licence": installation.LicenceNumber,
			"status":  installation.Status,
		}))
}

// RecordInstallationDeactivated records a client self-deactivation.
func RecordInstallationDeactivated(installation *models.Installation, actorUserNumber string) error {
	return recordInstallationFact(models.InstallationEventDeactivated, installation, actorUserNumber,
		WithClientContext("", ""),
		WithAfterState(map[string]interface{}{
			"number":         installation.Number,
			"licence":        installation.LicenceNumber,
			"status":         installation.Status,
			"deactivated_at": installation.ClientDeactivatedAt,
		}))
}

// RecordInstallationReleased records an owner release of capacity.
func RecordInstallationReleased(installation *models.Installation, actorUserNumber, reason string) error {
	return recordInstallationFact(models.InstallationEventReleased, installation, actorUserNumber,
		WithReason(reason),
		WithAfterState(map[string]interface{}{
			"number":  installation.Number,
			"licence": installation.LicenceNumber,
			"status":  installation.Status,
		}))
}

// RecordClientRenewalRequested records a client renewal request fact.
func RecordClientRenewalRequested(licence *models.Licence, actorUserNumber string) error {
	return recordLicenceFact(models.ClientEventRenewalRequested, licence, actorUserNumber,
		WithClientContext("", ""),
		WithAfterState(map[string]interface{}{"renewal_requested_at": licence.RenewalRequestedAt}))
}

func mustMarshal(v interface{}) string {
	if v == nil {
		return ""
	}
	data, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(data)
}

// ListLicenceEvents returns licensing facts for the control plane, newest
// first, with an optional filter. It is read-only and deliberately global
// (control plane); the caller is guarded by the registry authority.
func ListLicenceEvents(licenceNumber, organisationNumber, eventType string, limit, offset int) ([]responses.LicenceEventResponse, int64, error) {
	filter := models.LicenceEventFilter{Limit: limit, Offset: offset}
	if v := strings.TrimSpace(licenceNumber); v != "" {
		filter.LicenceNumber = &v
	}
	if v := strings.TrimSpace(organisationNumber); v != "" {
		filter.OrganisationNumber = &v
	}
	if v := strings.TrimSpace(eventType); v != "" {
		t := models.LicenceEventType(v)
		filter.EventType = &t
	}

	events, total, err := licenceRepositories.FindLicenceEvents(filter)
	if err != nil {
		return nil, 0, err
	}

	result := make([]responses.LicenceEventResponse, 0, len(events))
	for _, event := range events {
		item := responses.LicenceEventResponse{
			Number:             event.ID,
			EventType:          string(event.EventType),
			LicenceNumber:      event.LicenceNumber,
			OrganisationNumber: event.OrganisationNumber,
			ActorUserNumber:    event.ActorUserNumber,
			ClientInitiated:    event.ClientInitiated,
			CreatedAt:          event.CreatedAt,
		}
		if event.InstallationNumber != nil {
			item.InstallationNumber = *event.InstallationNumber
		}
		if event.InstallationIdentity != nil {
			item.InstallationIdentity = *event.InstallationIdentity
		}
		if event.Reason != nil {
			item.Reason = *event.Reason
		}
		result = append(result, item)
	}
	return result, total, nil
}
