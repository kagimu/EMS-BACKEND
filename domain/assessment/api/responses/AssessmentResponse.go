package responses

type NarrativeResponse struct {
	Number         string `json:"number"`
	Date           string `json:"date"`
	Teacher        string `json:"teacher"`
	AcademicYear   string `json:"academic_year"`
	StudyPeriod    string `json:"study_period"`
	Class          string `json:"class"`
	Narration      string `json:"narration"`
	AssessmentArea string `json:"assessment_area"`
}
