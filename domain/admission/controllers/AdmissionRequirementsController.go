package controllers

import (
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/admission/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/admission/services"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/gin-gonic/gin"
)

func NewAdmissionRequirementController(ctx *gin.Context) {
	user, err := commons.GetUserFromHeader(ctx)
	if err != nil || user == nil {
		ctx.JSON(http.StatusForbidden, gin.H{"message": "User not found"})
		return
	}

	var req requests.AdmissionRequirementsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request data: " + err.Error()})
		return
	}

	if err := req.Validate(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Validation failed: " + err.Error()})
		return
	}

	if err := services.SaveAdmissionRequirementsService(&req, user.Number); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save admission requirements: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Admission requirements saved successfully",
		"count":   len(req.ItemNumbers),
	})
}

// GetAdmissionRequirementsController returns admission requirements with optional query filters.
// Query params: school_number, campus_class_number, study_period_number, academic_year_number, item_number
func GetAdmissionRequirementsController(ctx *gin.Context) {
	if _, err := commons.GetUserFromHeader(ctx); err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"message": "User not found"})
		return
	}

	schoolNumber := ctx.Query("school_number")
	campusClassNumber := ctx.Query("campus_class_number")
	studyPeriodNumber := ctx.Query("study_period_number")
	academicYearNumber := ctx.Query("academic_year_number")
	itemNumber := ctx.Query("item_number")

	list, err := services.FetchAdmissionRequirementsService(schoolNumber, campusClassNumber, studyPeriodNumber, academicYearNumber, itemNumber)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch admission requirements: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Admission requirements fetched successfully",
		"data":    list,
		"count":   len(list),
	})
}
