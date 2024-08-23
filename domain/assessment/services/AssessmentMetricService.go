package services

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/repositories"
	repositories3 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/repositories"
	classService "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/services"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/pagination"
	repositories2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/assessmentArea/repositories"
	repositories5 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/interview/repositories"
	repositories4 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/studyPeriod/repositories"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"time"
)

var assessmentRepository = repositories.NewAssessmentMetricRepository()

func SaveAssessmentMetric(assessmentMetricRequest requests.AssessmentMetricRequest, currentUser string) error {

	formatDate, err := time.Parse(time.DateOnly, assessmentMetricRequest.ScheduledOn)
	if err != nil {
		return err
	}

	campusClass, err := classService.ResolveCampusClass(assessmentMetricRequest.ClassNumber, assessmentMetricRequest.CampusNumber)
	if err != nil {
		return err
	}

	saveAssessmentMetric := models.AssessmentMetric{
		Number:               uuid.NewV4().String(),
		StudyPeriodNumber:    assessmentMetricRequest.StudyPeriodNumber,
		ScheduledOn:          formatDate,
		CampusClassNumber:    campusClass.Number,
		MetricNumber:         assessmentMetricRequest.MetricNumber,
		AssessmentAreaNumber: assessmentMetricRequest.AssessmentAreaNumber,
		Description:          assessmentMetricRequest.Description,
		MaxScore:             assessmentMetricRequest.MaxScore,
		RequireComment:       assessmentMetricRequest.RequireComment,
		UserNumber:           currentUser,
	}

	err = assessmentRepository.SaveAssessmentMetric(saveAssessmentMetric)
	if err != nil {
		return errors.New("creating schedule error: " + err.Error())
	}

	return nil
}

func AllAssessmentsPaginated(assessmentAreaNumber, studyPeriodNumber, campusClassNumber, metricNumber, classNumber, campusNumber string, ctx *gin.Context) (pagination.PaginatedResponse, error) {

	var paginatedResponse = pagination.PaginatedResponse{}

	var paginate = pagination.Paginate(ctx, &paginatedResponse)

	assessments, err := assessmentRepository.AllAssessmentMetricPaginated(assessmentAreaNumber, studyPeriodNumber, campusClassNumber, metricNumber, classNumber, campusNumber, paginate)
	if err != nil {
		return paginatedResponse, errors.New("error getting list of assessments: " + err.Error())
	}

	var assessmentsResponse = make([]responses.AssessmentMetricResponse, 0)

	for _, assessment := range assessments {

		response, _ := CreateAssessmentsResponse(assessment)

		assessmentsResponse = append(assessmentsResponse, response)
	}

	paginatedResponse.Content = assessmentsResponse

	return paginatedResponse, nil
}

func AllAssessments(assessmentAreaNumber, studyPeriodNumber, campusClassNumber, metricNumber, classNumber, campusNumber string) ([]responses.AssessmentMetricResponse, error) {

	var assessmentsResponse = make([]responses.AssessmentMetricResponse, 0)

	assessments, err := assessmentRepository.AllAssessmentMetric(assessmentAreaNumber, studyPeriodNumber, campusClassNumber, metricNumber, classNumber, campusNumber)
	if err != nil {
		return assessmentsResponse, errors.New("error getting list of assessments: " + err.Error())
	}

	for _, assessment := range assessments {

		response, _ := CreateAssessmentsResponse(assessment)

		assessmentsResponse = append(assessmentsResponse, response)
	}

	return assessmentsResponse, nil
}

func CreateAssessmentsResponse(assessment models.AssessmentMetric) (responses.AssessmentMetricResponse, error) {

	var assessmentMetricResponse responses.AssessmentMetricResponse

	assessmentArea, err := repositories2.GetAssessmentAreaByNumber(assessment.AssessmentAreaNumber)
	if err != nil {
		return assessmentMetricResponse, errors.New("assessment area by number error: " + err.Error())
	}

	metric, err := repositories.GetMetricByNumber(assessment.MetricNumber)
	if err != nil {
		return assessmentMetricResponse, errors.New("getting metric error: " + err.Error())
	}

	class, err := repositories3.GetClassByCampusClassNumber(assessment.CampusClassNumber)
	if err != nil {
		return assessmentMetricResponse, errors.New("getting class error: " + err.Error())
	}

	campusClass, err := repositories3.GetCampusClassByNumber(assessment.CampusClassNumber)
	if err != nil {
		return assessmentMetricResponse, errors.New("failed to find campus class")
	}

	studyPeriod, err := repositories4.GetStudyPeriodByNumber(assessment.StudyPeriodNumber)
	if err != nil {
		return assessmentMetricResponse, errors.New("failed to find academic year")
	}

	assessmentMetricResponse.AssessmentArea = assessmentArea.Name
	assessmentMetricResponse.AssessmentAreaNumber = assessmentArea.Number
	assessmentMetricResponse.MetricName = metric.Name
	assessmentMetricResponse.MetricNumber = metric.Number
	assessmentMetricResponse.Number = assessment.Number
	assessmentMetricResponse.Description = assessment.Description
	assessmentMetricResponse.StudyPeriodNumber = assessment.StudyPeriodNumber
	assessmentMetricResponse.RequireComment = assessment.RequireComment
	assessmentMetricResponse.MaxScore = assessment.MaxScore
	assessmentMetricResponse.Class = class.Name
	assessmentMetricResponse.CampusLevelNumber = campusClass.CampusLevelNumber
	assessmentMetricResponse.AcademicYearNumber = studyPeriod.AcademicYearNumber
	assessmentMetricResponse.CampusClassNumber = assessment.CampusClassNumber
	assessmentMetricResponse.ScheduledOn = assessment.ScheduledOn.Format(time.DateOnly)

	return assessmentMetricResponse, nil
}

func FindInterviewAssessments(interviewNumber, campusClassNumber string) ([]responses.AssessmentMetricResponse, error) {
	var assessmentsResponse = make([]responses.AssessmentMetricResponse, 0)

	assessmentMetrics, err := repositories5.FindInterviewAssessmentMetricByInterviewNumberAndCampusClassNumber(interviewNumber, campusClassNumber)
	if err != nil {
		return assessmentsResponse, errors.New("error getting list of assessments: " + err.Error())
	}

	for _, assessmentMetric := range assessmentMetrics {
		assessment, err := assessmentRepository.FindAssessmentMetricByNumber(assessmentMetric.AssessmentMetricNumber)
		if err != nil {
			return assessmentsResponse, errors.New("error getting assessment details: " + err.Error())
		}

		response, err := CreateAssessmentsResponse(assessment)
		if err != nil {
			return assessmentsResponse, err
		}

		assessmentsResponse = append(assessmentsResponse, response)
	}

	return assessmentsResponse, nil
}

func UpdateAssessmentMetric(assessmentMetricRequest requests.UpdateAssessmentMetricRequest, currentUser string) error {

	formatDate, err := time.Parse(time.DateOnly, assessmentMetricRequest.ScheduledOn)
	if err != nil {
		return err
	}

	assessmentMetric, err := assessmentRepository.FindAssessmentMetricByNumber(assessmentMetricRequest.AssessmentMetricNumber)
	if err != nil {
		return errors.New("error getting assessment details: " + err.Error())
	}

	if assessmentMetric.Number == "" {
		return errors.New("invalid assessment metric number")
	}

	campusClass, err := classService.ResolveCampusClass(assessmentMetricRequest.ClassNumber, assessmentMetricRequest.CampusNumber)
	if err != nil {
		return err
	}

	assessmentMetric.UserNumber = currentUser
	assessmentMetric.MaxScore = assessmentMetricRequest.MaxScore
	assessmentMetric.RequireComment = assessmentMetricRequest.RequireComment
	assessmentMetric.Description = assessmentMetricRequest.Description
	assessmentMetric.AssessmentAreaNumber = assessmentMetricRequest.AssessmentAreaNumber
	assessmentMetric.StudyPeriodNumber = assessmentMetricRequest.StudyPeriodNumber
	assessmentMetric.MetricNumber = assessmentMetricRequest.MetricNumber
	assessmentMetric.ScheduledOn = formatDate
	assessmentMetric.CampusClassNumber = campusClass.Number

	err = assessmentRepository.SaveAssessmentMetric(assessmentMetric)
	if err != nil {
		return errors.New("updating schedule error: " + err.Error())
	}

	return nil

}
