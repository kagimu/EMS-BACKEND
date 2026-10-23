package services

import (
	"strings"
	"testing"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/stakeholders/models"
)

func ptrFloat(v float64) *float64    { return &v }
func ptrTime(t time.Time) *time.Time { return &t }

func TestIsValidStakeholderStatus(t *testing.T) {
	for _, valid := range []string{constants.StatusActive, constants.StatusInactive, StakeholderStatusEnded} {
		if !isValidStakeholderStatus(valid) {
			t.Errorf("%q must be a valid stakeholder status", valid)
		}
	}
	for _, invalid := range []string{"", "Active", "DELETED", "ended"} {
		if isValidStakeholderStatus(invalid) {
			t.Errorf("%q must not be a valid stakeholder status", invalid)
		}
	}
}

func TestIsValidSharePercentage(t *testing.T) {
	if !isValidSharePercentage(nil) {
		t.Fatal("nil share percentage must be allowed")
	}
	if !isValidSharePercentage(ptrFloat(0)) || !isValidSharePercentage(ptrFloat(100)) || !isValidSharePercentage(ptrFloat(25.5)) {
		t.Fatal("0, 100 and 25.5 must be allowed")
	}
	if isValidSharePercentage(ptrFloat(-1)) || isValidSharePercentage(ptrFloat(100.5)) {
		t.Fatal("negative and over-100 percentages must be rejected")
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

func TestStakeholderTypeCatalogue(t *testing.T) {
	for _, expected := range []string{
		models.StakeholderTypeOwner,
		models.StakeholderTypeShareholder,
		models.StakeholderTypeBoardMember,
		models.StakeholderTypeDirector,
		models.StakeholderTypeExecutive,
		models.StakeholderTypeStaff,
		models.StakeholderTypeConsultant,
		models.StakeholderTypePartner,
		models.StakeholderTypeOther,
	} {
		normalized, err := models.NormalizeStakeholderType(expected)
		if err != nil || normalized != expected {
			t.Errorf("%q must be a valid stakeholder type", expected)
		}
	}

	custom, err := models.NormalizeStakeholderType("superhero")
	if err != nil || custom != "superhero" {
		t.Fatal("custom stakeholder types must be accepted")
	}
	if _, err := models.NormalizeStakeholderType(""); err == nil {
		t.Fatal("empty stakeholder types must be rejected")
	}
	if _, err := models.NormalizeStakeholderType(strings.Repeat("a", models.StakeholderTypeMaxLength+1)); err == nil {
		t.Fatal("stakeholder types longer than the column must be rejected")
	}
}

func TestIsFinalActiveOwner(t *testing.T) {
	owner := models.StakeholderTypeOwner

	if isFinalActiveOwner(models.StakeholderTypeDirector, constants.StatusActive, constants.StatusInactive, 1) {
		t.Fatal("non-owner types must never be treated as the final owner")
	}
	if isFinalActiveOwner(owner, constants.StatusInactive, StakeholderStatusEnded, 1) {
		t.Fatal("an INACTIVE owner is not an ACTIVE owner")
	}
	if isFinalActiveOwner(owner, constants.StatusActive, constants.StatusActive, 1) {
		t.Fatal("reactivation must not be guarded")
	}

	if !isFinalActiveOwner(owner, constants.StatusActive, constants.StatusInactive, 1) {
		t.Fatal("the only ACTIVE owner must be protected when deactivated")
	}
	if !isFinalActiveOwner(owner, constants.StatusActive, StakeholderStatusEnded, 1) {
		t.Fatal("the only ACTIVE owner must be protected when ended")
	}

	if isFinalActiveOwner(owner, constants.StatusActive, StakeholderStatusEnded, 2) {
		t.Fatal("ending a non-final owner must be allowed")
	}
}
