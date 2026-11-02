package commons

import (
	"os"
	"strings"
)

// SoftwareOwnerUserNumbersEnv is the deployment-configured allow-list of user
// numbers that may exercise software-owner control-plane authority.
//
// This is deliberately not a user/identity system: the named users still
// authenticate through the normal JWT flow. It is an explicit, trusted
// deployment grant so owner access never depends on the first registered user,
// the Admin user type, or a role that a client administrator could create.
const SoftwareOwnerUserNumbersEnv = "EMS_SOFTWARE_OWNER_USER_NUMBERS"

// controlPlanePermissions are permissions whose authority spans organisations
// or the whole installation. They must never be delegated by an ordinary client
// administrator through organisation role management.
var controlPlanePermissions = map[string]struct{}{
	"manage_licence_registry": {},
	"manage_ems_audit":        {},
	"manage_ems_packages":     {},
	"manage_ems_modules":      {},
	"view_all_users":          {},
	"manage_support_tickets":  {},
	"assign_support_tickets":  {},
	"manage_support_settings": {},
}

// IsSoftwareOwner reports whether the authenticated user number is explicitly
// listed by the deployment as a software owner.
func IsSoftwareOwner(userNumber string) bool {
	userNumber = strings.TrimSpace(userNumber)
	if userNumber == "" {
		return false
	}
	for _, candidate := range strings.Split(os.Getenv(SoftwareOwnerUserNumbersEnv), ",") {
		if strings.TrimSpace(candidate) == userNumber {
			return true
		}
	}
	return false
}

// IsControlPlanePermission reports whether a permission grants installation-wide
// or cross-organisation authority.
func IsControlPlanePermission(permission string) bool {
	_, ok := controlPlanePermissions[permission]
	return ok
}

// ContainsControlPlanePermission reports whether any requested permission is a
// control-plane permission.
func ContainsControlPlanePermission(permissions []string) bool {
	for _, permission := range permissions {
		if IsControlPlanePermission(permission) {
			return true
		}
	}
	return false
}
