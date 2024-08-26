package controllers

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/services"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateAssessmentMetric(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var payload requests.AssessmentMetricRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input" + err.Error()})
		return
	}

	if err := payload.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err = services.SaveAssessmentMetric(payload, user.Number)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, commons.Feedback("Schedule created successfully"))
}

func FetchAssessmentMetrics(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	assessmentAreaNumber := c.Query("assessment_area_number")
	studyPeriodNumber := c.Query("study_period_number")
	campusClassNumber := c.Query("campus_class_number")
	metricNumber := c.Query("metric_number")
	classNumber := c.Query("class_number")
	campusNumber := c.Query("campus_number")

	var paginated = c.DefaultQuery("paginated", "no")

	if campusClassNumber == "" && (campusNumber == "" || classNumber == "") {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Please select at-least a campus and class to view assessment metrics"})
		return
	}

	/*if campusNumber != "" && classNumber != "" {
		campusClassNumber = ""
	}*/

	if paginated == "yes" {

		allAssessments, err := services.AllAssessmentsPaginated(assessmentAreaNumber, studyPeriodNumber, campusClassNumber, metricNumber, classNumber, campusNumber, c)
		if err != nil {
			c.JSON(http.StatusBadRequest, commons.Feedback(err.Error()))
			return
		}

		c.JSON(http.StatusOK, allAssessments)

	} else {

		assessmentMetrics, err := services.AllAssessments(assessmentAreaNumber, studyPeriodNumber, campusClassNumber, metricNumber, classNumber, campusNumber)
		if err != nil {
			c.JSON(http.StatusBadRequest, commons.Feedback(err.Error()))
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": assessmentMetrics})

	}
}

func UpdateAssessmentMetric(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var payload requests.UpdateAssessmentMetricRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input" + err.Error()})
		return
	}

	if err := payload.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err = services.UpdateAssessmentMetric(payload, user.Number)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, commons.Feedback("Schedule updated successfully"))
}
