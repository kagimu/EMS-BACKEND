package controllers

import (
	"errors"
	aggregateSaceltModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/aggregateScale/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/filters"
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	assessmentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/interview/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/interview/models"
	studentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

func CreateInterviewAssessmentMetric(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var payload requests.InterviewAssessmentMetricRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid payload"})
		return
	}

	if err := payload.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var interview models.Interview
	if err := config.DB.Where("number = ?", payload.InterviewNumber).First(&interview).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Interview schedule not found"})
		return
	}

	var assessmentMetric assessmentModel.AssessmentMetric
	if err := config.DB.Where("number = ?", payload.AssessmentMetricNumber).First(&assessmentMetric).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Assessment metric not found"})
		return
	}

	var interviewAssessmentMetric models.InterviewAssessmentMetric
	if err := config.DB.Where("interview_number = ? AND assessment_metric_number = ?", payload.InterviewNumber, payload.AssessmentMetricNumber).First(&interviewAssessmentMetric).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Interview assessment metric already exists"})
		return
	}

	interviewAssessmentMetric = models.InterviewAssessmentMetric{
		Number:                 uuid.NewV4().String(),
		InterviewNumber:        payload.InterviewNumber,
		AssessmentMetricNumber: payload.AssessmentMetricNumber,
		UserNumber:             user.Number,
	}

	if err := config.DB.Create(&interviewAssessmentMetric).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to create interview assessment metric"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": interviewAssessmentMetric, "message": "Interview assessment metric created successfully"})
}

func FetchInterviewAssessmentMetrics(c *gin.Context) {
	var interviewAssessmentMetrics []models.InterviewAssessmentMetric
	interviewNumber := c.Param("interview_number")
	if err := config.DB.Scopes(filters.ByInterviewNumber(interviewNumber)).Preload("AssessmentMetric").Preload("Interview").Find(&interviewAssessmentMetrics).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch interview assessment metrics"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": interviewAssessmentMetrics})
}

func FetchInterviewAssessmentMetricByStudent(c *gin.Context) {

	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	studentNumber := c.Query("student_number")
	var student studentModel.Student

	if studentNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Student number is missing"})
		return
	}

	if err := config.DB.First(&student, "number = ?", studentNumber).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": "Student not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		}
		return
	}

	var assessments []assessmentModel.Assessment
	if err := config.DB.Where("student_number = ?", studentNumber).Preload("AssessmentMetric").Preload("Student").Find(&assessments).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Assessment not found"})
		return
	}

	interviewNumber := c.Query("interview_number")
	var interview models.Interview

	if interviewNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Interview number is missing"})
		return
	}

	if err := config.DB.First(&interview, "number = ?", interviewNumber).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": "Interview not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		}
		return
	}

	type StudentInterviewResult struct {
		Name                  string
		Score                 int
		Comment               string
		GradingStandardNumber string
	}

	var assessmentMetrics []assessmentModel.AssessmentMetric

	var interviewAssessmentMetrics []models.InterviewAssessmentMetric
	if err := config.DB.
		Joins("join assessment_metrics on assessment_metrics.number = interview_assessment_metrics.assessment_metric_number AND assessment_metrics.campus_class_number=?", student.CampusClassNumber).
		Where("interview_number = ?", interviewNumber).Find(&interviewAssessmentMetrics).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch interview assessment metrics"})
		return
	}

	for _, v := range interviewAssessmentMetrics {
		var assessmentMetric assessmentModel.AssessmentMetric
		if err := config.DB.Where("number = ?", v.AssessmentMetricNumber).First(&assessmentMetric).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Assessment metric not found"})
			return
		}
		assessmentMetrics = append(assessmentMetrics, assessmentMetric)
	}

	var results []StudentInterviewResults

	for _, v := range assessmentMetrics {
		var aggregateScales []aggregateSaceltModel.AggregateScale
		var assessment assessmentModel.Assessment
		if err := config.DB.Where("student_number = ? AND assessment_metric_number = ?", studentNumber, v.Number).First(&assessment).Error; err != nil {
			//c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch assessment"})
			//return
		}

		if err := config.DB.Where("grading_standard_number = ?", v.GradingStandardNumber).Find(&aggregateScales).Error; err != nil {
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
			results = append(results, StudentInterviewResults{
				Name:    v.AssessmentAreaName,
				Comment: comment,
				Score:   int(assessment.ActualScore),
				Grade:   "Not Graded , no scales ",
			})

		} else {
			for _, val := range aggregateScales {
				if assessment.ActualScore >= val.MinValue && assessment.ActualScore <= val.MaxValue {
					results = append(results, StudentInterviewResults{
						Name:    v.AssessmentAreaName,
						Score:   int(assessment.ActualScore),
						Grade:   val.Name,
						Comment: comment,
					})
				}
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": results})
}

//var results []StudentInterviewResults
//var result []StudentInterviewResult
//config.DB.Raw("SELECT a.actual_score AS Score , a.comment , aa.name , gs.number  FROM assessments a,assessment_metrics am,assessment_areas aa , interview_assessment_metrics iam , interviews i ,grading_standards gs  WHERE a.student_number = ? AND a.assessment_metric_number = am.number AND am.number = iam.assessment_metric_number AND  am.assessment_area_number = aa.number AND am.grading_standard_number = gs.number AND  iam.interview_number = i.number AND i.number = ?", studentNumber, interviewNumber).Scan(&result)
//
//// loop through the result and get the grade
//
//for _, v := range result {
//
//	var aggregateScales []aggregateSaceltModel.AggregateScale
//
//	if err := config.DB.Where("grading_standard_number = ?", v.GradingStandardNumber).Find(&aggregateScales).Error; err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch aggregate scales"})
//		return
//	}
//
//	for _, val := range aggregateScales {
//		if v.Score >= val.MinValue && v.Score <= val.MaxValue {
//			results = append(results, StudentInterviewResults{
//				Name:    v.Name,
//				Score:   v.Score,
//				Grade:   val.Name,
//				Comment: v.Comment,
//			})
//		}
//	}
//
//}

func RemoveInterviewAssessmentMetric(c *gin.Context) {

	interviewAssessmentMetricNumber := c.Query("interview_assessment_metric_number")
	if interviewAssessmentMetricNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Interview assessment metric number is missing"})
		return
	}

	var interviewAssessmentMetric models.InterviewAssessmentMetric
	if err := config.DB.Where("number = ?", interviewAssessmentMetricNumber).First(&interviewAssessmentMetric).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Interview assessment metric not found"})
		return
	}

	if err := config.DB.Delete(&interviewAssessmentMetric).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to delete interview assessment metric"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Interview assessment metric deleted successfully"})
}
