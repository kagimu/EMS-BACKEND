package controllers

import (
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/aggregateScale/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/aggregateScale/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
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
		Number:                uuid.NewV4().String(),
		Name:                  payload.Name,
		Weight:                payload.Weight,
		MinValue:              payload.MinValue,
		MaxValue:              payload.MaxValue,
		Description:           payload.Description,
		SchoolNumber:          payload.SchoolNumber,
		ClassNumber:           payload.ClassNumber,
		CampusNumber:          payload.CampusNumber,
		UserNumber:            user.Number,
		GradingStandardNumber: payload.GradingStandardNumber,
	}

	tx := config.DB.Begin()

	if err := tx.Create(&saveAggregateScale).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"message": " Failed to save aggregate scale" + err.Error()})
		return
	}

	tx.Commit()

	c.JSON(http.StatusCreated, gin.H{"message": "Aggregate scale created successfully", "data": saveAggregateScale})
}

func FetchAggregateScales(c *gin.Context) {
	campusNumber := c.Query("campus_number")
	if campusNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Campus number is required"})
		return
	}

	var aggregateScales []models.AggregateScale
	if err := config.DB.Where("campus_number = ?", campusNumber).Find(&aggregateScales).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch aggregate scales"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": aggregateScales})
}

func FetchAggregateScaleBySchool(c *gin.Context) {

	campusNumber := c.Query("campus_number")
	if campusNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Campus number is required"})
		return
	}

	schoolNumber := c.Query("school_number")
	if schoolNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "School number is required"})
		return
	}

	var aggregateScales []models.AggregateScale
	if err := config.DB.Where("campus_number = ? AND school_number", campusNumber, schoolNumber).Find(&aggregateScales).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch aggregate scales"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": aggregateScales})
}

func EditAggregateScale(c *gin.Context) {
	_, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var payload requests.AggregateScaleRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input"})
		return
	}

	aggregateScaleNumber := c.Query("aggregate_scale_number")
	var aggregateScale models.AggregateScale
	if err := config.DB.First(&aggregateScale, "number = ?", aggregateScaleNumber).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "No aggregate scale found with the specified number"})
		return
	}

	updateAggregateScale := models.AggregateScale{
		Name:                  payload.Name,
		Weight:                payload.Weight,
		MinValue:              payload.MinValue,
		MaxValue:              payload.MaxValue,
		Description:           payload.Description,
		SchoolNumber:          payload.SchoolNumber,
		ClassNumber:           payload.ClassNumber,
		CampusNumber:          payload.CampusNumber,
		GradingStandardNumber: payload.GradingStandardNumber,
	}

	if err := config.DB.Model(&aggregateScale).Where("number  = ?", aggregateScaleNumber).Updates(updateAggregateScale).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update aggregate scale"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": aggregateScale})
}
