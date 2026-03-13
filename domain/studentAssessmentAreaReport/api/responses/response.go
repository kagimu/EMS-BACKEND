package responses

// AssessmentScheduleScore represents one assessment schedule (metric) with the student's score.
type AssessmentScheduleScore struct {
	AssessmentMetricNumber string  `json:"assessment_metric_number"`
	Description            string  `json:"description"`
	ScheduledOn            string  `json:"scheduled_on"`
	MaxScore               float64 `json:"max_score"`
	ActualScore            float64 `json:"actual_score"`
	Score                  float64 `json:"score"`
	Comment                string  `json:"comment"`
}

// StudentWithScores is a student with their scores per assessment schedule.
type StudentWithScores struct {
	StudentNumber   string                    `json:"student_number"`
	Surname         string                    `json:"surname"`
	SecondName      string                    `json:"second_name"`
	OtherName       string                    `json:"other_name"`
	ADN             string                    `json:"adn"`
	APN             string                    `json:"apn"`
	CampusClassName string                    `json:"campus_class_name,omitempty"`
	Scores          []AssessmentScheduleScore `json:"scores"`
}

// AssessmentScheduleSummary is a schedule (metric) summary for the report header.
type AssessmentScheduleSummary struct {
	AssessmentMetricNumber string  `json:"assessment_metric_number"`
	Description            string  `json:"description"`
	ScheduledOn            string  `json:"scheduled_on"`
	MaxScore               float64 `json:"max_score"`
}

// StudentAssessmentAreaReportResponse is the full report: assessment area, schedules, and students with scores.
type StudentAssessmentAreaReportResponse struct {
	AssessmentAreaNumber string                      `json:"assessment_area_number"`
	AssessmentAreaName   string                      `json:"assessment_area_name"`
	Schedules           []AssessmentScheduleSummary `json:"schedules"`
	Students             []StudentWithScores        `json:"students"`
}
