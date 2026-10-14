package repositories

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/organisation/models"
)

// FindAllByOrganisation returns units for one organisation, optionally filtered
// by a parent number and/or status, ordered by name.
func FindAllByOrganisation(organisationNumber, parentNumber, status string) ([]models.OrganisationUnit, error) {
	var units []models.OrganisationUnit

	query := config.DB.Where("organisation_number = ?", organisationNumber)
	if parentNumber != "" {
		query = query.Where("parent_number = ?", parentNumber)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Order("name asc").Find(&units).Error; err != nil {
		return nil, err
	}

	return units, nil
}

// FindByNumberAndOrganisation returns a unit only when it belongs to the given
// organisation. A unit from another organisation is never returned, which is
// how cross-organisation (IDOR) access is prevented.
func FindByNumberAndOrganisation(number, organisationNumber string) (*models.OrganisationUnit, error) {
	var unit models.OrganisationUnit
	if err := config.DB.
		Where("number = ? AND organisation_number = ?", number, organisationNumber).
		First(&unit).Error; err != nil {
		return nil, err
	}
	return &unit, nil
}

// CountChildren counts live child units of a unit within the same organisation.
func CountChildren(number, organisationNumber string) (int64, error) {
	var count int64
	if err := config.DB.Model(&models.OrganisationUnit{}).
		Where("parent_number = ? AND organisation_number = ?", number, organisationNumber).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// FindByCode returns a unit with the given non-empty code within an organisation.
func FindByCode(organisationNumber, code string) (*models.OrganisationUnit, error) {
	var unit models.OrganisationUnit
	if err := config.DB.
		Where("organisation_number = ? AND code = ?", organisationNumber, code).
		First(&unit).Error; err != nil {
		return nil, err
	}
	return &unit, nil
}

func Create(unit *models.OrganisationUnit) error {
	return config.DB.Create(unit).Error
}

func Save(unit *models.OrganisationUnit) error {
	return config.DB.Save(unit).Error
}

// Delete performs a soft delete (the model embeds gorm.DeletedAt).
func Delete(unit *models.OrganisationUnit) error {
	return config.DB.Delete(unit).Error
}
