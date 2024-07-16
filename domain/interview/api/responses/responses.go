package responses

type InterviewQuestionResponse struct {
	InterviewNumber string
	Question        string
	PossibleAnswers string
	CorrectAnswer   string
	Score           int
}

type StudentInterviewResults struct {
	Name    string
	Score   float64
	Grade   string
	Comment string
}

type Assessment struct {
	AssessmentAreaName     string  `json:"assessment_area_name"`
	Number                 string  `json:"number"`
	AssessmentMetricNumber string  `json:"assessment_metric_number"`
	Comment                string  `json:"comment"`
	ActualScore            float64 `json:"actual_score"`
	Score                  float64 `json:"score"`
}
