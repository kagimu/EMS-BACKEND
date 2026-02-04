package controllers

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/admission/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/admission/services"
	"github.com/gin-gonic/gin"
)

func NewAdmissionRequirementController(ctx *gin.Context) {
	var req requests.AdmissionRequirementsRequest
	err := ctx.BindJSON(&req)
	if err != nil {
		ctx.JSON(400, gin.H{"message": "Invalid request data " + err.Error()})
		return
	}

	if err := req.Validate(); err != nil {
		ctx.JSON(400, gin.H{"message": "Validation failed: " + err.Error()})
		return
	}

	if err := services.SaveAdmissionRequirementsService(req.ToModel()); err != nil {
		ctx.JSON(500, gin.H{"message": "Failed to save admission Requirement: " + err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"message": "Admission Requirement saved successfully"})

}
