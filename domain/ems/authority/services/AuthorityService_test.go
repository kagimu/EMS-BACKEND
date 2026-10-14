package services

import (
	"testing"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/ems/authority/models"
)

func ptrFloat(v float64) *float64    { return &v }
func ptrTime(t time.Time) *time.Time { return &t }

func TestIsValidAuthorityStatus(t *testing.T) {
	for _, valid := range []string{constants.StatusActive, constants.StatusInactive, AuthorityStatusEnded} {
		if !isValidAuthorityStatus(valid) {
			t.Errorf("%q must be a valid authority status", valid)
		}
	}
	for _, invalid := range []string{"", "Active", "DELETED", "ended"} {
		if isValidAuthorityStatus(invalid) {
			t.Errorf("%q must not be a valid authority status", invalid)
		}
	}
}

func TestIsValidAmountRange(t *testing.T) {
	if !isValidAmountRange(nil, nil) || !isValidAmountRange(ptrFloat(100), nil) || !isValidAmountRange(nil, ptrFloat(100)) {
		t.Fatal("partial or empty ranges must be allowed")
	}
	if !isValidAmountRange(ptrFloat(100), ptrFloat(100)) {
		t.Fatal("min == max must be allowed")
	}
	if !isValidAmountRange(ptrFloat(100), ptrFloat(200)) {
		t.Fatal("min < max must be allowed")
	}
	if isValidAmountRange(ptrFloat(200), ptrFloat(100)) {
		t.Fatal("min > max must be rejected")
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

// TestLevelCoversMaxAmountBoundary pins the exact amount boundary semantics:
// an amount equal to MaxAmount is covered; one unit above it is not; a NULL
// MaxAmount is unlimited.
func TestLevelCoversMaxAmountBoundary(t *testing.T) {
	limit := 5000000.0
	rows := []models.DecisionAuthority{{MaxAmount: ptrFloat(limit)}}

	if !levelCovers(rows, 0) {
		t.Fatal("zero must be covered by a finite limit")
	}
	if !levelCovers(rows, 100) {
		t.Fatal("an amount below the limit must be covered")
	}
	if !levelCovers(rows, limit-1) {
		t.Fatal("an amount just below the limit must be covered")
	}
	if levelCovers(rows, limit) {
		// equal is covered
	} else {
		t.Fatal("an amount exactly equal to MaxAmount must be covered")
	}
	if levelCovers(rows, limit+0.01) {
		t.Fatal("an amount one CENT above MaxAmount must NOT be covered")
	}
	if levelCovers(rows, limit+1) {
		t.Fatal("an amount one unit above MaxAmount must NOT be covered")
	}

	unlimited := []models.DecisionAuthority{{MaxAmount: nil}}
	if !levelCovers(unlimited, limit) || !levelCovers(unlimited, limit+1) {
		t.Fatal("a NULL MaxAmount must be unlimited")
	}

	if levelCovers([]models.DecisionAuthority{}, 10) {
		t.Fatal("an empty level must not cover any amount")
	}
}

func TestRowScopeApplies(t *testing.T) {
	ancestors := map[string]bool{"unit-a": true, "parent-1": true}

	if !rowScopeApplies(models.DecisionAuthority{OrganisationUnitNumber: ""}, ancestors) {
		t.Fatal("an organisation-wide row must always apply")
	}
	if !rowScopeApplies(models.DecisionAuthority{OrganisationUnitNumber: "unit-a"}, ancestors) {
		t.Fatal("a row scoped to the requested unit must apply")
	}
	if !rowScopeApplies(models.DecisionAuthority{OrganisationUnitNumber: "parent-1"}, ancestors) {
		t.Fatal("a row scoped to an ancestor unit must apply")
	}
	if rowScopeApplies(models.DecisionAuthority{OrganisationUnitNumber: "unrelated"}, ancestors) {
		t.Fatal("a row scoped to an unrelated unit must not apply")
	}
	if rowScopeApplies(models.DecisionAuthority{OrganisationUnitNumber: "unit-a"}, map[string]bool{}) {
		t.Fatal("with no requested unit, a unit-scoped row must not apply")
	}
}

// TestLevelResponseAnchorIdentity proves ratified D5: an anchor's identity is
// type + value + unit, so the same role in two different units produces two
// DISTINCT anchors and can never collapse merely because type and value match.
func TestLevelResponseAnchorIdentity(t *testing.T) {
	t.Run("same anchor different units stay distinct", func(t *testing.T) {
		level := levelResponse(1, []models.DecisionAuthority{
			{ApprovalLevel: 1, RoleNumber: "role-fin", OrganisationUnitNumber: "kampala"},
			{ApprovalLevel: 1, RoleNumber: "role-fin", OrganisationUnitNumber: "jinja"},
		})
		if len(level.Anchors) != 2 {
			t.Fatalf("expected 2 distinct anchors, got %d: %+v", len(level.Anchors), level.Anchors)
		}
		units := map[string]bool{}
		for _, anchor := range level.Anchors {
			if anchor.Type != "role" || anchor.Value != "role-fin" {
				t.Fatalf("unexpected anchor: %+v", anchor)
			}
			units[anchor.OrganisationUnitNumber] = true
		}
		if !units["kampala"] || !units["jinja"] {
			t.Fatalf("both units must be preserved on the anchors: %+v", level.Anchors)
		}
	})

	t.Run("same anchor same unit de-duplicates", func(t *testing.T) {
		level := levelResponse(1, []models.DecisionAuthority{
			{ApprovalLevel: 1, RoleNumber: "role-fin", OrganisationUnitNumber: "kampala"},
			{ApprovalLevel: 1, RoleNumber: "role-fin", OrganisationUnitNumber: "kampala"},
		})
		if len(level.Anchors) != 1 {
			t.Fatalf("expected 1 anchor, got %d: %+v", len(level.Anchors), level.Anchors)
		}
	})

	t.Run("organisation-wide and unit-scoped stay distinct", func(t *testing.T) {
		level := levelResponse(1, []models.DecisionAuthority{
			{ApprovalLevel: 1, RoleNumber: "role-fin"},
			{ApprovalLevel: 1, RoleNumber: "role-fin", OrganisationUnitNumber: "kampala"},
		})
		if len(level.Anchors) != 2 {
			t.Fatalf("expected 2 anchors, got %d: %+v", len(level.Anchors), level.Anchors)
		}
		orgWide := 0
		for _, anchor := range level.Anchors {
			if anchor.OrganisationUnitNumber == "" {
				orgWide++
			}
		}
		if orgWide != 1 {
			t.Fatalf("exactly one anchor must be organisation-wide: %+v", level.Anchors)
		}
	})
}
