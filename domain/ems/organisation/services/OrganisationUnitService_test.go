package services

import (
	"testing"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
)

func TestIsSelfParent(t *testing.T) {
	if !isSelfParent("unit-1", "unit-1") {
		t.Fatal("expected a unit to be detected as its own parent")
	}
	if isSelfParent("unit-1", "unit-2") {
		t.Fatal("did not expect a different unit to be treated as self-parent")
	}
	if isSelfParent("", "") {
		t.Fatal("an empty unit number must never be a self-parent")
	}
}

func TestIsValidStatus(t *testing.T) {
	if !isValidStatus(constants.StatusActive) {
		t.Errorf("%q must be a valid status", constants.StatusActive)
	}
	if !isValidStatus(constants.StatusInactive) {
		t.Errorf("%q must be a valid status", constants.StatusInactive)
	}
	if isValidStatus("DELETED") {
		t.Fatal("unexpected lifecycle state DELETED must not be valid")
	}
	if isValidStatus("") {
		t.Fatal("empty status must not be valid")
	}
}

func TestAncestryContains(t *testing.T) {
	ancestors := []string{"parent-1", "parent-2", "parent-3"}

	if !ancestryContains(ancestors, "parent-2") {
		t.Fatal("expected parent-2 to be found in the ancestor chain")
	}
	if ancestryContains(ancestors, "unrelated") {
		t.Fatal("did not expect an unrelated number to be found")
	}
	if ancestryContains(ancestors, "") {
		t.Fatal("an empty target must never match")
	}
	if ancestryContains(nil, "parent-1") {
		t.Fatal("an empty ancestor chain must never match")
	}
}
