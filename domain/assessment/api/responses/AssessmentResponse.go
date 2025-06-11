package responses

type NarrativeResponse struct {
	Number               string `json:"number"`
	AcademicYear         string `json:"academic_year"`
	StudyPeriod          string `json:"study_period"`
	Class                string `json:"class"`
	AssessmentArea       string `json:"assessment_area"`
	AssessmentAreaNumber string `json:"assessment_area_number"`
	NarrativeResponseRow
}

type NarrativeResponseRow struct {
	ReferenceDate string `json:"reference_date"`
	ReferenceTime string `json:"reference_time"`
	Narration     string `json:"narration"`
	CreatedAt     string `json:"created_at"`
	Teacher       string `json:"teacher"`
	Status        string `json:"status"`
	CurrentStatus string `json:"current_status"`
}

type AssessmentResponse struct {
	AssessmentAreaName     string  `json:"assessment_area_name"`
	Number                 string  `json:"number"`
	AssessmentMetricNumber string  `json:"assessment_metric_number"`
	Comment                string  `json:"comment"`
	ActualScore            float64 `json:"actual_score"`
	Score                  float64 `json:"score"`
}

type ScriptMarkingGuideResponse struct {
	AssessmentMetricNumber string `json:"assessment_metric_number"`
	ScriptNumber           string `json:"script_number"`
	MarkingGuideNumber     string `json:"marking_guide_number"`
	AddedBy                string `json:"added_by"`
	CreatedOn              string `json:"created_on"`
	UpdatedOn              string `json:"updated_on"`
}
