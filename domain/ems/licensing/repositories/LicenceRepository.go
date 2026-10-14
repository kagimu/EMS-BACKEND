package repositories

import (
	"errors"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/licensing/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Create inserts a licence. The caller may pass an explicit transaction.
// The unique index on active_org_number is the concurrency backstop for the
// one-active-licence rule (F-Q7).
func Create(licence *models.Licence, txn ...*gorm.DB) error {
	db := config.DB
	if len(txn) > 0 && txn[0] != nil {
		db = txn[0]
	}
	return db.Create(licence).Error
}

// Save persists changes to an existing licence.
func Save(licence *models.Licence, txn ...*gorm.DB) error {
	db := config.DB
	if len(txn) > 0 && txn[0] != nil {
		db = txn[0]
	}
	return db.Save(licence).Error
}

// FindByNumber returns a licence by business number, or (nil, nil) when it
// does not exist.
func FindByNumber(number string, txn ...*gorm.DB) (*models.Licence, error) {
	db := config.DB
	if len(txn) > 0 && txn[0] != nil {
		db = txn[0]
	}
	var licence models.Licence
	if err := db.Where("number = ?", number).First(&licence).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &licence, nil
}

// FindByKey returns the licence carrying the opaque key, or (nil, nil).
func FindByKey(licenceKey string) (*models.Licence, error) {
	var licence models.Licence
	if err := config.DB.Where("licence_key = ?", licenceKey).First(&licence).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &licence, nil
}

// FindByNumberForUpdate is FindByNumber with a row lock, for status changes
// and renewal inside a transaction.
func FindByNumberForUpdate(number string, txn ...*gorm.DB) (*models.Licence, error) {
	db := config.DB
	if len(txn) > 0 && txn[0] != nil {
		db = txn[0]
	}
	var licence models.Licence
	if err := db.Clauses(lockingUpdate()).Where("number = ?", number).First(&licence).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &licence, nil
}

// FindActiveByOrganisation returns the organisation's ACTIVE licence (the row
// whose active_org_number guard equals the organisation), or (nil, nil).
func FindActiveByOrganisation(organisationNumber string, txn ...*gorm.DB) (*models.Licence, error) {
	db := config.DB
	if len(txn) > 0 && txn[0] != nil {
		db = txn[0]
	}
	var licence models.Licence
	if err := db.Where("active_org_number = ?", organisationNumber).First(&licence).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &licence, nil
}

// FindLatestByOrganisation returns the organisation's most recently created
// licence, or (nil, nil). The client status view uses this so a suspended or
// revoked licence remains visible to its own organisation (the active guard
// is NULL for non-ACTIVE rows by design).
func FindLatestByOrganisation(organisationNumber string) (*models.Licence, error) {
	var licence models.Licence
	if err := config.DB.
		Where("organisation_number = ?", organisationNumber).
		Order("created_at desc").
		First(&licence).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &licence, nil
}

// FindAll returns licences, optionally filtered by organisation and stored
// status, newest first. This is a control-plane read (deliberately global).
func FindAll(organisationNumber, status string) ([]models.Licence, error) {
	var licences []models.Licence
	query := config.DB.Model(&models.Licence{})
	if organisationNumber != "" {
		query = query.Where("organisation_number = ?", organisationNumber)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if err := query.Order("created_at desc").Find(&licences).Error; err != nil {
		return nil, err
	}
	return licences, nil
}

// lockingUpdate returns a SELECT ... FOR UPDATE clause used by the finders
// above to serialise licence mutations inside transactions.
func lockingUpdate() clause.Locking {
	return clause.Locking{Strength: "UPDATE"}
}
