package service

import (
	"errors"
	repositories2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/interview/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/repositories"
)

func FetchInterviewAssessmentMetricByStudent(studentNumber string) ([]responses.StudentInterviewResults, error) {

	var studentInterviewResults []responses.StudentInterviewResults

	student, err := repositories.FindStudentByUserNumber(studentNumber)
	if err != nil {
		return studentInterviewResults, errors.New("checking student number error: " + err.Error())
	}

	if student.Number == "" {
		return studentInterviewResults, errors.New("invalid student number")
	}

	interview, err := GetInt

	assessments, err := repositories2.FindAssessmentByStudentNumberAndAssessmentMetricNumber(student.Number, "")
	if err != nil {
		return studentInterviewResults, errors.New("checking student assessments error: " + err.Error())
	}

}
