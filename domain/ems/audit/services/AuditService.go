package services

import (
	"encoding/json"
	"errors"
	"log"
	"strings"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/audit/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/audit/models"
	auditRepositories "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/audit/repositories"
	eosRepositories "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/eos/repositories"
)

var (
	ErrNoActiveOrganisation = errors.New("no active organisation is set for this user")
	ErrAuditLogNotFound     = errors.New("audit record not found")
	ErrInvalidAuditRecord   = errors.New("audit action and entity type are required")
)

// appendAuditLog is the persistence call used by Record. It is a variable so a
// unit test can simulate an audit-store failure and prove the non-fatal
// contract (Foundation G §3.5/§7).
var appendAuditLog = auditRepositories.CreateAuditLog

// AuditRecord is the input for appending an audit fact.
type AuditRecord struct {
	OrganisationNumber     string
	ActorUserNumber        string
	Action                 string
	EntityType             string
	EntityID               string
	OldValues              interface{}
	NewValues              interface{}
	IP                     string
	UserAgent              string
	CorrelationID          string
	WorkflowInstanceNumber string
}

// Record appends an immutable audit fact. It is best-effort: callers must treat
// an error as non-fatal so auditing never rolls back or alters an authoritative
// business action (Foundation G §3.5). G1 does not yet call this from
// Foundations A–E; that integration is G2.
func Record(record AuditRecord) error {
	action := strings.TrimSpace(record.Action)
	entityType := strings.TrimSpace(record.EntityType)
	if action == "" || entityType == "" {
		return ErrInvalidAuditRecord
	}

	log := models.AuditLog{
		OrganisationNumber: strings.TrimSpace(record.OrganisationNumber),
		ActorUserNumber:    strings.TrimSpace(record.ActorUserNumber),
		Action:             action,
		EntityType:         entityType,
		EntityID:           strings.TrimSpace(record.EntityID),
		OldValues:          toJSON(record.OldValues),
		NewValues:          toJSON(record.NewValues),
		IP:                 record.IP,
		UserAgent:          record.UserAgent,
		CorrelationID:      record.CorrelationID,
	}
	if w := strings.TrimSpace(record.WorkflowInstanceNumber); w != "" {
		log.WorkflowInstanceNumber = &w
	}
	return appendAuditLog(&log)
}

// RecordFact appends an audit fact best-effort. A failure is logged (without
// secrets) and swallowed so that auditing never changes the business result
// (Foundation G §3.5/§7). Callers must not treat an audit failure as a business
// failure, and must not retry, enqueue or background the write.
func RecordFact(record AuditRecord) {
	if err := Record(record); err != nil {
		log.Printf("audit: failed to record %q for %s: %v", record.Action, record.EntityType, err)
	}
}

func toJSON(v interface{}) string {
	if v == nil {
		return ""
	}
	data, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(data)
}

// ListAuditLogs returns audit facts for the control plane (empty
// scopedOrganisation, global filters) or for a single organisation
// (scopedOrganisation non-empty). The scoped organisation always overrides any
// caller-supplied organisation so client-plane requests can never widen scope.
func ListAuditLogs(scopedOrganisation, organisation, actor, action, entityType, entityID string, limit, offset int) ([]responses.AuditLogResponse, int64, error) {
	if scopedOrganisation != "" {
		organisation = scopedOrganisation
	}
	filter := models.AuditLogFilter{
		OrganisationNumber: strings.TrimSpace(organisation),
		ActorUserNumber:    strings.TrimSpace(actor),
		Action:             strings.TrimSpace(action),
		EntityType:         strings.TrimSpace(entityType),
		EntityID:           strings.TrimSpace(entityID),
		Limit:              limit,
		Offset:             offset,
	}

	logs, total, err := auditRepositories.FindAuditLogs(filter)
	if err != nil {
		return nil, 0, err
	}

	result := make([]responses.AuditLogResponse, 0, len(logs))
	for _, log := range logs {
		result = append(result, toResponse(log))
	}
	return result, total, nil
}

// GetAuditLog returns one audit fact. When scopedOrganisation is non-empty the
// record must belong to that organisation, otherwise it is reported as not
// found so a client can never probe another organisation's records.
func GetAuditLog(scopedOrganisation string, id uint) (*responses.AuditLogResponse, error) {
	log, err := auditRepositories.FindAuditLogByID(id)
	if err != nil {
		return nil, err
	}
	if log == nil {
		return nil, ErrAuditLogNotFound
	}
	if scopedOrganisation != "" && log.OrganisationNumber != scopedOrganisation {
		return nil, ErrAuditLogNotFound
	}
	response := toResponse(*log)
	return &response, nil
}

// ListOrganisationAuditLogs resolves the authenticated user's organisation and
// lists its audit records (client plane).
func ListOrganisationAuditLogs(userNumber, actor, action, entityType, entityID string, limit, offset int) ([]responses.AuditLogResponse, int64, error) {
	organisationNumber, err := currentOrganisationNumber(userNumber)
	if err != nil {
		return nil, 0, err
	}
	return ListAuditLogs(organisationNumber, "", actor, action, entityType, entityID, limit, offset)
}

// GetOrganisationAuditLog resolves the authenticated user's organisation and
// returns one of its audit records.
func GetOrganisationAuditLog(userNumber string, id uint) (*responses.AuditLogResponse, error) {
	organisationNumber, err := currentOrganisationNumber(userNumber)
	if err != nil {
		return nil, err
	}
	return GetAuditLog(organisationNumber, id)
}

func currentOrganisationNumber(userNumber string) (string, error) {
	active, err := eosRepositories.FindActiveUserEntityByUserNumber(userNumber)
	if err != nil {
		return "", err
	}
	if active.UserNumber == "" || active.EntityNumber == "" {
		return "", ErrNoActiveOrganisation
	}
	return active.EntityNumber, nil
}

func toResponse(log models.AuditLog) responses.AuditLogResponse {
	response := responses.AuditLogResponse{
		Number:             log.ID,
		OrganisationNumber: log.OrganisationNumber,
		ActorUserNumber:    log.ActorUserNumber,
		Action:             log.Action,
		EntityType:         log.EntityType,
		EntityID:           log.EntityID,
		IP:                 log.IP,
		UserAgent:          log.UserAgent,
		CorrelationID:      log.CorrelationID,
		CreatedAt:          log.CreatedAt,
	}
	if log.OldValues != "" {
		response.OldValues = json.RawMessage(log.OldValues)
	}
	if log.NewValues != "" {
		response.NewValues = json.RawMessage(log.NewValues)
	}
	if log.WorkflowInstanceNumber != nil {
		response.WorkflowInstanceNumber = *log.WorkflowInstanceNumber
	}
	return response
}
