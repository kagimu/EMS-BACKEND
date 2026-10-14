package repositories

import (
	"errors"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/assignments/models"
	"gorm.io/gorm"
)

// FindAllByOrganisation returns assignments of one organisation, optionally
// filtered by unit, user and/or status, ordered for stable pagination.
func FindAllByOrganisation(organisationNumber, unitNumber, userNumber, status string) ([]models.UserAssignment, error) {
	var assignments []models.UserAssignment

	query := config.DB.Where("organisation_number = ?", organisationNumber)
	if unitNumber != "" {
		query = query.Where("organisation_unit_number = ?", unitNumber)
	}
	if userNumber != "" {
		query = query.Where("user_number = ?", userNumber)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Order("created_at desc").Find(&assignments).Error; err != nil {
		return nil, err
	}

	return assignments, nil
}

// FindByNumberAndOrganisation returns an assignment only when it belongs to the
// given organisation, which prevents cross-organisation (IDOR) disclosure.
func FindByNumberAndOrganisation(number, organisationNumber string) (*models.UserAssignment, error) {
	var assignment models.UserAssignment
	if err := config.DB.
		Where("number = ? AND organisation_number = ?", number, organisationNumber).
		First(&assignment).Error; err != nil {
		return nil, err
	}
	return &assignment, nil
}

// FindActiveByOrgUserUnitType returns the ACTIVE assignment for a person at a
// unit and type, or (nil, nil) when none exists.
func FindActiveByOrgUserUnitType(organisationNumber, userNumber, unitNumber, assignmentType string) (*models.UserAssignment, error) {
	var assignment models.UserAssignment

	err := config.DB.
		Where("organisation_number = ? AND user_number = ? AND organisation_unit_number = ? AND assignment_type = ? AND status = ?",
			organisationNumber, userNumber, unitNumber, assignmentType, constants.StatusActive).
		First(&assignment).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &assignment, nil
}

// CountActiveByUnit counts ACTIVE assignments that reference a unit within an
// organisation. Used by the Foundation B unit-deletion integrity guard.
func CountActiveByUnit(unitNumber, organisationNumber string) (int64, error) {
	var count int64
	if err := config.DB.Model(&models.UserAssignment{}).
		Where("organisation_unit_number = ? AND organisation_number = ? AND status = ?",
			unitNumber, organisationNumber, constants.StatusActive).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func Create(assignment *models.UserAssignment) error {
	return config.DB.Create(assignment).Error
}

func Save(assignment *models.UserAssignment) error {
	return config.DB.Save(assignment).Error
}
