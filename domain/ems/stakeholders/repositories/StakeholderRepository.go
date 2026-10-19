package repositories

import (
	"errors"
	"strings"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/stakeholders/models"
	"gorm.io/gorm"
)

// FindAllByOrganisation returns stakeholders of one organisation, optionally
// filtered by stakeholder type and/or status, ordered by type then title.
func FindAllByOrganisation(organisationNumber, stakeholderType, status string) ([]models.Stakeholder, error) {
	var stakeholders []models.Stakeholder

	query := config.DB.Where("organisation_number = ?", organisationNumber)
	if stakeholderType != "" {
		query = query.Where("stakeholder_type = ?", stakeholderType)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Order("stakeholder_type asc, title asc").Find(&stakeholders).Error; err != nil {
		return nil, err
	}

	return stakeholders, nil
}

// FindByNumberAndOrganisation returns a stakeholder only when it belongs to the
// given organisation, which prevents cross-organisation (IDOR) disclosure.
func FindByNumberAndOrganisation(number, organisationNumber string, tx ...*gorm.DB) (*models.Stakeholder, error) {
	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}

	var stakeholder models.Stakeholder
	if err := db.
		Where("number = ? AND organisation_number = ?", number, organisationNumber).
		First(&stakeholder).Error; err != nil {
		return nil, err
	}
	return &stakeholder, nil
}

// FindActiveByOrgUserType returns the ACTIVE relationship for a person and type
// within an organisation, or (nil, nil) when none exists.
func FindActiveByOrgUserType(organisationNumber, userNumber, stakeholderType string, tx ...*gorm.DB) (*models.Stakeholder, error) {
	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}

	var stakeholder models.Stakeholder

	err := db.
		Where("organisation_number = ? AND user_number = ? AND stakeholder_type = ? AND status = ?",
			organisationNumber, userNumber, stakeholderType, constants.StatusActive).
		First(&stakeholder).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &stakeholder, nil
}

// CountEffectiveOwners counts the organisation's effective owners at now.
//
// This is the ratified effective-owner predicate (Foundation C/E2, D10):
// stakeholder_type='owner', status='ACTIVE', and the validity window contains
// now (valid_to inclusive; NULL on either side is unbounded). It is the only
// predicate used by the last-owner invariant, so an ACTIVE but expired or
// future-dated owner never satisfies it.
func CountEffectiveOwners(organisationNumber string, now time.Time, tx ...*gorm.DB) (int64, error) {
	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}

	var count int64
	if err := db.Model(&models.Stakeholder{}).
		Where("organisation_number = ? AND stakeholder_type = ? AND status = ?",
			organisationNumber, models.StakeholderTypeOwner, constants.StatusActive).
		Where("(valid_from IS NULL OR valid_from <= ?)", now).
		Where("(valid_to IS NULL OR valid_to >= ?)", now).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func Create(stakeholder *models.Stakeholder, tx ...*gorm.DB) error {
	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}
	return db.Create(stakeholder).Error
}

func Save(stakeholder *models.Stakeholder, tx ...*gorm.DB) error {
	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}
	return db.Save(stakeholder).Error
}

// FindUserDisplayNames returns the display name for each supplied user number in
// one batched query. It is used only for the organisation-scoped stakeholder
// display_name field and returns no contact or credential data.
func FindUserDisplayNames(userNumbers []string) (map[string]string, error) {
	result := make(map[string]string)
	unique := make([]string, 0, len(userNumbers))
	seen := make(map[string]bool, len(userNumbers))
	for _, number := range userNumbers {
		if number == "" || seen[number] {
			continue
		}
		seen[number] = true
		unique = append(unique, number)
	}
	if len(unique) == 0 {
		return result, nil
	}

	var users []authModel.User
	if err := config.DB.Where("number IN ?", unique).Find(&users).Error; err != nil {
		return nil, err
	}
	for _, user := range users {
		parts := make([]string, 0, 3)
		for _, part := range []string{user.FirstName, user.SecondName, user.LastName} {
			if trimmed := strings.TrimSpace(part); trimmed != "" {
				parts = append(parts, trimmed)
			}
		}
		result[user.Number] = strings.Join(parts, " ")
	}
	return result, nil
}

// FindAuthenticatableUserNumbers reports which of the supplied user numbers
// currently have a usable login, defined as: an Email contact AND a password
// row AND an account status of "Active".
//
// Foundation C2.3: this is a read-only derivation used for the optional
// can_authenticate response field. It performs three batched queries (never one
// per user), exposes no contact or credential data, and does not change
// authentication behaviour in any way.
//
// "Active" is the Title Case literal used by the existing login flow
// (see AuthController); it is intentionally NOT the EMS constants.StatusActive
// ("ACTIVE"), which is a different axis.
func FindAuthenticatableUserNumbers(userNumbers []string) (map[string]bool, error) {
	result := make(map[string]bool)
	if len(userNumbers) == 0 {
		return result, nil
	}

	const authActiveStatus = "Active"

	activeUsers := make(map[string]bool)
	var users []authModel.User
	if err := config.DB.Where("number IN ?", userNumbers).Find(&users).Error; err != nil {
		return nil, err
	}
	for _, user := range users {
		if user.Status == authActiveStatus {
			activeUsers[user.Number] = true
		}
	}

	hasPassword := make(map[string]bool)
	var passwords []authModel.UserPassword
	if err := config.DB.Where("user_number IN ?", userNumbers).Where("deleted_at IS NULL").Find(&passwords).Error; err != nil {
		return nil, err
	}
	for _, password := range passwords {
		hasPassword[password.UserNumber] = true
	}

	hasEmail := make(map[string]bool)
	var contacts []authModel.Contact
	if err := config.DB.Where("user_number IN ?", userNumbers).Where("contact_type = ?", "Email").Where("deleted_at IS NULL").Find(&contacts).Error; err != nil {
		return nil, err
	}
	for _, contact := range contacts {
		hasEmail[contact.UserNumber] = true
	}

	for _, number := range userNumbers {
		if activeUsers[number] && hasPassword[number] && hasEmail[number] {
			result[number] = true
		}
	}

	return result, nil
}
