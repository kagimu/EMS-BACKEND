package services

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/addons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/licensing/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/licensing/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/licensing/models"
	licenceRepositories "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/licensing/repositories"
	eosRepositories "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/eos/repositories"
	eosServices "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/eos/services"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

var (
	ErrNoActiveOrganisation    = errors.New("no active organisation is set for this user")
	ErrOrganisationNotFound    = errors.New("organisation not found")
	ErrLicenceNotFound         = errors.New("licence not found")
	ErrActiveLicenceExists     = errors.New("an active licence already exists for this organisation")
	ErrInvalidValidity         = errors.New("invalid validity period")
	ErrInvalidMaxInstallations = errors.New("max_installations must be at least 1")
	ErrUnknownAddon            = errors.New("unknown add-on code")
	ErrInvalidStatus           = errors.New("invalid status")
	ErrInvalidTransition       = errors.New("licence status does not permit this operation")
	ErrValidToNotExtended      = errors.New("renewal must extend the validity period")
)

// currentOrganisationNumber resolves the authoritative organisation for the
// authenticated user through the existing eos mechanism (User.Number ->
// ActiveUserEntity -> Entity.Number). No fallback organisation is invented.
func currentOrganisationNumber(userNumber string) (string, error) {
	entityService := eosServices.EntityService{}
	active, err := entityService.GetCurrentEntity(userNumber)
	if err != nil {
		return "", err
	}
	if active.Number == "" {
		return "", ErrNoActiveOrganisation
	}
	return active.Number, nil
}

// OrganisationNumberForUser resolves the authenticated user's organisation
// through the existing eos mechanism. It is the exported counterpart used by
// the licence enforcement middleware; licensing decisions never trust a
// client-supplied organisation number.
func OrganisationNumberForUser(userNumber string) (string, error) {
	return currentOrganisationNumber(userNumber)
}

// EffectiveStatus derives the read-time licence status. Derived states are
// never persisted. A future-dated licence (now < valid_from) is a window
// condition: it reports ACTIVE here while access predicates deny it.
func EffectiveStatus(licence models.Licence, now time.Time) string {
	switch licence.Status {
	case models.StatusSuspended:
		return models.EffectiveStatusSuspended
	case models.StatusRevoked:
		return models.EffectiveStatusRevoked
	}

	if now.Before(licence.ValidTo) {
		return models.EffectiveStatusActive
	}
	if now.Before(graceEnd(licence)) {
		return models.EffectiveStatusInGrace
	}
	return models.EffectiveStatusExpired
}

// LicencePermitsAccess reports whether the licence permits module access at
// the given time: stored ACTIVE, within [valid_from, valid_to + grace).
func LicencePermitsAccess(licence models.Licence, now time.Time) bool {
	if licence.Status != models.StatusActive {
		return false
	}
	if now.Before(licence.ValidFrom) {
		return false
	}
	return now.Before(graceEnd(licence))
}

// LicenceContainsAddon reports whether the licence entitlements include the
// add-on code.
func LicenceContainsAddon(licence models.Licence, addonCode string) bool {
	for _, code := range decodeAddons(licence.LicensedAddons) {
		if code == addonCode {
			return true
		}
	}
	return false
}

// NormaliseAddons trims, de-duplicates and validates entitlement codes
// against the existing module catalogue. It never duplicates the catalogue.
func NormaliseAddons(codes []string) ([]string, error) {
	seen := make(map[string]bool)
	result := make([]string, 0, len(codes))
	for _, code := range codes {
		code = strings.TrimSpace(code)
		if code == "" {
			continue
		}
		if _, ok := addons.FindModuleDef(code); !ok {
			return nil, ErrUnknownAddon
		}
		if !seen[code] {
			seen[code] = true
			result = append(result, code)
		}
	}
	return result, nil
}

// generateLicenceKey returns an opaque, cryptographically random key that is
// not derived from any organisation information (F-Q1).
func generateLicenceKey() (string, error) {
	buffer := make([]byte, 32)
	if _, err := rand.Read(buffer); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buffer), nil
}

func encodeAddons(codes []string) (string, error) {
	normalised, err := NormaliseAddons(codes)
	if err != nil {
		return "", err
	}
	raw, err := json.Marshal(normalised)
	if err != nil {
		return "", err
	}
	return string(raw), nil
}

func decodeAddons(raw string) []string {
	var codes []string
	if err := json.Unmarshal([]byte(raw), &codes); err != nil {
		return []string{}
	}
	return codes
}

func graceEnd(licence models.Licence) time.Time {
	return licence.ValidTo.Add(time.Duration(licence.GraceDays) * 24 * time.Hour)
}

func validateValidity(validFrom, validTo time.Time) error {
	if validFrom.IsZero() || validTo.IsZero() {
		return ErrInvalidValidity
	}
	if !validTo.After(validFrom) {
		return ErrInvalidValidity
	}
	return nil
}

func validStoredStatus(status string) bool {
	return status == models.StatusActive || status == models.StatusSuspended || status == models.StatusRevoked
}

// IssueLicence creates and assigns a licence with a server-generated opaque
// key. F-Q7 is enforced by a transaction plus the unique active_org_number
// guard, so concurrent issuance cannot leave two ACTIVE licences.
func IssueLicence(actorUserNumber string, request requests.IssueLicenceRequest) (*responses.IssuedLicenceResponse, error) {
	organisationNumber := strings.TrimSpace(request.OrganisationNumber)
	if organisationNumber == "" {
		return nil, ErrOrganisationNotFound
	}
	entity, err := eosRepositories.FindEntityByNumber(organisationNumber)
	if err != nil {
		return nil, err
	}
	if entity.Number == "" {
		// FindEntityByNumber returns an empty entity (nil error) when the
		// organisation does not exist.
		return nil, ErrOrganisationNotFound
	}
	if err := validateValidity(request.ValidFrom, request.ValidTo); err != nil {
		return nil, err
	}
	if request.MaxInstallations < 1 {
		return nil, ErrInvalidMaxInstallations
	}
	encodedAddons, err := encodeAddons(request.LicensedAddons)
	if err != nil {
		return nil, err
	}
	licenceKey, err := generateLicenceKey()
	if err != nil {
		return nil, err
	}

	activeOrg := organisationNumber
	licence := models.Licence{
		Number:             uuid.NewV4().String(),
		LicenceKey:         licenceKey,
		OrganisationNumber: organisationNumber,
		ValidFrom:          request.ValidFrom,
		ValidTo:            request.ValidTo,
		GraceDays:          models.DefaultGraceDays,
		MaxInstallations:   request.MaxInstallations,
		LicensedAddons:     encodedAddons,
		Status:             models.StatusActive,
		ActiveOrgNumber:    &activeOrg,
		IssuedBy:           actorUserNumber,
		LastActedBy:        &actorUserNumber,
	}

	err = config.DB.Transaction(func(tx *gorm.DB) error {
		if existing, err := licenceRepositories.FindActiveByOrganisation(organisationNumber, tx); err != nil {
			return err
		} else if existing != nil {
			return ErrActiveLicenceExists
		}
		return licenceRepositories.Create(&licence, tx)
	})
	if err != nil {
		// The unique active_org_number index is the concurrency backstop: when
		// two issuances race, the loser is reported as an active-licence
		// conflict rather than a raw database error.
		if existing, checkErr := licenceRepositories.FindActiveByOrganisation(organisationNumber); checkErr == nil && existing != nil {
			return nil, ErrActiveLicenceExists
		}
		return nil, err
	}

	var updated models.Licence
	updated = licence

	// Record licence issuance event for Foundation G audit
	RecordLicenceIssued(&updated, actorUserNumber)

	response := toIssuedResponse(updated, time.Now())
	return &response, nil
}

// ListLicences returns licences filtered by organisation and/or stored status.
func ListLicences(organisationNumber, status string) ([]responses.LicenceResponse, error) {
	if status != "" && !validStoredStatus(status) {
		return nil, ErrInvalidStatus
	}
	licences, err := licenceRepositories.FindAll(strings.TrimSpace(organisationNumber), status)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	result := make([]responses.LicenceResponse, 0, len(licences))
	for _, licence := range licences {
		result = append(result, toResponse(licence, now))
	}
	return result, nil
}

// GetLicence returns one licence, or ErrLicenceNotFound.
func GetLicence(number string) (*responses.LicenceResponse, error) {
	licence, err := licenceRepositories.FindByNumber(number)
	if err != nil {
		return nil, err
	}
	if licence == nil {
		return nil, ErrLicenceNotFound
	}
	response := toResponse(*licence, time.Now())
	return &response, nil
}

// ConfigureLicence updates validity, installation allowance and entitlements
// of an ACTIVE licence. Renewal, suspension and revocation have their own
// operations.
func ConfigureLicence(actorUserNumber, number string, request requests.ConfigureLicenceRequest) (*responses.LicenceResponse, error) {
	var updated models.Licence
	var beforeState map[string]interface{}
	err := config.DB.Transaction(func(tx *gorm.DB) error {
		licence, err := licenceRepositories.FindByNumberForUpdate(number, tx)
		if err != nil {
			return err
		}
		if licence == nil {
			return ErrLicenceNotFound
		}
		if licence.Status != models.StatusActive {
			return ErrInvalidTransition
		}

		beforeState = map[string]interface{}{
			"valid_from":        licence.ValidFrom,
			"valid_to":          licence.ValidTo,
			"max_installations": licence.MaxInstallations,
			"licensed_addons":   licence.LicensedAddons,
		}

		validFrom := licence.ValidFrom
		if request.ValidFrom != nil {
			validFrom = *request.ValidFrom
		}
		validTo := licence.ValidTo
		if request.ValidTo != nil {
			validTo = *request.ValidTo
		}
		if err := validateValidity(validFrom, validTo); err != nil {
			return err
		}

		if request.MaxInstallations != nil {
			if *request.MaxInstallations < 1 {
				return ErrInvalidMaxInstallations
			}
			licence.MaxInstallations = *request.MaxInstallations
		}
		if request.LicensedAddons != nil {
			encoded, err := encodeAddons(*request.LicensedAddons)
			if err != nil {
				return err
			}
			licence.LicensedAddons = encoded
		}

		licence.ValidFrom = validFrom
		licence.ValidTo = validTo
		licence.LastActedBy = &actorUserNumber

		if err := licenceRepositories.Save(licence, tx); err != nil {
			return err
		}
		updated = *licence
		return nil
	})
	if err != nil {
		return nil, err
	}

	// Record durable configuration fact for Foundation G (§18).
	RecordLicenceConfigured(&updated, actorUserNumber, beforeState, map[string]interface{}{
		"valid_from":        updated.ValidFrom,
		"valid_to":          updated.ValidTo,
		"max_installations": updated.MaxInstallations,
		"licensed_addons":   updated.LicensedAddons,
	})

	response := toResponse(updated, time.Now())
	return &response, nil
}

// RenewLicence extends the validity of an ACTIVE licence, preserving the
// licence number, key, organisation, entitlements and installation allowance
// (F-Q6). It never creates a new licence identity.
func RenewLicence(actorUserNumber, number string, request requests.RenewLicenceRequest) (*responses.LicenceResponse, error) {
	var updated models.Licence
	var previousValidTo time.Time
	err := config.DB.Transaction(func(tx *gorm.DB) error {
		licence, err := licenceRepositories.FindByNumberForUpdate(number, tx)
		if err != nil {
			return err
		}
		if licence == nil {
			return ErrLicenceNotFound
		}
		if licence.Status != models.StatusActive {
			return ErrInvalidTransition
		}
		if !request.ValidTo.After(licence.ValidTo) || !request.ValidTo.After(licence.ValidFrom) {
			return ErrValidToNotExtended
		}
		previousValidTo = licence.ValidTo
		licence.ValidTo = request.ValidTo
		licence.LastActedBy = &actorUserNumber
		if err := licenceRepositories.Save(licence, tx); err != nil {
			return err
		}
		updated = *licence
		return nil
	})
	if err != nil {
		return nil, err
	}

	// Record durable renewal fact for Foundation G (§18). Ordinary renewal
	// preserves the licence number, key, organisation, installations and
	// entitlements; only validity is extended.
	RecordLicenceRenewed(&updated, actorUserNumber, previousValidTo)

	response := toResponse(updated, time.Now())
	return &response, nil
}

// ReplaceLicence issues a brand-new licence that supersedes an ACTIVE licence.
// The operation is atomic: the superseded licence is revoked (clearing the
// one-active guard), every capacity-holding installation of the old licence
// (client-deactivated included) is explicitly released with the actor,
// timestamp and mandatory reason, and the new licence becomes the
// organisation's single ACTIVE licence. Installations are never inherited,
// moved or reattached, and capacity is never transferred (F-Q6, F-Q7, §13,
// §24 detail 1).
func ReplaceLicence(actorUserNumber, supersededNumber string, request requests.ReplaceLicenceRequest) (*responses.IssuedLicenceResponse, error) {
	organisationNumber := strings.TrimSpace(request.OrganisationNumber)
	if organisationNumber == "" {
		return nil, ErrOrganisationNotFound
	}
	reason := strings.TrimSpace(request.Reason)
	if reason == "" {
		return nil, ErrReleaseReasonRequired
	}
	entity, err := eosRepositories.FindEntityByNumber(organisationNumber)
	if err != nil {
		return nil, err
	}
	if entity.Number == "" {
		return nil, ErrOrganisationNotFound
	}
	if err := validateValidity(request.ValidFrom, request.ValidTo); err != nil {
		return nil, err
	}
	if request.MaxInstallations < 1 {
		return nil, ErrInvalidMaxInstallations
	}
	encodedAddons, err := encodeAddons(request.LicensedAddons)
	if err != nil {
		return nil, err
	}
	licenceKey, err := generateLicenceKey()
	if err != nil {
		return nil, err
	}

	var created models.Licence
	var superseded models.Licence
	err = config.DB.Transaction(func(tx *gorm.DB) error {
		old, err := licenceRepositories.FindByNumberForUpdate(supersededNumber, tx)
		if err != nil {
			return err
		}
		if old == nil {
			return ErrLicenceNotFound
		}
		if old.Status != models.StatusActive {
			return ErrInvalidTransition
		}
		superseded = *old

		now := time.Now()

		// 1. Explicitly release every capacity-holding installation of the
		// superseded licence. Client-deactivated installations still consume
		// capacity, so they are released here too. The licence association is
		// never edited; the records remain historical.
		if _, err := licenceRepositories.ReleaseActiveByLicence(
			old.Number, actorUserNumber, reason, now, tx); err != nil {
			return err
		}

		// 2. Supersede the old licence. It is revoked and its one-active guard
		// cleared INSIDE the transaction so the unique active_org_number index
		// admits the new active row without ever allowing two active licences.
		newNumber := uuid.NewV4().String()
		old.Status = models.StatusRevoked
		old.ActiveOrgNumber = nil
		old.ReplacedByNumber = &newNumber
		old.LastActedBy = &actorUserNumber
		if err := licenceRepositories.Save(old, tx); err != nil {
			return err
		}

		// 3. Create the new ACTIVE licence with its own opaque key, validity,
		// entitlements and max_installations. No capacity is transferred.
		activeOrg := organisationNumber
		newLicence := models.Licence{
			Number:             newNumber,
			LicenceKey:         licenceKey,
			OrganisationNumber: organisationNumber,
			ValidFrom:          request.ValidFrom,
			ValidTo:            request.ValidTo,
			GraceDays:          models.DefaultGraceDays,
			MaxInstallations:   request.MaxInstallations,
			LicensedAddons:     encodedAddons,
			Status:             models.StatusActive,
			ActiveOrgNumber:    &activeOrg,
			IssuedBy:           actorUserNumber,
			LastActedBy:        &actorUserNumber,
		}
		if err := licenceRepositories.Create(&newLicence, tx); err != nil {
			return err
		}
		created = newLicence
		return nil
	})
	if err != nil {
		// Preserve explicit domain failures. A duplicate active_org_number from
		// a concurrent issuance is reported as the active-licence conflict.
		if errors.Is(err, ErrLicenceNotFound) || errors.Is(err, ErrInvalidTransition) ||
			errors.Is(err, ErrOrganisationNotFound) || errors.Is(err, ErrInvalidValidity) ||
			errors.Is(err, ErrInvalidMaxInstallations) || errors.Is(err, ErrUnknownAddon) ||
			errors.Is(err, ErrReleaseReasonRequired) {
			return nil, err
		}
		if existing, checkErr := licenceRepositories.FindActiveByOrganisation(organisationNumber); checkErr == nil && existing != nil {
			return nil, ErrActiveLicenceExists
		}
		return nil, err
	}

	// Record durable supersession fact for Foundation G (§18). Installations
	// are never transferred; the fact captures the old/new licence identities.
	RecordLicenceReplaced(&superseded, &created, actorUserNumber, reason)

	response := toIssuedResponse(created, time.Now())
	return &response, nil
}

// SuspendLicence suspends an ACTIVE licence and clears the one-active guard.
func SuspendLicence(actorUserNumber, number string) (*responses.LicenceResponse, error) {
	return transitionLicence(actorUserNumber, number, models.LicenceEventSuspended, func(licence *models.Licence) error {
		if licence.Status != models.StatusActive {
			return ErrInvalidTransition
		}
		licence.Status = models.StatusSuspended
		licence.ActiveOrgNumber = nil
		return nil
	})
}

// RevokeLicence revokes an ACTIVE or SUSPENDED licence and clears the
// one-active guard. Revocation is terminal.
func RevokeLicence(actorUserNumber, number string) (*responses.LicenceResponse, error) {
	return transitionLicence(actorUserNumber, number, models.LicenceEventRevoked, func(licence *models.Licence) error {
		if licence.Status != models.StatusActive && licence.Status != models.StatusSuspended {
			return ErrInvalidTransition
		}
		licence.Status = models.StatusRevoked
		licence.ActiveOrgNumber = nil
		return nil
	})
}

func transitionLicence(actorUserNumber, number string, eventType models.LicenceEventType, mutate func(*models.Licence) error) (*responses.LicenceResponse, error) {
	var updated models.Licence
	var previousStatus string
	err := config.DB.Transaction(func(tx *gorm.DB) error {
		licence, err := licenceRepositories.FindByNumberForUpdate(number, tx)
		if err != nil {
			return err
		}
		if licence == nil {
			return ErrLicenceNotFound
		}
		previousStatus = licence.Status
		if err := mutate(licence); err != nil {
			return err
		}
		licence.LastActedBy = &actorUserNumber
		if err := licenceRepositories.Save(licence, tx); err != nil {
			return err
		}
		updated = *licence
		return nil
	})
	if err != nil {
		return nil, err
	}
	// Record durable status-transition fact for Foundation G (§18).
	RecordLicenceTransition(eventType, &updated, actorUserNumber, previousStatus)
	response := toResponse(updated, time.Now())
	return &response, nil
}

// ClientLicenceStatus returns the authenticated organisation's latest licence
// for the client plane, so SUSPENDED/REVOKED states remain visible to their
// own organisation. It never exposes another organisation's licence.
func ClientLicenceStatus(userNumber string) (*responses.LicenceStatusResponse, error) {
	organisationNumber, err := currentOrganisationNumber(userNumber)
	if err != nil {
		return nil, err
	}
	licence, err := licenceRepositories.FindLatestByOrganisation(organisationNumber)
	if err != nil {
		return nil, err
	}
	if licence == nil {
		return nil, ErrLicenceNotFound
	}
	response := toLicenceStatusResponse(*licence, time.Now())
	return &response, nil
}

// RequestRenewal records the authenticated organisation's renewal request on
// its latest licence (actor + timestamp). It never changes validity, capacity
// or any other licence state; the control plane performs the renewal
// (Foundation F §18/§19).
func RequestRenewal(userNumber string) (*responses.LicenceStatusResponse, error) {
	organisationNumber, err := currentOrganisationNumber(userNumber)
	if err != nil {
		return nil, err
	}
	licence, err := licenceRepositories.FindLatestByOrganisation(organisationNumber)
	if err != nil {
		return nil, err
	}
	if licence == nil {
		return nil, ErrLicenceNotFound
	}
	if licence.Status == models.StatusRevoked {
		// A revoked licence is terminal; renewal requires a new licence.
		return nil, ErrInvalidTransition
	}
	now := time.Now()
	licence.RenewalRequestedBy = &userNumber
	licence.RenewalRequestedAt = &now
	if err := licenceRepositories.Save(licence); err != nil {
		return nil, err
	}
	// Record durable client renewal-request fact for Foundation G (§18).
	RecordClientRenewalRequested(licence, userNumber)
	response := toLicenceStatusResponse(*licence, now)
	return &response, nil
}

func toLicenceStatusResponse(licence models.Licence, now time.Time) responses.LicenceStatusResponse {
	return responses.LicenceStatusResponse{
		Number:             licence.Number,
		EffectiveStatus:    EffectiveStatus(licence, now),
		ValidFrom:          licence.ValidFrom,
		ValidTo:            licence.ValidTo,
		GraceDays:          licence.GraceDays,
		LicensedAddons:     decodeAddons(licence.LicensedAddons),
		MaxInstallations:   licence.MaxInstallations,
		RenewalRequestedAt: licence.RenewalRequestedAt,
	}
}

// IsAddonLicensed reports whether the organisation's active licence permits
// access AND contains the add-on code.
func IsAddonLicensed(organisationNumber, addonCode string) (bool, error) {
	licence, err := licenceRepositories.FindActiveByOrganisation(organisationNumber)
	if err != nil {
		return false, err
	}
	if licence == nil {
		return false, nil
	}
	if !LicencePermitsAccess(*licence, time.Now()) {
		return false, nil
	}
	return LicenceContainsAddon(*licence, addonCode), nil
}

func toResponse(licence models.Licence, now time.Time) responses.LicenceResponse {
	response := responses.LicenceResponse{
		Number:             licence.Number,
		OrganisationNumber: licence.OrganisationNumber,
		ValidFrom:          licence.ValidFrom,
		ValidTo:            licence.ValidTo,
		GraceDays:          licence.GraceDays,
		MaxInstallations:   licence.MaxInstallations,
		LicensedAddons:     decodeAddons(licence.LicensedAddons),
		Status:             licence.Status,
		EffectiveStatus:    EffectiveStatus(licence, now),
		IssuedBy:           licence.IssuedBy,
		CreatedAt:          licence.CreatedAt,
		UpdatedAt:          licence.UpdatedAt,
	}
	if licence.ReplacedByNumber != nil {
		response.ReplacedByNumber = *licence.ReplacedByNumber
	}
	if licence.LastActedBy != nil {
		response.LastActedBy = *licence.LastActedBy
	}
	if licence.RenewalRequestedBy != nil {
		response.RenewalRequestedBy = *licence.RenewalRequestedBy
	}
	if licence.RenewalRequestedAt != nil {
		response.RenewalRequestedAt = licence.RenewalRequestedAt
	}
	return response
}

func toIssuedResponse(licence models.Licence, now time.Time) responses.IssuedLicenceResponse {
	return responses.IssuedLicenceResponse{
		LicenceResponse: toResponse(licence, now),
		LicenceKey:      licence.LicenceKey,
	}
}
