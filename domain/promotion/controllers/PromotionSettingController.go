package controllers

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/promotion/services"
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/promotion/api/requests"
	"github.com/gin-gonic/gin"
)

var promotionService = services.NewPromotionService()

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

	err = promotionService.NewPromotionSetting(payload, user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to create promotion settings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Promotion settings created successfully"})
}

func FetchPromotionSettings(c *gin.Context) {
	if _, err := commons.GetUserFromHeader(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User not found"})
		return
	}

	studyPeriodNumber := c.Query("study_period_number")
	nextCampusClassNumber := c.Query("next_campus_class_number")
	settingsType := c.Query("type")
	if studyPeriodNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Choose a study period"})
		return
	}

	promotionSettings, err := promotionService.AllPromotionSettings(studyPeriodNumber, nextCampusClassNumber, settingsType)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Error fetching promotion settings: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": promotionSettings})
}

func FetchPromotionSettingDetails(c *gin.Context) {
	if _, err := commons.GetUserFromHeader(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User not found"})
		return
	}

	promotionSettingNumber := c.Query("promotion_setting_number")

	promotionSettings, err := promotionService.PromotionSettingDetails(promotionSettingNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Error fetching promotion setting details error: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": promotionSettings})
}

func FetchPromotionSettingGrades(c *gin.Context) {
	if _, err := commons.GetUserFromHeader(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User not found"})
		return
	}

	promotionSettingNumber := c.Query("promotion_setting_number")

	promotionSettings, err := promotionService.PromotionSettingGrades(promotionSettingNumber, c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Error fetching promotion setting grades error: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": promotionSettings})
}

func UpdatePromotionSettings(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	promotionSettingNumber := c.Param("promotion_setting_number")

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

	err = promotionService.UpdatePromotionSetting(payload, user, promotionSettingNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to update promotion settings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Promotion settings updated successfully"})
}
