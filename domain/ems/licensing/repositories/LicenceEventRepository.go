package repositories

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/licensing/models"
	"gorm.io/gorm"
)

// CreateLicenceEvent inserts a licence event. The caller may pass an explicit transaction.
func CreateLicenceEvent(event *models.LicenceEvent, txn ...*gorm.DB) error {
	db := config.DB
	if len(txn) > 0 && txn[0] != nil {
		db = txn[0]
	}
	return db.Create(event).Error
}

// FindLicenceEvents returns licence events matching the filter.
// Events are returned newest first. The caller must provide a valid filter.
func FindLicenceEvents(filter models.LicenceEventFilter) ([]models.LicenceEvent, int64, error) {
	db := config.DB.Model(&models.LicenceEvent{})

	if filter.EventType != nil {
		filterVal := *filter.EventType
		// Use GORM's Where with the string value
		if filterVal != "" {
			db = db.Where("event_type = ?", filterVal)
		}
	}
	if filter.LicenceNumber != nil {
		db = db.Where("licence_number = ?", *filter.LicenceNumber)
	}
	if filter.OrganisationNumber != nil {
		db = db.Where("organisation_number = ?", *filter.OrganisationNumber)
	}
	if filter.ActorUserNumber != nil {
		db = db.Where("actor_user_number = ?", *filter.ActorUserNumber)
	}
	if filter.InstallationNumber != nil {
		db = db.Where("installation_number = ?", *filter.InstallationNumber)
	}
	if filter.Since != nil {
		db = db.Where("created_at >= ?", *filter.Since)
	}
	if filter.Until != nil {
		db = db.Where("created_at <= ?", *filter.Until)
	}

	var total int64
	if err := db.Model(&models.LicenceEvent{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var events []models.LicenceEvent
	db = db.Order("created_at desc")
	if filter.Limit > 0 {
		db = db.Limit(filter.Limit)
	}
	if filter.Offset > 0 {
		db = db.Offset(filter.Offset)
	}

	if err := db.Find(&events).Error; err != nil {
		return nil, 0, err
	}

	return events, total, nil
}

// FindLicenceEventsByLicence returns all events for a specific licence,
// newest first. Convenience wrapper for common queries.
func FindLicenceEventsByLicence(licenceNumber string, limit, offset int) ([]models.LicenceEvent, int64, error) {
	filter := models.LicenceEventFilter{
		LicenceNumber: func() *string { s := licenceNumber; return &s }(),
		Limit:         limit,
		Offset:        offset,
	}
	return FindLicenceEvents(filter)
}

// FindOrganisationLicenceEvents returns all events for an organisation's licences.
func FindOrganisationLicenceEvents(organisationNumber string, limit, offset int) ([]models.LicenceEvent, int64, error) {
	filter := models.LicenceEventFilter{
		OrganisationNumber: func() *string { s := organisationNumber; return &s }(),
		Limit:              limit,
		Offset:             offset,
	}
	return FindLicenceEvents(filter)
}

// FindEventsByActor returns events initiated by a specific actor.
func FindEventsByActor(actorUserNumber string, limit, offset int) ([]models.LicenceEvent, int64, error) {
	filter := models.LicenceEventFilter{
		ActorUserNumber: func() *string { s := actorUserNumber; return &s }(),
		Limit:           limit,
		Offset:          offset,
	}
	return FindLicenceEvents(filter)
}
