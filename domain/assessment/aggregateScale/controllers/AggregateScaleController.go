package controllers

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/filters"
	"net/http"
	"os"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/aggregateScale/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/aggregateScale/api/responses"
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
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process payload" + err.Error()})
		return
	}

	if err := payload.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var fileName string

	_, err = c.FormFile("illustration")
	if err == nil {
		file, err := commons.HandleFileUpload(c, "illustration", "/fileRepo/illustrations")
		fileName = file
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to handle file upload"})
			return
		}
	}

	saveAggregateScale := models.AggregateScale{
		Code:                  payload.Code,
		Number:                uuid.NewV4().String(),
		Name:                  payload.Name,
		Weight:                payload.Weight,
		MinValue:              payload.MinValue,
		MaxValue:              payload.MaxValue,
		Illustration:          fileName,
		Description:           payload.Description,
		CampusClassNumber:     payload.CampusClassNumber,
		CampusNumber:          payload.CampusNumber,
		UserNumber:            user.Number,
		GradingStandardNumber: payload.GradingStandardNumber,
	}

	var duplicate models.AggregateScale
	if err := config.DB.Where("name = ?", saveAggregateScale.Name).First(&duplicate).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Duplicate aggregate scale name"})
		return
	}

	if err := config.DB.Where("code = ?", saveAggregateScale.Code).First(&duplicate).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Duplicate aggregate scale code"})
		return
	}

	tx := config.DB.Begin()

	if err := tx.Create(&saveAggregateScale).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"message": " Failed to save aggregate scale"})
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

	campusClassNumber := c.Query("campus_class_number")

	var aggregateScales []models.AggregateScale
	if err := config.DB.
		Scopes(filters.ByCampusClassNumber(campusClassNumber), filters.ByCampusNumber(campusNumber)).
		Preload("Campus").Preload("CampusClass").Preload("GradingStandard").
		Find(&aggregateScales).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch aggregate scales"})
		return
	}
	var response []responses.AggregateScaleResponse
	baseURL := os.Getenv("BASE_URL")
	for _, aggregateScale := range aggregateScales {
		response = append(response, responses.AggregateScaleResponse{
			Code:                aggregateScale.Code,
			Number:              aggregateScale.Number,
			MinValue:            aggregateScale.MinValue,
			MaxValue:            aggregateScale.MaxValue,
			Name:                aggregateScale.Name,
			Weight:              aggregateScale.Weight,
			Illustration:        baseURL + aggregateScale.Illustration,
			Description:         aggregateScale.Description,
			ClassName:           aggregateScale.CampusClass.ClassName,
			CampusName:          aggregateScale.Campus.Name,
			GradingStandardName: aggregateScale.GradingStandard.Name,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
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

	var fileName string
	_, err = c.FormFile("illustration")
	if err == nil {
		file, err := commons.HandleFileUpload(c, "illustration", "/fileRepo/illustrations")
		fileName = file
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to handle file upload"})
			return
		}
	}

	updateAggregateScale := models.AggregateScale{
		Code:                  payload.Code,
		Name:                  payload.Name,
		Weight:                payload.Weight,
		MinValue:              payload.MinValue,
		MaxValue:              payload.MaxValue,
		Description:           payload.Description,
		Illustration:          fileName,
		CampusClassNumber:     payload.CampusClassNumber,
		CampusNumber:          payload.CampusNumber,
		GradingStandardNumber: payload.GradingStandardNumber,
	}

	if err := config.DB.Model(&aggregateScale).Where("number  = ?", aggregateScaleNumber).Updates(updateAggregateScale).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update aggregate scale"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": aggregateScale})
}
