package models

import "gorm.io/gorm"

// OrganisationUnit is a generic EMS organisational unit belonging to one
// organisation (the existing eos.Entity, referenced by value).
//
// Foundation B: this is the only table introduced by this release
// (ems_organisation_units). It does NOT replace the Education hierarchy
// (School/Campus/Level/Class/Stream) and has no foreign key into any existing
// table — references are stored by value.
type OrganisationUnit struct {
	gorm.Model
	Number             string `json:"number" gorm:"size:36;unique;not null"`
	OrganisationNumber string `json:"organisation_number" gorm:"size:36;not null;index;index:idx_ems_organisation_units_org_status,priority:1"`
	ParentNumber       string `json:"parent_number" gorm:"size:36;index"`
	UnitType           string `json:"unit_type" gorm:"size:50;not null"`
	Name               string `json:"name" gorm:"size:128;not null"`
	Code               string `json:"code" gorm:"size:64;index"`
	Status             string `json:"status" gorm:"size:20;not null;index:idx_ems_organisation_units_org_status,priority:2"`
	AddedByNumber      string `json:"added_by_number" gorm:"size:36;not null"`
}

// TableName pins the table to ems_organisation_units.
func (OrganisationUnit) TableName() string {
	return "ems_organisation_units"
}
