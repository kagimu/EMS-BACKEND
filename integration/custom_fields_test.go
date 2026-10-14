package integration

import (
	"errors"
	"fmt"
	"testing"
	"time"

	customfieldModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/customfields/models"
	customfieldServices "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/customfields/services"
)

// TestOrganisationCustomFields exercises the data-driven custom-field backend:
// create definitions, save values, validation, and cross-organisation safety.
// It writes only to the dedicated onboarding database and skips otherwise.
func TestOrganisationCustomFields(t *testing.T) {
	db := requireOnboardingDB(t)

	suffix := time.Now().UnixNano() % 1000000000
	orgA := fmt.Sprintf("cf-org-a-%d", suffix)
	orgB := fmt.Sprintf("cf-org-b-%d", suffix)
	scope := customfieldModels.ScopeOrganisation

	cleanup := func() {
		db.Unscoped().Where("organisation_number IN ?", []string{orgA, orgB}).Delete(&customfieldModels.FieldValue{})
		db.Unscoped().Where("organisation_number IN ?", []string{orgA, orgB}).Delete(&customfieldModels.FieldDefinition{})
	}
	cleanup()
	t.Cleanup(cleanup)

	// 1. Create a custom text field.
	textField, err := customfieldServices.CreateDefinition(orgA, scope, orgA, customfieldServices.DefinitionInput{
		Label: "Registration Number", Key: "registration_number", Type: customfieldModels.TypeText, Required: false, Active: true,
	})
	if err != nil {
		t.Fatalf("create text field: %v", err)
	}

	// 2. Duplicate key rejected.
	if _, err := customfieldServices.CreateDefinition(orgA, scope, orgA, customfieldServices.DefinitionInput{
		Label: "Dup", Key: "registration_number", Type: customfieldModels.TypeText, Active: true,
	}); !errors.Is(err, customfieldServices.ErrKeyExists) {
		t.Fatalf("duplicate key not rejected: %v", err)
	}

	// 3. Select field requires options.
	if _, err := customfieldServices.CreateDefinition(orgA, scope, orgA, customfieldServices.DefinitionInput{
		Label: "Category", Type: customfieldModels.TypeSelect, Active: true,
	}); !errors.Is(err, customfieldServices.ErrSelectNeedsOptions) {
		t.Fatalf("select without options not rejected: %v", err)
	}

	// 4. Required field + a select field with options.
	reqField, err := customfieldServices.CreateDefinition(orgA, scope, orgA, customfieldServices.DefinitionInput{
		Label: "Tax ID", Key: "tax_id", Type: customfieldModels.TypeText, Required: true, Active: true,
	})
	if err != nil {
		t.Fatalf("create required field: %v", err)
	}
	selectField, err := customfieldServices.CreateDefinition(orgA, scope, orgA, customfieldServices.DefinitionInput{
		Label: "Category", Key: "category", Type: customfieldModels.TypeSelect, Options: []string{"Government", "Private"}, Active: true,
	})
	if err != nil {
		t.Fatalf("create select field: %v", err)
	}

	// 5. Save a valid value.
	if err := customfieldServices.SaveValues(orgA, scope, orgA, map[string]string{
		textField.Number:   "REG-001",
		reqField.Number:    "TIN-999",
		selectField.Number: "Private",
	}); err != nil {
		t.Fatalf("save values: %v", err)
	}

	// 6. Values persisted and retrievable (simulating another authorised user).
	values, err := customfieldServices.ListValues(orgA, scope)
	if err != nil {
		t.Fatalf("list values: %v", err)
	}
	if values[textField.Number] != "REG-001" || values[selectField.Number] != "Private" {
		t.Fatalf("values not persisted: %#v", values)
	}

	// 7. Required field rejects an empty value.
	if err := customfieldServices.SaveValues(orgA, scope, orgA, map[string]string{reqField.Number: ""}); !errors.Is(err, customfieldServices.ErrRequiredField) {
		t.Fatalf("required empty not rejected: %v", err)
	}

	// 8. Invalid dropdown value rejected by the backend.
	if err := customfieldServices.SaveValues(orgA, scope, orgA, map[string]string{selectField.Number: "Unknown"}); !errors.Is(err, customfieldServices.ErrInvalidDropdown) {
		t.Fatalf("invalid dropdown not rejected: %v", err)
	}

	// 9. Another organisation's definition cannot be used (cross-org safety).
	if err := customfieldServices.SaveValues(orgB, scope, orgB, map[string]string{textField.Number: "X"}); !errors.Is(err, customfieldServices.ErrDefinitionNotFound) {
		t.Fatalf("cross-organisation definition not rejected: %v", err)
	}

	// 10. Deactivating a field must not destroy stored values.
	if _, err := customfieldServices.UpdateDefinition(orgA, textField.Number, customfieldServices.DefinitionInput{
		Label: "Registration Number", Type: customfieldModels.TypeText, Active: false,
	}); err != nil {
		t.Fatalf("deactivate field: %v", err)
	}
	afterDeactivate, err := customfieldServices.ListValues(orgA, scope)
	if err != nil {
		t.Fatalf("list values after deactivate: %v", err)
	}
	if afterDeactivate[textField.Number] != "REG-001" {
		t.Fatalf("deactivating a field destroyed its value: %#v", afterDeactivate)
	}

	// 11. Changing a field type while values exist is refused.
	if _, err := customfieldServices.UpdateDefinition(orgA, reqField.Number, customfieldServices.DefinitionInput{
		Label: "Tax ID", Type: customfieldModels.TypeNumber, Active: true,
	}); !errors.Is(err, customfieldServices.ErrTypeChangeWithValue) {
		t.Fatalf("type change with values not rejected: %v", err)
	}
}
