package repositories

import (
	"errors"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/workflow/models"
	"gorm.io/gorm"
)

// FindDeadlineConfig returns the organisation's workflow deadline configuration,
// or (nil, nil) when none has been set. The organisation is always the
// authoritative one supplied by the caller.
func FindDeadlineConfig(organisationNumber string) (*models.WorkflowDeadlineConfig, error) {
	var record models.WorkflowDeadlineConfig
	err := config.DB.Where("organisation_number = ?", organisationNumber).First(&record).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &record, nil
}

// UpsertDeadlineConfig inserts or updates the single configuration row for the
// record's organisation.
func UpsertDeadlineConfig(record *models.WorkflowDeadlineConfig) error {
	existing, err := FindDeadlineConfig(record.OrganisationNumber)
	if err != nil {
		return err
	}
	if existing == nil {
		return config.DB.Create(record).Error
	}
	existing.DefaultDeadlineHours = record.DefaultDeadlineHours
	return config.DB.Save(existing).Error
}
