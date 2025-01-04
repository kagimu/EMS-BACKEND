package controllers

import (
	services2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/services"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/promotion/services"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/studyPeriod/repositories"
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/promotion/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/promotion/models"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
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

	nextCampusClass, err := services2.ResolveCampusClass(payload.NextClassNumber, payload.CampusNumber, payload.CurriculumNumber, payload.CampusCurriculumNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	prevCampusClass, err := services2.ResolveCampusClass(payload.PreviousClassNumber, payload.CampusNumber, payload.CurriculumNumber, payload.CampusCurriculumNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	studyPeriod, err := repositories.GetStudyPeriodByNumber(payload.StudyPeriodNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if studyPeriod.Number == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid study period number"})
		return
	}

	savePromotionSettings := models.PromotionSetting{
		Number:                    uuid.NewV4().String(),
		Name:                      payload.Name,
		Type:                      payload.Type,
		Status:                    payload.Status,
		NextCampusClassNumber:     nextCampusClass.Number,
		PreviousCampusClassNumber: prevCampusClass.Number,
		UserNumber:                user.Number,
		StudyPeriodNumber:         studyPeriod.Number,
	}

	if err := config.DB.Create(&savePromotionSettings).Error; err != nil {
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

	promotionSettings, err := promotionService.AllPromotionSettings(studyPeriodNumber, nextCampusClassNumber, settingsType)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Error fetching promotion settings: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": promotionSettings})
}
