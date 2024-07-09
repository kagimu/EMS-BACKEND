package controllers

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/aggregateScale/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/aggregateScale/api/responses"
	models2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/aggregateScale/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/services"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/gin-gonic/gin"
	"net/http"
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

	err = services.SaveAggregateScale(payload, user.Number)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Aggregate scale created successfully"})
}

func FetchAggregateScales(c *gin.Context) {
	campusNumber := c.Query("campus_number")
	if campusNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Campus number is required"})
		return
	}

	campusClassNumber := c.Query("campus_class_number")

	if campusClassNumber == "" {
		c.JSON(http.StatusOK, gin.H{"data": make([]responses.AggregateScaleResponse, 0)})
		return
	}

	//baseURL := os.Getenv("BASE_URL")

	gradingStandardNumber := c.Query("grading_standard_number")

	assessmentAreaNumber := c.Query("assessment_area_number")

	response, err := services.FetchAggregateScales(campusClassNumber, campusNumber, gradingStandardNumber, assessmentAreaNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

func UpdateAggregateScale(c *gin.Context) {
	user, err := commons.GetUserFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var payload requests.EditAggregateScaleRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process input"})
		return
	}

	if err = payload.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	scaleNumber := payload.ScaleNumber
	var scale models2.Scale
	if err := config.DB.First(&scale, "number = ?", scaleNumber).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "No aggregate scale found with the specified number"})
		return
	}

	/*var fileName string
	_, err = c.FormFile("illustration")
	if err == nil {
		file, err := commons.HandleFileUpload(c, "illustration", "/fileRepo/illustrations")
		fileName = file
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to handle file upload"})
			return
		}
	}*/

	scale.Name = payload.ScaleName
	scale.UserNumber = user.Number

	if err := config.DB.Model(&scale).Where("number  = ?", scale.Number).Updates(scale).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update aggregate scale"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": scale})
}

func FetchAllScales(context *gin.Context) {
	schoolNumber := context.Query("school_number")
	if schoolNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Choose a school to fetch scales"})
		return
	}

	scales, err := services.FindAllScales(schoolNumber)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"data": scales})

}
