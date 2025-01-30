package controllers

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/aggregateScale/api/requests"
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

	err = services.UpdateAggregateScale(payload, user.Number)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, commons.Feedback("Aggregate scale updated successfully"))
}

func FetchAllDetailedScales(context *gin.Context) {
	schoolNumber := context.Query("school_number")
	if schoolNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Choose a school to fetch scales"})
		return
	}

	scales, err := services.FindAllDetailsScales(schoolNumber)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"data": scales})

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

func RemoveAggregateScale(context *gin.Context) {
	scaleNumber := context.Param("scale_number")
	if scaleNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Choose a scale to delete"})
		return
	}

	err := services.RemoveScale(scaleNumber)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, commons.Feedback("Scale removed successfully"))
}

func UploadAggregateScaleIllustration(c *gin.Context) {
	aggregateScaleNumber := c.PostForm("aggregate_scale_number")

	if err := services.UploadAggregateScaleIllustration(aggregateScaleNumber, c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to upload aggregate scale photo" + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Aggregate scale photo uploaded successfully"})

}
