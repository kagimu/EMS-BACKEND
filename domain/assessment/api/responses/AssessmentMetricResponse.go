package responses

type AssessmentMetricResponse struct {
	MetricName           string  `json:"metric_name"`
	Number               string  `json:"number"`
	Description          string  `json:"description"`
	MaxScore             float64 `json:"max_score"`
	ScheduledOn          string  `json:"scheduled_on"`
	AssessmentAreaNumber string  `json:"assessment_area_number"`
	AssessmentArea       string  `json:"assessment_area"`
	MetricNumber         string  `json:"metric_number"`
	StudyPeriodNumber    string  `json:"study_period_number"`
	CampusClassNumber    string  `json:"campus_class_number"`
	Class                string  `json:"class"`
	RequireComment       string  `json:"require_comment"`
	CampusLevelNumber    string  `json:"campus_level_number"`
	AcademicYearNumber   string  `json:"academic_year_number"`
	GradeStandardNumber  string  `json:"grade_standard_number"`
}
