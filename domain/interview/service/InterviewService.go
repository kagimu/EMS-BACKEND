package service

import (
	"errors"
	"fmt"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	aggregateSaceltModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/aggregateScale/models"
	assessmentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/models"
	repositories2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/interview/api/responses"
	repositories3 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/interview/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/repositories"
	"strconv"
)

func FetchInterviewAssessmentMetricByStudent(studentNumber, interviewNumber string) ([]responses.StudentInterviewResults, error) {

	var studentInterviewResults []responses.StudentInterviewResults

	student, err := repositories.FindStudentByUserNumber(studentNumber)
	if err != nil {
		return studentInterviewResults, errors.New("checking student number error: " + err.Error())
	}

	if student.Number == "" {
		return studentInterviewResults, errors.New("invalid student number")
	}

	assessments, err := repositories2.FindAssessmentByStudentNumberAndAssessmentMetricNumber(student.Number, "")
	if err != nil {
		return studentInterviewResults, errors.New("checking student assessments error: " + err.Error())
	}

	fmt.Println("Assessments attempted: " + strconv.Itoa(len(assessments)))

	var assessmentMetrics []assessmentModel.AssessmentMetric

	interviewAssessmentMetrics, err := repositories3.FindInterviewAssessmentMetricByInterviewNumberAndCampusClassNumber(interviewNumber, student.CampusClassNumber)
	if err != nil {
		return studentInterviewResults, errors.New("getting interview assessment metric error: " + err.Error())
	}

	for _, v := range interviewAssessmentMetrics {
		assessmentMetric, err := repositories2.FindAssessmentMetricByNumber(v.AssessmentMetricNumber)
		if err != nil {
			return studentInterviewResults, errors.New("getting assessment metric error: " + err.Error())
		}

		assessmentMetrics = append(assessmentMetrics, assessmentMetric)
	}

	for _, assessmentMetric := range assessmentMetrics {
		var aggregateScales []aggregateSaceltModel.AggregateScale
		assessment, err := repositories2.GetAssessmentByStudentNumberAndAssessmentMetricNumber(studentNumber, assessmentMetric.Number)
		if err != nil {
			return studentInterviewResults, errors.New("error getting learner scores: " + err.Error())
		}

		if err := config.DB.Where("grading_standard_number = ?", assessmentMetric.GradingStandardNumber).Find(&aggregateScales).Error; err != nil {
			//c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch aggregate scales"})
			//return
		}
		var comment string
		if assessment.Comment == "" {
			comment = "No comment"
		} else {
			comment = assessment.Comment
		}

		if len(aggregateScales) == 0 {
			studentInterviewResults = append(studentInterviewResults, responses.StudentInterviewResults{
				Name:    assessmentMetric.AssessmentAreaName,
				Comment: comment,
				Score:   int(assessment.ActualScore),
				Grade:   "Not Graded , no scales ",
			})

		} else {
			for _, val := range aggregateScales {
				if assessment.ActualScore >= val.MinValue && assessment.ActualScore <= val.MaxValue {
					studentInterviewResults = append(studentInterviewResults, responses.StudentInterviewResults{
						Name:    assessmentMetric.AssessmentAreaName,
						Score:   int(assessment.ActualScore),
						Grade:   val.Name,
						Comment: comment,
					})
				}
			}
		}
	}

}
