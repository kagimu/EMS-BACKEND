package services

import (
	"errors"
	"strings"
	"time"

	authRepository "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/repository"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/assignments/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/assignments/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/assignments/models"
	assignmentRepositories "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/assignments/repositories"
	organisationRepositories "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/organisation/repositories"
	organisationServices "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/organisation/services"
	uuid "github.com/satori/go.uuid"
)

// AssignmentStatusEnded terminates an assignment while retaining it for history.
const AssignmentStatusEnded = "ENDED"

var (
	ErrAssignmentNotFound        = errors.New("user assignment not found")
	ErrDuplicateActiveAssignment = errors.New("an active assignment of this type already exists for this person at this organisation unit")
	ErrInvalidAssignmentType     = errors.New("invalid assignment type")
	ErrInvalidStatus             = errors.New("invalid status")
	ErrInvalidDateRange          = errors.New("valid_from cannot be after valid_to")
	ErrUserRequired              = errors.New("user_number is required")
	ErrUserNotFound              = errors.New("referenced user was not found")
	ErrUnitRequired              = errors.New("organisation_unit_number is required")
	ErrUnitNotFound              = errors.New("organisation unit was not found in this organisation")
	ErrUnitNotActive             = errors.New("organisation unit is inactive")
)

func isValidAssignmentStatus(status string) bool {
	return status == constants.StatusActive ||
		status == constants.StatusInactive ||
		status == AssignmentStatusEnded
}

func isValidDateRange(from, to *time.Time) bool {
	if from == nil || to == nil {
		return true
	}
	return !from.After(*to)
}

// validateUnit enforces the organisation-unit rule: the unit must exist in the
// SAME organisation (a cross-organisation or soft-deleted unit is treated as
// absent) and must be ACTIVE before it can receive a new assignment.
func validateUnit(organisationNumber, unitNumber string) error {
	if unitNumber == "" {
		return ErrUnitRequired
	}

	unit, err := organisationRepositories.FindByNumberAndOrganisation(unitNumber, organisationNumber)
	if err != nil {
		return ErrUnitNotFound
	}
	if unit.Status != constants.StatusActive {
		return ErrUnitNotActive
	}

	return nil
}

// ListAssignments returns the assignments of an organisation, optionally
// filtered by unit, user and/or status. A single query is used (no N+1).
func ListAssignments(organisationNumber, unitNumber, userNumber, status string) ([]responses.UserAssignmentResponse, error) {
	if status != "" && !isValidAssignmentStatus(status) {
		return nil, ErrInvalidStatus
	}

	assignments, err := assignmentRepositories.FindAllByOrganisation(organisationNumber, unitNumber, userNumber, status)
	if err != nil {
		return nil, err
	}

	result := make([]responses.UserAssignmentResponse, 0, len(assignments))
	for _, assignment := range assignments {
		result = append(result, toResponse(assignment))
	}
	return result, nil
}

// GetAssignment returns one assignment within the organisation.
func GetAssignment(organisationNumber, number string) (*responses.UserAssignmentResponse, error) {
	assignment, err := assignmentRepositories.FindByNumberAndOrganisation(number, organisationNumber)
	if err != nil {
		return nil, ErrAssignmentNotFound
	}
	response := toResponse(*assignment)
	return &response, nil
}

// CreateAssignment creates an assignment for an existing person at an existing
// ACTIVE unit of the same organisation.
//
// It creates no user, unit, role or authority, and it grants no access: EMS
// operation still requires EntityUser membership plus an active organisation
// context. Stakeholder status is never consulted — stakeholders and assignments
// are independent concepts.
func CreateAssignment(actorNumber, organisationNumber string, request requests.CreateUserAssignmentRequest) (*responses.UserAssignmentResponse, error) {
	userNumber := strings.TrimSpace(request.UserNumber)
	unitNumber := strings.TrimSpace(request.OrganisationUnitNumber)
	assignmentType, typeErr := models.NormalizeAssignmentType(request.AssignmentType)

	if typeErr != nil {
		return nil, ErrInvalidAssignmentType
	}
	if !isValidDateRange(request.ValidFrom, request.ValidTo) {
		return nil, ErrInvalidDateRange
	}
	if userNumber == "" {
		return nil, ErrUserRequired
	}

	user, err := authRepository.GetUserByNumber(userNumber)
	if err != nil {
		return nil, err
	}
	if user.Number == "" {
		return nil, ErrUserNotFound
	}

	if err := validateUnit(organisationNumber, unitNumber); err != nil {
		return nil, err
	}

	existing, err := assignmentRepositories.FindActiveByOrgUserUnitType(organisationNumber, userNumber, unitNumber, assignmentType)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrDuplicateActiveAssignment
	}

	assignment := models.UserAssignment{
		Number:                 uuid.NewV4().String(),
		OrganisationNumber:     organisationNumber,
		UserNumber:             userNumber,
		OrganisationUnitNumber: unitNumber,
		AssignmentType:         assignmentType,
		PositionTitle:          strings.TrimSpace(request.PositionTitle),
		ValidFrom:              request.ValidFrom,
		ValidTo:                request.ValidTo,
		Status:                 constants.StatusActive,
		AddedByNumber:          actorNumber,
	}

	if err := assignmentRepositories.Create(&assignment); err != nil {
		return nil, err
	}

	response := toResponse(assignment)
	return &response, nil
}

// UpdateAssignment updates placement and lifecycle fields. The person is
// immutable. The unit is re-validated (same organisation, ACTIVE).
func UpdateAssignment(organisationNumber, number string, request requests.UpdateUserAssignmentRequest) (*responses.UserAssignmentResponse, error) {
	assignment, err := assignmentRepositories.FindByNumberAndOrganisation(number, organisationNumber)
	if err != nil {
		return nil, ErrAssignmentNotFound
	}

	unitNumber := strings.TrimSpace(request.OrganisationUnitNumber)
	assignmentType, typeErr := models.NormalizeAssignmentType(request.AssignmentType)

	if typeErr != nil {
		return nil, ErrInvalidAssignmentType
	}
	if !isValidDateRange(request.ValidFrom, request.ValidTo) {
		return nil, ErrInvalidDateRange
	}
	if err := validateUnit(organisationNumber, unitNumber); err != nil {
		return nil, err
	}

	if unitNumber != assignment.OrganisationUnitNumber || assignmentType != assignment.AssignmentType {
		existing, err := assignmentRepositories.FindActiveByOrgUserUnitType(organisationNumber, assignment.UserNumber, unitNumber, assignmentType)
		if err != nil {
			return nil, err
		}
		if existing != nil && existing.Number != assignment.Number {
			return nil, ErrDuplicateActiveAssignment
		}
	}

	assignment.OrganisationUnitNumber = unitNumber
	assignment.AssignmentType = assignmentType
	assignment.PositionTitle = strings.TrimSpace(request.PositionTitle)
	assignment.ValidFrom = request.ValidFrom
	assignment.ValidTo = request.ValidTo

	if err := assignmentRepositories.Save(assignment); err != nil {
		return nil, err
	}

	response := toResponse(*assignment)
	return &response, nil
}

// SetAssignmentStatus applies the ACTIVE / INACTIVE / ENDED lifecycle.
//
// Reactivating re-checks the single-ACTIVE rule. Moving to ENDED stamps ValidTo
// when it is not already set; the assignment is retained for history. Changing
// an assignment's status never changes a user's authentication status or a
// stakeholder's status.
func SetAssignmentStatus(organisationNumber, number, status string) (*responses.UserAssignmentResponse, error) {
	if !isValidAssignmentStatus(status) {
		return nil, ErrInvalidStatus
	}

	assignment, err := assignmentRepositories.FindByNumberAndOrganisation(number, organisationNumber)
	if err != nil {
		return nil, ErrAssignmentNotFound
	}

	if status == constants.StatusActive {
		existing, err := assignmentRepositories.FindActiveByOrgUserUnitType(
			organisationNumber, assignment.UserNumber, assignment.OrganisationUnitNumber, assignment.AssignmentType)
		if err != nil {
			return nil, err
		}
		if existing != nil && existing.Number != assignment.Number {
			return nil, ErrDuplicateActiveAssignment
		}
	}

	assignment.Status = status
	if status == AssignmentStatusEnded && assignment.ValidTo == nil {
		now := time.Now()
		assignment.ValidTo = &now
	}

	if err := assignmentRepositories.Save(assignment); err != nil {
		return nil, err
	}

	response := toResponse(*assignment)
	return &response, nil
}

// CurrentOrganisationNumber is re-exported so callers depend on one resolver.
var CurrentOrganisationNumber = organisationServices.CurrentOrganisationNumber

func toResponse(assignment models.UserAssignment) responses.UserAssignmentResponse {
	return responses.UserAssignmentResponse{
		Number:                 assignment.Number,
		OrganisationNumber:     assignment.OrganisationNumber,
		UserNumber:             assignment.UserNumber,
		OrganisationUnitNumber: assignment.OrganisationUnitNumber,
		AssignmentType:         assignment.AssignmentType,
		PositionTitle:          assignment.PositionTitle,
		ValidFrom:              assignment.ValidFrom,
		ValidTo:                assignment.ValidTo,
		Status:                 assignment.Status,
		AddedByNumber:          assignment.AddedByNumber,
		CreatedAt:              assignment.CreatedAt,
		UpdatedAt:              assignment.UpdatedAt,
	}
}
