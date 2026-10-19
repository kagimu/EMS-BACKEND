package services

import (
	"errors"
	"strings"
	"time"

	authRepository "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/repository"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
	organisationServices "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/organisation/services"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/stakeholders/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/stakeholders/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/stakeholders/models"
	stakeholderRepositories "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/stakeholders/repositories"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// StakeholderStatusEnded terminates a relationship while retaining it for
// history. It is deliberately distinct from auth.User's Title Case
// "Active"/"Inactive" status values, which are NOT used for EMS stakeholders.
const StakeholderStatusEnded = "ENDED"

var (
	ErrStakeholderNotFound        = errors.New("stakeholder relationship not found")
	ErrDuplicateActiveStakeholder = errors.New("an active relationship of this type already exists for this person in this organisation")
	ErrInvalidStakeholderType     = errors.New("invalid stakeholder type")
	ErrInvalidStatus              = errors.New("invalid status")
	ErrInvalidSharePercentage     = errors.New("share percentage must be between 0 and 100")
	ErrInvalidDateRange           = errors.New("valid_from cannot be after valid_to")
	ErrUserRequired               = errors.New("user_number is required")
	ErrUserNotFound               = errors.New("referenced user was not found")
	ErrLastActiveOwner            = errors.New("an organisation must retain at least one active owner")
)

// ErrOwnershipChangeBusy is re-exported so controllers can map the per-
// organisation ownership lock timeout without importing the repository package.
var ErrOwnershipChangeBusy = stakeholderRepositories.ErrOwnershipChangeBusy

func isValidStakeholderStatus(status string) bool {
	return status == constants.StatusActive ||
		status == constants.StatusInactive ||
		status == StakeholderStatusEnded
}

func isValidSharePercentage(value *float64) bool {
	if value == nil {
		return true
	}
	return *value >= 0 && *value <= 100
}

func isValidDateRange(from, to *time.Time) bool {
	if from == nil || to == nil {
		return true
	}
	return !from.After(*to)
}

// isFinalActiveOwner reports whether moving an ACTIVE owner relationship to a
// non-active status would remove the organisation's final ACTIVE owner
// (Foundation C4 last-owner invariant).
//
// Only ACTIVE owner relationships count; INACTIVE and ENDED owners do not.
// The rule is actor-independent: it makes no reference to who is performing the
// operation.
func isFinalActiveOwner(stakeholderType, currentStatus, targetStatus string, activeOwnerCount int64) bool {
	if stakeholderType != models.StakeholderTypeOwner {
		return false
	}
	if currentStatus != constants.StatusActive {
		return false
	}
	if targetStatus != constants.StatusInactive && targetStatus != StakeholderStatusEnded {
		return false
	}
	return activeOwnerCount <= 1
}

// isEffectiveOwner applies the ratified effective-owner predicate to one row at
// a single instant: owner, ACTIVE, and the validity window contains now
// (valid_to inclusive; NULL on either side is unbounded). It is used to decide
// whether a mutation actually removes an owner from the effective set.
func isEffectiveOwner(stakeholder models.Stakeholder, now time.Time) bool {
	if stakeholder.StakeholderType != models.StakeholderTypeOwner || stakeholder.Status != constants.StatusActive {
		return false
	}
	if stakeholder.ValidFrom != nil && stakeholder.ValidFrom.After(now) {
		return false
	}
	if stakeholder.ValidTo != nil && stakeholder.ValidTo.Before(now) {
		return false
	}
	return true
}

// authenticatableUserNumbers derives, in one batched query for the whole set,
// which referenced users currently have a usable login (Foundation C2.3).
func authenticatableUserNumbers(userNumbers []string) (map[string]bool, error) {
	unique := make([]string, 0, len(userNumbers))
	seen := make(map[string]bool)
	for _, number := range userNumbers {
		if number == "" || seen[number] {
			continue
		}
		seen[number] = true
		unique = append(unique, number)
	}
	return stakeholderRepositories.FindAuthenticatableUserNumbers(unique)
}

// ListStakeholders returns the stakeholders of an organisation, optionally
// filtered by stakeholder type and/or status.
func ListStakeholders(organisationNumber, stakeholderType, status string) ([]responses.StakeholderResponse, error) {
	if stakeholderType != "" && !models.IsValidStakeholderType(stakeholderType) {
		return nil, ErrInvalidStakeholderType
	}
	if status != "" && !isValidStakeholderStatus(status) {
		return nil, ErrInvalidStatus
	}

	stakeholders, err := stakeholderRepositories.FindAllByOrganisation(organisationNumber, stakeholderType, status)
	if err != nil {
		return nil, err
	}

	userNumbers := make([]string, 0, len(stakeholders))
	for _, stakeholder := range stakeholders {
		userNumbers = append(userNumbers, stakeholder.UserNumber)
	}
	canAuthenticate, err := authenticatableUserNumbers(userNumbers)
	if err != nil {
		return nil, err
	}
	displayNames, err := stakeholderRepositories.FindUserDisplayNames(userNumbers)
	if err != nil {
		return nil, err
	}

	result := make([]responses.StakeholderResponse, 0, len(stakeholders))
	for _, stakeholder := range stakeholders {
		result = append(result, toResponse(
			stakeholder,
			canAuthenticate[stakeholder.UserNumber],
			displayNames[stakeholder.UserNumber],
		))
	}
	return result, nil
}

// GetStakeholder returns one stakeholder within the organisation.
func GetStakeholder(organisationNumber, number string) (*responses.StakeholderResponse, error) {
	stakeholder, err := stakeholderRepositories.FindByNumberAndOrganisation(number, organisationNumber)
	if err != nil {
		return nil, ErrStakeholderNotFound
	}
	return responseFor(*stakeholder)
}

// CreateStakeholder creates a generic relationship for an EXISTING person.
//
// Foundation C2 integration rules encoded here:
//   - V2: the referenced user must exist (a soft-deleted user is treated as
//     absent, because the lookup excludes soft-deleted rows).
//   - V3: the referenced user does NOT have to be Active — a relationship may be
//     recorded for a person without a usable login.
//   - V4: the referenced user does NOT have to have EntityUser membership in
//     this organisation — external/governance relationships are valid records.
//   - V8: this operation grants NO system access. Authenticated EMS operation
//     still requires EntityUser membership plus an active organisation context,
//     which are enforced by the existing mechanisms, not here.
//
// It never creates a user, contact, password, or login, and never writes to any
// authentication table.
func CreateStakeholder(actorNumber, organisationNumber string, request requests.CreateStakeholderRequest) (*responses.StakeholderResponse, error) {
	userNumber := strings.TrimSpace(request.UserNumber)
	stakeholderType := strings.TrimSpace(request.StakeholderType)

	if userNumber == "" {
		return nil, ErrUserRequired
	}
	if !models.IsValidStakeholderType(stakeholderType) {
		return nil, ErrInvalidStakeholderType
	}
	if !isValidSharePercentage(request.SharePercentage) {
		return nil, ErrInvalidSharePercentage
	}
	if !isValidDateRange(request.ValidFrom, request.ValidTo) {
		return nil, ErrInvalidDateRange
	}

	user, err := authRepository.GetUserByNumber(userNumber)
	if err != nil {
		return nil, err
	}
	if user.Number == "" {
		return nil, ErrUserNotFound
	}

	stakeholder := models.Stakeholder{
		Number:             uuid.NewV4().String(),
		OrganisationNumber: organisationNumber,
		UserNumber:         userNumber,
		StakeholderType:    stakeholderType,
		Title:              strings.TrimSpace(request.Title),
		SharePercentage:    request.SharePercentage,
		Shares:             request.Shares,
		ValidFrom:          request.ValidFrom,
		ValidTo:            request.ValidTo,
		Status:             constants.StatusActive,
		AddedByNumber:      actorNumber,
	}

	// The duplicate check and the insert are serialised per organisation so two
	// concurrent requests cannot both create an ACTIVE relationship.
	if err := stakeholderRepositories.WithOwnershipTransaction(organisationNumber, func(tx *gorm.DB) error {
		existing, err := stakeholderRepositories.FindActiveByOrgUserType(organisationNumber, userNumber, stakeholderType, tx)
		if err != nil {
			return err
		}
		if existing != nil {
			return ErrDuplicateActiveStakeholder
		}
		return stakeholderRepositories.Create(&stakeholder, tx)
	}); err != nil {
		return nil, err
	}

	return responseFor(stakeholder)
}

// UpdateStakeholder updates descriptive and lifecycle fields. The person and
// the stakeholder type are intentionally immutable.
func UpdateStakeholder(organisationNumber, number string, request requests.UpdateStakeholderRequest) (*responses.StakeholderResponse, error) {
	if !isValidSharePercentage(request.SharePercentage) {
		return nil, ErrInvalidSharePercentage
	}
	if !isValidDateRange(request.ValidFrom, request.ValidTo) {
		return nil, ErrInvalidDateRange
	}

	var saved models.Stakeholder
	if err := stakeholderRepositories.WithOwnershipTransaction(organisationNumber, func(tx *gorm.DB) error {
		stakeholder, err := stakeholderRepositories.FindByNumberAndOrganisation(number, organisationNumber, tx)
		if err != nil {
			return ErrStakeholderNotFound
		}

		now := time.Now()
		wasEffective := isEffectiveOwner(*stakeholder, now)

		stakeholder.Title = strings.TrimSpace(request.Title)
		stakeholder.SharePercentage = request.SharePercentage
		stakeholder.Shares = request.Shares
		stakeholder.ValidFrom = request.ValidFrom
		stakeholder.ValidTo = request.ValidTo

		// A validity-window change can remove the final effective owner without
		// changing status (the verified Stage 0C bypass). Apply the same
		// invariant as an ACTIVE -> INACTIVE/ENDED transition.
		if wasEffective && !isEffectiveOwner(*stakeholder, now) {
			effectiveOwners, err := stakeholderRepositories.CountEffectiveOwners(organisationNumber, now, tx)
			if err != nil {
				return err
			}
			if effectiveOwners <= 1 {
				return ErrLastActiveOwner
			}
		}

		if err := stakeholderRepositories.Save(stakeholder, tx); err != nil {
			return err
		}
		saved = *stakeholder
		return nil
	}); err != nil {
		return nil, err
	}

	return responseFor(saved)
}

// SetStakeholderStatus applies the ACTIVE / INACTIVE / ENDED lifecycle.
//
// Reactivating re-checks the single-ACTIVE rule. Moving to ENDED stamps
// ValidTo when it is not already set; the relationship is retained for history.
// Changing a stakeholder's status never changes the user's authentication
// status (the two axes are independent).
func SetStakeholderStatus(organisationNumber, number, status string) (*responses.StakeholderResponse, error) {
	if !isValidStakeholderStatus(status) {
		return nil, ErrInvalidStatus
	}

	var saved models.Stakeholder
	if err := stakeholderRepositories.WithOwnershipTransaction(organisationNumber, func(tx *gorm.DB) error {
		stakeholder, err := stakeholderRepositories.FindByNumberAndOrganisation(number, organisationNumber, tx)
		if err != nil {
			return ErrStakeholderNotFound
		}

		if status == constants.StatusActive {
			existing, err := stakeholderRepositories.FindActiveByOrgUserType(organisationNumber, stakeholder.UserNumber, stakeholder.StakeholderType, tx)
			if err != nil {
				return err
			}
			if existing != nil && existing.Number != stakeholder.Number {
				return ErrDuplicateActiveStakeholder
			}
		}

		// Foundation C4 / Stage 0C-1: last-owner invariant using the ratified
		// effective-owner predicate (status + validity window). A row that is not
		// currently effective does not reduce effective ownership, so moving it
		// to INACTIVE/ENDED is not blocked.
		now := time.Now()
		if isEffectiveOwner(*stakeholder, now) &&
			(status == constants.StatusInactive || status == StakeholderStatusEnded) {

			effectiveOwners, err := stakeholderRepositories.CountEffectiveOwners(organisationNumber, now, tx)
			if err != nil {
				return err
			}
			if isFinalActiveOwner(stakeholder.StakeholderType, stakeholder.Status, status, effectiveOwners) {
				return ErrLastActiveOwner
			}
		}

		stakeholder.Status = status
		if status == StakeholderStatusEnded && stakeholder.ValidTo == nil {
			endedAt := now
			stakeholder.ValidTo = &endedAt
		}

		if err := stakeholderRepositories.Save(stakeholder, tx); err != nil {
			return err
		}
		saved = *stakeholder
		return nil
	}); err != nil {
		return nil, err
	}

	return responseFor(saved)
}

// CurrentOrganisationNumber is re-exported so callers depend on one resolver.
var CurrentOrganisationNumber = organisationServices.CurrentOrganisationNumber

func responseFor(stakeholder models.Stakeholder) (*responses.StakeholderResponse, error) {
	canAuthenticate, err := authenticatableUserNumbers([]string{stakeholder.UserNumber})
	if err != nil {
		return nil, err
	}
	displayNames, err := stakeholderRepositories.FindUserDisplayNames([]string{stakeholder.UserNumber})
	if err != nil {
		return nil, err
	}
	response := toResponse(
		stakeholder,
		canAuthenticate[stakeholder.UserNumber],
		displayNames[stakeholder.UserNumber],
	)
	return &response, nil
}

func toResponse(stakeholder models.Stakeholder, canAuthenticate bool, displayName string) responses.StakeholderResponse {
	return responses.StakeholderResponse{
		Number:             stakeholder.Number,
		OrganisationNumber: stakeholder.OrganisationNumber,
		UserNumber:         stakeholder.UserNumber,
		DisplayName:        displayName,
		StakeholderType:    stakeholder.StakeholderType,
		Title:              stakeholder.Title,
		SharePercentage:    stakeholder.SharePercentage,
		Shares:             stakeholder.Shares,
		ValidFrom:          stakeholder.ValidFrom,
		ValidTo:            stakeholder.ValidTo,
		Status:             stakeholder.Status,
		CanAuthenticate:    canAuthenticate,
		AddedByNumber:      stakeholder.AddedByNumber,
		CreatedAt:          stakeholder.CreatedAt,
		UpdatedAt:          stakeholder.UpdatedAt,
	}
}
