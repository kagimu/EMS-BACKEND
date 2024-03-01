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

func CreateAggregateScale(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var payload requests.AggregateScaleRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process payload"})
		return
	}

	if err := payload.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	saveAggregateScale := models.AggregateScale{
		Number:            uuid.NewV4().String(),
		Name:              payload.Name,
		Scale:             payload.Scale,
		MinValue:          payload.MinValue,
		MaxValue:          payload.MaxValue,
		Description:       payload.Description,
		LevelNumber:       payload.LevelNumber,
		SchoolNumber:      payload.SchoolNumber,
		StudyPeriodNumber: payload.StudyPeriodNumber,
		ClassNumber:       payload.ClassNumber,
		CampusNumber:      payload.CampusNumber,
		UserNumber:        user.Number,
	}

	tx := config.DB.Begin()

	if err := tx.Create(&saveAggregateScale).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"message": " Failed to save aggregatet scale"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Aggregate scale created successfully", "data": saveAggregateScale})
}

func FetchAggregateScales(c *gin.Context) {
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

	var aggregateScales []models.AggregateScale
	if err := config.DB.Where("campus_number = ?", user.Number).Find(&aggregateScales).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch aggregate scales"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": aggregateScales})
}
