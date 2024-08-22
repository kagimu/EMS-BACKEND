package controllers

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	requests2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/school/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/school/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

var modeOfStudyPeriod = services.NewModeOfStudyService()

func NewModeOfStudy(context *gin.Context) {
	user, err := commons.GetUserFromHeader(context)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var payload requests2.ModeOfStudyRequest
	if err := context.ShouldBindJSON(&payload); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process request input: " + err.Error()})
		return
	}

	if err := payload.Validate(); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err = modeOfStudyPeriod.NewModeOfStudy(payload, user.Number)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Added successfully"})

}

func AllModeOfStudy(context *gin.Context) {
	schoolNumber := context.Query("school_number")
	if schoolNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Choose school"})
		return
	}

	all, err := modeOfStudyPeriod.AllModesOfStudy(schoolNumber)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Study modes", "data": all})
}

func UpdateModeOfStudy(context *gin.Context) {
	user, err := commons.GetUserFromHeader(context)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var payload requests2.UpdateModeOfStudyRequest
	if err := context.ShouldBindJSON(&payload); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process request input: " + err.Error()})
		return
	}

	if err := payload.Validate(); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err = modeOfStudyPeriod.UpdateModeOfStudy(payload, user.Number)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Updated successfully"})
}

func DeleteModeOfStudy(context *gin.Context) {
	modeOfStudyNumber := context.Param("mode_of_study_number")
	if modeOfStudyNumber == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Choose mode of study to delete"})
		return
	}

	err := modeOfStudyPeriod.DeleteModeOfStudy(modeOfStudyNumber)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Removed successfully"})
}
