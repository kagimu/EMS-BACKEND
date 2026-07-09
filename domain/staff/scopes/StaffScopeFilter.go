package scopes

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/staff/models"
	staffMembersModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/staffMembers/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ResolveStaffDeployments checks if the caller is an active Staff member in staff_deployments table.
// It matches against both StaffDeployment.Number and StaffDeployment.UserNumber.
func ResolveStaffDeployments(callerNumber string) ([]staffMembersModels.StaffDeployment, bool) {
	if callerNumber == "" {
		return nil, false
	}
	var deployments []staffMembersModels.StaffDeployment
	err := config.DB.Where("number = ? OR user_number = ?", callerNumber, callerNumber).Find(&deployments).Error
	if err != nil || len(deployments) == 0 {
		return nil, false
	}
	return deployments, true
}

// GetStaffAllowedSchools retrieves all unique SchoolNumbers assigned to a staff member
// across both staff_scopes and stream_teachers tables.
func GetStaffAllowedSchools(deployments []staffMembersModels.StaffDeployment) []string {
	var ids []string
	seenId := make(map[string]bool)

	for _, dep := range deployments {
		if dep.Number != "" && !seenId[dep.Number] {
			seenId[dep.Number] = true
			ids = append(ids, dep.Number)
		}
		if dep.UserNumber != "" && !seenId[dep.UserNumber] {
			seenId[dep.UserNumber] = true
			ids = append(ids, dep.UserNumber)
		}
	}

	if len(ids) == 0 {
		return nil
	}

	var staffScopes []models.StaffScope
	config.DB.Where("staff_number IN (?) OR user_number IN (?)", ids, ids).Find(&staffScopes)

	var st []models.StreamTeacher
	config.DB.Where("user_number IN (?)", ids).Find(&st)

	seen := make(map[string]bool)
	var allowed []string

	for _, sc := range staffScopes {
		if sc.SchoolNumber != "" && !seen[sc.SchoolNumber] {
			seen[sc.SchoolNumber] = true
			allowed = append(allowed, sc.SchoolNumber)
		}
	}
	for _, teacher := range st {
		if teacher.SchoolNumber != "" && !seen[teacher.SchoolNumber] {
			seen[teacher.SchoolNumber] = true
			allowed = append(allowed, teacher.SchoolNumber)
		}
	}

	return allowed
}

// FilterSchoolsByStaffScope is a GORM scope that strictly enforces school access for staff members.
// - If the caller is not found in staff_deployments (e.g. Entity Owner / Admin), it allows all schools.
// - If the caller is found in staff_deployments with 0 allowed schools, it returns WHERE 1 = 0 (0 schools).
// - If the caller is found in staff_deployments with allowed schools, it returns WHERE number IN (...).
func FilterSchoolsByStaffScope(ctx *gin.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		callerNumber := ctx.GetHeader("User-Number")
		if callerNumber == "" {
			return db
		}

		deployments, isStaff := ResolveStaffDeployments(callerNumber)
		if !isStaff {
			// Not in staff_deployments -> allow all schools (Entity Owner / Admin)
			return db
		}

		allowedSchools := GetStaffAllowedSchools(deployments)
		if len(allowedSchools) == 0 {
			// Staff member has no schools assigned -> strictly return 0 schools
			return db.Where("1 = 0")
		}

		return db.Where("number IN (?)", allowedSchools)
	}
}

// GetDebugInfo returns diagnostic information about the caller and their resolved scope.
func GetDebugInfo(ctx *gin.Context) gin.H {
	callerNumber := ctx.GetHeader("User-Number")
	deployments, isStaff := ResolveStaffDeployments(callerNumber)
	if !isStaff {
		return gin.H{
			"caller_number": callerNumber,
			"is_staff":      false,
			"debug_message": "Caller was not found in 'staff_deployments' table matching number or user_number.",
		}
	}

	var depNumbers []string
	var depUserNumbers []string
	for _, d := range deployments {
		depNumbers = append(depNumbers, d.Number)
		depUserNumbers = append(depUserNumbers, d.UserNumber)
	}

	allowed := GetStaffAllowedSchools(deployments)
	return gin.H{
		"caller_number":             callerNumber,
		"is_staff":                  true,
		"staff_deployment_numbers":  depNumbers,
		"staff_user_numbers":        depUserNumbers,
		"allowed_schools":           allowed,
		"debug_message":             "Caller identified in 'staff_deployments' table.",
	}
}
