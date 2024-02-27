package controllers

import (
	"errors"
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/aggregateScale/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/aggregateScale/models"
	campusModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

func CreateAlternateAggregateScale(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var payload requests.AlternativeAggregateScaleRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process payload"})
		return
	}

	if err := payload.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	saveAlternativeAggregateScale := models.AlternativeAggregateScale{
		Number:            uuid.NewV4().String(),
		Name:              payload.Name,
		Factors:           payload.Factors,
		Description:       payload.Description,
		MinValue:          payload.MinValue,
		MaxValue:          payload.MaxValue,
		LevelNumber:       payload.LevelNumber,
		StudyPeriodNumber: payload.StudyPeriodNumber,
		ClassNumber:       payload.ClassNumber,
		CampusNumber:      payload.CampusNumber,
		UserNumber:        user.Number,
	}

	tx := config.DB.Begin()

	if err := tx.Create(&saveAlternativeAggregateScale).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"message": " Failed to save alternative aggregatet scale"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Alternative aggregate scale created successfully", "data": saveAlternativeAggregateScale})
}

func FetchAlternativeAggregateScale(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	campusNumber := c.Query("campus_number")
	if campusNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Campus number is required"})
		return
	}

	var campus campusModel.Campus

	if err := config.DB.Where("number = ?", campusNumber).First(&campus).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": "Campus not found"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch campus"})
		}
		return
	}

	var alternativeAggregateScales []models.AlternativeAggregateScale
	if err := config.DB.Where("campus_number = ?", user.Number).Find(&alternativeAggregateScales).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch alternative aggregate scales"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": alternativeAggregateScales})
}
