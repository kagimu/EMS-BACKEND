package models

import (
	"errors"
	"strings"

	"gorm.io/gorm"
)

const (
	CommercialProductStatusDraft    = "DRAFT"
	CommercialProductStatusActive   = "ACTIVE"
	CommercialProductStatusInactive = "INACTIVE"
)

const (
	CommercialProductTypeCore   = "CORE"
	CommercialProductTypeAddon  = "ADDON"
	CommercialProductTypeBundle = "BUNDLE"
)

const (
	CommercialComponentStatusActive   = "ACTIVE"
	CommercialComponentStatusInactive = "INACTIVE"
)

var CommercialProductStatuses = []string{
	CommercialProductStatusDraft,
	CommercialProductStatusActive,
	CommercialProductStatusInactive,
}

var CommercialProductTypes = []string{
	CommercialProductTypeCore,
	CommercialProductTypeAddon,
	CommercialProductTypeBundle,
}

// CommercialProduct is a sellable catalogue entry. It is not an addon and does
// not grant, enable or alter any licence entitlement.
type CommercialProduct struct {
	gorm.Model
	Number          string `json:"number" gorm:"size:36;unique;not null"`
	ProductCode     string `json:"product_code" gorm:"size:60;unique;not null"`
	Name            string `json:"name" gorm:"size:191;not null"`
	Description     string `json:"description" gorm:"type:text"`
	ProductType     string `json:"product_type" gorm:"size:20;not null;index"`
	Status          string `json:"status" gorm:"size:20;not null;index"`
	DisplayOrder    int    `json:"display_order" gorm:"not null;default:0"`
	CreatedByNumber string `json:"created_by_number" gorm:"size:36"`
	UpdatedByNumber string `json:"updated_by_number" gorm:"size:36"`
}

func (CommercialProduct) TableName() string { return "ems_commercial_products" }

// CommercialProductComponent maps a commercial product to an existing addon
// catalogue code. The snapshot preserves the display name at configuration time.
type CommercialProductComponent struct {
	gorm.Model
	Number          string `json:"number" gorm:"size:36;unique;not null"`
	ProductNumber   string `json:"product_number" gorm:"size:36;not null;uniqueIndex:idx_commercial_product_component"`
	ComponentCode   string `json:"component_code" gorm:"size:60;not null;uniqueIndex:idx_commercial_product_component"`
	ComponentName   string `json:"component_name" gorm:"size:191"`
	IsIncluded      bool   `json:"is_included" gorm:"not null;default:false"`
	IsOptional      bool   `json:"is_optional" gorm:"not null;default:false"`
	Status          string `json:"status" gorm:"size:20;not null;index"`
	DisplayOrder    int    `json:"display_order" gorm:"not null;default:0"`
	CreatedByNumber string `json:"created_by_number" gorm:"size:36"`
	UpdatedByNumber string `json:"updated_by_number" gorm:"size:36"`
}

func (CommercialProductComponent) TableName() string { return "ems_commercial_product_components" }

func NormaliseCommercialProductStatus(value string) (string, error) {
	value = strings.ToUpper(strings.TrimSpace(value))
	for _, candidate := range CommercialProductStatuses {
		if value == candidate {
			return value, nil
		}
	}
	return "", errors.New("invalid commercial product status")
}

func NormaliseCommercialProductType(value string) (string, error) {
	value = strings.ToUpper(strings.TrimSpace(value))
	for _, candidate := range CommercialProductTypes {
		if value == candidate {
			return value, nil
		}
	}
	return "", errors.New("invalid commercial product type")
}

func NormaliseCommercialComponentStatus(value string) (string, error) {
	value = strings.ToUpper(strings.TrimSpace(value))
	if value == CommercialComponentStatusActive || value == CommercialComponentStatusInactive {
		return value, nil
	}
	return "", errors.New("invalid commercial component status")
}

func ValidateProductStatusTransition(from, to string) error {
	if from == to {
		return nil
	}
	allowed := map[string]map[string]bool{
		CommercialProductStatusDraft:    {CommercialProductStatusActive: true, CommercialProductStatusInactive: true},
		CommercialProductStatusActive:   {CommercialProductStatusInactive: true},
		CommercialProductStatusInactive: {CommercialProductStatusActive: true},
	}
	if !allowed[from][to] {
		return errors.New("invalid commercial product status transition")
	}
	return nil
}
