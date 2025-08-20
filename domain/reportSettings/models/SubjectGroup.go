package models

type SubjectGroup struct {
	Number                 string `json:"number" gorm:"unique"`
	SubjectGroupName       string `json:"subject_group_name"`
	ReportSettingsNumber   string `json:"report_settings_number"`
	ClassNumber            string `json:"class_number"`
	CampusNumber           string `json:"campus_number"`
	CurriculumNumber       string `json:"curriculum_number"`
	CampusCurriculumNumber string `json:"campus_curriculum_number"`
}

type GroupAssessmentArea struct {
	AssessmentAreaNumber string `json:"assessment_area_number"`
	GroupNumber          string `json:"group_number"`
	SubjectNameOnReport  string `json:"subject_name_on_report"`
}
