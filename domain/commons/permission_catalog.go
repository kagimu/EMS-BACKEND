package commons

import "sync"

// emsPermissionCatalogue is the authoritative list of EMS Core and EOS
// permission strings used by the route table. Education permissions live in
// PermissionGroups(); together they form the only permission strings a role may
// contain.
var emsPermissionCatalogue = []string{
	"can_update_entity",
	"can_delete_entity",

	"view_ems_modules",
	"manage_ems_modules",
	"view_ems_organisation",
	"manage_ems_organisation",
	"view_ems_stakeholders",
	"manage_ems_stakeholders",
	"view_ems_assignments",
	"manage_ems_assignments",
	"view_ems_authority",
	"manage_ems_authority",
	"view_ems_workflow",
	"manage_ems_workflow",
	"act_ems_approval",
	"view_licence_status",
	"activate_licence_installation",
	"deactivate_own_installation",
	"request_licence_renewal",
	"manage_licence_registry",
	"view_ems_audit",
	"manage_ems_audit",
	"view_ems_notifications",
	"manage_ems_packages",
}

var (
	knownPermissionsOnce sync.Once
	knownPermissions     map[string]struct{}
)

// KnownPermissions returns the authoritative set of permission strings that a
// role may contain. It reuses the existing Education catalogue rather than
// duplicating the role system.
func KnownPermissions() map[string]struct{} {
	knownPermissionsOnce.Do(func() {
		knownPermissions = make(map[string]struct{})
		for _, group := range PermissionGroups() {
			for _, permissions := range group {
				for _, permission := range permissions {
					knownPermissions[permission] = struct{}{}
				}
			}
		}
		for _, permission := range emsPermissionCatalogue {
			knownPermissions[permission] = struct{}{}
		}
	})
	return knownPermissions
}

// IsKnownPermission reports whether permission is part of the authoritative
// catalogue.
func IsKnownPermission(permission string) bool {
	_, ok := KnownPermissions()[permission]
	return ok
}
