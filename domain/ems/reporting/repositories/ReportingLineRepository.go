package repositories

import (
	"errors"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/reporting/models"
	"gorm.io/gorm"
)

// FindAllByOrganisation returns reporting lines of one organisation, optionally
// filtered by subordinate assignment, manager assignment and/or status.
func FindAllByOrganisation(organisationNumber, subordinateAssignment, managerAssignment, status string) ([]models.ReportingLine, error) {
	var lines []models.ReportingLine

	query := config.DB.Where("organisation_number = ?", organisationNumber)
	if subordinateAssignment != "" {
		query = query.Where("subordinate_assignment_number = ?", subordinateAssignment)
	}
	if managerAssignment != "" {
		query = query.Where("manager_assignment_number = ?", managerAssignment)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Order("created_at asc").Find(&lines).Error; err != nil {
		return nil, err
	}
	return lines, nil
}

// FindByNumberAndOrganisation returns a line only when it belongs to the given
// organisation, preventing cross-organisation disclosure.
func FindByNumberAndOrganisation(number, organisationNumber string, tx ...*gorm.DB) (*models.ReportingLine, error) {
	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}

	var line models.ReportingLine
	if err := db.
		Where("number = ? AND organisation_number = ?", number, organisationNumber).
		First(&line).Error; err != nil {
		return nil, err
	}
	return &line, nil
}

// FindActiveBySubordinate returns the ACTIVE reporting line for an assignment,
// or (nil, nil) when none exists. Used by the single-active-manager rule and the
// cycle walk.
func FindActiveBySubordinate(organisationNumber, subordinateAssignment string, tx ...*gorm.DB) (*models.ReportingLine, error) {
	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}

	var line models.ReportingLine
	err := db.
		Where("organisation_number = ? AND subordinate_assignment_number = ? AND status = ?",
			organisationNumber, subordinateAssignment, models.ReportingStatusActive).
		First(&line).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &line, nil
}

// FindAllActive returns every ACTIVE reporting line of one organisation. Used
// by the overlap-aware duplicate and cycle checks.
func FindAllActive(organisationNumber string, tx ...*gorm.DB) ([]models.ReportingLine, error) {
	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}

	var lines []models.ReportingLine
	if err := db.
		Where("organisation_number = ? AND status = ?", organisationNumber, models.ReportingStatusActive).
		Find(&lines).Error; err != nil {
		return nil, err
	}
	return lines, nil
}

func Create(line *models.ReportingLine, tx ...*gorm.DB) error {
	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}
	return db.Create(line).Error
}

func Save(line *models.ReportingLine, tx ...*gorm.DB) error {
	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}
	return db.Save(line).Error
}
