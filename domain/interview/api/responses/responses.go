package responses

import "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/api/responses"

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

type InterviewResponse struct {
	InterviewNumber   string `json:"interview_number"`
	Message           string `json:"message"`
	Venue             string `json:"venue"`
	Name              string `json:"name"`
	InterviewDate     string `json:"interview_date"`
	CampusClassNumber string `json:"campus_class_number"`
	Class             string `json:"class"`
	SchoolNumber      string `json:"school_number"`
	CampusNumber      string `json:"campus_number"`
}

type InterviewResponseDetails struct {
	InterviewResponse
	AssessmentSchedules []responses.AssessmentMetricResponse `json:"assessment_schedules"`
}
