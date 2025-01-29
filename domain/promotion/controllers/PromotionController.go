package controllers

import (
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	academicYearModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/academicYear/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/promotion/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/promotion/models"
	studentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/models"

	"github.com/gin-gonic/gin"
)

func PromoteStudent(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
		return
	}

	var payload requests.PromotionRequest
	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input"})
		return
	}

	if err := payload.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	studentNumber := payload.StudentNumber
	var student studentModel.Student
	if err := config.DB.Where("number = ?", studentNumber).First(&student).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Student not found"})
		return
	}

	academicYearNumber := payload.CurrentAcademicYearNumber
	var academicYear academicYearModel.AcademicYear
	if err := config.DB.Where("number = ?", academicYearNumber).First(&academicYear).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Academic year not found"})
		return
	}

	promotionSettingNumber := payload.PromotionSettingNumber
	var promotionSetting models.PromotionSetting
	if err := config.DB.Where("number = ?", promotionSettingNumber).First(&promotionSetting).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Promotion setting not found"})
		return
	}

	/*savePromotion := models.Promotion{
		Number:                    uuid.NewV4().String(),
		PromotionSettingNumber:    payload.PromotionSettingNumber,
		StudentNumber:             studentNumber,
		Status:                    payload.Status,
		Reason:                    payload.Reason,
		CurrentAcademicYearNumber: payload.CurrentAcademicYearNumber,
		UserNumber:                user.Number,
	}

	if err := config.DB.Create(&savePromotion).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to create promotion"})
		return
	}*/

	c.JSON(http.StatusOK, gin.H{"message": "Promotion created successfully"})
}

func FetchAllPromotions(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var promotions []models.Promotion
	if err := config.DB.Preload("CurrentAcademicYear").Preload("Student").Preload("PromotionSetting").Preload("PromotionSetting.NextClass").Preload("PromotionSetting.PreviousClass").Preload("Student.PreferredCampus").Find(&promotions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch promotions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": promotions})
}

func FetchPromotionsByAcademicYear(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	academicYearNumber := c.Query("academic_year_number")
	var promotions []models.Promotion
	if err := config.DB.Where("current_academic_year_number = ?", academicYearNumber).Preload("CurrentAcademicYear").Preload("Student").Preload("PromotionSetting").Preload("PromotionSetting.NextClass").Preload("PromotionSetting.PreviousClass").Preload("Student.PreferredCampus").Find(&promotions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch promotions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": promotions})
}

func EditPromotion(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Unauthorized"})
		return
	}

	var payload requests.PromotionRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input"})
		return
	}

	payload.Prepare()

	var promotion models.Promotion
	promotionNumber := c.Param("promotion_number")
	if err := config.DB.First(&promotion, "number = ?", promotionNumber).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "No promotion found with the specified number"})
		return
	}

	/*updatedPromotion := models.Promotion{
		PromotionSettingNumber:    payload.PromotionSettingNumber,
		StudentNumber:             payload.StudentNumber,
		Status:                    payload.Status,
		Reason:                    payload.Reason,
		CurrentAcademicYearNumber: payload.CurrentAcademicYearNumber,
		UserNumber:                user.Number,
	}

	if err := config.DB.Model(&promotion).Where("number = ?", promotionNumber).Updates(updatedPromotion).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update promotion"})
		return
	}*/

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": promotion})
}
