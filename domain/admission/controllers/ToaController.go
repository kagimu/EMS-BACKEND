package controllers

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/admission/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/admission/services"
	"github.com/gin-gonic/gin"
)

// SaveToaController handles POST /main/admissions/toa
func SaveToaController(ctx *gin.Context) {
	var req requests.ToaRequest
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"message": "Invalid request data: " + err.Error()})
		return
	}

	if err := req.Validate(); err != nil {
		ctx.JSON(400, gin.H{"message": "Validation failed: " + err.Error()})
		return
	}

	if err := services.SaveToaService(req.ToModel()); err != nil {
		ctx.JSON(500, gin.H{"message": "Failed to save TOA: " + err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"message": "Terms of Admission saved successfully"})
}

// GetToaController handles GET /main/admissions/toa/:school_number
func GetToaController(ctx *gin.Context) {
	schoolNumber := ctx.Param("school_number")
	if schoolNumber == "" {
		ctx.JSON(400, gin.H{"message": "School number is required"})
		return
	}

	toa, err := services.GetToaService(schoolNumber)
	if err != nil {
		// Return empty content gracefully when no record yet
		ctx.JSON(200, gin.H{"content": nil})
		return
	}

	ctx.JSON(200, toa)
}
