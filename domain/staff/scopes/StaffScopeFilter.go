package scopes

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/staff/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ResolveStaffInfo checks if the caller is an active Staff member.
// It matches against both Staff.Number and Staff.UserNumber.
func ResolveStaffInfo(callerNumber string) (*models.Staff, bool) {
	if callerNumber == "" {
		return nil, false
	}
	var staff models.Staff
	err := config.DB.Where("number = ? OR user_number = ?", callerNumber, callerNumber).First(&staff).Error
	if err != nil || staff.Number == "" {
		return nil, false
	}
	return &staff, true
}

// GetStaffAllowedSchools retrieves all unique SchoolNumbers assigned to a staff member
// across both staff_scopes and stream_teachers tables.
func GetStaffAllowedSchools(staff *models.Staff) []string {
	ids := []string{staff.Number}
	if staff.UserNumber != "" {
		ids = append(ids, staff.UserNumber)
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
// - If the caller is not a staff member (e.g. Entity Owner / Admin), it allows all schools.
// - If the caller is a staff member with 0 allowed schools, it returns WHERE 1 = 0 (0 schools).
// - If the caller is a staff member with allowed schools, it returns WHERE number IN (...).
func FilterSchoolsByStaffScope(ctx *gin.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		callerNumber := ctx.GetHeader("User-Number")
		if callerNumber == "" {
			return db
		}

		staff, isStaff := ResolveStaffInfo(callerNumber)
		if !isStaff {
			// Not a staff member -> allow all schools
			return db
		}

		allowedSchools := GetStaffAllowedSchools(staff)
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
	staff, isStaff := ResolveStaffInfo(callerNumber)
	if !isStaff {
		return gin.H{
			"caller_number": callerNumber,
			"is_staff":      false,
			"debug_message": "Caller was not found in 'staffs' table matching number or user_number.",
		}
	}
	allowed := GetStaffAllowedSchools(staff)
	return gin.H{
		"caller_number":     callerNumber,
		"is_staff":          true,
		"staff_number":      staff.Number,
		"staff_user_number": staff.UserNumber,
		"allowed_schools":   allowed,
		"debug_message":     "Caller identified as staff member.",
	}
}

