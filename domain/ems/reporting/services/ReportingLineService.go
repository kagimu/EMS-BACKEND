package services

import (
	"errors"
	"strings"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
	assignmentRepositories "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/assignments/repositories"
	organisationServices "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/organisation/services"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/reporting/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/reporting/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/reporting/models"
	reportingRepositories "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/reporting/repositories"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

const ReportingStatusEnded = "ENDED"

// Effective-state values exposed to clients. They are derived, never stored.
const (
	EffectiveStateEffective          = "EFFECTIVE"
	EffectiveStateUpcoming           = "UPCOMING"
	EffectiveStateExpired            = "EXPIRED"
	EffectiveStateInactive           = "INACTIVE"
	EffectiveStateAssignmentInactive = "ASSIGNMENT_INACTIVE"
)

var (
	ErrReportingLineNotFound        = errors.New("reporting line not found")
	ErrSubordinateRequired          = errors.New("subordinate_assignment_number is required")
	ErrManagerRequired              = errors.New("manager_assignment_number is required")
	ErrSelfReporting                = errors.New("an assignment cannot report to itself")
	ErrReportingCycle               = errors.New("this reporting line would create a cycle")
	ErrAssignmentNotFound           = errors.New("assignment was not found in this organisation")
	ErrAssignmentNotActive          = errors.New("both assignments must be active for an effective reporting line")
	ErrDuplicateActiveReportingLine = errors.New("this assignment already has an effective reporting line for the requested period")
	ErrInvalidStatus                = errors.New("invalid status")
	ErrInvalidDateRange             = errors.New("valid_from cannot be after valid_to")
	ErrReactivationRequiresDates    = errors.New("this reporting line has expired; update the validity dates before reactivating it")
)

// ErrReportingChangeBusy is re-exported so controllers can map the lock timeout.
var ErrReportingChangeBusy = reportingRepositories.ErrReportingChangeBusy

func isValidReportingStatus(status string) bool {
	return status == constants.StatusActive ||
		status == constants.StatusInactive ||
		status == ReportingStatusEnded
}

func isValidDateRange(from, to *time.Time) bool {
	if from == nil || to == nil {
		return true
	}
	return !from.After(*to)
}

// effectiveState is the single derived definition of whether a reporting line is
// currently effective:
//
//   - line status must be ACTIVE;
//   - both assignments must currently be ACTIVE;
//   - valid_from (if set) must not be in the future;
//   - valid_to (if set) must not be in the past.
//
// The organisation is already enforced by the query. It is never stored.
func effectiveState(line models.ReportingLine, subordinateStatus, managerStatus string, now time.Time) (bool, string) {
	if line.Status != constants.StatusActive {
		return false, EffectiveStateInactive
	}
	if subordinateStatus != constants.StatusActive || managerStatus != constants.StatusActive {
		return false, EffectiveStateAssignmentInactive
	}
	if line.ValidFrom != nil && line.ValidFrom.After(now) {
		return false, EffectiveStateUpcoming
	}
	if line.ValidTo != nil && line.ValidTo.Before(now) {
		return false, EffectiveStateExpired
	}
	return true, EffectiveStateEffective
}

// windowsOverlap reports whether an existing line's validity window overlaps the
// candidate window. An existing nil valid_from is unbounded in the past and a nil
// valid_to is unbounded in the future. A candidate nil valid_from starts now (a
// new relationship has no historical start), and a nil valid_to is unbounded.
// Boundaries are exclusive: a line ending exactly when another starts does not
// overlap.
func windowsOverlap(lineFrom, lineTo, candidateFrom, candidateTo *time.Time, now time.Time) bool {
	lineStart := time.Time{}
	if lineFrom != nil {
		lineStart = *lineFrom
	}
	candidateStart := now
	if candidateFrom != nil {
		candidateStart = *candidateFrom
	}

	if lineTo != nil && !lineTo.After(candidateStart) {
		return false
	}
	if candidateTo != nil && !candidateTo.After(lineStart) {
		return false
	}
	return true
}

// ListReportingLines returns the reporting lines of one organisation, optionally
// filtered by subordinate assignment, manager assignment and/or status. Each row
// carries the derived effective/effective_state fields.
func ListReportingLines(organisationNumber, subordinateAssignment, managerAssignment, status string) ([]responses.ReportingLineResponse, error) {
	if status != "" && !isValidReportingStatus(status) {
		return nil, ErrInvalidStatus
	}

	lines, err := reportingRepositories.FindAllByOrganisation(organisationNumber, subordinateAssignment, managerAssignment, status)
	if err != nil {
		return nil, err
	}

	assignments, err := assignmentRepositories.FindAllByOrganisation(organisationNumber, "", "", "")
	if err != nil {
		return nil, err
	}
	statusByAssignment := make(map[string]string, len(assignments))
	for _, assignment := range assignments {
		statusByAssignment[assignment.Number] = assignment.Status
	}

	now := time.Now()
	result := make([]responses.ReportingLineResponse, 0, len(lines))
	for _, line := range lines {
		result = append(result, toResponse(
			line,
			statusByAssignment[line.SubordinateAssignmentNumber],
			statusByAssignment[line.ManagerAssignmentNumber],
			now,
		))
	}
	return result, nil
}

// GetReportingLine returns one reporting line within the organisation.
func GetReportingLine(organisationNumber, number string) (*responses.ReportingLineResponse, error) {
	line, err := reportingRepositories.FindByNumberAndOrganisation(number, organisationNumber)
	if err != nil {
		return nil, ErrReportingLineNotFound
	}

	subordinateStatus, managerStatus := assignmentStatuses(organisationNumber, line.SubordinateAssignmentNumber, line.ManagerAssignmentNumber)
	response := toResponse(*line, subordinateStatus, managerStatus, time.Now())
	return &response, nil
}

// CreateReportingLine creates a reporting line between two existing ACTIVE
// assignments of the same organisation.
func CreateReportingLine(actorNumber, organisationNumber string, request requests.CreateReportingLineRequest) (*responses.ReportingLineResponse, error) {
	subordinate := strings.TrimSpace(request.SubordinateAssignmentNumber)
	manager := strings.TrimSpace(request.ManagerAssignmentNumber)

	if subordinate == "" {
		return nil, ErrSubordinateRequired
	}
	if manager == "" {
		return nil, ErrManagerRequired
	}
	if subordinate == manager {
		return nil, ErrSelfReporting
	}
	if !isValidDateRange(request.ValidFrom, request.ValidTo) {
		return nil, ErrInvalidDateRange
	}

	subordinateAssignment, err := assignmentRepositories.FindByNumberAndOrganisation(subordinate, organisationNumber)
	if err != nil {
		return nil, ErrAssignmentNotFound
	}
	managerAssignment, err := assignmentRepositories.FindByNumberAndOrganisation(manager, organisationNumber)
	if err != nil {
		return nil, ErrAssignmentNotFound
	}
	if subordinateAssignment.Status != constants.StatusActive ||
		managerAssignment.Status != constants.StatusActive {
		return nil, ErrAssignmentNotActive
	}

	var created models.ReportingLine
	if err := reportingRepositories.WithReportingTransaction(organisationNumber, func(tx *gorm.DB) error {
		existing, err := overlappingActiveLineForSubordinate(organisationNumber, subordinate, request.ValidFrom, request.ValidTo, "", tx)
		if err != nil {
			return err
		}
		if existing != nil {
			return ErrDuplicateActiveReportingLine
		}

		cycle, err := wouldCreateCycle(organisationNumber, subordinate, manager, request.ValidFrom, request.ValidTo, "", tx)
		if err != nil {
			return err
		}
		if cycle {
			return ErrReportingCycle
		}

		created = models.ReportingLine{
			Number:                      uuid.NewV4().String(),
			OrganisationNumber:          organisationNumber,
			SubordinateAssignmentNumber: subordinate,
			ManagerAssignmentNumber:     manager,
			ValidFrom:                   request.ValidFrom,
			ValidTo:                     request.ValidTo,
			Status:                      constants.StatusActive,
			AddedByNumber:               actorNumber,
		}
		return reportingRepositories.Create(&created, tx)
	}); err != nil {
		return nil, err
	}

	response := toResponse(created, subordinateAssignment.Status, managerAssignment.Status, time.Now())
	return &response, nil
}

// UpdateReportingLine updates the manager and lifecycle fields. The subordinate
// is immutable so history is preserved.
func UpdateReportingLine(organisationNumber, number string, request requests.UpdateReportingLineRequest) (*responses.ReportingLineResponse, error) {
	line, err := reportingRepositories.FindByNumberAndOrganisation(number, organisationNumber)
	if err != nil {
		return nil, ErrReportingLineNotFound
	}

	manager := strings.TrimSpace(request.ManagerAssignmentNumber)
	if manager == "" {
		return nil, ErrManagerRequired
	}
	if !isValidDateRange(request.ValidFrom, request.ValidTo) {
		return nil, ErrInvalidDateRange
	}
	if manager == line.SubordinateAssignmentNumber {
		return nil, ErrSelfReporting
	}

	subordinateAssignment, err := assignmentRepositories.FindByNumberAndOrganisation(line.SubordinateAssignmentNumber, organisationNumber)
	if err != nil {
		return nil, ErrAssignmentNotFound
	}
	managerAssignment, err := assignmentRepositories.FindByNumberAndOrganisation(manager, organisationNumber)
	if err != nil {
		return nil, ErrAssignmentNotFound
	}
	if managerAssignment.Status != constants.StatusActive {
		return nil, ErrAssignmentNotActive
	}
	if line.Status == constants.StatusActive && subordinateAssignment.Status != constants.StatusActive {
		return nil, ErrAssignmentNotActive
	}

	var saved models.ReportingLine
	if err := reportingRepositories.WithReportingTransaction(organisationNumber, func(tx *gorm.DB) error {
		current, err := reportingRepositories.FindByNumberAndOrganisation(number, organisationNumber, tx)
		if err != nil {
			return ErrReportingLineNotFound
		}

		if current.Status == constants.StatusActive {
			existing, err := overlappingActiveLineForSubordinate(
				organisationNumber, current.SubordinateAssignmentNumber,
				request.ValidFrom, request.ValidTo, current.Number, tx,
			)
			if err != nil {
				return err
			}
			if existing != nil {
				return ErrDuplicateActiveReportingLine
			}

			cycle, err := wouldCreateCycle(
				organisationNumber, current.SubordinateAssignmentNumber, manager,
				request.ValidFrom, request.ValidTo, current.Number, tx,
			)
			if err != nil {
				return err
			}
			if cycle {
				return ErrReportingCycle
			}
		}

		current.ManagerAssignmentNumber = manager
		current.ValidFrom = request.ValidFrom
		current.ValidTo = request.ValidTo

		if err := reportingRepositories.Save(current, tx); err != nil {
			return err
		}
		saved = *current
		return nil
	}); err != nil {
		return nil, err
	}

	response := toResponse(saved, subordinateAssignment.Status, managerAssignment.Status, time.Now())
	return &response, nil
}

// SetReportingLineStatus applies the ACTIVE / INACTIVE / ENDED lifecycle.
//
// Reactivation is refused when the existing validity window has already expired:
// an administrator must explicitly update the dates first, so an expired
// relationship is never silently revived with a historical window. Reactivation
// also re-checks assignment status, the single-effective-manager rule and cycles.
func SetReportingLineStatus(organisationNumber, number, status string) (*responses.ReportingLineResponse, error) {
	if !isValidReportingStatus(status) {
		return nil, ErrInvalidStatus
	}

	line, err := reportingRepositories.FindByNumberAndOrganisation(number, organisationNumber)
	if err != nil {
		return nil, ErrReportingLineNotFound
	}

	subordinateStatus, managerStatus := assignmentStatuses(organisationNumber, line.SubordinateAssignmentNumber, line.ManagerAssignmentNumber)

	if status == constants.StatusActive {
		if line.ValidTo != nil && line.ValidTo.Before(time.Now()) {
			return nil, ErrReactivationRequiresDates
		}
		if subordinateStatus != constants.StatusActive || managerStatus != constants.StatusActive {
			return nil, ErrAssignmentNotActive
		}
	}

	var saved models.ReportingLine
	if err := reportingRepositories.WithReportingTransaction(organisationNumber, func(tx *gorm.DB) error {
		current, err := reportingRepositories.FindByNumberAndOrganisation(number, organisationNumber, tx)
		if err != nil {
			return ErrReportingLineNotFound
		}

		if status == constants.StatusActive {
			if current.ValidTo != nil && current.ValidTo.Before(time.Now()) {
				return ErrReactivationRequiresDates
			}

			existing, err := overlappingActiveLineForSubordinate(
				organisationNumber, current.SubordinateAssignmentNumber,
				current.ValidFrom, current.ValidTo, current.Number, tx,
			)
			if err != nil {
				return err
			}
			if existing != nil {
				return ErrDuplicateActiveReportingLine
			}

			cycle, err := wouldCreateCycle(
				organisationNumber, current.SubordinateAssignmentNumber, current.ManagerAssignmentNumber,
				current.ValidFrom, current.ValidTo, current.Number, tx,
			)
			if err != nil {
				return err
			}
			if cycle {
				return ErrReportingCycle
			}
		}

		current.Status = status
		if status == ReportingStatusEnded && current.ValidTo == nil {
			now := time.Now()
			current.ValidTo = &now
		}

		if err := reportingRepositories.Save(current, tx); err != nil {
			return err
		}
		saved = *current
		return nil
	}); err != nil {
		return nil, err
	}

	response := toResponse(saved, subordinateStatus, managerStatus, time.Now())
	return &response, nil
}

// overlappingActiveLineForSubordinate returns another ACTIVE line for the same
// subordinate whose validity window overlaps the candidate window and whose two
// assignments are currently ACTIVE. Expired/future/non-overlapping lines and
// lines whose assignments are no longer active are ignored.
func overlappingActiveLineForSubordinate(organisationNumber, subordinate string, from, to *time.Time, excludeNumber string, tx *gorm.DB) (*models.ReportingLine, error) {
	lines, err := reportingRepositories.FindAllActive(organisationNumber, tx)
	if err != nil {
		return nil, err
	}
	activeAssignments, err := assignmentRepositories.FindAllByOrganisation(organisationNumber, "", "", constants.StatusActive)
	if err != nil {
		return nil, err
	}
	activeSet := make(map[string]bool, len(activeAssignments))
	for _, assignment := range activeAssignments {
		activeSet[assignment.Number] = true
	}

	now := time.Now()
	for i := range lines {
		line := lines[i]
		if line.Number == excludeNumber || line.SubordinateAssignmentNumber != subordinate {
			continue
		}
		if !activeSet[line.SubordinateAssignmentNumber] || !activeSet[line.ManagerAssignmentNumber] {
			continue
		}
		if windowsOverlap(line.ValidFrom, line.ValidTo, from, to, now) {
			return &line, nil
		}
	}
	return nil, nil
}

// wouldCreateCycle walks the overlapping ACTIVE manager chain from manager to the
// root. Historical (expired), future-only, non-overlapping and assignment-inactive
// lines are excluded, so a completed historical relationship never blocks a new
// valid structure. A pre-existing cycle in the active graph fails closed.
func wouldCreateCycle(organisationNumber, subordinate, manager string, from, to *time.Time, excludeNumber string, tx *gorm.DB) (bool, error) {
	lines, err := reportingRepositories.FindAllActive(organisationNumber, tx)
	if err != nil {
		return false, err
	}
	activeAssignments, err := assignmentRepositories.FindAllByOrganisation(organisationNumber, "", "", constants.StatusActive)
	if err != nil {
		return false, err
	}
	activeSet := make(map[string]bool, len(activeAssignments))
	for _, assignment := range activeAssignments {
		activeSet[assignment.Number] = true
	}

	now := time.Now()
	managerBySubordinate := make(map[string]string)
	for i := range lines {
		line := lines[i]
		if line.Number == excludeNumber {
			continue
		}
		if !activeSet[line.SubordinateAssignmentNumber] || !activeSet[line.ManagerAssignmentNumber] {
			continue
		}
		if !windowsOverlap(line.ValidFrom, line.ValidTo, from, to, now) {
			continue
		}
		managerBySubordinate[line.SubordinateAssignmentNumber] = line.ManagerAssignmentNumber
	}

	current := manager
	visited := map[string]bool{}
	for current != "" {
		if current == subordinate {
			return true, nil
		}
		if visited[current] {
			return true, nil
		}
		visited[current] = true
		current = managerBySubordinate[current]
	}

	return false, nil
}

func assignmentStatuses(organisationNumber, subordinateAssignment, managerAssignment string) (string, string) {
	var statuses [2]string
	for index, number := range []string{subordinateAssignment, managerAssignment} {
		assignment, err := assignmentRepositories.FindByNumberAndOrganisation(number, organisationNumber)
		if err == nil && assignment != nil {
			statuses[index] = assignment.Status
		}
	}
	return statuses[0], statuses[1]
}

// CurrentOrganisationNumber is re-exported so callers depend on one resolver.
var CurrentOrganisationNumber = organisationServices.CurrentOrganisationNumber

func toResponse(line models.ReportingLine, subordinateStatus, managerStatus string, now time.Time) responses.ReportingLineResponse {
	effective, effectiveState := effectiveState(line, subordinateStatus, managerStatus, now)
	return responses.ReportingLineResponse{
		Number:                      line.Number,
		OrganisationNumber:          line.OrganisationNumber,
		SubordinateAssignmentNumber: line.SubordinateAssignmentNumber,
		ManagerAssignmentNumber:     line.ManagerAssignmentNumber,
		ValidFrom:                   line.ValidFrom,
		ValidTo:                     line.ValidTo,
		Status:                      line.Status,
		Effective:                   effective,
		EffectiveState:              effectiveState,
		AddedByNumber:               line.AddedByNumber,
		CreatedAt:                   line.CreatedAt,
		UpdatedAt:                   line.UpdatedAt,
	}
}
