package controllers

import (
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/studentAssessmentAreaReport/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/studentAssessmentAreaReport/services"
	"github.com/gin-gonic/gin"
)

// GetStudentAssessmentAreaReport returns students with their scores per assessment schedule, filtered by assessment_area_number.
// Query params: assessment_area_number (required), study_period_number (optional), campus_class_number (optional).
func GetStudentAssessmentAreaReport(c *gin.Context) {
	if _, err := commons.GetUserFromHeader(c); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var req requests.StudentAssessmentAreaReportRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "assessment_area_number is required", "error": err.Error()})
		return
	}

	report, err := services.GetReport(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": report})
}

func GetAllStudentAssessmentAreaReportsByClassAndStudyPeriod(c *gin.Context) {
	if _, err := commons.GetUserFromHeader(c); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	report, err := services.GetAllStudentScoresReportByClassNumberAndStudyPeriod(c.Query("campus_class_number"), c.Query("study_period_number"), c.Query("campus_number"), c.Query("curriculum_number"), c.Query("campus_curriculum_number"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": report})
}
