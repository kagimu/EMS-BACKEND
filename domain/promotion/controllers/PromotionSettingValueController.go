package controllers

import (
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/promotion/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/promotion/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func CreatePromotionSettingValue(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var payload requests.PromotionSettingValue
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	payload.Prepare()
	if err := payload.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	savePromotionSettings := models.PromotionSettingsValue{
		Number:     uuid.NewV4().String(),
		Type:       payload.Type,
		Score:      payload.Score,
		UserNumber: user.Number,
	}

	if err := config.DB.Create(&savePromotionSettings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create promotion setting value"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Promotion setting value created successfully"})

}

func FetchPromotionSettingsValues(c *gin.Context) {
	if _, err := commons.GetUserFromHeader(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User not found"})
		return
	}

	var promotionSettingsValues []models.PromotionSettingsValue
	err := config.DB.Find(&promotionSettingsValues).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch promotion settings values"})
	}

	c.JSON(http.StatusOK, gin.H{"data": promotionSettingsValues})
}
