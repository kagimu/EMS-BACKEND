package service

import (
	"errors"
	"fmt"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/grades/repository"
	assessmentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/models"
	assessmentRepositories "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/repositories"
	classRepositories "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/repositories"
	assessmentAreaRepositories "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/assessmentArea/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/interview/api/requests"
	interviewResponses "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/interview/api/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/interview/models"
	interviewRepositories "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/interview/repositories"
	repository2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/interview/repository"
	studentRepositories "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/repositories"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"strconv"
)

var assessmentMetricRepository = assessmentRepositories.NewAssessmentMetricRepository()

var interviewRepository = repository2.NewInterviewRepoImpl(config.DB)

func FetchInterviewAssessmentMetricByStudent(studentNumber, interviewNumber string) ([]interviewResponses.StudentInterviewResults, error) {

	var studentInterviewResults []interviewResponses.StudentInterviewResults

	student, err := studentRepositories.FindStudentByUserNumber(studentNumber)
	if err != nil {
		return studentInterviewResults, errors.New("checking student number error: " + err.Error())
	}

	if student.Number == "" {
		return studentInterviewResults, errors.New("invalid student number")
	}

	assessments, err := assessmentRepositories.FindAssessmentByStudentNumberAndAssessmentMetricNumber(student.Number, "")
	if err != nil {
		return studentInterviewResults, errors.New("checking student assessments error: " + err.Error())
	}

	fmt.Println("Assessments attempted: " + strconv.Itoa(len(assessments)))

	interviewAssessmentMetrics, err := interviewRepositories.FindInterviewAssessmentMetricByInterviewNumberAndCampusClassNumber(interviewNumber, student.CampusClassNumber)
	if err != nil {
		return studentInterviewResults, errors.New("getting interview assessment metric error: " + err.Error())
	}

	for _, interviewAssessmentMetric := range interviewAssessmentMetrics {

		assessmentMetric, err := assessmentMetricRepository.FindAssessmentMetricByNumber(interviewAssessmentMetric.AssessmentMetricNumber)
		if err != nil {
			return studentInterviewResults, errors.New("getting assessment metric error: " + err.Error())
		}

		assessmentAreaName, err := assessmentAreaRepositories.GetAssessmentAreaByNumber(assessmentMetric.AssessmentAreaNumber)
		if err != nil {
			return studentInterviewResults, errors.New("error getting assessment area: " + err.Error())
		}

		if assessmentAreaName.Number == "" {
			return studentInterviewResults, errors.New("no assessment area found")
		}

		class, err := classRepositories.GetClassByCampusClassNumber(assessmentMetric.CampusClassNumber)
		if err != nil {
			return studentInterviewResults, errors.New("error getting class: " + err.Error())
		}

		if class.Number == "" {
			return studentInterviewResults, errors.New("campus class number " + assessmentMetric.CampusClassNumber + " has no class")
		}

		assessment, err := assessmentRepositories.GetAssessmentByStudentNumberAndAssessmentMetricNumber(studentNumber, assessmentMetric.Number)
		if err != nil {
			return studentInterviewResults, errors.New("error getting learner scores: " + err.Error())
		}

		gradingStandards, err := assessmentRepositories.FindGradingStandardByCampusNumberAndCampusClassNumber("", class.Number, "")
		if err != nil {
			return studentInterviewResults, errors.New("error getting " + class.Name + "'s grading standard: " + err.Error())
		}

		if len(gradingStandards) == 0 {
			return studentInterviewResults, errors.New("no grading standard set for: " + class.Name)
		}

		if len(gradingStandards) > 1 {
			return studentInterviewResults, errors.New("Class " + class.Name + " has more than one standard set, choose the default standard to use.")
		}

		aggregateScales, err := assessmentRepositories.FindAggregateScaleByClassNumberAndGradingStandardNumber(class.Number, gradingStandards[0].Number)
		if err != nil {
			return studentInterviewResults, errors.New("error getting aggregate scale: " + err.Error())
		}

		var comment string
		if assessment.Comment == "" {
			comment = "No comment"
		} else {
			comment = assessment.Comment
		}

		if len(aggregateScales) == 0 {
			studentInterviewResults = append(studentInterviewResults, createInterviewResult(assessmentAreaName.Name, comment, "Scale not set", assessment.ActualScore))

		} else {
			for _, aggregateScale := range aggregateScales {

				if assessment.ActualScore >= aggregateScale.MinValue && assessment.ActualScore <= aggregateScale.MaxValue {

					grade, err := repository.FindGradeByNumber(aggregateScale.GradeNumber)

					if err != nil {
						return studentInterviewResults, errors.New("error resolving assessment grade: " + err.Error())
					}

					if grade.Number == "" {
						return studentInterviewResults, errors.New("invalid assessment grade number: " + aggregateScale.GradeNumber)
					}

					studentInterviewResults = append(studentInterviewResults, createInterviewResult(assessmentAreaName.Name, comment, grade.Name, assessment.ActualScore))
				}
			}
		}
	}

	return studentInterviewResults, nil
}

func createInterviewResult(name, comment, grade string, score float64) interviewResponses.StudentInterviewResults {
	return interviewResponses.StudentInterviewResults{
		Name:    name,
		Comment: comment,
		Score:   score,
		Grade:   grade,
	}
}

func CreateInterview(payload requests.InterviewRequest, currentUser string) error {

	tx := config.DB.Begin()

	var repo = repository2.NewInterviewRepoImpl(tx)

	interviewNumber := uuid.NewV4().String()
	saveInterview := models.Interview{
		Number:            interviewNumber,
		Message:           payload.Message,
		Venue:             payload.Venue,
		Name:              payload.Name,
		InterviewDate:     payload.InterviewDate,
		UserNumber:        currentUser,
		CampusClassNumber: payload.CampusClassNumber,
	}

	if err := repo.CreateInterview(&saveInterview); err != nil {
		tx.Rollback()
		return errors.New("creating interview error: " + err.Error())
	}

	var interview models.Interview
	if err := tx.Where("number = ?", interviewNumber).First(&interview).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Interview schedule not found"})
		return
	}

	for _, assessmentMetricNumber := range payload.AssessmentMetrics {

		var assessmentMetric assessmentModel.AssessmentMetric
		if err := tx.Where("number = ?", assessmentMetricNumber).First(&assessmentMetric).Error; err != nil {
			continue
		}

		var interviewAssessmentMetric models.InterviewAssessmentMetric
		if err := tx.Where("interview_number = ? AND assessment_metric_number = ?", interviewNumber, assessmentMetricNumber).First(&interviewAssessmentMetric).Error; err == nil {
			continue
		}

		interviewAssessmentMetric = models.InterviewAssessmentMetric{
			Number:                 uuid.NewV4().String(),
			InterviewNumber:        interviewNumber,
			AssessmentMetricNumber: assessmentMetricNumber,
			UserNumber:             user.Number,
		}

		if err := tx.Create(&interviewAssessmentMetric).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to create interview assessment metric"})
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to commit transaction"})
		return

	}

}

func FindInterviewStudentsByInterViewNumber(interviewNumber string) error {

}
