package services

import (
	"errors"
	"testing"

	customfieldModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/customfields/models"
)

func TestValidType(t *testing.T) {
	for _, kind := range []string{"text", "multiline", "number", "date", "email", "url", "select", "boolean"} {
		if !ValidType(kind) {
			t.Fatalf("type %q should be valid", kind)
		}
	}
	if ValidType("richtext") || ValidType("") {
		t.Fatal("unsupported types must be rejected")
	}
}

func TestNormaliseKey(t *testing.T) {
	cases := map[string]string{
		"Registration Number": "registration_number",
		"Tax ID":              "tax_id",
		"":                    "",
	}
	for in, want := range cases {
		if got := normaliseKey(in, ""); got != want {
			t.Fatalf("normaliseKey(%q) = %q, want %q", in, got, want)
		}
	}
	if got := normaliseKey("ignored", "Custom_Key"); got != "custom_key" {
		t.Fatalf("explicit key not normalised: %q", got)
	}
}

func TestValidateValue(t *testing.T) {
	selectDef := customfieldModels.FieldDefinition{Type: customfieldModels.TypeSelect, OptionsJSON: `["A","B"]`}
	cases := []struct {
		name    string
		def     customfieldModels.FieldDefinition
		value   string
		wantErr error
	}{
		{"required empty rejected", customfieldModels.FieldDefinition{Type: customfieldModels.TypeText, Required: true}, "", ErrRequiredField},
		{"optional empty ok", customfieldModels.FieldDefinition{Type: customfieldModels.TypeText}, "", nil},
		{"valid number", customfieldModels.FieldDefinition{Type: customfieldModels.TypeNumber}, "12.5", nil},
		{"invalid number", customfieldModels.FieldDefinition{Type: customfieldModels.TypeNumber}, "abc", ErrInvalidValue},
		{"valid date", customfieldModels.FieldDefinition{Type: customfieldModels.TypeDate}, "2026-01-31", nil},
		{"invalid date", customfieldModels.FieldDefinition{Type: customfieldModels.TypeDate}, "31/01/2026", ErrInvalidValue},
		{"valid email", customfieldModels.FieldDefinition{Type: customfieldModels.TypeEmail}, "a@b.com", nil},
		{"invalid email", customfieldModels.FieldDefinition{Type: customfieldModels.TypeEmail}, "nope", ErrInvalidValue},
		{"valid url", customfieldModels.FieldDefinition{Type: customfieldModels.TypeURL}, "https://example.com", nil},
		{"invalid url", customfieldModels.FieldDefinition{Type: customfieldModels.TypeURL}, "example", ErrInvalidValue},
		{"valid option", selectDef, "A", nil},
		{"invalid option", selectDef, "C", ErrInvalidDropdown},
		{"valid boolean", customfieldModels.FieldDefinition{Type: customfieldModels.TypeBoolean}, "true", nil},
		{"invalid boolean", customfieldModels.FieldDefinition{Type: customfieldModels.TypeBoolean}, "yes", ErrInvalidValue},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateValue(tc.def, tc.value)
			if tc.wantErr == nil {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("err = %v, want %v", err, tc.wantErr)
			}
		})
	}
}
