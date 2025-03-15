package responses

type ScheduledNarrativeResponse struct {
	Number               string `json:"number"`
	ScheduledOn          string `json:"scheduled_on"`
	AssessmentAreaNumber string `json:"assessment_area_number"`
	CampusClassNumber    string `json:"campus_class_number"`
	StudyPeriodNumber    string `json:"study_period_number"`
	Description          string `json:"description"`
}
