package services

import (
	"errors"
	"strings"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/licensing/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/licensing/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/licensing/models"
	licensingRepositories "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/licensing/repositories"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

var (
	ErrInvalidLicenceKey              = errors.New("invalid licence key for this organisation")
	ErrInvalidInstallationIdentity    = errors.New("invalid installation identity")
	ErrLicenceNotActive               = errors.New("licence is not active")
	ErrLicenceNotYetValid             = errors.New("licence is not yet valid")
	ErrLicenceExpired                 = errors.New("licence has expired")
	ErrCapacityReached                = errors.New("installation capacity reached")
	ErrAddonNotLicensed               = errors.New("add-on is not licensed")
	ErrInstallationNotFound           = errors.New("installation not found")
	ErrInstallationReleased           = errors.New("installation has been released and cannot be reactivated")
	ErrInstallationNotActive          = errors.New("installation is not active")
	ErrInstallationAlreadyDeactivated = errors.New("installation is already deactivated")
	ErrReleaseReasonRequired          = errors.New("a release reason is required")
	ErrInvalidReplacementLink         = errors.New("invalid replacement installation link")
)

// validInstallationIdentity accepts the opaque, client-generated installation
// identity. It is a recognition value only and never grants authority.
func validInstallationIdentity(identity string) bool {
	if len(identity) < 8 || len(identity) > 128 {
		return false
	}
	return !strings.ContainsAny(identity, " \t\r\n")
}

// installationCountsTowardCapacity reports whether the installation consumes
// licence capacity. Client self-deactivation does NOT free capacity: the row
// stays ACTIVE until an explicit control-plane release (F-Q3A, §24 detail 2).
func installationCountsTowardCapacity(installation models.Installation) bool {
	return installation.Status == models.InstallationStatusActive
}

func validInstallationStatus(status string) bool {
	return status == models.InstallationStatusActive ||
		status == models.InstallationStatusReleased ||
		status == models.InstallationStatusRevoked
}

// activationBlockedError reports why the licence cannot activate an
// installation at the given time. Activation is allowed while ACTIVE or
// IN_GRACE, and fails closed otherwise (F-Q2, §7).
func activationBlockedError(licence models.Licence, now time.Time) error {
	switch {
	case licence.Status != models.StatusActive:
		return ErrLicenceNotActive
	case now.Before(licence.ValidFrom):
		return ErrLicenceNotYetValid
	case !now.Before(graceEnd(licence)):
		return ErrLicenceExpired
	}
	return nil
}

// ActivateInstallation registers (or reactivates) an installation against the
// caller's organisation licence. Capacity allocation is concurrency-safe: the
// licence row is locked for the duration of the transaction, so concurrent
// activations can never exceed max_installations.
func ActivateInstallation(userNumber string, request requests.ActivateInstallationRequest) (*responses.ActivationResponse, error) {
	organisationNumber, err := currentOrganisationNumber(userNumber)
	if err != nil {
		return nil, err
	}

	licenceKey := strings.TrimSpace(request.LicenceKey)
	identity := strings.TrimSpace(request.InstallationIdentity)
	if !validInstallationIdentity(identity) {
		return nil, ErrInvalidInstallationIdentity
	}
	addonCode := strings.TrimSpace(request.AddonCode)
	replacementNumber := strings.TrimSpace(request.ReplacesInstallationNumber)

	licence, err := licensingRepositories.FindByKey(licenceKey)
	if err != nil {
		return nil, err
	}
	if licence == nil || licence.OrganisationNumber != organisationNumber {
		// Key unknown and key belonging to another organisation are reported
		// identically: activation must not disclose another organisation's
		// licensing information.
		return nil, ErrInvalidLicenceKey
	}
	if addonCode != "" && !LicenceContainsAddon(*licence, addonCode) {
		return nil, ErrAddonNotLicensed
	}
	if err := activationBlockedError(*licence, time.Now()); err != nil {
		return nil, err
	}

	var activated models.Installation
	var reactivated bool
	err = config.DB.Transaction(func(tx *gorm.DB) error {
		locked, err := licensingRepositories.FindByNumberForUpdate(licence.Number, tx)
		if err != nil {
			return err
		}
		if locked == nil {
			return ErrLicenceNotFound
		}
		if addonCode != "" && !LicenceContainsAddon(*locked, addonCode) {
			return ErrAddonNotLicensed
		}
		if err := activationBlockedError(*locked, time.Now()); err != nil {
			return err
		}

		existing, err := licensingRepositories.FindLatestInstallationByIdentity(locked.Number, identity, tx)
		if err != nil {
			return err
		}
		if existing != nil {
			switch existing.Status {
			case models.InstallationStatusReleased, models.InstallationStatusRevoked:
				// A released/revoked identity is not silently re-permitted;
				// history stays bound to the old record (fail safely).
				return ErrInstallationReleased
			}
			// Same identity: idempotent activation, or reactivation after a
			// permissible client self-deactivation. Both are capacity-neutral.
			now := time.Now()
			if existing.ClientDeactivatedAt != nil {
				existing.ClientDeactivatedAt = nil
				existing.ClientDeactivatedBy = nil
				reactivated = true
			}
			existing.LastSeenAt = &now
			if err := licensingRepositories.SaveInstallation(existing, tx); err != nil {
				return err
			}
			activated = *existing
			return nil
		}

		capacity, err := licensingRepositories.CountCapacityHolding(locked.Number, tx)
		if err != nil {
			return err
		}
		if capacity >= int64(locked.MaxInstallations) {
			return ErrCapacityReached
		}

		now := time.Now()
		actor := userNumber
		installation := models.Installation{
			Number:               uuid.NewV4().String(),
			InstallationIdentity: identity,
			OrganisationNumber:   organisationNumber,
			LicenceNumber:        locked.Number,
			ActivatedAt:          now,
			ActivatedBy:          &actor,
			LastSeenAt:           &now,
			Status:               models.InstallationStatusActive,
		}
		if err := licensingRepositories.CreateInstallation(&installation, tx); err != nil {
			return err
		}

		if replacementNumber != "" {
			old, err := licensingRepositories.FindInstallationByNumberForUpdate(replacementNumber, tx)
			if err != nil {
				return err
			}
			// The link is history only: same licence, same organisation,
			// owner-released, and not already linked. No record is transferred
			// or rewritten beyond the forward link.
			if old == nil || old.LicenceNumber != locked.Number || old.OrganisationNumber != organisationNumber ||
				old.Status != models.InstallationStatusReleased || old.ReplacedByNumber != nil {
				return ErrInvalidReplacementLink
			}
			old.ReplacedByNumber = &installation.Number
			if err := licensingRepositories.SaveInstallation(old, tx); err != nil {
				return err
			}
		}

		activated = installation
		return nil
	})
	if err != nil {
		return nil, err
	}

	current, err := licensingRepositories.FindByNumber(licence.Number)
	if err != nil {
		return nil, err
	}
	if current == nil {
		current = licence
	}
	capacityUsed, err := licensingRepositories.CountCapacityHolding(current.Number)
	if err != nil {
		return nil, err
	}
	// Record durable installation fact for Foundation G (§17/§18).
	if reactivated {
		RecordInstallationReactivated(&activated, userNumber)
	} else {
		RecordInstallationActivated(&activated, userNumber)
	}
	response := toActivationResponse(activated, *current, capacityUsed)
	return &response, nil
}

// ListClientInstallations returns the authenticated organisation's
// installations. It never reaches another organisation.
func ListClientInstallations(userNumber string) ([]responses.InstallationResponse, error) {
	organisationNumber, err := currentOrganisationNumber(userNumber)
	if err != nil {
		return nil, err
	}
	installations, err := licensingRepositories.FindInstallations("", organisationNumber, "")
	if err != nil {
		return nil, err
	}
	result := make([]responses.InstallationResponse, 0, len(installations))
	for _, installation := range installations {
		result = append(result, toInstallationResponse(installation))
	}
	return result, nil
}

// DeactivateOwnInstallation is the permissible client self-deactivation: it
// terminates access, records the client actor and timestamp, and holds the
// capacity slot until an explicit control-plane release (§24 detail 2).
func DeactivateOwnInstallation(userNumber, installationNumber string) (*responses.InstallationResponse, error) {
	organisationNumber, err := currentOrganisationNumber(userNumber)
	if err != nil {
		return nil, err
	}
	installation, err := licensingRepositories.FindInstallationByNumber(installationNumber)
	if err != nil {
		return nil, err
	}
	if installation == nil || installation.OrganisationNumber != organisationNumber {
		return nil, ErrInstallationNotFound
	}
	if installation.Status != models.InstallationStatusActive {
		return nil, ErrInstallationNotActive
	}
	if installation.ClientDeactivatedAt != nil {
		return nil, ErrInstallationAlreadyDeactivated
	}

	now := time.Now()
	actor := userNumber
	installation.ClientDeactivatedAt = &now
	installation.ClientDeactivatedBy = &actor
	if err := licensingRepositories.SaveInstallation(installation); err != nil {
		return nil, err
	}
	// Record durable client-deactivation fact for Foundation G (§17/§18).
	RecordInstallationDeactivated(installation, userNumber)
	response := toInstallationResponse(*installation)
	return &response, nil
}

// ListInstallations returns installations for the control plane, optionally
// filtered (deliberately global).
func ListInstallations(licenceNumber, organisationNumber, status string) ([]responses.InstallationResponse, error) {
	if status != "" && !validInstallationStatus(status) {
		return nil, ErrInvalidStatus
	}
	installations, err := licensingRepositories.FindInstallations(
		strings.TrimSpace(licenceNumber), strings.TrimSpace(organisationNumber), status)
	if err != nil {
		return nil, err
	}
	result := make([]responses.InstallationResponse, 0, len(installations))
	for _, installation := range installations {
		result = append(result, toInstallationResponse(installation))
	}
	return result, nil
}

// GetInstallation returns one installation for the control plane.
func GetInstallation(number string) (*responses.InstallationResponse, error) {
	installation, err := licensingRepositories.FindInstallationByNumber(number)
	if err != nil {
		return nil, err
	}
	if installation == nil {
		return nil, ErrInstallationNotFound
	}
	response := toInstallationResponse(*installation)
	return &response, nil
}

// ReleaseInstallation is the explicit control-plane release — the only action
// that makes installation capacity available again (F-Q3A). It records the
// actor, timestamp and reason.
func ReleaseInstallation(actorUserNumber, installationNumber, reason string) (*responses.InstallationResponse, error) {
	reason = strings.TrimSpace(reason)
	if reason == "" {
		return nil, ErrReleaseReasonRequired
	}

	var released models.Installation
	err := config.DB.Transaction(func(tx *gorm.DB) error {
		installation, err := licensingRepositories.FindInstallationByNumberForUpdate(installationNumber, tx)
		if err != nil {
			return err
		}
		if installation == nil {
			return ErrInstallationNotFound
		}
		if installation.Status != models.InstallationStatusActive {
			return ErrInstallationNotActive
		}

		now := time.Now()
		actor := actorUserNumber
		installation.Status = models.InstallationStatusReleased
		installation.ReleasedBy = &actor
		installation.ReleasedAt = &now
		installation.ReleaseReason = &reason
		if err := licensingRepositories.SaveInstallation(installation, tx); err != nil {
			return err
		}
		released = *installation
		return nil
	})
	if err != nil {
		return nil, err
	}
	// Record durable owner-release fact for Foundation G (§17/§18).
	RecordInstallationReleased(&released, actorUserNumber, reason)
	response := toInstallationResponse(released)
	return &response, nil
}

func toInstallationResponse(installation models.Installation) responses.InstallationResponse {
	response := responses.InstallationResponse{
		Number:               installation.Number,
		InstallationIdentity: installation.InstallationIdentity,
		OrganisationNumber:   installation.OrganisationNumber,
		LicenceNumber:        installation.LicenceNumber,
		Status:               installation.Status,
		ActivatedAt:          installation.ActivatedAt,
		LastSeenAt:           installation.LastSeenAt,
		ClientDeactivatedAt:  installation.ClientDeactivatedAt,
		CreatedAt:            installation.CreatedAt,
		UpdatedAt:            installation.UpdatedAt,
	}
	if installation.ActivatedBy != nil {
		response.ActivatedBy = *installation.ActivatedBy
	}
	if installation.ClientDeactivatedBy != nil {
		response.ClientDeactivatedBy = *installation.ClientDeactivatedBy
	}
	if installation.ReleasedBy != nil {
		response.ReleasedBy = *installation.ReleasedBy
	}
	if installation.ReleasedAt != nil {
		response.ReleasedAt = installation.ReleasedAt
	}
	if installation.ReleaseReason != nil {
		response.ReleaseReason = *installation.ReleaseReason
	}
	if installation.ReplacedByNumber != nil {
		response.ReplacedByNumber = *installation.ReplacedByNumber
	}
	return response
}

func toActivationResponse(installation models.Installation, licence models.Licence, capacityUsed int64) responses.ActivationResponse {
	return responses.ActivationResponse{
		InstallationNumber:   installation.Number,
		InstallationIdentity: installation.InstallationIdentity,
		OrganisationNumber:   installation.OrganisationNumber,
		LicenceNumber:        installation.LicenceNumber,
		Status:               installation.Status,
		LicenceStatus:        licence.Status,
		EffectiveStatus:      EffectiveStatus(licence, time.Now()),
		ActivatedAt:          installation.ActivatedAt,
		LastSeenAt:           installation.LastSeenAt,
		LicensedAddons:       decodeAddons(licence.LicensedAddons),
		MaxInstallations:     licence.MaxInstallations,
		CapacityUsed:         int(capacityUsed),
	}
}
