package controllers

import (
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/admission/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/admission/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func CreateAdmissionOfferCondition(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var payload requests.AdmissionOfferConditionRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to validate input"})
		return
	}

	saveAdmission := models.AdmissionOfferCondition{
		Number:      uuid.NewV4().String(),
		Description: payload.Description,
		UserNumber:  user.Number,
	}

	if err := config.DB.Create(&saveAdmission).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create admission offer condition"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Admission offer conditions created successfully"})
}

func FetchAdmissionOfferConditions(c *gin.Context) {
	if _, err := commons.GetUserFromHeader(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User not found"})
		return
	}

	var admissions []models.AdmissionOfferCondition
	query := config.DB.Model(&models.AdmissionOfferCondition{})

	if err := query.Find(&admissions).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Admission offer conditions fetching failed"})
		return
	}

	c.JSON(http.StatusOK, admissions)
}
