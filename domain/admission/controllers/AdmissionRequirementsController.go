package controllers

import (
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/admission/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/admission/services"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
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
// Query params: school_number, campus_class_number, class_number, campus_number, campus_curriculum_number, campus_level_number,
// study_period_number, academic_year_number, item_number, mode_of_study_number
func GetAdmissionRequirementsController(ctx *gin.Context) {
	if _, err := commons.GetUserFromHeader(ctx); err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"message": "User not found"})
		return
	}

	schoolNumber := ctx.Query("school_number")
	campusClassNumber := ctx.Query("campus_class_number")
	classNumber := ctx.Query("class_number")
	campusNumber := ctx.Query("campus_number")
	campusCurriculumNumber := ctx.Query("campus_curriculum_number")
	campusLevelNumber := ctx.Query("campus_level_number")
	studyPeriodNumber := ctx.Query("study_period_number")
	academicYearNumber := ctx.Query("academic_year_number")
	itemNumber := ctx.Query("item_number")
	modeOfStudyNumber := ctx.Query("mode_of_study_number")

	list, err := services.FetchAdmissionRequirementsService(
		schoolNumber,
		campusClassNumber,
		classNumber,
		campusNumber,
		campusCurriculumNumber,
		campusLevelNumber,
		studyPeriodNumber,
		academicYearNumber,
		itemNumber,
		modeOfStudyNumber,
	)
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

func DeleteAdmissionRequirementController(ctx *gin.Context) {
	if _, err := commons.GetUserFromHeader(ctx); err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"message": "User not found"})
		return
	}

	number := ctx.Param("number")
	if err := services.DeleteAdmissionRequirementService(number); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Failed to delete admission requirement: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Admission requirement deleted successfully"})
}

type UpdateAdmissionRequirementStatusRequest struct {
	Status string `json:"status"`
}

func UpdateAdmissionRequirementStatusController(ctx *gin.Context) {
	if _, err := commons.GetUserFromHeader(ctx); err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"message": "User not found"})
		return
	}

	number := ctx.Param("number")
	var req UpdateAdmissionRequirementStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request data: " + err.Error()})
		return
	}
	if req.Status == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "status is required"})
		return
	}
	if req.Status != constants.StatusActive && req.Status != constants.StatusInactive {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "status must be ACTIVE or INACTIVE"})
		return
	}

	if err := services.UpdateAdmissionRequirementStatusService(number, req.Status); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Failed to update admission requirement status: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Admission requirement status updated successfully"})
}
