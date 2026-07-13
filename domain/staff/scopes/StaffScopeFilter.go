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

// GetStaffAllowedCampuses retrieves all unique CampusNumbers assigned to any of the resolved staff IDs.
func GetStaffAllowedCampuses(staffIDs []string) []string {
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
		if sc.CampusNumber != "" && !seen[sc.CampusNumber] {
			seen[sc.CampusNumber] = true
			allowed = append(allowed, sc.CampusNumber)
		}
	}
	for _, teacher := range st {
		if teacher.CampusNumber != "" && !seen[teacher.CampusNumber] {
			seen[teacher.CampusNumber] = true
			allowed = append(allowed, teacher.CampusNumber)
		}
	}

	return allowed
}

// GetStaffAllowedCurriculums retrieves all unique Curriculum numbers (CampusCurriculum.Number) allowed for any of the resolved staff IDs.
func GetStaffAllowedCurriculums(staffIDs []string) []string {
	if len(staffIDs) == 0 {
		return nil
	}

	var staffScopes []models.StaffScope
	config.DB.Where("staff_number IN (?) OR user_number IN (?)", staffIDs, staffIDs).Find(&staffScopes)

	var st []models.StreamTeacher
	config.DB.Where("user_number IN (?)", staffIDs).Find(&st)

	seen := make(map[string]bool)
	var allowed []string

	addCurriculum := func(num string) {
		if num != "" && !seen[num] {
			seen[num] = true
			allowed = append(allowed, num)
		}
	}

	// 1. Direct curriculum numbers from staff_scopes or broader scopes
	for _, sc := range staffScopes {
		if sc.CurriculumNumber != "" {
			addCurriculum(sc.CurriculumNumber)
		} else if sc.CampusNumber != "" {
			type ccResult struct {
				Number string
			}
			var ccs []ccResult
			config.DB.Table("campus_curriculums").Where("campus_number = ? AND deleted_at IS NULL", sc.CampusNumber).Select("number").Scan(&ccs)
			for _, cc := range ccs {
				addCurriculum(cc.Number)
			}
		} else if sc.SchoolNumber != "" {
			type ccResult struct {
				Number string
			}
			var ccs []ccResult
			config.DB.Table("campus_curriculums").
				Joins("JOIN campuses ON campuses.number = campus_curriculums.campus_number AND campuses.deleted_at IS NULL").
				Where("campuses.school_number = ? AND campus_curriculums.deleted_at IS NULL", sc.SchoolNumber).
				Select("campus_curriculums.number").Scan(&ccs)
			for _, cc := range ccs {
				addCurriculum(cc.Number)
			}
		}
	}

	// 2. Direct curriculum numbers from stream_teachers
	for _, teacher := range st {
		if teacher.CurriculumNumber != "" {
			addCurriculum(teacher.CurriculumNumber)
		}
	}

	// 3. Ensure both CampusCurriculum.Number and SchoolCurriculum/Curriculum IDs match correctly
	if len(allowed) > 0 {
		type ccResult struct {
			Number string
		}
		var extra []ccResult
		config.DB.Table("campus_curriculums").
			Joins("JOIN school_curriculums ON school_curriculums.number = campus_curriculums.school_curriculum_number AND school_curriculums.deleted_at IS NULL").
			Where("campus_curriculums.deleted_at IS NULL AND (campus_curriculums.number IN (?) OR campus_curriculums.school_curriculum_number IN (?) OR school_curriculums.curriculum_number IN (?))", allowed, allowed, allowed).
			Select("campus_curriculums.number").Scan(&extra)
		for _, ex := range extra {
			addCurriculum(ex.Number)
		}
	}

	return allowed
}

// GetStaffAllowedLevels retrieves all unique CampusLevel numbers (CampusLevel.Number) allowed for any of the resolved staff IDs.
func GetStaffAllowedLevels(staffIDs []string) []string {
	if len(staffIDs) == 0 {
		return nil
	}

	var staffScopes []models.StaffScope
	config.DB.Where("staff_number IN (?) OR user_number IN (?)", staffIDs, staffIDs).Find(&staffScopes)

	var st []models.StreamTeacher
	config.DB.Where("user_number IN (?)", staffIDs).Find(&st)

	seen := make(map[string]bool)
	var allowed []string

	addLevel := func(num string) {
		if num != "" && !seen[num] {
			seen[num] = true
			allowed = append(allowed, num)
		}
	}

	// 1. Levels from staff_scopes or broader scopes
	for _, sc := range staffScopes {
		if sc.CampusLevelNumber != "" {
			addLevel(sc.CampusLevelNumber)
		} else if sc.CurriculumNumber != "" {
			type clResult struct {
				Number string
			}
			var cls []clResult
			config.DB.Table("campus_levels").Where("campus_curriculum_number = ? AND deleted_at IS NULL", sc.CurriculumNumber).Select("number").Scan(&cls)
			for _, cl := range cls {
				addLevel(cl.Number)
			}
		} else if sc.CampusNumber != "" {
			type clResult struct {
				Number string
			}
			var cls []clResult
			config.DB.Table("campus_levels").
				Joins("JOIN campus_curriculums ON campus_curriculums.number = campus_levels.campus_curriculum_number AND campus_curriculums.deleted_at IS NULL").
				Where("campus_curriculums.campus_number = ? AND campus_levels.deleted_at IS NULL", sc.CampusNumber).
				Select("campus_levels.number").Scan(&cls)
			for _, cl := range cls {
				addLevel(cl.Number)
			}
		} else if sc.SchoolNumber != "" {
			type clResult struct {
				Number string
			}
			var cls []clResult
			config.DB.Table("campus_levels").
				Joins("JOIN campus_curriculums ON campus_curriculums.number = campus_levels.campus_curriculum_number AND campus_curriculums.deleted_at IS NULL").
				Joins("JOIN campuses ON campuses.number = campus_curriculums.campus_number AND campuses.deleted_at IS NULL").
				Where("campuses.school_number = ? AND campus_levels.deleted_at IS NULL", sc.SchoolNumber).
				Select("campus_levels.number").Scan(&cls)
			for _, cl := range cls {
				addLevel(cl.Number)
			}
		}
	}

	// 2. Levels from stream_teachers
	for _, teacher := range st {
		if teacher.CampusClassNumber != "" {
			type clResult struct {
				CampusLevelNumber string
			}
			var res clResult
			if err := config.DB.Table("campus_classes").Where("number = ? AND deleted_at IS NULL", teacher.CampusClassNumber).Select("campus_level_number").First(&res).Error; err == nil {
				addLevel(res.CampusLevelNumber)
			}
		}
	}

	return allowed
}

// GetStaffAllowedClasses retrieves all unique CampusClass numbers (CampusClass.Number) allowed for any of the resolved staff IDs.
func GetStaffAllowedClasses(staffIDs []string) []string {
	if len(staffIDs) == 0 {
		return nil
	}

	var staffScopes []models.StaffScope
	config.DB.Where("staff_number IN (?) OR user_number IN (?)", staffIDs, staffIDs).Find(&staffScopes)

	var st []models.StreamTeacher
	config.DB.Where("user_number IN (?)", staffIDs).Find(&st)

	seen := make(map[string]bool)
	var allowed []string

	addClass := func(num string) {
		if num != "" && !seen[num] {
			seen[num] = true
			allowed = append(allowed, num)
		}
	}

	// 1. Classes from staff_scopes or broader scopes
	for _, sc := range staffScopes {
		if sc.ClassNumber != "" {
			addClass(sc.ClassNumber)
		} else if sc.CampusLevelNumber != "" {
			type ccResult struct {
				Number string
			}
			var ccs []ccResult
			config.DB.Table("campus_classes").Where("campus_level_number = ? AND deleted_at IS NULL", sc.CampusLevelNumber).Select("number").Scan(&ccs)
			for _, cc := range ccs {
				addClass(cc.Number)
			}
		} else if sc.CurriculumNumber != "" {
			type ccResult struct {
				Number string
			}
			var ccs []ccResult
			config.DB.Table("campus_classes").
				Joins("JOIN campus_levels ON campus_levels.number = campus_classes.campus_level_number AND campus_levels.deleted_at IS NULL").
				Where("campus_levels.campus_curriculum_number = ? AND campus_classes.deleted_at IS NULL", sc.CurriculumNumber).
				Select("campus_classes.number").Scan(&ccs)
			for _, cc := range ccs {
				addClass(cc.Number)
			}
		} else if sc.CampusNumber != "" {
			type ccResult struct {
				Number string
			}
			var ccs []ccResult
			config.DB.Table("campus_classes").Where("campus_number = ? AND deleted_at IS NULL", sc.CampusNumber).Select("number").Scan(&ccs)
			for _, cc := range ccs {
				addClass(cc.Number)
			}
		} else if sc.SchoolNumber != "" {
			type ccResult struct {
				Number string
			}
			var ccs []ccResult
			config.DB.Table("campus_classes").
				Joins("JOIN campuses ON campuses.number = campus_classes.campus_number AND campuses.deleted_at IS NULL").
				Where("campuses.school_number = ? AND campus_classes.deleted_at IS NULL", sc.SchoolNumber).
				Select("campus_classes.number").Scan(&ccs)
			for _, cc := range ccs {
				addClass(cc.Number)
			}
		}
	}

	// 2. Classes from stream_teachers
	for _, teacher := range st {
		if teacher.CampusClassNumber != "" {
			addClass(teacher.CampusClassNumber)
		}
	}

	// Expand allowed list to include all matching CampusClass.Number, Class.Number, and CurriculumClass.Number
	if len(allowed) > 0 {
		type classTrio struct {
			CampusClassNumber     string
			ClassNumber           string
			CurriculumClassNumber string
		}
		var trios []classTrio
		config.DB.Table("campus_classes").
			Joins("JOIN curriculum_classes ON curriculum_classes.number = campus_classes.curriculum_class_number AND curriculum_classes.deleted_at IS NULL").
			Where("campus_classes.number IN (?) OR curriculum_classes.class_number IN (?) OR curriculum_classes.number IN (?)", allowed, allowed, allowed).
			Select("campus_classes.number AS campus_class_number, curriculum_classes.class_number AS class_number, curriculum_classes.number AS curriculum_class_number").
			Scan(&trios)
		for _, t := range trios {
			addClass(t.CampusClassNumber)
			addClass(t.ClassNumber)
			addClass(t.CurriculumClassNumber)
		}
	}

	return allowed
}

// GetStaffAllowedStreams retrieves all unique Stream numbers (Stream.Number) allowed for any of the resolved staff IDs.
func GetStaffAllowedStreams(staffIDs []string) []string {
	if len(staffIDs) == 0 {
		return nil
	}

	var staffScopes []models.StaffScope
	config.DB.Where("staff_number IN (?) OR user_number IN (?)", staffIDs, staffIDs).Find(&staffScopes)

	var st []models.StreamTeacher
	config.DB.Where("user_number IN (?)", staffIDs).Find(&st)

	seen := make(map[string]bool)
	var allowed []string

	addStream := func(num string) {
		if num != "" && !seen[num] {
			seen[num] = true
			allowed = append(allowed, num)
		}
	}

	// 1. Streams from staff_scopes or broader scopes
	for _, sc := range staffScopes {
		if sc.StreamNumber != "" {
			addStream(sc.StreamNumber)
		} else if sc.ClassNumber != "" {
			type sResult struct {
				Number string
			}
			var ss []sResult
			config.DB.Table("streams").Where("campus_class_number = ? AND deleted_at IS NULL", sc.ClassNumber).Select("number").Scan(&ss)
			for _, s := range ss {
				addStream(s.Number)
			}
		} else if sc.CampusLevelNumber != "" {
			type sResult struct {
				Number string
			}
			var ss []sResult
			config.DB.Table("streams").
				Joins("JOIN campus_classes ON campus_classes.number = streams.campus_class_number AND campus_classes.deleted_at IS NULL").
				Where("campus_classes.campus_level_number = ? AND streams.deleted_at IS NULL", sc.CampusLevelNumber).
				Select("streams.number").Scan(&ss)
			for _, s := range ss {
				addStream(s.Number)
			}
		} else if sc.CurriculumNumber != "" {
			type sResult struct {
				Number string
			}
			var ss []sResult
			config.DB.Table("streams").
				Joins("JOIN campus_classes ON campus_classes.number = streams.campus_class_number AND campus_classes.deleted_at IS NULL").
				Joins("JOIN campus_levels ON campus_levels.number = campus_classes.campus_level_number AND campus_levels.deleted_at IS NULL").
				Where("campus_levels.campus_curriculum_number = ? AND streams.deleted_at IS NULL", sc.CurriculumNumber).
				Select("streams.number").Scan(&ss)
			for _, s := range ss {
				addStream(s.Number)
			}
		} else if sc.CampusNumber != "" {
			type sResult struct {
				Number string
			}
			var ss []sResult
			config.DB.Table("streams").
				Joins("JOIN campus_classes ON campus_classes.number = streams.campus_class_number AND campus_classes.deleted_at IS NULL").
				Where("campus_classes.campus_number = ? AND streams.deleted_at IS NULL", sc.CampusNumber).
				Select("streams.number").Scan(&ss)
			for _, s := range ss {
				addStream(s.Number)
			}
		} else if sc.SchoolNumber != "" {
			type sResult struct {
				Number string
			}
			var ss []sResult
			config.DB.Table("streams").
				Joins("JOIN campus_classes ON campus_classes.number = streams.campus_class_number AND campus_classes.deleted_at IS NULL").
				Joins("JOIN campuses ON campuses.number = campus_classes.campus_number AND campuses.deleted_at IS NULL").
				Where("campuses.school_number = ? AND streams.deleted_at IS NULL", sc.SchoolNumber).
				Select("streams.number").Scan(&ss)
			for _, s := range ss {
				addStream(s.Number)
			}
		}
	}

	// 2. Streams from stream_teachers
	for _, teacher := range st {
		if teacher.StreamNumber != "" {
			addStream(teacher.StreamNumber)
		} else if teacher.CampusClassNumber != "" {
			type sResult struct {
				Number string
			}
			var ss []sResult
			config.DB.Table("streams").Where("campus_class_number = ? AND deleted_at IS NULL", teacher.CampusClassNumber).Select("number").Scan(&ss)
			for _, s := range ss {
				addStream(s.Number)
			}
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

// FilterCampusesByStaffScope is a GORM scope that enforces campus access for staff members.
func FilterCampusesByStaffScope(ctx *gin.Context) func(db *gorm.DB) *gorm.DB {
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

		allowedCampuses := GetStaffAllowedCampuses(staffIDs)
		if len(allowedCampuses) == 0 {
			// Staff member has no campuses assigned -> allow bypassing
			return db
		}

		return db.Where("number IN (?)", allowedCampuses)
	}
}

// FilterCurriculumsByStaffScope is a GORM scope that enforces curriculum access for staff members.
func FilterCurriculumsByStaffScope(ctx *gin.Context) func(db *gorm.DB) *gorm.DB {
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
			return db
		}

		staffIDs, isStaff := ResolveAllStaffIDs(callerNumber)
		if !isStaff {
			return db
		}

		allowedCurriculums := GetStaffAllowedCurriculums(staffIDs)
		if len(allowedCurriculums) == 0 {
			return db
		}

		return db.Where("number IN (?)", allowedCurriculums)
	}
}

// FilterLevelsByStaffScope is a GORM scope that enforces campus level access for staff members.
func FilterLevelsByStaffScope(ctx *gin.Context) func(db *gorm.DB) *gorm.DB {
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
			return db
		}

		staffIDs, isStaff := ResolveAllStaffIDs(callerNumber)
		if !isStaff {
			return db
		}

		allowedLevels := GetStaffAllowedLevels(staffIDs)
		if len(allowedLevels) == 0 {
			return db
		}

		return db.Where("number IN (?)", allowedLevels)
	}
}

// FilterClassesByStaffScope is a GORM scope that enforces campus class access for staff members.
func FilterClassesByStaffScope(ctx *gin.Context, tablePrefixes ...string) func(db *gorm.DB) *gorm.DB {
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
			return db
		}

		staffIDs, isStaff := ResolveAllStaffIDs(callerNumber)
		if !isStaff {
			return db
		}

		allowedClasses := GetStaffAllowedClasses(staffIDs)
		if len(allowedClasses) == 0 {
			return db
		}

		if len(tablePrefixes) == 0 || (len(tablePrefixes) == 1 && tablePrefixes[0] == "") {
			return db.Where("number IN (?)", allowedClasses)
		}

		var conds []string
		var args []interface{}
		for _, prefix := range tablePrefixes {
			if prefix != "" {
				conds = append(conds, prefix+".number IN (?)")
				args = append(args, allowedClasses)
			}
		}
		if len(conds) == 0 {
			return db.Where("number IN (?)", allowedClasses)
		}
		query := ""
		for i, c := range conds {
			if i > 0 {
				query += " OR "
			}
			query += c
		}
		return db.Where("("+query+")", args...)
	}
}

// FilterStreamsByStaffScope is a GORM scope that enforces stream access for staff members.
func FilterStreamsByStaffScope(ctx *gin.Context, tablePrefixes ...string) func(db *gorm.DB) *gorm.DB {
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
			return db
		}

		staffIDs, isStaff := ResolveAllStaffIDs(callerNumber)
		if !isStaff {
			return db
		}

		allowedStreams := GetStaffAllowedStreams(staffIDs)
		if len(allowedStreams) == 0 {
			return db
		}

		if len(tablePrefixes) == 0 || (len(tablePrefixes) == 1 && tablePrefixes[0] == "") {
			return db.Where("number IN (?)", allowedStreams)
		}

		var conds []string
		var args []interface{}
		for _, prefix := range tablePrefixes {
			if prefix != "" {
				conds = append(conds, prefix+".number IN (?)")
				args = append(args, allowedStreams)
			}
		}
		if len(conds) == 0 {
			return db.Where("number IN (?)", allowedStreams)
		}
		query := ""
		for i, c := range conds {
			if i > 0 {
				query += " OR "
			}
			query += c
		}
		return db.Where("("+query+")", args...)
	}
}

// GetDebugInfo returns diagnostic information about the caller and their resolved school/campus scope.
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

	allowedSchools := GetStaffAllowedSchools(staffIDs)
	allowedCampuses := GetStaffAllowedCampuses(staffIDs)
	allowedCurriculums := GetStaffAllowedCurriculums(staffIDs)
	allowedLevels := GetStaffAllowedLevels(staffIDs)
	allowedClasses := GetStaffAllowedClasses(staffIDs)
	allowedStreams := GetStaffAllowedStreams(staffIDs)
	isScoped := len(allowedSchools) > 0 || len(allowedCampuses) > 0 || len(allowedCurriculums) > 0 || len(allowedLevels) > 0 || len(allowedClasses) > 0 || len(allowedStreams) > 0
	msg := "Caller is Staff and scoped."
	if !isScoped {
		msg = "Scoping bypassed: Staff member has no schools, campuses, curriculums, levels, classes, or streams assigned in their scope."
	}

	return gin.H{
		"caller_number":       callerNumber,
		"user_type":           userType,
		"is_staff":            true,
		"is_scoped":           isScoped,
		"resolved_staff_ids":  staffIDs,
		"allowed_schools":     allowedSchools,
		"allowed_campuses":    allowedCampuses,
		"allowed_curriculums": allowedCurriculums,
		"allowed_levels":      allowedLevels,
		"allowed_classes":     allowedClasses,
		"allowed_streams":     allowedStreams,
		"debug_message":       msg,
	}
}


// GetCampusDebugInfo returns diagnostic information about the caller and their resolved campus scope.
func GetCampusDebugInfo(ctx *gin.Context) gin.H {
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
			"debug_message": "Campus scoping bypassed: caller is not Staff usertype (user_type: " + userType + ").",
		}
	}

	staffIDs, isStaff := ResolveAllStaffIDs(callerNumber)
	if !isStaff {
		return gin.H{
			"caller_number": callerNumber,
			"user_type":     userType,
			"is_staff":      false,
			"is_scoped":     false,
			"debug_message": "Campus scoping bypassed: Caller is Staff usertype but has no staff records setup.",
		}
	}

	allowedCampuses := GetStaffAllowedCampuses(staffIDs)
	isScoped := len(allowedCampuses) > 0
	msg := "Caller identified as Staff usertype and scoped to assigned campuses."
	if !isScoped {
		msg = "Campus scoping bypassed: Staff member has no campuses assigned in their scope."
	}

	return gin.H{
		"caller_number":      callerNumber,
		"user_type":          userType,
		"is_staff":           true,
		"is_scoped":          isScoped,
		"resolved_staff_ids": staffIDs,
		"allowed_campuses":   allowedCampuses,
		"debug_message":      msg,
	}
}

// GetCurriculumDebugInfo returns diagnostic information about the caller and their resolved curriculum scope.
func GetCurriculumDebugInfo(ctx *gin.Context) gin.H {
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
			"debug_message": "Curriculum scoping bypassed: caller is not Staff usertype (user_type: " + userType + ").",
		}
	}

	staffIDs, isStaff := ResolveAllStaffIDs(callerNumber)
	if !isStaff {
		return gin.H{
			"caller_number": callerNumber,
			"user_type":     userType,
			"is_staff":      false,
			"is_scoped":     false,
			"debug_message": "Curriculum scoping bypassed: Caller is Staff usertype but has no staff records setup.",
		}
	}

	allowedCurriculums := GetStaffAllowedCurriculums(staffIDs)
	isScoped := len(allowedCurriculums) > 0
	msg := "Caller identified as Staff usertype and scoped to assigned curriculums."
	if !isScoped {
		msg = "Curriculum scoping bypassed: Staff member has no curriculums assigned in their scope."
	}

	return gin.H{
		"caller_number":       callerNumber,
		"user_type":           userType,
		"is_staff":            true,
		"is_scoped":           isScoped,
		"resolved_staff_ids":  staffIDs,
		"allowed_curriculums": allowedCurriculums,
		"debug_message":       msg,
	}
}

// GetLevelDebugInfo returns diagnostic information about the caller and their resolved campus level scope.
func GetLevelDebugInfo(ctx *gin.Context) gin.H {
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
			"debug_message": "Level scoping bypassed: caller is not Staff usertype (user_type: " + userType + ").",
		}
	}

	staffIDs, isStaff := ResolveAllStaffIDs(callerNumber)
	if !isStaff {
		return gin.H{
			"caller_number": callerNumber,
			"user_type":     userType,
			"is_staff":      false,
			"is_scoped":     false,
			"debug_message": "Level scoping bypassed: Caller is Staff usertype but has no staff records setup.",
		}
	}

	allowedLevels := GetStaffAllowedLevels(staffIDs)
	isScoped := len(allowedLevels) > 0
	msg := "Caller identified as Staff usertype and scoped to assigned campus levels."
	if !isScoped {
		msg = "Level scoping bypassed: Staff member has no campus levels assigned in their scope."
	}

	return gin.H{
		"caller_number":      callerNumber,
		"user_type":          userType,
		"is_staff":           true,
		"is_scoped":          isScoped,
		"resolved_staff_ids": staffIDs,
		"allowed_levels":     allowedLevels,
		"debug_message":      msg,
	}
}


// GetClassDebugInfo returns diagnostic information about the caller and their resolved class scope.
func GetClassDebugInfo(ctx *gin.Context) gin.H {
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
			"debug_message": "Class scoping bypassed: caller is not Staff usertype (user_type: " + userType + ").",
		}
	}

	staffIDs, isStaff := ResolveAllStaffIDs(callerNumber)
	if !isStaff {
		return gin.H{
			"caller_number": callerNumber,
			"user_type":     userType,
			"is_staff":      false,
			"is_scoped":     false,
			"debug_message": "Class scoping bypassed: Caller is Staff usertype but has no staff records setup.",
		}
	}

	allowedClasses := GetStaffAllowedClasses(staffIDs)
	isScoped := len(allowedClasses) > 0
	msg := "Caller identified as Staff usertype and scoped to assigned campus classes."
	if !isScoped {
		msg = "Class scoping bypassed: Staff member has no campus classes assigned in their scope."
	}

	return gin.H{
		"caller_number":      callerNumber,
		"user_type":          userType,
		"is_staff":           true,
		"is_scoped":          isScoped,
		"resolved_staff_ids": staffIDs,
		"allowed_classes":    allowedClasses,
		"debug_message":      msg,
	}
}

// GetStreamDebugInfo returns diagnostic information about the caller and their resolved stream scope.
func GetStreamDebugInfo(ctx *gin.Context) gin.H {
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
			"debug_message": "Stream scoping bypassed: caller is not Staff usertype (user_type: " + userType + ").",
		}
	}

	staffIDs, isStaff := ResolveAllStaffIDs(callerNumber)
	if !isStaff {
		return gin.H{
			"caller_number": callerNumber,
			"user_type":     userType,
			"is_staff":      false,
			"is_scoped":     false,
			"debug_message": "Stream scoping bypassed: Caller is Staff usertype but has no staff records setup.",
		}
	}

	allowedStreams := GetStaffAllowedStreams(staffIDs)
	isScoped := len(allowedStreams) > 0
	msg := "Caller identified as Staff usertype and scoped to assigned streams."
	if !isScoped {
		msg = "Stream scoping bypassed: Staff member has no streams assigned in their scope."
	}

	return gin.H{
		"caller_number":      callerNumber,
		"user_type":          userType,
		"is_staff":           true,
		"is_scoped":          isScoped,
		"resolved_staff_ids": staffIDs,
		"allowed_streams":    allowedStreams,
		"debug_message":      msg,
	}
}

