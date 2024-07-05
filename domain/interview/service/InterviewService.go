package service

import (
	"errors"
	"fmt"
	assessmentRepositories "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/repositories"
	classRepositories "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/repositories"
	assessmentAreaRepositories "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/assessmentArea/repositories"
	interviewResponses "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/interview/api/responses"
	interviewRepositories "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/interview/repositories"
	studentRepositories "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/repositories"
	"strconv"
)

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

		assessmentMetric, err := assessmentRepositories.FindAssessmentMetricByNumber(interviewAssessmentMetric.AssessmentMetricNumber)
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

		aggregateScales, err := assessmentRepositories.FindAggregateScaleByClassNumberAndGradingStandardNumber(class.Number, assessmentMetric.GradingStandardNumber)
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
					assessmentGrade, err := assessmentRepositories.GetAssessmentGradeByNumber(aggregateScale.AssessmentGradeNumber)

					if err != nil {
						return studentInterviewResults, errors.New("error resolving assessment grade: " + err.Error())
					}

					if assessmentGrade.Number == "" {
						return studentInterviewResults, errors.New("invalid assessment grade number: " + aggregateScale.AssessmentGradeNumber)
					}

					studentInterviewResults = append(studentInterviewResults, createInterviewResult(assessmentAreaName.Name, comment, assessmentGrade.Name, assessment.ActualScore))
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
