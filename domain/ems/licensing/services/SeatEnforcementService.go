package services

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	licenseModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/licensing/models"
	eosModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/eos/models"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrSeatCapacityReached = errors.New("user seat capacity reached")
	ErrSeatScopeInactive   = errors.New("the USER_SEAT scope is inactive or expired")
	ErrSiteCapacityReached = errors.New("site capacity reached")
	ErrSiteScopeInactive   = errors.New("the SITE scope is inactive or expired")
)

func activeScopeFor(organisationNumber, scopeType string) (*licenseModels.LicenceScope, error) {
	return activeScopeForDB(config.DB, organisationNumber, scopeType)
}

func activeScopeForDB(db *gorm.DB, organisationNumber, scopeType string) (*licenseModels.LicenceScope, error) {
	now := time.Now()
	var licence licenseModels.Licence
	if err := db.Where("organisation_number = ? AND status = ?", organisationNumber, licenseModels.StatusActive).
		Order("created_at DESC").First(&licence).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	if now.Before(licence.ValidFrom) || now.After(licence.ValidTo.Add(time.Duration(licence.GraceDays)*24*time.Hour)) {
		return nil, ErrSeatScopeInactive
	}
	var scope licenseModels.LicenceScope
	if err := db.Where("licence_number = ? AND scope_type = ?", licence.Number, scopeType).First(&scope).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	if scope.Status != licenseModels.FlexStatusActive || now.Before(scope.ValidFrom) || now.After(scope.ValidTo) {
		return nil, ErrSeatScopeInactive
	}
	return &scope, nil
}

func activeUnits(tx *gorm.DB, scopeNumber string) (int64, error) {
	// A locking read is required: under REPEATABLE READ a non-locking COUNT
	// can keep the transaction snapshot and miss a concurrently committed unit.
	var units []licenseModels.LicensedUnit
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("scope_number = ? AND status = ?", scopeNumber, licenseModels.UnitStatusActive).
		Find(&units).Error
	return int64(len(units)), err
}

func recordUnitEvent(tx *gorm.DB, actorNumber, organisationNumber, action string, unit licenseModels.LicensedUnit, reason string) {
	encoded, _ := json.Marshal(unit)
	reasonValue := action
	_ = tx.Create(&licenseModels.LicenceEvent{
		EventType:          licenseModels.LicenceEventConfigured,
		LicenceNumber:      unit.LicenceNumber,
		OrganisationNumber: organisationNumber,
		ActorUserNumber:    actorNumber,
		AfterState:         string(encoded),
		Reason:             &reasonValue,
	}).Error
}

func EnsureMembershipSeat(tx *gorm.DB, actorNumber, organisationNumber, userNumber string) error {
	scope, err := activeScopeFor(organisationNumber, licenseModels.ScopeUserSeat)
	if err != nil {
		return err
	}
	if scope == nil {
		return nil
	}
	var existing licenseModels.LicensedUnit
	existingErr := tx.Where("licence_number = ? AND unit_type = ? AND unit_identifier = ?",
		scope.LicenceNumber, licenseModels.ScopeUserSeat, userNumber).First(&existing).Error
	if existingErr == nil && existing.Status == licenseModels.UnitStatusActive {
		return nil
	}
	if existingErr != nil && !errors.Is(existingErr, gorm.ErrRecordNotFound) {
		return existingErr
	}
	var locked licenseModels.LicenceScope
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("number = ?", scope.Number).First(&locked).Error; err != nil {
		return err
	}
	used, err := activeUnits(tx, locked.Number)
	if err != nil {
		return err
	}
	if locked.Quantity > 0 && used >= int64(locked.Quantity) {
		return ErrSeatCapacityReached
	}
	unit := existing
	if unit.Number == "" {
		unit = licenseModels.LicensedUnit{
			Number: uuid.NewV4().String(), LicenceNumber: locked.LicenceNumber, ScopeNumber: locked.Number,
			UnitType: licenseModels.ScopeUserSeat, UnitIdentifier: userNumber,
			OrganisationNumber: organisationNumber, Status: licenseModels.UnitStatusActive,
			RegisteredAt: time.Now(), RegisteredBy: actorNumber,
		}
		if err := tx.Create(&unit).Error; err != nil {
			return err
		}
	} else {
		unit.ScopeNumber = locked.Number
		unit.Status = licenseModels.UnitStatusActive
		unit.RegisteredAt = time.Now()
		unit.RegisteredBy = actorNumber
		unit.ReleasedAt = nil
		unit.ReleasedBy = nil
		unit.ReleaseReason = ""
		if err := tx.Save(&unit).Error; err != nil {
			return err
		}
	}
	recordUnitEvent(tx, actorNumber, organisationNumber, "USER_SEAT_REGISTERED", unit, "")
	return nil
}

func ReleaseMembershipSeat(tx *gorm.DB, actorNumber, organisationNumber, userNumber string) error {
	scope, err := activeScopeFor(organisationNumber, licenseModels.ScopeUserSeat)
	if err != nil {
		if errors.Is(err, ErrSeatScopeInactive) {
			return nil
		}
		return err
	}
	if scope == nil {
		return nil
	}
	var unit licenseModels.LicensedUnit
	if err := tx.Where("licence_number = ? AND unit_type = ? AND unit_identifier = ? AND status = ?",
		scope.LicenceNumber, licenseModels.ScopeUserSeat, userNumber, licenseModels.UnitStatusActive).First(&unit).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}
	now := time.Now()
	unit.Status = licenseModels.UnitStatusReleased
	unit.ReleasedAt = &now
	unit.ReleasedBy = &actorNumber
	unit.ReleaseReason = "Membership removed"
	if err := tx.Save(&unit).Error; err != nil {
		return err
	}
	recordUnitEvent(tx, actorNumber, organisationNumber, "USER_SEAT_RELEASED", unit, "Membership removed")
	return nil
}

func AddEntityMembership(actorNumber, organisationNumber, userNumber string) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		var existing eosModels.EntityUser
		if err := tx.Where("user_number = ? AND entity_number = ?", userNumber, organisationNumber).First(&existing).Error; err == nil {
			return nil
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if err := EnsureMembershipSeat(tx, actorNumber, organisationNumber, userNumber); err != nil {
			return err
		}
		return tx.Create(&eosModels.EntityUser{
			EntityNumber: organisationNumber, UserNumber: userNumber, AddedByNumber: actorNumber,
		}).Error
	})
}

// RemoveEntityMembership only releases a seat when the membership actually
// existed in that organisation. Repeated removals are idempotent.
func RemoveEntityMembership(actorNumber, organisationNumber, userNumber string) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		var membership eosModels.EntityUser
		if err := tx.Where("user_number = ? AND entity_number = ?", userNumber, organisationNumber).First(&membership).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil
			}
			return err
		}
		if err := tx.Delete(&membership).Error; err != nil {
			return err
		}
		return ReleaseMembershipSeat(tx, actorNumber, organisationNumber, userNumber)
	})
}

func MembershipSeatUsage(organisationNumber string) (purchased, used, available int64, hasScope bool, err error) {
	scope, err := activeScopeFor(organisationNumber, licenseModels.ScopeUserSeat)
	if err != nil && !errors.Is(err, ErrSeatScopeInactive) {
		return 0, 0, 0, false, err
	}
	if scope == nil {
		return 0, 0, 0, false, nil
	}
	used, err = activeUnits(config.DB, scope.Number)
	if err != nil {
		return 0, 0, 0, true, err
	}
	purchased = int64(scope.Quantity)
	if purchased > 0 {
		available = purchased - used
		if available < 0 {
			available = 0
		}
	}
	return purchased, used, available, true, nil
}

func RegisterSiteSeat(actorNumber, organisationNumber, siteNumber string) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		return registerSiteSeatTx(tx, actorNumber, organisationNumber, siteNumber)
	})
}

// RegisterSiteSeatTx registers a SITE unit inside an existing transaction.
func RegisterSiteSeatTx(tx *gorm.DB, actorNumber, organisationNumber, siteNumber string) error {
	return registerSiteSeatTx(tx, actorNumber, organisationNumber, siteNumber)
}

func registerSiteSeatTx(tx *gorm.DB, actorNumber, organisationNumber, siteNumber string) error {
	scope, err := activeScopeForDB(tx, organisationNumber, licenseModels.ScopeSite)
	if err != nil {
		return err
	}
	if scope == nil {
		return nil
	}
	var existing licenseModels.LicensedUnit
	existingErr := tx.Where("licence_number = ? AND unit_type = ? AND unit_identifier = ?",
		scope.LicenceNumber, licenseModels.ScopeSite, siteNumber).First(&existing).Error
	if existingErr == nil && existing.Status == licenseModels.UnitStatusActive {
		return nil
	}
	if existingErr != nil && !errors.Is(existingErr, gorm.ErrRecordNotFound) {
		return existingErr
	}
	var locked licenseModels.LicenceScope
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("number = ?", scope.Number).First(&locked).Error; err != nil {
		return err
	}
	used, err := activeUnits(tx, locked.Number)
	if err != nil {
		return err
	}
	if locked.Quantity > 0 && used >= int64(locked.Quantity) {
		return ErrSiteCapacityReached
	}
	unit := existing
	if unit.Number == "" {
		unit = licenseModels.LicensedUnit{
			Number: uuid.NewV4().String(), LicenceNumber: locked.LicenceNumber, ScopeNumber: locked.Number,
			UnitType: licenseModels.ScopeSite, UnitIdentifier: siteNumber, SiteNumber: siteNumber,
			OrganisationNumber: organisationNumber, Status: licenseModels.UnitStatusActive,
			RegisteredAt: time.Now(), RegisteredBy: actorNumber,
		}
		if err := tx.Create(&unit).Error; err != nil {
			return err
		}
	} else {
		unit.ScopeNumber = locked.Number
		unit.Status = licenseModels.UnitStatusActive
		unit.RegisteredAt = time.Now()
		unit.RegisteredBy = actorNumber
		unit.ReleasedAt = nil
		unit.ReleasedBy = nil
		unit.ReleaseReason = ""
		if err := tx.Save(&unit).Error; err != nil {
			return err
		}
	}
	recordUnitEvent(tx, actorNumber, organisationNumber, "SITE_REGISTERED", unit, "")
	return nil
}

func ReleaseSiteSeat(actorNumber, organisationNumber, siteNumber, reason string) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		scope, err := activeScopeFor(organisationNumber, licenseModels.ScopeSite)
		if err != nil {
			if errors.Is(err, ErrSeatScopeInactive) {
				return nil
			}
			return err
		}
		if scope == nil {
			return nil
		}
		var unit licenseModels.LicensedUnit
		if err := tx.Where("licence_number = ? AND unit_type = ? AND unit_identifier = ? AND status = ?",
			scope.LicenceNumber, licenseModels.ScopeSite, siteNumber, licenseModels.UnitStatusActive).First(&unit).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil
			}
			return err
		}
		now := time.Now()
		unit.Status = licenseModels.UnitStatusReleased
		unit.ReleasedAt = &now
		unit.ReleasedBy = &actorNumber
		unit.ReleaseReason = strings.TrimSpace(reason)
		if err := tx.Save(&unit).Error; err != nil {
			return err
		}
		recordUnitEvent(tx, actorNumber, organisationNumber, "SITE_RELEASED", unit, reason)
		return nil
	})
}

// ReleaseAllUnitsForLicence releases active licensed units of all types when a
// licence is replaced. Historical rows are retained and a traceable event is
// recorded for each release.
func ReleaseAllUnitsForLicence(tx *gorm.DB, licenceNumber, actorNumber, organisationNumber, reason string) error {
	var units []licenseModels.LicensedUnit
	if err := tx.Where("licence_number = ? AND status = ?", licenceNumber, licenseModels.UnitStatusActive).Find(&units).Error; err != nil {
		return err
	}
	now := time.Now()
	for index := range units {
		units[index].Status = licenseModels.UnitStatusReleased
		units[index].ReleasedAt = &now
		units[index].ReleasedBy = &actorNumber
		units[index].ReleaseReason = reason
		if err := tx.Save(&units[index]).Error; err != nil {
			return err
		}
		recordUnitEvent(tx, actorNumber, organisationNumber, "UNIT_RELEASED_ON_REPLACEMENT", units[index], reason)
	}
	return nil
}

func IsSiteUnitType(unitType string) bool {
	switch strings.ToLower(strings.TrimSpace(unitType)) {
	case "site", "branch", "campus", "location":
		return true
	}
	return false
}
