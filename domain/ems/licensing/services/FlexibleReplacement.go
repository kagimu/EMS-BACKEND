package services

import (
	"encoding/json"
	"time"

	licenseModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/licensing/models"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// copyFlexibleScopesTx copies the flexible scope and entitlement configuration
// from a superseded licence to its replacement. It deliberately does not copy
// licensed units; active units are released by ReleaseAllUnitsForLicence.
func copyFlexibleScopesTx(tx *gorm.DB, old, replacement licenseModels.Licence) error {
	var scopes []licenseModels.LicenceScope
	if err := tx.Where("licence_number = ?", old.Number).Find(&scopes).Error; err != nil {
		return err
	}
	if len(scopes) == 0 {
		for _, scopeType := range []string{licenseModels.ScopeOrganisation, licenseModels.ScopeInstallation} {
			quantity := 1
			if scopeType == licenseModels.ScopeInstallation {
				quantity = old.MaxInstallations
			}
			scope := licenseModels.LicenceScope{
				Number: uuid.NewV4().String(), LicenceNumber: replacement.Number, ScopeType: scopeType,
				Quantity: quantity, Status: licenseModels.FlexStatusActive,
				ValidFrom: replacement.ValidFrom, ValidTo: replacement.ValidTo,
			}
			if err := tx.Create(&scope).Error; err != nil {
				return err
			}
		}
	} else {
		for _, scope := range scopes {
			copyScope := licenseModels.LicenceScope{
				Number: uuid.NewV4().String(), LicenceNumber: replacement.Number, ScopeType: scope.ScopeType,
				ScopeKey: scope.ScopeKey, Quantity: scope.Quantity, Status: scope.Status,
				ValidFrom: replacement.ValidFrom, ValidTo: replacement.ValidTo, ConfigJSON: scope.ConfigJSON,
			}
			if err := tx.Create(&copyScope).Error; err != nil {
				return err
			}
		}
	}

	var entitlements []licenseModels.LicenceEntitlement
	if err := tx.Where("licence_number = ?", old.Number).Find(&entitlements).Error; err != nil {
		return err
	}
	if len(entitlements) == 0 {
		var codes []string
		_ = json.Unmarshal([]byte(old.LicensedAddons), &codes)
		for _, code := range codes {
			entitlement := licenseModels.LicenceEntitlement{
				Number: uuid.NewV4().String(), LicenceNumber: replacement.Number, Code: code,
				Quantity: 1, Status: licenseModels.FlexStatusActive,
				ValidFrom: replacement.ValidFrom, ValidTo: replacement.ValidTo,
			}
			if err := tx.Create(&entitlement).Error; err != nil {
				return err
			}
		}
	} else {
		for _, entitlement := range entitlements {
			copyEntitlement := licenseModels.LicenceEntitlement{
				Number: uuid.NewV4().String(), LicenceNumber: replacement.Number, Code: entitlement.Code,
				Quantity: entitlement.Quantity, Status: entitlement.Status,
				ValidFrom: replacement.ValidFrom, ValidTo: replacement.ValidTo, ConfigJSON: entitlement.ConfigJSON,
			}
			if err := tx.Create(&copyEntitlement).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

var _ = time.Now

// ReplacementCapacityError reports replacement capacity that cannot hold the
// active units being migrated.
type ReplacementCapacityError struct {
	ScopeType string
	Used      int64
	Quantity  int64
}

func (e *ReplacementCapacityError) Error() string {
	return "replacement licence has insufficient " + e.ScopeType + " capacity: used " +
		itoa(e.Used) + ", purchased " + itoa(e.Quantity)
}

func itoa(value int64) string {
	if value == 0 {
		return "0"
	}
	negative := value < 0
	if negative {
		value = -value
	}
	buf := make([]byte, 0, 20)
	for value > 0 {
		buf = append([]byte{byte('0' + value%10)}, buf...)
		value /= 10
	}
	if negative {
		buf = append([]byte{'-'}, buf...)
	}
	return string(buf)
}

// migrateReplacementUnitsTx migrates active USER_SEAT and SITE units to the
// replacement licence. It validates capacity before changing anything. Other
// unit types are released separately by ReleaseAllUnitsForLicence.
func migrateReplacementUnitsTx(tx *gorm.DB, old, replacement licenseModels.Licence, actorNumber, reason string) error {
	var oldScopes []licenseModels.LicenceScope
	if err := tx.Where("licence_number = ?", old.Number).Find(&oldScopes).Error; err != nil {
		return err
	}
	var newScopes []licenseModels.LicenceScope
	if err := tx.Where("licence_number = ?", replacement.Number).Find(&newScopes).Error; err != nil {
		return err
	}
	newByKey := map[string]licenseModels.LicenceScope{}
	for _, scope := range newScopes {
		newByKey[scope.ScopeType+"|"+scope.ScopeKey] = scope
	}

	for _, oldScope := range oldScopes {
		if oldScope.ScopeType != licenseModels.ScopeUserSeat && oldScope.ScopeType != licenseModels.ScopeSite {
			continue
		}
		var units []licenseModels.LicensedUnit
		if err := tx.Where("licence_number = ? AND scope_number = ? AND status = ?", old.Number, oldScope.Number, licenseModels.UnitStatusActive).Find(&units).Error; err != nil {
			return err
		}
		if len(units) == 0 {
			continue
		}
		newScope, ok := newByKey[oldScope.ScopeType+"|"+oldScope.ScopeKey]
		if !ok {
			return &ReplacementCapacityError{ScopeType: oldScope.ScopeType, Used: int64(len(units)), Quantity: 0}
		}
		if newScope.Quantity > 0 && int64(len(units)) > int64(newScope.Quantity) {
			return &ReplacementCapacityError{ScopeType: oldScope.ScopeType, Used: int64(len(units)), Quantity: int64(newScope.Quantity)}
		}
		for index := range units {
			units[index].LicenceNumber = replacement.Number
			units[index].ScopeNumber = newScope.Number
			if err := tx.Save(&units[index]).Error; err != nil {
				return err
			}
			recordUnitEvent(tx, actorNumber, replacement.OrganisationNumber, "UNIT_MIGRATED_TO_REPLACEMENT", units[index], reason)
		}
	}
	return nil
}
