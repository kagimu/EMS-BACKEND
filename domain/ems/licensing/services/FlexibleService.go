package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/addons"
	licenseRequests "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/licensing/api/requests"
	licenseResponses "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/licensing/api/responses"
	licenseModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/licensing/models"
	organisationModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/organisation/models"
	eosModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/eos/models"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

var (
	ErrFlexibleValidation      = errors.New("invalid licence scope or entitlement")
	ErrScopeQuantityBelowUsage = errors.New("scope quantity cannot be below current usage")
	ErrScopeNotFound           = errors.New("licence scope not found")
	ErrUnitDuplicate           = errors.New("licensed unit is already registered")
	ErrUnitCapacityReached     = errors.New("licensed unit capacity reached")
	ErrUnitCrossOrganisation   = errors.New("licensed unit does not belong to this organisation")
)

func ensureFlexible(licence *licenseModels.Licence) error {
	if licence == nil {
		return nil
	}
	var scopeCount int64
	config.DB.Model(&licenseModels.LicenceScope{}).Where("licence_number = ?", licence.Number).Count(&scopeCount)
	if scopeCount == 0 {
		for _, scope := range []licenseModels.LicenceScope{
			{ScopeType: licenseModels.ScopeOrganisation, Quantity: 1},
			{ScopeType: licenseModels.ScopeInstallation, Quantity: licence.MaxInstallations},
		} {
			scope.Number = uuid.NewV4().String()
			scope.LicenceNumber = licence.Number
			scope.Status = licenseModels.FlexStatusActive
			scope.ValidFrom = licence.ValidFrom
			scope.ValidTo = licence.ValidTo
			if err := config.DB.Create(&scope).Error; err != nil {
				return err
			}
		}
	} else {
		var installationScope licenseModels.LicenceScope
		if err := config.DB.Where("licence_number = ? AND scope_type = ?", licence.Number, licenseModels.ScopeInstallation).First(&installationScope).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			scope := licenseModels.LicenceScope{Number: uuid.NewV4().String(), LicenceNumber: licence.Number, ScopeType: licenseModels.ScopeInstallation, Quantity: licence.MaxInstallations, Status: licenseModels.FlexStatusActive, ValidFrom: licence.ValidFrom, ValidTo: licence.ValidTo}
			config.DB.Create(&scope)
		}
	}

	var entitlementCount int64
	config.DB.Model(&licenseModels.LicenceEntitlement{}).Where("licence_number = ?", licence.Number).Count(&entitlementCount)
	if entitlementCount == 0 {
		for _, code := range decodeAddons(licence.LicensedAddons) {
			entitlement := licenseModels.LicenceEntitlement{
				Number: uuid.NewV4().String(), LicenceNumber: licence.Number, Code: code, Quantity: 1,
				Status: licenseModels.FlexStatusActive, ValidFrom: licence.ValidFrom, ValidTo: licence.ValidTo,
			}
			config.DB.Create(&entitlement)
		}
	}

	var installationScope licenseModels.LicenceScope
	if err := config.DB.Where("licence_number = ? AND scope_type = ?", licence.Number, licenseModels.ScopeInstallation).First(&installationScope).Error; err == nil {
		var installations []licenseModels.Installation
		config.DB.Where("licence_number = ? AND status = ?", licence.Number, licenseModels.InstallationStatusActive).Find(&installations)
		for _, installation := range installations {
			var count int64
			config.DB.Model(&licenseModels.LicensedUnit{}).Where("licence_number = ? AND unit_type = ? AND unit_identifier = ?", licence.Number, licenseModels.ScopeInstallation, installation.InstallationIdentity).Count(&count)
			if count == 0 {
				config.DB.Create(&licenseModels.LicensedUnit{
					Number: uuid.NewV4().String(), LicenceNumber: licence.Number, ScopeNumber: installationScope.Number,
					UnitType: licenseModels.ScopeInstallation, UnitIdentifier: installation.InstallationIdentity,
					OrganisationNumber: licence.OrganisationNumber, Status: licenseModels.UnitStatusActive,
					RegisteredAt: installation.ActivatedAt, RegisteredBy: "SYSTEM_MIGRATION",
				})
			}
		}
	}
	return nil
}

func scopeUsed(scope licenseModels.LicenceScope) int64 {
	var count int64
	config.DB.Model(&licenseModels.LicensedUnit{}).
		Where("scope_number = ? AND status = ?", scope.Number, licenseModels.UnitStatusActive).Count(&count)
	return count
}

func flexibleView(number string) (*licenseResponses.FlexibleLicenceResponse, error) {
	var licence licenseModels.Licence
	if err := config.DB.Where("number = ?", number).First(&licence).Error; err != nil {
		return nil, ErrLicenceNotFound
	}
	if err := ensureFlexible(&licence); err != nil {
		return nil, err
	}
	var scopes []licenseModels.LicenceScope
	config.DB.Where("licence_number = ?", number).Order("scope_type ASC, created_at ASC").Find(&scopes)
	var entitlements []licenseModels.LicenceEntitlement
	config.DB.Where("licence_number = ?", number).Order("code ASC").Find(&entitlements)
	var units []licenseModels.LicensedUnit
	config.DB.Where("licence_number = ?", number).Order("registered_at DESC").Find(&units)

	result := &licenseResponses.FlexibleLicenceResponse{LicenceNumber: number, Scopes: make([]licenseResponses.ScopeResponse, 0, len(scopes)), Entitlements: make([]licenseResponses.EntitlementResponse, 0, len(entitlements)), Units: make([]licenseResponses.LicensedUnitResponse, 0, len(units))}
	for _, scope := range scopes {
		used := scopeUsed(scope)
		result.Scopes = append(result.Scopes, licenseResponses.ScopeResponse{
			Number: scope.Number, ScopeType: scope.ScopeType, ScopeKey: scope.ScopeKey, Quantity: scope.Quantity,
			Used: used, Available: maxInt64(0, int64(scope.Quantity)-used), Status: scope.Status,
			ValidFrom: scope.ValidFrom, ValidTo: scope.ValidTo, ConfigJSON: scope.ConfigJSON,
		})
	}
	for _, entitlement := range entitlements {
		result.Entitlements = append(result.Entitlements, licenseResponses.EntitlementResponse{
			Number: entitlement.Number, Code: entitlement.Code, Quantity: entitlement.Quantity, Used: 0,
			Status: entitlement.Status, ValidFrom: entitlement.ValidFrom, ValidTo: entitlement.ValidTo, ConfigJSON: entitlement.ConfigJSON,
		})
	}
	for _, unit := range units {
		response := licenseResponses.LicensedUnitResponse{
			Number: unit.Number, LicenceNumber: unit.LicenceNumber, ScopeNumber: unit.ScopeNumber,
			UnitType: unit.UnitType, UnitIdentifier: unit.UnitIdentifier, OrganisationNumber: unit.OrganisationNumber,
			SiteNumber: unit.SiteNumber, Status: unit.Status, RegisteredAt: unit.RegisteredAt,
			ReleasedAt: unit.ReleasedAt, ReleaseReason: unit.ReleaseReason,
		}
		result.Units = append(result.Units, response)
	}
	return result, nil
}

func maxInt64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

// FlexibleLicenceView exposes the authoritative scopes, entitlements and units.
func FlexibleLicenceView(number string) (*licenseResponses.FlexibleLicenceResponse, error) {
	return flexibleView(number)
}

// ConfigureFlexible replaces/upserts scopes and entitlements. Reductions below
// current usage are rejected; existing units are never silently released.
func ConfigureFlexible(actorNumber, number string, request licenseRequests.ConfigureFlexibleRequest) (*licenseResponses.FlexibleLicenceResponse, error) {
	var licence licenseModels.Licence
	if err := config.DB.Where("number = ?", number).First(&licence).Error; err != nil {
		return nil, ErrLicenceNotFound
	}
	if err := ensureFlexible(&licence); err != nil {
		return nil, err
	}
	before, _ := flexibleView(number)

	err := config.DB.Transaction(func(tx *gorm.DB) error {
		for _, input := range request.Scopes {
			scopeType := strings.ToUpper(strings.TrimSpace(input.ScopeType))
			if !licenseModels.IsScopeType(scopeType) || input.Quantity < 0 {
				return ErrFlexibleValidation
			}
			validFrom := input.ValidFrom
			validTo := input.ValidTo
			if validFrom.IsZero() {
				validFrom = licence.ValidFrom
			}
			if validTo.IsZero() {
				validTo = licence.ValidTo
			}
			if !validTo.After(validFrom) {
				return ErrFlexibleValidation
			}
			status := strings.ToUpper(strings.TrimSpace(input.Status))
			if status == "" {
				status = licenseModels.FlexStatusActive
			}
			var existing licenseModels.LicenceScope
			err := tx.Where("licence_number = ? AND scope_type = ? AND scope_key = ?", number, scopeType, strings.TrimSpace(input.ScopeKey)).First(&existing).Error
			if err == nil {
				if int64(input.Quantity) < scopeUsed(existing) {
					return ErrScopeQuantityBelowUsage
				}
				existing.Quantity = input.Quantity
				existing.Status = status
				existing.ValidFrom = validFrom
				existing.ValidTo = validTo
				existing.ConfigJSON = input.ConfigJSON
				if err := tx.Save(&existing).Error; err != nil {
					return err
				}
			} else {
				scope := licenseModels.LicenceScope{Number: uuid.NewV4().String(), LicenceNumber: number, ScopeType: scopeType, ScopeKey: strings.TrimSpace(input.ScopeKey), Quantity: input.Quantity, Status: status, ValidFrom: validFrom, ValidTo: validTo, ConfigJSON: input.ConfigJSON}
				if err := tx.Create(&scope).Error; err != nil {
					return err
				}
			}
		}
		for _, input := range request.Entitlements {
			code := strings.ToLower(strings.TrimSpace(input.Code))
			if _, ok := addons.FindModuleDef(code); !ok || input.Quantity < 0 {
				return ErrFlexibleValidation
			}
			validFrom := input.ValidFrom
			validTo := input.ValidTo
			if validFrom.IsZero() {
				validFrom = licence.ValidFrom
			}
			if validTo.IsZero() {
				validTo = licence.ValidTo
			}
			if !validTo.After(validFrom) {
				return ErrFlexibleValidation
			}
			status := strings.ToUpper(strings.TrimSpace(input.Status))
			if status == "" {
				status = licenseModels.FlexStatusActive
			}
			var existing licenseModels.LicenceEntitlement
			if err := tx.Where("licence_number = ? AND code = ?", number, code).First(&existing).Error; err == nil {
				existing.Quantity = input.Quantity
				existing.Status = status
				existing.ValidFrom = validFrom
				existing.ValidTo = validTo
				existing.ConfigJSON = input.ConfigJSON
				if err := tx.Save(&existing).Error; err != nil {
					return err
				}
			} else {
				entitlement := licenseModels.LicenceEntitlement{Number: uuid.NewV4().String(), LicenceNumber: number, Code: code, Quantity: input.Quantity, Status: status, ValidFrom: validFrom, ValidTo: validTo, ConfigJSON: input.ConfigJSON}
				if err := tx.Create(&entitlement).Error; err != nil {
					return err
				}
			}
		}

		// Keep legacy fields as a projection for compatibility.
		var installationScope licenseModels.LicenceScope
		if err := tx.Where("licence_number = ? AND scope_type = ?", number, licenseModels.ScopeInstallation).First(&installationScope).Error; err == nil {
			if installationScope.Quantity > 0 {
				licence.MaxInstallations = installationScope.Quantity
			}
		}
		var entitlements []licenseModels.LicenceEntitlement
		tx.Where("licence_number = ? AND status = ?", number, licenseModels.FlexStatusActive).Find(&entitlements)
		codes := make([]string, 0, len(entitlements))
		for _, entitlement := range entitlements {
			codes = append(codes, entitlement.Code)
		}
		encoded, _ := json.Marshal(codes)
		licence.LicensedAddons = string(encoded)
		licence.LastActedBy = &actorNumber
		return tx.Save(&licence).Error
	})
	if err != nil {
		return nil, err
	}
	after, _ := flexibleView(number)
	_ = RecordLicenceConfigured(&licence, actorNumber, map[string]interface{}{"scopes": before.Scopes, "entitlements": before.Entitlements}, map[string]interface{}{"scopes": after.Scopes, "entitlements": after.Entitlements, "reason": request.Reason})
	return after, nil
}

func RegisterLicensedUnit(actorNumber, number string, request licenseRequests.RegisterUnitRequest) (*licenseResponses.LicensedUnitResponse, error) {
	var licence licenseModels.Licence
	if err := config.DB.Where("number = ?", number).First(&licence).Error; err != nil {
		return nil, ErrLicenceNotFound
	}
	if err := ensureFlexible(&licence); err != nil {
		return nil, err
	}
	unitType := strings.ToUpper(strings.TrimSpace(request.ScopeType))
	identifier := strings.TrimSpace(request.UnitIdentifier)
	if !licenseModels.IsScopeType(unitType) || identifier == "" {
		return nil, ErrFlexibleValidation
	}

	var unit licenseModels.LicensedUnit
	err := config.DB.Transaction(func(tx *gorm.DB) error {
		var scope licenseModels.LicenceScope
		if err := tx.Set("gorm:query_option", "FOR UPDATE").Where("licence_number = ? AND scope_type = ?", number, unitType).First(&scope).Error; err != nil {
			return ErrScopeNotFound
		}
		if scope.Status != licenseModels.FlexStatusActive {
			return ErrScopeNotFound
		}
		used := scopeUsed(scope)
		if used >= int64(scope.Quantity) {
			return ErrUnitCapacityReached
		}
		var duplicate int64
		tx.Model(&licenseModels.LicensedUnit{}).Where("licence_number = ? AND unit_type = ? AND unit_identifier = ?", number, unitType, identifier).Count(&duplicate)
		if duplicate > 0 {
			return ErrUnitDuplicate
		}

		siteNumber := strings.TrimSpace(request.SiteNumber)
		switch unitType {
		case licenseModels.ScopeSite:
			var count int64
			tx.Model(&organisationModels.OrganisationUnit{}).Where("number = ? AND organisation_number = ?", siteNumber, licence.OrganisationNumber).Count(&count)
			if count == 0 {
				return ErrUnitCrossOrganisation
			}
		case licenseModels.ScopeUserSeat:
			var count int64
			tx.Model(&eosModels.EntityUser{}).Where("user_number = ? AND entity_number = ?", identifier, licence.OrganisationNumber).Count(&count)
			if count == 0 {
				return ErrUnitCrossOrganisation
			}
		}

		unit = licenseModels.LicensedUnit{
			Number: uuid.NewV4().String(), LicenceNumber: number, ScopeNumber: scope.Number,
			UnitType: unitType, UnitIdentifier: identifier, OrganisationNumber: licence.OrganisationNumber,
			SiteNumber: siteNumber, Status: licenseModels.UnitStatusActive, RegisteredAt: time.Now(),
			RegisteredBy: actorNumber, ConfigJSON: request.ConfigJSON,
		}
		return tx.Create(&unit).Error
	})
	if err != nil {
		return nil, err
	}
	response := licenseResponses.LicensedUnitResponse{Number: unit.Number, LicenceNumber: unit.LicenceNumber, ScopeNumber: unit.ScopeNumber, UnitType: unit.UnitType, UnitIdentifier: unit.UnitIdentifier, OrganisationNumber: unit.OrganisationNumber, SiteNumber: unit.SiteNumber, Status: unit.Status, RegisteredAt: unit.RegisteredAt}
	return &response, nil
}

func ReleaseLicensedUnit(actorNumber, number, reason string) error {
	reason = strings.TrimSpace(reason)
	if reason == "" {
		return ErrFlexibleValidation
	}
	var unit licenseModels.LicensedUnit
	if err := config.DB.Where("number = ?", number).First(&unit).Error; err != nil {
		return ErrUnitDuplicate
	}
	if unit.Status != licenseModels.UnitStatusActive {
		return ErrFlexibleValidation
	}
	now := time.Now()
	unit.Status = licenseModels.UnitStatusReleased
	unit.ReleasedAt = &now
	unit.ReleasedBy = &actorNumber
	unit.ReleaseReason = reason
	return config.DB.Save(&unit).Error
}

// FlexibleAddonLicensed enforces active entitlement validity. If no entitlement
// rows exist, legacy licensed_addons remains authoritative.
func FlexibleAddonLicensed(licence *licenseModels.Licence, code string) bool {
	var count int64
	config.DB.Model(&licenseModels.LicenceEntitlement{}).Where("licence_number = ?", licence.Number).Count(&count)
	if count == 0 {
		return LicenceContainsAddon(*licence, code)
	}
	now := time.Now()
	var entitlement licenseModels.LicenceEntitlement
	if err := config.DB.Where("licence_number = ? AND code = ? AND status = ? AND valid_from <= ? AND valid_to >= ?", licence.Number, strings.ToLower(code), licenseModels.FlexStatusActive, now, now).First(&entitlement).Error; err != nil {
		return false
	}
	return true
}

func InstallationScopeLimit(licence *licenseModels.Licence) int {
	var scope licenseModels.LicenceScope
	if err := config.DB.Where("licence_number = ? AND scope_type = ? AND status = ?", licence.Number, licenseModels.ScopeInstallation, licenseModels.FlexStatusActive).First(&scope).Error; err == nil {
		return scope.Quantity
	}
	return licence.MaxInstallations
}

func CreateInstallationUnit(licence *licenseModels.Licence, identity string) error {
	var scope licenseModels.LicenceScope
	if err := config.DB.Where("licence_number = ? AND scope_type = ?", licence.Number, licenseModels.ScopeInstallation).First(&scope).Error; err != nil {
		return nil
	}
	var count int64
	config.DB.Model(&licenseModels.LicensedUnit{}).Where("licence_number = ? AND unit_type = ? AND unit_identifier = ?", licence.Number, licenseModels.ScopeInstallation, identity).Count(&count)
	if count > 0 {
		return nil
	}
	return config.DB.Create(&licenseModels.LicensedUnit{
		Number: uuid.NewV4().String(), LicenceNumber: licence.Number, ScopeNumber: scope.Number,
		UnitType: licenseModels.ScopeInstallation, UnitIdentifier: identity,
		OrganisationNumber: licence.OrganisationNumber, Status: licenseModels.UnitStatusActive, RegisteredAt: time.Now(),
	}).Error
}

func FlexibleErrorText(err error) string {
	if err == nil {
		return ""
	}
	return fmt.Sprintf("%v", err)
}
