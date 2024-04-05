package controllers

import (
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	campusClassModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/promotion/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/promotion/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func CreatePromotionSettings(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var payload requests.PromotionSettingRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	payload.Prepare()
	if err := payload.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var campusClass campusClassModel.CampusClass
	if err := config.DB.Where("number = ?", payload.NextClassNumber).Find(&campusClass).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch next class"})
		return
	}

	if err := config.DB.Where("number = ?", payload.PreviousClassNumber).Find(&campusClass).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch previous class"})
		return
	}

	savePromotionSettings := models.PromotionSetting{
		Number:              uuid.NewV4().String(),
		Name:                payload.Name,
		Type:                payload.Type,
		Status:              payload.Status,
		NextClassNumber:     payload.NextClassNumber,
		PreviousClassNumber: payload.PreviousClassNumber,
		UserNumber:          user.Number,
	}

	if err := config.DB.Create(&savePromotionSettings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create promotion settings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Promotion settings created successfully"})
}

func FetchPromotionSettings(c *gin.Context) {
	if _, err := commons.GetUserFromHeader(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User not found"})
		return
	}

	var promotionSettings []models.PromotionSetting
	err := config.DB.Preload("NextClass").Preload("PreviousClass").Find(&promotionSettings).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch promotion settings"})
	}

	c.JSON(http.StatusOK, gin.H{"data": promotionSettings})
}
