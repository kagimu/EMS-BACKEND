package commons

import "testing"

func TestKnownPermissionsIncludesEducationAndEms(t *testing.T) {
	for _, permission := range []string{
		"view_students",
		"update_permission",
		"can_update_entity",
		"manage_ems_organisation",
		"act_ems_approval",
		"manage_licence_registry",
	} {
		if !IsKnownPermission(permission) {
			t.Fatalf("expected %q to be a known permission", permission)
		}
	}

	if IsKnownPermission("totally_unknown_permission") {
		t.Fatalf("unknown permission was accepted")
	}
}
