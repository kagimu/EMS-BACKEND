package services

import (
	"encoding/json"
	"errors"
	"net/mail"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	customfieldModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/customfields/models"
	uuid "github.com/satori/go.uuid"
)

var (
	ErrInvalidFieldType    = errors.New("unsupported custom field type")
	ErrLabelRequired       = errors.New("custom field label is required")
	ErrKeyRequired         = errors.New("custom field key is required")
	ErrKeyExists           = errors.New("a custom field with this key already exists")
	ErrDefinitionNotFound  = errors.New("custom field definition not found")
	ErrSelectNeedsOptions  = errors.New("a select field requires at least one option")
	ErrRequiredField       = errors.New("a required custom field cannot be empty")
	ErrInvalidValue        = errors.New("the value is not valid for this field type")
	ErrInvalidDropdown     = errors.New("the value is not one of the allowed options")
	ErrTypeChangeWithValue = errors.New("cannot change a field type while values exist; deactivate it instead")
)

var keyPattern = regexp.MustCompile(`^[a-z0-9_]+$`)

// ValidType reports whether a field type is supported.
func ValidType(kind string) bool {
	switch kind {
	case customfieldModels.TypeText, customfieldModels.TypeMultiline, customfieldModels.TypeNumber,
		customfieldModels.TypeDate, customfieldModels.TypeEmail, customfieldModels.TypeURL,
		customfieldModels.TypeSelect, customfieldModels.TypeBoolean:
		return true
	}
	return false
}

// DefinitionInput is an administrator-supplied field definition.
type DefinitionInput struct {
	Label        string
	Key          string
	Type         string
	Description  string
	Required     bool
	DefaultValue string
	Options      []string
	DisplayOrder int
	Active       bool
}

func normaliseKey(label, key string) string {
	k := strings.ToLower(strings.TrimSpace(key))
	if k == "" {
		k = strings.ToLower(strings.TrimSpace(label))
		k = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(k, "_")
		k = strings.Trim(k, "_")
	}
	return k
}

// ListDefinitions returns the definitions for an organisation+scope by order.
func ListDefinitions(organisationNumber, scope string) ([]customfieldModels.FieldDefinition, error) {
	var defs []customfieldModels.FieldDefinition
	err := config.DB.Where("organisation_number = ? AND scope = ?", organisationNumber, scope).
		Order("display_order asc, id asc").Find(&defs).Error
	return defs, err
}

// ListValues returns definition_number → value for an organisation+scope.
func ListValues(organisationNumber, scope string) (map[string]string, error) {
	var values []customfieldModels.FieldValue
	err := config.DB.Where("organisation_number = ? AND scope = ?", organisationNumber, scope).Find(&values).Error
	if err != nil {
		return nil, err
	}
	out := make(map[string]string, len(values))
	for _, v := range values {
		out[v.DefinitionNumber] = v.Value
	}
	return out, nil
}

// CreateDefinition validates and persists a new field definition.
func CreateDefinition(organisationNumber, scope, actor string, in DefinitionInput) (*customfieldModels.FieldDefinition, error) {
	label := strings.TrimSpace(in.Label)
	if label == "" {
		return nil, ErrLabelRequired
	}
	if !ValidType(in.Type) {
		return nil, ErrInvalidFieldType
	}
	key := normaliseKey(label, in.Key)
	if key == "" || !keyPattern.MatchString(key) {
		return nil, ErrKeyRequired
	}
	if in.Type == customfieldModels.TypeSelect && len(in.Options) == 0 {
		return nil, ErrSelectNeedsOptions
	}

	var existing customfieldModels.FieldDefinition
	if err := config.DB.Where("organisation_number = ? AND scope = ? AND field_key = ?", organisationNumber, scope, key).
		First(&existing).Error; err == nil {
		return nil, ErrKeyExists
	}

	optionsJSON := "[]"
	if len(in.Options) > 0 {
		data, _ := json.Marshal(in.Options)
		optionsJSON = string(data)
	}

	def := customfieldModels.FieldDefinition{
		Number:             uuid.NewV4().String(),
		OrganisationNumber: organisationNumber,
		Scope:              scope,
		Label:              label,
		Key:                key,
		Type:               in.Type,
		Description:        strings.TrimSpace(in.Description),
		Required:           in.Required,
		DefaultValue:       in.DefaultValue,
		OptionsJSON:        optionsJSON,
		DisplayOrder:       in.DisplayOrder,
		Active:             in.Active,
		CreatedByNumber:    actor,
	}
	if err := config.DB.Create(&def).Error; err != nil {
		return nil, err
	}
	return &def, nil
}

// UpdateDefinition edits a definition within the organisation. Changing the
// type is refused while values exist (values would be destroyed); deactivate or
// create a new field instead.
func UpdateDefinition(organisationNumber, definitionNumber string, in DefinitionInput) (*customfieldModels.FieldDefinition, error) {
	var def customfieldModels.FieldDefinition
	if err := config.DB.Where("number = ? AND organisation_number = ? AND scope = ?", definitionNumber, organisationNumber, customfieldModels.ScopeOrganisation).
		First(&def).Error; err != nil {
		return nil, ErrDefinitionNotFound
	}

	if strings.TrimSpace(in.Label) != "" {
		def.Label = strings.TrimSpace(in.Label)
	}
	def.Description = strings.TrimSpace(in.Description)
	def.Required = in.Required
	def.DisplayOrder = in.DisplayOrder
	def.Active = in.Active
	if in.DefaultValue != "" {
		def.DefaultValue = in.DefaultValue
	}

	if in.Type != "" && in.Type != def.Type {
		if !ValidType(in.Type) {
			return nil, ErrInvalidFieldType
		}
		var count int64
		config.DB.Model(&customfieldModels.FieldValue{}).
			Where("organisation_number = ? AND definition_number = ?", organisationNumber, def.Number).Count(&count)
		if count > 0 {
			return nil, ErrTypeChangeWithValue
		}
		def.Type = in.Type
	}

	if def.Type == customfieldModels.TypeSelect {
		if len(in.Options) == 0 {
			return nil, ErrSelectNeedsOptions
		}
		data, _ := json.Marshal(in.Options)
		def.OptionsJSON = string(data)
	}

	if err := config.DB.Save(&def).Error; err != nil {
		return nil, err
	}
	return &def, nil
}

// ValidateValue applies type/required/option rules for one definition.
func ValidateValue(def customfieldModels.FieldDefinition, value string) error {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		if def.Required {
			return ErrRequiredField
		}
		return nil
	}
	switch def.Type {
	case customfieldModels.TypeNumber:
		if _, err := strconv.ParseFloat(trimmed, 64); err != nil {
			return ErrInvalidValue
		}
	case customfieldModels.TypeDate:
		if _, err := time.Parse("2006-01-02", trimmed); err != nil {
			return ErrInvalidValue
		}
	case customfieldModels.TypeEmail:
		if _, err := mail.ParseAddress(trimmed); err != nil {
			return ErrInvalidValue
		}
	case customfieldModels.TypeURL:
		u, err := url.ParseRequestURI(trimmed)
		if err != nil || u.Scheme == "" || u.Host == "" {
			return ErrInvalidValue
		}
	case customfieldModels.TypeSelect:
		var options []string
		_ = json.Unmarshal([]byte(def.OptionsJSON), &options)
		allowed := false
		for _, option := range options {
			if option == trimmed {
				allowed = true
				break
			}
		}
		if !allowed {
			return ErrInvalidDropdown
		}
	case customfieldModels.TypeBoolean:
		if trimmed != "true" && trimmed != "false" {
			return ErrInvalidValue
		}
	}
	return nil
}

// SaveValues validates and upserts a batch of values for the organisation. It
// rejects unknown/inactive definitions (so a foreign organisation's definition
// number cannot be used).
func SaveValues(organisationNumber, scope, actor string, values map[string]string) error {
	for definitionNumber, value := range values {
		var def customfieldModels.FieldDefinition
		if err := config.DB.Where("number = ? AND organisation_number = ? AND scope = ?", definitionNumber, organisationNumber, scope).
			First(&def).Error; err != nil {
			return ErrDefinitionNotFound
		}
		if !def.Active {
			// Inactive fields are not editable; keep their stored value.
			continue
		}
		if err := ValidateValue(def, value); err != nil {
			return err
		}

		var existing customfieldModels.FieldValue
		err := config.DB.Where("organisation_number = ? AND definition_number = ?", organisationNumber, definitionNumber).First(&existing).Error
		if err == nil {
			existing.Value = value
			existing.UpdatedByNumber = actor
			if saveErr := config.DB.Save(&existing).Error; saveErr != nil {
				return saveErr
			}
			continue
		}
		record := customfieldModels.FieldValue{
			Number:             uuid.NewV4().String(),
			OrganisationNumber: organisationNumber,
			Scope:              scope,
			DefinitionNumber:   definitionNumber,
			Value:              value,
			UpdatedByNumber:    actor,
		}
		if createErr := config.DB.Create(&record).Error; createErr != nil {
			return createErr
		}
	}
	return nil
}
