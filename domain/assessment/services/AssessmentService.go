package services

import (
	"errors"
	"fmt"
	repositories3 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/academicYear/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/models"
	assessmentMetricRepositories "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/repositories"
	userModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	repository3 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/repository"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
	repositories4 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/assessmentArea/repositories"
	repository2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/enrollment/repository"
	studentRepositories "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/repositories"
	repositories2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/studyPeriod/repositories"
	uuid "github.com/satori/go.uuid"
	"strconv"
	"time"
)

var assessmentMetricRepository = assessmentMetricRepositories.NewAssessmentMetricRepository()

func SaveAssessment(assessmentRequest requests.AssessmentRequest, currentUserNumber string) error {

	student, err := studentRepositories.FindStudentByNumber(assessmentRequest.StudentNumber)
	if err != nil {
		return errors.New("checking student error: " + err.Error())
	}

	if student.Number == "" {
		return errors.New("student not found")
	}

	assessmentMetric, err := assessmentMetricRepository.FindAssessmentMetricByNumber(assessmentRequest.AssessmentMetricNumber)
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
	enrollment, err := repository2.GetEnrollmentByStudentNumber(student.Number)
	if err != nil {
		return errors.New("checking student enrollment error: " + err.Error())
	}

	if enrollment.Number == "" {

		if student.CampusClassNumber != assessmentMetric.CampusClassNumber {
			return errors.New("You can not add results on this schedule because this assessment is not meant  for this student")
		}
		//return errors.New("student is not enrolled yet")
	} else {
		if enrollment.CampusClassNumber != assessmentMetric.CampusClassNumber {
			return errors.New("You can not add results on this schedule because this assessment is not meant  for this student on the current enrollment")
		}

	}

	//find if assessment with student number and assessment metric  already exists
	currentAssessment, err := assessmentMetricRepositories.GetAssessmentByStudentNumberAndAssessmentMetricNumber(student.Number, assessmentMetric.Number)
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

func FetchLearnerNarrative(studentNumber string) ([]responses.NarrativeResponse, error) {

	student, err := studentRepositories.FindStudentByNumber(studentNumber)
	if err != nil {
		return nil, errors.New("checking student number error: " + err.Error())
	}

	if student.Number == "" {
		return nil, errors.New("student not found")
	}

	var assessmentsResponse = make([]responses.NarrativeResponse, 0)

	var assessments []models.Assessment

	assessments, err = assessmentMetricRepositories.FindAssessmentByStudentNumberAndAssessmentMetricNumber(studentNumber, "")
	if err != nil {
		return assessmentsResponse, errors.New("checking student assessments error: " + err.Error())
	}

	for _, assessment := range assessments {
		assessmentsResponse = append(assessmentsResponse, buildNarrativeResponse(assessment))
	}

	return assessmentsResponse, nil
}

func buildNarrativeResponse(assessment models.Assessment) responses.NarrativeResponse {

	user, err := repository3.GetUserByNumber(assessment.UserNumber)
	if err != nil {
		fmt.Println("user not found")
		user = userModels.User{}
	}

	assessmentMetric, err := assessmentMetricRepository.FindAssessmentMetricByNumber(assessment.AssessmentMetricNumber)
	if err != nil {
		fmt.Println("assessment metric not found")
		return responses.NarrativeResponse{}
	}

	if assessmentMetric.Number == "" {
		return responses.NarrativeResponse{}
	}

	assessmentArea, err := repositories4.GetAssessmentAreaByNumber(assessmentMetric.AssessmentAreaNumber)

	class, err := repositories.GetCampusClassByNumber(assessmentMetric.CampusClassNumber)
	if err != nil {
		return responses.NarrativeResponse{}
	}

	studyPeriod, err := repositories2.GetStudyPeriodByNumber(assessmentMetric.StudyPeriodNumber)
	if err != nil {
		return responses.NarrativeResponse{}
	}

	academicYear, err := repositories3.GetAcademicYearByNumber(studyPeriod.AcademicYearNumber)
	if err != nil {
		return responses.NarrativeResponse{}
	}

	studyPeriodName, _ := repositories2.GetStudyPeriodNameByNumber(studyPeriod.StudyPeriodNameNumber)

	return responses.NarrativeResponse{
		Date:           assessment.CreatedAt.Format(time.DateOnly),
		Number:         assessment.Number,
		Narration:      assessment.Comment,
		Teacher:        user.FirstName + " " + user.LastName,
		AssessmentArea: assessmentArea.Name,
		Class:          class.ClassName,
		StudyPeriod:    studyPeriodName.Name,
		AcademicYear:   academicYear.AcademicYearName,
	}
}
