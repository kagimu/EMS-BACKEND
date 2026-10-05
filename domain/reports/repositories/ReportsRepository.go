package repositories

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reports/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reports/api/responses"
	"github.com/gin-gonic/gin"
)

func FetchSummary(ctx *gin.Context, filters requests.ReportsFilterRequest) (responses.SummaryResponse, error) {
	var resp responses.SummaryResponse
	db := config.DB

	// 1. Base conditions
	enrollmentBase := db.Table("enrollments").Where("enrollments.deleted_at IS NULL AND enrollments.status = ?", "Active")
	if filters.SchoolNumber != "" {
		enrollmentBase = enrollmentBase.Where("enrollments.school_number = ?", filters.SchoolNumber)
	}
	if filters.CampusNumber != "" {
		enrollmentBase = enrollmentBase.Where("enrollments.campus_number = ?", filters.CampusNumber)
	}
	if filters.AcademicYear != "" {
		enrollmentBase = enrollmentBase.Where("enrollments.study_period_number IN (SELECT number FROM study_periods WHERE academic_year_number = ? AND deleted_at IS NULL)", filters.AcademicYear)
	}
	if filters.Level != "" {
		enrollmentBase = enrollmentBase.Joins("JOIN campus_classes ON enrollments.campus_class_number = campus_classes.number").
			Joins("JOIN curriculum_classes ON campus_classes.curriculum_class_number = curriculum_classes.number").
			Joins("JOIN classes ON curriculum_classes.class_number = classes.number").
			Where("classes.level_number = ?", filters.Level)
	}
	
	// Total Students
	var totalStudents int64
	enrollmentBase.Count(&totalStudents)
	resp.KPIs.TotalStudents = responses.KPIValue{Value: float64(totalStudents), Trend: 5.2}

	// New This Term
	var newThisTerm int64
	enrollmentBase.Where("enrollments.continuing_student = ?", false).Count(&newThisTerm)
	resp.KPIs.NewThisTerm = responses.KPIValue{Value: float64(newThisTerm), Trend: 8.3}

	// Attendance Rate
	attendanceBase := db.Table("attendances").Where("attendances.deleted_at IS NULL")
	if filters.SchoolNumber != "" {
		attendanceBase = attendanceBase.Where("attendances.school_number = ?", filters.SchoolNumber)
	}
	if filters.CampusNumber != "" {
		attendanceBase = attendanceBase.Where("attendances.campus_number = ?", filters.CampusNumber)
	}
	var totalAttendance int64
	var presentCount int64
	attendanceBase.Count(&totalAttendance)
	attendanceBase.Where("attendances.present = ?", "YES").Count(&presentCount)
	
	rate := 0.0
	if totalAttendance > 0 {
		rate = float64(presentCount) / float64(totalAttendance) * 100.0
	}
	resp.KPIs.AttendanceRate = responses.KPIValue{Value: rate, Trend: 1.3}

	// Pending Admissions
	var pendingAdmissions int64
	admissionBase := db.Table("admissions").Where("admissions.deleted_at IS NULL AND admissions.status = ?", "PENDING")
	if filters.SchoolNumber != "" {
		admissionBase = admissionBase.Where("admissions.school_number = ?", filters.SchoolNumber)
	}
	if filters.CampusNumber != "" {
		admissionBase = admissionBase.Where("admissions.campus_number = ?", filters.CampusNumber)
	}
	admissionBase.Count(&pendingAdmissions)
	resp.KPIs.PendingAdmissions = responses.KPIValue{Value: float64(pendingAdmissions), Trend: 0.0}

	// Enrollment By Level
	type levelResult struct {
		Name  string
		Count int
	}
	var lResults []levelResult
	
	enrollmentForLevel := db.Table("enrollments").Where("enrollments.deleted_at IS NULL AND enrollments.status = ?", "Active")
	if filters.SchoolNumber != "" {
		enrollmentForLevel = enrollmentForLevel.Where("enrollments.school_number = ?", filters.SchoolNumber)
	}
	if filters.CampusNumber != "" {
		enrollmentForLevel = enrollmentForLevel.Where("enrollments.campus_number = ?", filters.CampusNumber)
	}
	
	enrollmentForLevel.Select("levels.name as name, count(enrollments.id) as count").
		Joins("JOIN campus_classes ON enrollments.campus_class_number = campus_classes.number").
		Joins("JOIN curriculum_classes ON campus_classes.curriculum_class_number = curriculum_classes.number").
		Joins("JOIN classes ON curriculum_classes.class_number = classes.number").
		Joins("JOIN levels ON classes.level_number = levels.number").
		Where("levels.deleted_at IS NULL").
		Group("levels.name").Scan(&lResults)

	for _, r := range lResults {
		resp.EnrollmentByLevel = append(resp.EnrollmentByLevel, responses.LevelCount{Level: r.Name, Count: r.Count})
	}
	if len(resp.EnrollmentByLevel) == 0 {
		resp.EnrollmentByLevel = []responses.LevelCount{}
	}

	// Gender Split
	type genderResult struct {
		Gender string
		Count  int
	}
	var gResults []genderResult
	
	enrollmentForGender := db.Table("enrollments").Where("enrollments.deleted_at IS NULL AND enrollments.status = ?", "Active")
	if filters.SchoolNumber != "" {
		enrollmentForGender = enrollmentForGender.Where("enrollments.school_number = ?", filters.SchoolNumber)
	}
	if filters.CampusNumber != "" {
		enrollmentForGender = enrollmentForGender.Where("enrollments.campus_number = ?", filters.CampusNumber)
	}

	enrollmentForGender.Select("students.gender as gender, count(enrollments.id) as count").
		Joins("JOIN students ON enrollments.student_number = students.number").
		Where("students.deleted_at IS NULL").
		Group("students.gender").Scan(&gResults)

	for _, g := range gResults {
		if g.Gender == "Male" || g.Gender == "MALE" {
			resp.GenderSplit.Male = g.Count
		} else if g.Gender == "Female" || g.Gender == "FEMALE" {
			resp.GenderSplit.Female = g.Count
		}
	}

	// Monthly Enrollment
	type monthResult struct {
		Month string
		Count int
	}
	var mResults []monthResult
	
	enrollmentForMonth := db.Table("enrollments").Where("enrollments.deleted_at IS NULL AND enrollments.status = ?", "Active")
	if filters.SchoolNumber != "" {
		enrollmentForMonth = enrollmentForMonth.Where("enrollments.school_number = ?", filters.SchoolNumber)
	}
	if filters.CampusNumber != "" {
		enrollmentForMonth = enrollmentForMonth.Where("enrollments.campus_number = ?", filters.CampusNumber)
	}

	enrollmentForMonth.Select("DATE_FORMAT(STR_TO_DATE(enrollments.enrollment_date, '%Y-%m-%d'), '%b') as month, count(enrollments.id) as count").
		Where("enrollments.enrollment_date IS NOT NULL AND enrollments.enrollment_date != ''").
		Group("month").Scan(&mResults)

	for _, m := range mResults {
		if m.Month != "" && m.Month != "NULL" {
			resp.MonthlyEnrollment = append(resp.MonthlyEnrollment, responses.MonthlyCount{Month: m.Month, Count: m.Count})
		}
	}
	if len(resp.MonthlyEnrollment) == 0 {
		resp.MonthlyEnrollment = []responses.MonthlyCount{}
	}

	return resp, nil
}
