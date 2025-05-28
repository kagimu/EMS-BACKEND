package responses

type ScheduledNarrativeResponse struct {
	Number               string `json:"number"`
	ScheduledOn          string `json:"scheduled_on"`
	AssessmentAreaNumber string `json:"assessment_area_number"`
	AssessmentArea       string `json:"assessment_area"`
	CampusClassNumber    string `json:"campus_class_number"`
	Class                string `json:"class"`
	StudyPeriodNumber    string `json:"study_period_number"`
	StudyPeriod          string `json:"study_period"`
	Description          string `json:"description"`
}
