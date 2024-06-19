package services

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/models"
	assessmentMetricRepositories "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/enrollment/repository"
	studentRepositories "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/repositories"
	uuid "github.com/satori/go.uuid"
	"strconv"
)

func SaveAssessment(assessmentRequest requests.AssessmentRequest, currentUserNumber string) error {

	student, err := studentRepositories.FindStudentByNumber(assessmentRequest.StudentNumber)
	if err != nil {
		return errors.New("checking student error: " + err.Error())
	}

	if student.Number == "" {
		return errors.New("student not found")
	}

	assessmentMetric, err := assessmentMetricRepositories.FindAssessmentMetricByNumber(assessmentRequest.AssessmentMetricNumber)
	if err != nil {
		return errors.New("checking assessment metric error: " + err.Error())
	}

	if assessmentMetric.Number == "" {
		return errors.New("assessment metric not found")
	}

	if assessmentRequest.ActualScore > assessmentMetric.MaxScore {
		return errors.New("score exceeds the max value set at: " + strconv.Itoa(int(assessmentMetric.MaxScore)))
	}

	var actualScore float64
	if assessmentMetric.MaxScore > 0 && assessmentRequest.ActualScore > 0 {
		actualScore = (assessmentRequest.ActualScore / assessmentMetric.MaxScore) * 100
	} else {
		actualScore = 0
	}

	if assessmentMetric.RequireComment == constants.StatusYes && assessmentRequest.Comment == "" {
		return errors.New("this assessment requires a comment")
	}

	//check if there is existing enrollment for the student
	enrollment, err := repository.FindEnrollmentByStudentNumber(student.Number)
	if err != nil {
		return errors.New("checking student enrollment error: " + err.Error())
	}

	if enrollment.Number == "" {
		return errors.New("student is not enrolled yet")
	}

	if student.CampusClassNumber != assessmentMetric.CampusClassNumber {
		return errors.New("this assessment is not for this student")
	}

	//find if assessment with student number and assessment metric  already exists
	currentAssessment, err := assessmentMetricRepositories.FindAssessmentByStudentNumberAndAssessmentMetricNumber(student.Number, assessmentMetric.Number)
	if err != nil {
		return errors.New("checking student current score error: " + err.Error())
	}

	if currentAssessment.Number != "" {
		currentAssessment.Comment = assessmentRequest.Comment
		currentAssessment.ActualScore = actualScore
		currentAssessment.UserNumber = currentUserNumber
		currentAssessment.Score = assessmentRequest.ActualScore

		err = assessmentMetricRepositories.SaveAssessment(currentAssessment)
		if err != nil {
			return errors.New("updating student score error: " + err.Error())
		}

	} else {

		saveAssessment := models.Assessment{
			Number:                 uuid.NewV4().String(),
			UserNumber:             currentUserNumber,
			Comment:                assessmentRequest.Comment,
			AssessmentMetricNumber: assessmentRequest.AssessmentMetricNumber,
			StudentNumber:          assessmentRequest.StudentNumber,
			ActualScore:            actualScore,
			Score:                  assessmentRequest.ActualScore,
		}

		err = assessmentMetricRepositories.SaveAssessment(saveAssessment)
		if err != nil {
			return errors.New("saving student score error: " + err.Error())
		}
	}

	return nil
}
