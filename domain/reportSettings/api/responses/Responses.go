package responses

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/models"
)

type ReportSettingsResponse struct {
	Number       string `json:"number"`
	Name         string `json:"name"`
	Columns      int64  `json:"columns"`
	Status       string `json:"status"`
	Code         int64  `json:"code"`
	StudyPeriod  string `json:"study_period"`
	AcademicYear string `json:"academic_year"`
	Class        string `json:"class"`
}

type ReportColumnResponse struct {
	Number            string                           `json:"number"`
	Name              string                           `json:"name"`
	ReportSettingName string                           `json:"report_setting_name"`
	PercentageShare   float64                          `json:"percentage_share"`
	Code              int64                            `json:"code"`
	Assessments       []ReportColumnAssessmentResponse `json:"assessments"`
}

type ReportColumnAssessmentResponse struct {
	Number           string `json:"number"`
	AssessmentArea   string `json:"assessment_area"`
	AssessmentMetric string `json:"assessment_metric"`
	Code             int64  `json:"code"`
	ReportColumn     string `json:"report_column"`
	ScheduledOn      string `json:"scheduled_on"`
}

type SummarizedFinalReportResponse struct {
	Student           string `json:"student"`
	StudentNumber     string `json:"student_number"`
	Grade             string `json:"grade"`
	Aggregate         string `json:"aggregate"`
	Class             string `json:"class"`
	CampusClassNumber string `json:"campus_class_number"`
	EnrollmentNumber  string `json:"enrollment_number"`
}

type LeanReportSummary struct {
	Number        string `json:"number"`
	StudentName   string `json:"student_name"`
	StudentNumber string `json:"student_number"`
	Grade         string `json:"grade"`
	StudyPeriod   string `json:"study_period"`
	Stream        string `json:"stream"`
	Class         string `json:"class"`
}

type LeanStudent struct {
	Number     string `json:"number"`
	Surname    string `json:"surname"`
	SecondName string `json:"second_name"`
	OtherName  string `json:"other_name"`
	Gender     string `json:"gender"`
}

type Assessments struct {
	Name               string                    `json:"name"`
	ReportScores       []models.ReportScore      `json:"report_scores"`
	ReportScoreSummary models.ReportScoreSummary `json:"report_score_summary"`
}

type LearningArea struct {
	Name        string        `json:"name"`
	Assessments []Assessments `json:"assessments"`
}

type AppliedScale struct {
	ShortName string  `json:"short_name"`
	Name      string  `json:"name"`
	Min       float64 `json:"min"`
	Max       float64 `json:"max"`
}
type FullReportSummary struct {
	Student             LeanStudent    `json:"student"`
	Grade               string         `json:"grade"`
	LearningAreas       []LearningArea `json:"learning_areas"`
	Class               string         `json:"class"`
	ClassAbbreviation   string         `json:"class_abbreviation"`
	Stream              string         `json:"stream"`
	SchoolName          string         `json:"school_name"`
	CampusName          string         `json:"campus_name"`
	LevelName           string         `json:"level_name"`
	School              LeanSchool     `json:"school"`
	Campus              LeanCampus     `json:"campus"`
	StudyPeriod         string         `json:"study_period"`
	AcademicYear        string         `json:"academic_year"`
	ReportColumns       []string       `json:"report_columns"`
	ReportTemplateName  string         `json:"report_template_name"`
	TeacherComment      string         `json:"teacher_comment"`
	HeadTeachersComment string         `json:"head_teachers_comment"`
	AggregateScale      []AppliedScale `json:"aggregate_scale"`
}

type LeanSchool interface{}

type LeanCampus interface{}

type HeadTeacherCommentResponse struct {
	Comment      string  `json:"comment"`
	AcademicYear string  `json:"academic_year"`
	MinAggregate float64 `json:"min_aggregate"`
	MaxAggregate float64 `json:"max_aggregate"`
	StudyPeriod  string  `json:"study_period"`
	Class        string  `json:"class"`
}
