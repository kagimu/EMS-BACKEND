package services

import (
	"strings"
	"testing"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/assignments/models"
)

func ptrTime(t time.Time) *time.Time { return &t }

func TestIsValidAssignmentStatus(t *testing.T) {
	for _, valid := range []string{constants.StatusActive, constants.StatusInactive, AssignmentStatusEnded} {
		if !isValidAssignmentStatus(valid) {
			t.Errorf("%q must be a valid assignment status", valid)
		}
	}
	for _, invalid := range []string{"", "Active", "DELETED", "ended"} {
		if isValidAssignmentStatus(invalid) {
			t.Errorf("%q must not be a valid assignment status", invalid)
		}
	}
}

func TestIsValidDateRange(t *testing.T) {
	earlier := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	later := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)

	if !isValidDateRange(nil, nil) || !isValidDateRange(ptrTime(earlier), nil) || !isValidDateRange(nil, ptrTime(later)) {
		t.Fatal("partial or empty ranges must be allowed")
	}
	if !isValidDateRange(ptrTime(earlier), ptrTime(later)) {
		t.Fatal("an ascending range must be allowed")
	}
	if isValidDateRange(ptrTime(later), ptrTime(earlier)) {
		t.Fatal("a descending range must be rejected")
	}
}

func TestAssignmentTypeCatalogue(t *testing.T) {
	for _, expected := range []string{
		models.AssignmentTypeMember,
		models.AssignmentTypeManager,
		models.AssignmentTypeHead,
		models.AssignmentTypeStaff,
	} {
		normalized, err := models.NormalizeAssignmentType(expected)
		if err != nil || normalized != expected {
			t.Errorf("%q must be a valid assignment type", expected)
		}
	}

	custom, err := models.NormalizeAssignmentType("superhero")
	if err != nil || custom != "superhero" {
		t.Fatal("custom assignment types must be accepted")
	}
	if _, err := models.NormalizeAssignmentType(""); err == nil {
		t.Fatal("empty assignment types must be rejected")
	}
	if _, err := models.NormalizeAssignmentType(strings.Repeat("a", models.AssignmentTypeMaxLength+1)); err == nil {
		t.Fatal("assignment types longer than the column must be rejected")
	}
}
