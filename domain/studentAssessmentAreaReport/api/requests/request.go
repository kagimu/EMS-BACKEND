package requests

// StudentAssessmentAreaReportRequest holds query params for the report.
// AssessmentAreaNumber is required; others narrow the result.
type StudentAssessmentAreaReportRequest struct {
	AssessmentAreaNumber string `form:"assessment_area_number" binding:"required"`
	StudyPeriodNumber    string `form:"study_period_number"`
	CampusClassNumber    string `form:"campus_class_number"`
}
