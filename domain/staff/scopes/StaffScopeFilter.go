package scopes

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	authModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/staff/models"
	staffMembersModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/staffMembers/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ResolveAllStaffIDs looks up a staff member across staff_deployments, staffs, job_offers, and job_applications.
// Returns all known staff identifiers (deployment numbers, staff numbers, user numbers) for the caller.
func ResolveAllStaffIDs(callerNumber string) ([]string, bool) {
	if callerNumber == "" {
		return nil, false
	}
	seen := make(map[string]bool)
	var ids []string
	addID := func(id string) {
		if id != "" && !seen[id] {
			seen[id] = true
			ids = append(ids, id)
		}
	}

	// 1. Check staff_deployments directly
	var deployments []staffMembersModels.StaffDeployment
	config.DB.Where("number = ? OR user_number = ?", callerNumber, callerNumber).Find(&deployments)
	for _, d := range deployments {
		addID(d.Number)
		addID(d.UserNumber)
	}

	// 2. Check staffs table directly
	var staffs []models.Staff
	config.DB.Where("number = ? OR user_number = ?", callerNumber, callerNumber).Find(&staffs)
	for _, s := range staffs {
		addID(s.Number)
		addID(s.UserNumber)
	}

	// 3. Check job_offers by user_number
	var offers []staffMembersModels.JobOffer
	config.DB.Where("user_number = ?", callerNumber).Find(&offers)
	for _, o := range offers {
		var deps []staffMembersModels.StaffDeployment
		config.DB.Where("job_offer_number = ?", o.Number).Find(&deps)
		for _, d := range deps {
			addID(d.Number)
			addID(d.UserNumber)
		}
	}

	// 4. Check job_applications by user_number
	var apps []staffMembersModels.JobApplication
	config.DB.Where("user_number = ?", callerNumber).Find(&apps)
	for _, a := range apps {
		var offs []staffMembersModels.JobOffer
		config.DB.Where("job_application_number = ?", a.Number).Find(&offs)
		for _, o := range offs {
			var deps []staffMembersModels.StaffDeployment
			config.DB.Where("job_offer_number = ?", o.Number).Find(&deps)
			for _, d := range deps {
				addID(d.Number)
				addID(d.UserNumber)
			}
		}
	}

	if len(ids) == 0 {
		return nil, false
	}
	return ids, true
}

// GetStaffAllowedSchools retrieves all unique SchoolNumbers assigned to any of the resolved staff IDs.
func GetStaffAllowedSchools(staffIDs []string) []string {
	if len(staffIDs) == 0 {
		return nil
	}

	var staffScopes []models.StaffScope
	config.DB.Where("staff_number IN (?) OR user_number IN (?)", staffIDs, staffIDs).Find(&staffScopes)

	var st []models.StreamTeacher
	config.DB.Where("user_number IN (?)", staffIDs).Find(&st)

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

// GetCallerUserType inspects the gin context or users table to determine the caller's User.Type.
func GetCallerUserType(ctx *gin.Context, callerNumber string) string {
	if ctx != nil {
		if userObj, exists := ctx.Get("user"); exists {
			if user, ok := userObj.(authModels.User); ok {
				return user.Type
			}
		}
	}

	if callerNumber != "" {
		var u authModels.User
		if err := config.DB.Where("number = ?", callerNumber).First(&u).Error; err == nil {
			return u.Type
		}
	}

	return ""
}

// FilterSchoolsByStaffScope is a GORM scope that enforces school access for staff members.
func FilterSchoolsByStaffScope(ctx *gin.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		callerNumber := ""
		if ctx != nil {
			callerNumber = ctx.GetHeader("User-Number")
		}
		if callerNumber == "" {
			return db
		}

		userType := GetCallerUserType(ctx, callerNumber)
		if userType != "Staff" {
			// Scoping should not affect Admin or any non-Staff usertype
			return db
		}

		staffIDs, isStaff := ResolveAllStaffIDs(callerNumber)
		if !isStaff {
			// Caller is Staff usertype but has no staff deployment records -> allow bypassing
			return db
		}

		allowedSchools := GetStaffAllowedSchools(staffIDs)
		if len(allowedSchools) == 0 {
			// Staff member has no schools assigned -> allow bypassing
			return db
		}

		return db.Where("number IN (?)", allowedSchools)
	}
}

// GetDebugInfo returns diagnostic information about the caller and their resolved scope.
func GetDebugInfo(ctx *gin.Context) gin.H {
	callerNumber := ""
	if ctx != nil {
		callerNumber = ctx.GetHeader("User-Number")
	}
	userType := GetCallerUserType(ctx, callerNumber)

	if userType != "Staff" {
		return gin.H{
			"caller_number": callerNumber,
			"user_type":     userType,
			"is_staff":      false,
			"is_scoped":     false,
			"debug_message": "Scoping bypassed: caller is not Staff usertype (user_type: " + userType + ").",
		}
	}

	staffIDs, isStaff := ResolveAllStaffIDs(callerNumber)
	if !isStaff {
		return gin.H{
			"caller_number": callerNumber,
			"user_type":     userType,
			"is_staff":      false,
			"is_scoped":     false,
			"debug_message": "Scoping bypassed: Caller is Staff usertype but has no staff records setup.",
		}
	}

	allowed := GetStaffAllowedSchools(staffIDs)
	isScoped := len(allowed) > 0
	msg := "Caller identified as Staff usertype and scoped to assigned schools."
	if !isScoped {
		msg = "Scoping bypassed: Staff member has no schools assigned in their scope."
	}

	return gin.H{
		"caller_number":      callerNumber,
		"user_type":          userType,
		"is_staff":           true,
		"is_scoped":          isScoped,
		"resolved_staff_ids": staffIDs,
		"allowed_schools":    allowed,
		"debug_message":      msg,
	}
}
