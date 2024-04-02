package controllers

import (
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	assessmentGradeModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/assessmentGrade/models"
	assessmentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/promotion/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/promotion/models"
	studentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/models"
	uuid "github.com/satori/go.uuid"

	"github.com/gin-gonic/gin"
)

func PromoteStudent(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "Unauthorizr=ed"})
		return
	}

	var payload requests.PromotionRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input"})
		return
	}

	if err := payload.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	assessmentGradeNumber := payload.AssessmentGradeNumber
	var assessmentGrade assessmentGradeModel.AssessmentGrade
	if err := config.DB.Where("number = ?", assessmentGradeNumber).First(&assessmentGrade).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Assessment grade not found"})
		return
	}

	assessmentNumber := payload.AssessmentNumber
	var assessment assessmentModel.Assessment
	if err := config.DB.Where("number = ?", assessmentNumber).First(&assessment).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Assessment not found"})
		return
	}

	studentNumber := payload.StudentNumber
	var student studentModel.Student
	if err := config.DB.Where("number = ?", studentNumber).First(&student).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Student not found"})
		return
	}

	savePromotion := models.Promotion{
		Number:                    uuid.NewV4().String(),
		StudentNumber:             studentNumber,
		AssessmentNumber:          assessmentNumber,
		AssessmentGradeNumber:     assessmentGradeNumber,
		Status:                    payload.Status,
		Type:                      payload.Type,
		PreviousCampusClassNumber: payload.PreviousCampusClassNumber,
		NextCampusClassNumber:     payload.NextCampusClassNumber,
		Condition:                 payload.Condition,
		ReasonForAltering:         payload.ReasonForAltering,
		Authorizer:                payload.Authorizer,
		UserNumber:                user.Number,
	}

	if err := config.DB.Create(&savePromotion).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to create promotion"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Promotion created successfully"})
}

func FetchAllPromotedStudents(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var promotions []models.Promotion
	if err := config.DB.Find(&promotions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch promotions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": promotions})
}
