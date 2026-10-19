package integration

import (
	"testing"

	permissionModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/permissions/models"
	"gorm.io/gorm"
)

// cloneRoleForOrganisation creates a copy of sourceNumber in organisation with
// the same permissions. Fixtures use it so a client receives its role in the
// organisation it actually operates in, rather than relying on a permission
// granted in a different organisation.
func cloneRoleForOrganisation(t *testing.T, db *gorm.DB, sourceNumber, newNumber, organisation, owner string) {
	t.Helper()

	var source permissionModel.Role
	if err := db.Where("number = ?", sourceNumber).First(&source).Error; err != nil {
		t.Fatalf("loading source role %s: %v", sourceNumber, err)
	}

	copyRole := permissionModel.Role{
		Number:       newNumber,
		Name:         newNumber,
		EntityNumber: organisation,
		Description:  source.Description,
		Status:       "ACTIVE",
		UserNumber:   owner,
	}
	if err := db.Omit("User", "Entity", "Permissions").Create(&copyRole).Error; err != nil {
		t.Fatalf("seeding role copy %s: %v", newNumber, err)
	}

	var sourcePermissions []permissionModel.RolePermission
	if err := db.Where("role_number = ?", sourceNumber).Find(&sourcePermissions).Error; err != nil {
		t.Fatalf("loading permissions for %s: %v", sourceNumber, err)
	}
	for _, permission := range sourcePermissions {
		copyPermission := permissionModel.RolePermission{
			RoleNumber:    newNumber,
			Permission:    permission.Permission,
			AddedByNumber: owner,
		}
		if err := db.Omit("AddedBy", "Role").Create(&copyPermission).Error; err != nil {
			t.Fatalf("seeding permission %s for %s: %v", permission.Permission, newNumber, err)
		}
	}
}
