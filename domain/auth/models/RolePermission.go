package models

// RolePermission stores the relationship between roles and permissions
type RolePermission struct {
	ID           uint
	RoleID       uint
	PermissionID uint
}
