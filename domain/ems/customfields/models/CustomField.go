package models

import "gorm.io/gorm"

// Organisation-scoped custom fields (data-driven). Two tables:
//   - FieldDefinition: an administrator-defined field (label, key, type, …).
//   - FieldValue: the value of one definition for one organisation.
//
// This is a small, additive EAV that can later be reused for other EMS records
// by changing Scope; it does NOT add a database column per customer field and
// does NOT create a duplicate organisation hierarchy.

// ScopeOrganisation is the only scope implemented in this phase.
const ScopeOrganisation = "organisation"

// Supported field types (only types we can validate and render reliably).
const (
	TypeText      = "text"
	TypeMultiline = "multiline"
	TypeNumber    = "number"
	TypeDate      = "date"
	TypeEmail     = "email"
	TypeURL       = "url"
	TypeSelect    = "select"
	TypeBoolean   = "boolean"
)

// FieldDefinition is an administrator-defined custom field for a scope.
type FieldDefinition struct {
	gorm.Model
	Number             string `json:"number" gorm:"uniqueIndex;size:36"`
	OrganisationNumber string `json:"organisation_number" gorm:"index;size:36"`
	Scope              string `json:"scope" gorm:"index;size:32"`
	Label              string `json:"label" gorm:"size:191"`
	Key                string `json:"key" gorm:"column:field_key;index;size:64"`
	Type               string `json:"type" gorm:"size:32"`
	Description        string `json:"description" gorm:"type:text"`
	Required           bool   `json:"required"`
	DefaultValue       string `json:"default_value" gorm:"type:text"`
	OptionsJSON        string `json:"options_json" gorm:"type:text"` // JSON array of strings (select)
	DisplayOrder       int    `json:"display_order"`
	Active             bool   `json:"active" gorm:"default:true"`
	CreatedByNumber    string `json:"created_by_number" gorm:"size:36"`
}

// FieldValue is one definition's value for one organisation.
type FieldValue struct {
	gorm.Model
	Number             string `json:"number" gorm:"uniqueIndex;size:36"`
	OrganisationNumber string `json:"organisation_number" gorm:"index;size:36"`
	Scope              string `json:"scope" gorm:"index;size:32"`
	DefinitionNumber   string `json:"definition_number" gorm:"index;size:36"`
	Value              string `json:"value" gorm:"type:text"`
	UpdatedByNumber    string `json:"updated_by_number" gorm:"size:36"`
}

func (FieldDefinition) TableName() string { return "ems_custom_field_definitions" }
func (FieldValue) TableName() string      { return "ems_custom_field_values" }
