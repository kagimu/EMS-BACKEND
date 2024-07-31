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

type AssessmentResponse struct {
	AssessmentAreaName     string  `json:"assessment_area_name"`
	Number                 string  `json:"number"`
	AssessmentMetricNumber string  `json:"assessment_metric_number"`
	Comment                string  `json:"comment"`
	ActualScore            float64 `json:"actual_score"`
	Score                  float64 `json:"score"`
}
