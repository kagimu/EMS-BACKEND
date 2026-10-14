package services

import (
	"errors"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/authority/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/authority/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/authority/models"
	authorityRepositories "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/authority/repositories"
	organisationRepositories "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/organisation/repositories"
	organisationServices "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/organisation/services"
	stakeholderModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/stakeholders/models"
	uuid "github.com/satori/go.uuid"
)

const AuthorityStatusEnded = "ENDED"

var (
	ErrAuthorityNotFound      = errors.New("decision authority not found")
	ErrInvalidDecisionType    = errors.New("decision_type is required")
	ErrInvalidApprovalLevel   = errors.New("approval_level must be at least 1")
	ErrInvalidAmount          = errors.New("amount must not be negative")
	ErrInvalidAmountRange     = errors.New("min_amount cannot be greater than max_amount")
	ErrInvalidAnchor          = errors.New("exactly one authority anchor (role_number or stakeholder_type) is required")
	ErrInvalidStakeholderType = errors.New("invalid stakeholder type")
	ErrRoleNotFound           = errors.New("referenced role was not found in this organisation")
	ErrInvalidStatus          = errors.New("invalid status")
	ErrInvalidDateRange       = errors.New("valid_from cannot be after valid_to")
	ErrUnitNotFound           = errors.New("organisation unit was not found in this organisation")
	ErrUnitNotActive          = errors.New("organisation unit is inactive")
)

func isValidAuthorityStatus(status string) bool {
	return status == constants.StatusActive ||
		status == constants.StatusInactive ||
		status == AuthorityStatusEnded
}

func isValidDateRange(from, to *time.Time) bool {
	if from == nil || to == nil {
		return true
	}
	return !from.After(*to)
}

func isNegativeAmount(value *float64) bool {
	return value != nil && *value < 0
}

// isValidAmountRange enforces that MinAmount <= MaxAmount when both are present.
func isValidAmountRange(min, max *float64) bool {
	if min == nil || max == nil {
		return true
	}
	return *min <= *max
}

// toMinorUnits converts a two-decimal monetary value into integer minor units
// (cents) so authority coverage can be compared WITHOUT floating-point
// comparison error.
//
// The database column remains DECIMAL(18,2); only the comparison is done in
// integer minor units.
func toMinorUnits(value float64) int64 {
	return int64(math.Round(value * 100))
}

// levelCovers reports whether any authority row at a level covers the amount.
//
// Ratified rule (BD-1): coverage is determined by MaxAmount ONLY — a row covers
// the amount when MaxAmount is NULL (unlimited) or the amount is less than or
// equal to MaxAmount. MinAmount is stored and validated as a declared lower
// bound but is NOT used to determine coverage.
//
// Amounts are compared in integer minor units (cents) to avoid unsafe
// floating-point comparison.
func levelCovers(rows []models.DecisionAuthority, amount float64) bool {
	amountMinor := toMinorUnits(amount)
	for _, row := range rows {
		if row.MaxAmount == nil || amountMinor <= toMinorUnits(*row.MaxAmount) {
			return true
		}
	}
	return false
}

// rowScopeApplies reports whether an authority row applies to the requested unit.
// An empty OrganisationUnitNumber is organisation-wide; otherwise the row's unit
// must be the requested unit or one of its ancestors.
func rowScopeApplies(row models.DecisionAuthority, ancestorSet map[string]bool) bool {
	if row.OrganisationUnitNumber == "" {
		return true
	}
	return ancestorSet[row.OrganisationUnitNumber]
}

// validateUnit is the WRITE-path unit rule (BD-3): a unit referenced when an
// authority is created or updated must exist, belong to the same organisation,
// and be ACTIVE.
func validateUnit(organisationNumber, unitNumber string) error {
	if strings.TrimSpace(unitNumber) == "" {
		return nil
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

// validateUnitExists is the RESOLUTION-path unit rule (BD-3): the requested unit
// must exist in the same organisation, but it does NOT have to remain ACTIVE.
// Deactivating a unit must not silently revoke authority that was validly
// configured against it; revocation is explicit, via the authority lifecycle.
func validateUnitExists(organisationNumber, unitNumber string) error {
	if strings.TrimSpace(unitNumber) == "" {
		return nil
	}

	if _, err := organisationRepositories.FindByNumberAndOrganisation(unitNumber, organisationNumber); err != nil {
		return ErrUnitNotFound
	}

	return nil
}

// validateAnchor enforces exactly one anchor and validates it.
func validateAnchor(organisationNumber, roleNumber, stakeholderType string) error {
	role := strings.TrimSpace(roleNumber)
	stakeholder := strings.TrimSpace(stakeholderType)

	if (role == "") == (stakeholder == "") {
		return ErrInvalidAnchor
	}

	if stakeholder != "" && !stakeholderModels.IsValidStakeholderType(stakeholder) {
		return ErrInvalidStakeholderType
	}

	if role != "" {
		exists, err := authorityRepositories.RoleExistsInOrganisation(role, organisationNumber)
		if err != nil {
			return err
		}
		if !exists {
			return ErrRoleNotFound
		}
	}

	return nil
}

func validateAuthorityInput(organisationNumber, decisionType string, approvalLevel int, min, max *float64,
	roleNumber, stakeholderType, unitNumber string, from, to *time.Time) error {

	if strings.TrimSpace(decisionType) == "" {
		return ErrInvalidDecisionType
	}
	if approvalLevel < 1 {
		return ErrInvalidApprovalLevel
	}
	if isNegativeAmount(min) || isNegativeAmount(max) {
		return ErrInvalidAmount
	}
	if !isValidAmountRange(min, max) {
		return ErrInvalidAmountRange
	}
	if !isValidDateRange(from, to) {
		return ErrInvalidDateRange
	}
	if err := validateAnchor(organisationNumber, roleNumber, stakeholderType); err != nil {
		return err
	}
	return validateUnit(organisationNumber, unitNumber)
}

// ListDecisionAuthorities returns the authority rows of an organisation.
func ListDecisionAuthorities(organisationNumber, decisionType, status string) ([]responses.DecisionAuthorityResponse, error) {
	if status != "" && !isValidAuthorityStatus(status) {
		return nil, ErrInvalidStatus
	}

	authorities, err := authorityRepositories.FindAllByOrganisation(organisationNumber, decisionType, status)
	if err != nil {
		return nil, err
	}

	result := make([]responses.DecisionAuthorityResponse, 0, len(authorities))
	for _, authority := range authorities {
		result = append(result, toResponse(authority))
	}
	return result, nil
}

// GetDecisionAuthority returns one authority row within the organisation.
func GetDecisionAuthority(organisationNumber, number string) (*responses.DecisionAuthorityResponse, error) {
	authority, err := authorityRepositories.FindByNumberAndOrganisation(number, organisationNumber)
	if err != nil {
		return nil, ErrAuthorityNotFound
	}
	response := toResponse(*authority)
	return &response, nil
}

// CreateDecisionAuthority creates an authority definition.
func CreateDecisionAuthority(actorNumber, organisationNumber string, request requests.CreateDecisionAuthorityRequest) (*responses.DecisionAuthorityResponse, error) {
	if err := validateAuthorityInput(organisationNumber, request.DecisionType, request.ApprovalLevel,
		request.MinAmount, request.MaxAmount, request.RoleNumber, request.StakeholderType,
		request.OrganisationUnitNumber, request.ValidFrom, request.ValidTo); err != nil {
		return nil, err
	}

	authority := models.DecisionAuthority{
		Number:                 uuid.NewV4().String(),
		OrganisationNumber:     organisationNumber,
		DecisionType:           strings.TrimSpace(request.DecisionType),
		ApprovalLevel:          request.ApprovalLevel,
		MinAmount:              request.MinAmount,
		MaxAmount:              request.MaxAmount,
		RoleNumber:             strings.TrimSpace(request.RoleNumber),
		StakeholderType:        strings.TrimSpace(request.StakeholderType),
		OrganisationUnitNumber: strings.TrimSpace(request.OrganisationUnitNumber),
		Status:                 constants.StatusActive,
		ValidFrom:              request.ValidFrom,
		ValidTo:                request.ValidTo,
		AddedByNumber:          actorNumber,
	}

	if err := authorityRepositories.Create(&authority); err != nil {
		return nil, err
	}

	response := toResponse(authority)
	return &response, nil
}

// UpdateDecisionAuthority replaces the editable fields of an authority row.
func UpdateDecisionAuthority(organisationNumber, number string, request requests.UpdateDecisionAuthorityRequest) (*responses.DecisionAuthorityResponse, error) {
	authority, err := authorityRepositories.FindByNumberAndOrganisation(number, organisationNumber)
	if err != nil {
		return nil, ErrAuthorityNotFound
	}

	if err := validateAuthorityInput(organisationNumber, request.DecisionType, request.ApprovalLevel,
		request.MinAmount, request.MaxAmount, request.RoleNumber, request.StakeholderType,
		request.OrganisationUnitNumber, request.ValidFrom, request.ValidTo); err != nil {
		return nil, err
	}

	authority.DecisionType = strings.TrimSpace(request.DecisionType)
	authority.ApprovalLevel = request.ApprovalLevel
	authority.MinAmount = request.MinAmount
	authority.MaxAmount = request.MaxAmount
	authority.RoleNumber = strings.TrimSpace(request.RoleNumber)
	authority.StakeholderType = strings.TrimSpace(request.StakeholderType)
	authority.OrganisationUnitNumber = strings.TrimSpace(request.OrganisationUnitNumber)
	authority.ValidFrom = request.ValidFrom
	authority.ValidTo = request.ValidTo

	if err := authorityRepositories.Save(authority); err != nil {
		return nil, err
	}

	response := toResponse(*authority)
	return &response, nil
}

// SetDecisionAuthorityStatus applies the ACTIVE / INACTIVE / ENDED lifecycle.
func SetDecisionAuthorityStatus(organisationNumber, number, status string) (*responses.DecisionAuthorityResponse, error) {
	if !isValidAuthorityStatus(status) {
		return nil, ErrInvalidStatus
	}

	authority, err := authorityRepositories.FindByNumberAndOrganisation(number, organisationNumber)
	if err != nil {
		return nil, ErrAuthorityNotFound
	}

	authority.Status = status

	if err := authorityRepositories.Save(authority); err != nil {
		return nil, err
	}

	response := toResponse(*authority)
	return &response, nil
}

// ResolveAuthority is Foundation D's core deliverable: a deterministic,
// read-only computation of the required approval path for a decision type and
// amount within the authenticated organisation.
//
// It loads applicable ACTIVE authority rows in a single query, applies unit
// scope in memory (no per-row queries), groups rows by approval level, and walks
// the levels in ascending order until one covers the amount. It never creates
// approvals, never routes anything, and never writes.
//
// Three states are distinguished for Foundation E (BD-2):
//
//	State 1: configured = false                      -> nothing is configured for
//	                                                    this decision type.
//	State 2: configured = true, non-empty path       -> governed and applicable.
//	State 3: configured = true, empty path           -> governed, but no authority
//	                                                    applies to this scope.
//
// `configured` is therefore scope-INDEPENDENT. A unit-less request applies only
// organisation-wide rows (BD-4), and a unit that has since become INACTIVE still
// resolves (BD-3).
func ResolveAuthority(organisationNumber, decisionType, unitNumber string, amount float64) (*responses.AuthorityResolutionResponse, error) {
	decision := strings.TrimSpace(decisionType)
	if decision == "" {
		return nil, ErrInvalidDecisionType
	}
	// Non-finite values (NaN, +Inf, -Inf) are not valid monetary amounts and must
	// never silently enter authority coverage evaluation.
	if math.IsNaN(amount) || math.IsInf(amount, 0) {
		return nil, ErrInvalidAmount
	}
	if amount < 0 {
		return nil, ErrInvalidAmount
	}

	unit := strings.TrimSpace(unitNumber)
	if unit != "" {
		// BD-3: existence and same-organisation are required; ACTIVE is not.
		if err := validateUnitExists(organisationNumber, unit); err != nil {
			return nil, err
		}
	}

	rows, err := authorityRepositories.FindActiveForResolution(organisationNumber, decision, time.Now())
	if err != nil {
		return nil, err
	}

	ancestors, err := ancestorSet(organisationNumber, unit)
	if err != nil {
		return nil, err
	}

	applicable := make([]models.DecisionAuthority, 0, len(rows))
	for _, row := range rows {
		if rowScopeApplies(row, ancestors) {
			applicable = append(applicable, row)
		}
	}

	// BD-2: `configured` is scope-INDEPENDENT — it reports whether this decision
	// type is governed by any ACTIVE authority definition in the organisation,
	// before unit-scope filtering.
	response := &responses.AuthorityResolutionResponse{
		Configured:             len(rows) > 0,
		DecisionType:           decision,
		Amount:                 amount,
		OrganisationUnitNumber: unit,
		RequiredPath:           []responses.AuthorityLevelResponse{},
	}
	if len(applicable) == 0 {
		// State 3: authority exists for this decision type, but none applies to
		// this scope/request. exceeds_declared_authority remains false because
		// there is nothing declared for THIS scope to exceed.
		return response, nil
	}

	levels := make(map[int][]models.DecisionAuthority)
	for _, row := range applicable {
		levels[row.ApprovalLevel] = append(levels[row.ApprovalLevel], row)
	}

	ordered := make([]int, 0, len(levels))
	for level := range levels {
		ordered = append(ordered, level)
	}
	sort.Ints(ordered)

	terminal := false
	for _, level := range ordered {
		response.RequiredPath = append(response.RequiredPath, levelResponse(level, levels[level]))
		if levelCovers(levels[level], amount) {
			terminal = true
			break
		}
	}
	// Applicable authority exists here, so a non-terminal walk genuinely
	// exceeds the declared limits for this scope.
	response.ExceedsDeclaredAuthority = !terminal

	return response, nil
}

// ancestorSet returns the requested unit and all of its ancestors within the
// organisation, reading the whole unit set once (no per-level queries).
func ancestorSet(organisationNumber, unitNumber string) (map[string]bool, error) {
	set := map[string]bool{}
	if unitNumber == "" {
		return set, nil
	}

	units, err := organisationRepositories.FindAllByOrganisation(organisationNumber, "", "")
	if err != nil {
		return nil, err
	}

	parents := make(map[string]string, len(units))
	for _, unit := range units {
		parents[unit.Number] = unit.ParentNumber
	}

	current := unitNumber
	for current != "" && !set[current] {
		set[current] = true
		current = parents[current]
	}

	return set, nil
}

// levelResponse converts the authority rows of one approval level into the
// anchors that level requires.
//
// Foundation D (ratified D5): an anchor's identity is **type + value + unit**, so
// the de-duplication key includes the row's organisation-unit scope. Two rows at
// the same level with the same role but different units — for example a
// Kampala-scoped and a Jinja-scoped Finance Manager — therefore produce two
// DISTINCT anchors and can never collapse into one. An empty unit represents
// organisation-wide scope.
func levelResponse(level int, rows []models.DecisionAuthority) responses.AuthorityLevelResponse {
	anchors := []responses.AuthorityAnchorResponse{}
	seen := map[string]bool{}

	for _, row := range rows {
		unit := row.OrganisationUnitNumber
		if row.RoleNumber != "" {
			key := "role|" + row.RoleNumber + "|" + unit
			if !seen[key] {
				seen[key] = true
				anchors = append(anchors, responses.AuthorityAnchorResponse{
					Type: "role", Value: row.RoleNumber, OrganisationUnitNumber: unit,
				})
			}
		}
		if row.StakeholderType != "" {
			key := "stakeholder_type|" + row.StakeholderType + "|" + unit
			if !seen[key] {
				seen[key] = true
				anchors = append(anchors, responses.AuthorityAnchorResponse{
					Type: "stakeholder_type", Value: row.StakeholderType, OrganisationUnitNumber: unit,
				})
			}
		}
	}

	return responses.AuthorityLevelResponse{ApprovalLevel: level, Anchors: anchors}
}

// CurrentOrganisationNumber is re-exported so callers depend on one resolver.
var CurrentOrganisationNumber = organisationServices.CurrentOrganisationNumber

func toResponse(authority models.DecisionAuthority) responses.DecisionAuthorityResponse {
	return responses.DecisionAuthorityResponse{
		Number:                 authority.Number,
		OrganisationNumber:     authority.OrganisationNumber,
		DecisionType:           authority.DecisionType,
		ApprovalLevel:          authority.ApprovalLevel,
		MinAmount:              authority.MinAmount,
		MaxAmount:              authority.MaxAmount,
		RoleNumber:             authority.RoleNumber,
		StakeholderType:        authority.StakeholderType,
		OrganisationUnitNumber: authority.OrganisationUnitNumber,
		Status:                 authority.Status,
		ValidFrom:              authority.ValidFrom,
		ValidTo:                authority.ValidTo,
		AddedByNumber:          authority.AddedByNumber,
		CreatedAt:              authority.CreatedAt,
		UpdatedAt:              authority.UpdatedAt,
	}
}
