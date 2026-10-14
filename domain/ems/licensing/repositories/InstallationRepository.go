package repositories

import (
	"errors"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/licensing/models"
	"gorm.io/gorm"
)

// CreateInstallation inserts an installation. The caller may pass an explicit
// transaction.
func CreateInstallation(installation *models.Installation, txn ...*gorm.DB) error {
	db := config.DB
	if len(txn) > 0 && txn[0] != nil {
		db = txn[0]
	}
	return db.Create(installation).Error
}

// SaveInstallation persists changes to an existing installation.
func SaveInstallation(installation *models.Installation, txn ...*gorm.DB) error {
	db := config.DB
	if len(txn) > 0 && txn[0] != nil {
		db = txn[0]
	}
	return db.Save(installation).Error
}

// FindInstallationByNumber returns an installation by business number, or
// (nil, nil) when it does not exist.
func FindInstallationByNumber(number string, txn ...*gorm.DB) (*models.Installation, error) {
	db := config.DB
	if len(txn) > 0 && txn[0] != nil {
		db = txn[0]
	}
	var installation models.Installation
	if err := db.Where("number = ?", number).First(&installation).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &installation, nil
}

// FindInstallationByNumberForUpdate is FindInstallationByNumber with a row
// lock.
func FindInstallationByNumberForUpdate(number string, txn ...*gorm.DB) (*models.Installation, error) {
	db := config.DB
	if len(txn) > 0 && txn[0] != nil {
		db = txn[0]
	}
	var installation models.Installation
	if err := db.Clauses(lockingUpdate()).Where("number = ?", number).First(&installation).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &installation, nil
}

// FindLatestInstallationByIdentity returns the most recent installation record
// for the licence + identity pair, or (nil, nil). History is preserved: a
// released record is never edited to look like a new activation.
func FindLatestInstallationByIdentity(licenceNumber, installationIdentity string, txn ...*gorm.DB) (*models.Installation, error) {
	db := config.DB
	if len(txn) > 0 && txn[0] != nil {
		db = txn[0]
	}
	var installation models.Installation
	if err := db.
		Where("licence_number = ? AND installation_identity = ?", licenceNumber, installationIdentity).
		Order("created_at desc").
		First(&installation).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &installation, nil
}

// CountCapacityHolding counts installations that consume licence capacity:
// status ACTIVE, including client-deactivated ones. Released/revoked records
// do not count (F-Q3A, §24 detail 2).
func CountCapacityHolding(licenceNumber string, txn ...*gorm.DB) (int64, error) {
	db := config.DB
	if len(txn) > 0 && txn[0] != nil {
		db = txn[0]
	}
	var count int64
	if err := db.Model(&models.Installation{}).
		Where("licence_number = ? AND status = ?", licenceNumber, models.InstallationStatusActive).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// ReleaseActiveByLicence releases every capacity-holding installation of the
// licence (status ACTIVE, including client-deactivated ones) as part of an
// authorised control-plane replacement. It never edits the licence
// association: history stays bound to the original licence (§13, §24 detail 1).
func ReleaseActiveByLicence(licenceNumber, actorUserNumber, reason string, releasedAt time.Time, txn ...*gorm.DB) (int64, error) {
	db := config.DB
	if len(txn) > 0 && txn[0] != nil {
		db = txn[0]
	}
	result := db.Model(&models.Installation{}).
		Where("licence_number = ? AND status = ?", licenceNumber, models.InstallationStatusActive).
		Updates(map[string]interface{}{
			"status":         models.InstallationStatusReleased,
			"released_by":    actorUserNumber,
			"released_at":    releasedAt,
			"release_reason": reason,
			"updated_at":     releasedAt,
		})
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

// FindInstallations returns installations, optionally filtered by licence,
// organisation and status, newest first.
func FindInstallations(licenceNumber, organisationNumber, status string) ([]models.Installation, error) {
	var installations []models.Installation
	query := config.DB.Model(&models.Installation{})
	if licenceNumber != "" {
		query = query.Where("licence_number = ?", licenceNumber)
	}
	if organisationNumber != "" {
		query = query.Where("organisation_number = ?", organisationNumber)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if err := query.Order("created_at desc").Find(&installations).Error; err != nil {
		return nil, err
	}
	return installations, nil
}
