package controllers

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/promotion/services"
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/promotion/api/requests"
	"github.com/gin-gonic/gin"
)

var promotionValueService = services.NewPromotionService()

func CreatePromotionSettingValue(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var payload requests.PromotionSettingValueRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	payload.Prepare()
	if err := payload.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err = promotionValueService.SavePromotionSettingValue(payload, user)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Promotion setting value created successfully"})

}

func FetchPromotionSettingsValues(c *gin.Context) {
	if _, err := commons.GetUserFromHeader(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User not found"})
		return
	}

	promotionSettingNumber := c.Query("promotion_setting_number")
	if promotionSettingNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "promotion setting number is required"})
		return
	}

	data, err := promotionValueService.PromotionSettingValue(promotionSettingNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"data": data})
}
