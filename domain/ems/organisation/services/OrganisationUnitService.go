package services

import (
	"errors"
	"strings"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
	assignmentRepositories "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/assignments/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/organisation/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/organisation/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/organisation/models"
	organisationRepositories "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/organisation/repositories"
	eosRepositories "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/eos/repositories"
	uuid "github.com/satori/go.uuid"
)

var (
	ErrNoActiveOrganisation = errors.New("no active organisation is set for this user")
	ErrUnitNotFound         = errors.New("organisation unit not found")
	ErrParentNotFound       = errors.New("parent organisation unit was not found in this organisation")
	ErrParentInactive       = errors.New("parent organisation unit is inactive")
	ErrSelfParent           = errors.New("an organisation unit cannot be its own parent")
	ErrCircularHierarchy    = errors.New("circular organisation-unit hierarchy is not allowed")
	ErrHasChildren          = errors.New("organisation unit has child units and cannot be deleted")
	ErrHasActiveAssignments = errors.New("organisation unit has active assignments and cannot be deleted")
	ErrInvalidStatus        = errors.New("invalid status")
	ErrDuplicateCode        = errors.New("code already exists in this organisation")
	ErrNameRequired         = errors.New("name is required")
	ErrUnitTypeRequired     = errors.New("unit type is required")
)

// CurrentOrganisationNumber resolves the authoritative organisation for the
// authenticated user through the existing eos mechanism:
//
//	User.Number -> ActiveUserEntity -> Entity.Number
//
// It reuses the existing repository; it does not create a parallel resolver.
// When no active organisation exists it returns ErrNoActiveOrganisation
// (mapped to HTTP 409). No fallback organisation is ever invented.
func CurrentOrganisationNumber(userNumber string) (string, error) {
	active, err := eosRepositories.FindActiveUserEntityByUserNumber(userNumber)
	if err != nil {
		return "", err
	}
	if active.UserNumber == "" || active.EntityNumber == "" {
		return "", ErrNoActiveOrganisation
	}
	return active.EntityNumber, nil
}

func isSelfParent(unitNumber, parentNumber string) bool {
	return unitNumber != "" && unitNumber == parentNumber
}

func isActiveStatus(status string) bool {
	return status == constants.StatusActive
}

func isValidStatus(status string) bool {
	return status == constants.StatusActive || status == constants.StatusInactive
}

// ancestryContains reports whether target appears in the supplied ancestor chain.
func ancestryContains(ancestorNumbers []string, target string) bool {
	if target == "" {
		return false
	}
	for _, ancestor := range ancestorNumbers {
		if ancestor == target {
			return true
		}
	}
	return false
}

// collectAncestors walks upwards from startParentNumber and returns the chain
// (including the start). Walking stops if a link is missing.
func collectAncestors(organisationNumber, startParentNumber string) []string {
	var ancestors []string
	current := startParentNumber
	for current != "" {
		ancestors = append(ancestors, current)
		parent, err := organisationRepositories.FindByNumberAndOrganisation(current, organisationNumber)
		if err != nil {
			break
		}
		current = parent.ParentNumber
	}
	return ancestors
}

// validateParent enforces the hierarchy rules for a proposed parent.
func validateParent(organisationNumber, unitNumber, parentNumber string) error {
	if parentNumber == "" {
		return nil
	}
	if isSelfParent(unitNumber, parentNumber) {
		return ErrSelfParent
	}

	parent, err := organisationRepositories.FindByNumberAndOrganisation(parentNumber, organisationNumber)
	if err != nil {
		return ErrParentNotFound
	}
	if !isActiveStatus(parent.Status) {
		return ErrParentInactive
	}

	if unitNumber != "" && ancestryContains(collectAncestors(organisationNumber, parentNumber), unitNumber) {
		return ErrCircularHierarchy
	}

	return nil
}

// ListUnits returns the units of an organisation, optionally filtered.
func ListUnits(organisationNumber, parentNumber, status string) ([]responses.OrganisationUnitResponse, error) {
	if status != "" && !isValidStatus(status) {
		return nil, ErrInvalidStatus
	}

	units, err := organisationRepositories.FindAllByOrganisation(organisationNumber, parentNumber, status)
	if err != nil {
		return nil, err
	}

	result := make([]responses.OrganisationUnitResponse, 0, len(units))
	for _, unit := range units {
		result = append(result, toResponse(unit))
	}
	return result, nil
}

// GetUnit returns a single unit within the organisation, or ErrUnitNotFound.
func GetUnit(organisationNumber, number string) (*responses.OrganisationUnitResponse, error) {
	unit, err := organisationRepositories.FindByNumberAndOrganisation(number, organisationNumber)
	if err != nil {
		return nil, ErrUnitNotFound
	}
	response := toResponse(*unit)
	return &response, nil
}

// CreateUnit creates a unit in the organisation. The caller supplies the
// authenticated actor; the organisation comes from the authenticated context.
func CreateUnit(actorUserNumber, organisationNumber string, request requests.CreateOrganisationUnitRequest) (*responses.OrganisationUnitResponse, error) {
	name := strings.TrimSpace(request.Name)
	unitType := strings.TrimSpace(request.UnitType)
	code := strings.TrimSpace(request.Code)
	parentNumber := strings.TrimSpace(request.ParentNumber)

	if name == "" {
		return nil, ErrNameRequired
	}
	if unitType == "" {
		return nil, ErrUnitTypeRequired
	}
	if err := validateParent(organisationNumber, "", parentNumber); err != nil {
		return nil, err
	}
	if code != "" {
		if existing, err := organisationRepositories.FindByCode(organisationNumber, code); err == nil && existing != nil {
			return nil, ErrDuplicateCode
		}
	}

	unit := models.OrganisationUnit{
		Number:             uuid.NewV4().String(),
		OrganisationNumber: organisationNumber,
		ParentNumber:       parentNumber,
		UnitType:           unitType,
		Name:               name,
		Code:               code,
		Status:             constants.StatusActive,
		AddedByNumber:      actorUserNumber,
	}

	if err := organisationRepositories.Create(&unit); err != nil {
		return nil, err
	}

	response := toResponse(unit)
	return &response, nil
}

// UpdateUnit updates editable fields and re-validates hierarchy rules when the
// parent changes.
func UpdateUnit(organisationNumber, number string, request requests.UpdateOrganisationUnitRequest) (*responses.OrganisationUnitResponse, error) {
	unit, err := organisationRepositories.FindByNumberAndOrganisation(number, organisationNumber)
	if err != nil {
		return nil, ErrUnitNotFound
	}

	name := strings.TrimSpace(request.Name)
	unitType := strings.TrimSpace(request.UnitType)
	code := strings.TrimSpace(request.Code)
	parentNumber := strings.TrimSpace(request.ParentNumber)

	if name == "" {
		return nil, ErrNameRequired
	}
	if unitType == "" {
		return nil, ErrUnitTypeRequired
	}

	if parentNumber != unit.ParentNumber {
		if err := validateParent(organisationNumber, unit.Number, parentNumber); err != nil {
			return nil, err
		}
	}

	if code != "" && code != unit.Code {
		if existing, err := organisationRepositories.FindByCode(organisationNumber, code); err == nil && existing != nil {
			return nil, ErrDuplicateCode
		}
	}

	unit.Name = name
	unit.UnitType = unitType
	unit.Code = code
	unit.ParentNumber = parentNumber

	if err := organisationRepositories.Save(unit); err != nil {
		return nil, err
	}

	response := toResponse(*unit)
	return &response, nil
}

// SetUnitStatus activates or deactivates a unit.
func SetUnitStatus(organisationNumber, number, status string) (*responses.OrganisationUnitResponse, error) {
	if !isValidStatus(status) {
		return nil, ErrInvalidStatus
	}

	unit, err := organisationRepositories.FindByNumberAndOrganisation(number, organisationNumber)
	if err != nil {
		return nil, ErrUnitNotFound
	}

	unit.Status = status

	if err := organisationRepositories.Save(unit); err != nil {
		return nil, err
	}

	response := toResponse(*unit)
	return &response, nil
}

// DeleteUnit soft-deletes a unit, refusing when it still has children.
//
// Foundation C3 integrity enhancement (backward compatible): deletion is also
// refused when the unit has ACTIVE assignments, so assignments are never
// silently orphaned. Historical INACTIVE/ENDED assignments do NOT block
// deletion, and no assignment is cascaded, deleted, or automatically ended.
func DeleteUnit(organisationNumber, number string) error {
	unit, err := organisationRepositories.FindByNumberAndOrganisation(number, organisationNumber)
	if err != nil {
		return ErrUnitNotFound
	}

	children, err := organisationRepositories.CountChildren(number, organisationNumber)
	if err != nil {
		return err
	}
	if children > 0 {
		return ErrHasChildren
	}

	activeAssignments, err := assignmentRepositories.CountActiveByUnit(number, organisationNumber)
	if err != nil {
		return err
	}
	if activeAssignments > 0 {
		return ErrHasActiveAssignments
	}

	return organisationRepositories.Delete(unit)
}

func toResponse(unit models.OrganisationUnit) responses.OrganisationUnitResponse {
	return responses.OrganisationUnitResponse{
		Number:             unit.Number,
		OrganisationNumber: unit.OrganisationNumber,
		ParentNumber:       unit.ParentNumber,
		UnitType:           unit.UnitType,
		Name:               unit.Name,
		Code:               unit.Code,
		Status:             unit.Status,
		AddedByNumber:      unit.AddedByNumber,
		CreatedAt:          unit.CreatedAt,
		UpdatedAt:          unit.UpdatedAt,
	}
}
