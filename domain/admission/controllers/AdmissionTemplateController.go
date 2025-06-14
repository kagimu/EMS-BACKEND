package controllers

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/admission/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/admission/services"
	"github.com/gin-gonic/gin"
)

func NewAdmissionTemplateController(ctx *gin.Context) {
	var req requests.AdmissionLetterTemplateRequest
	err := ctx.BindJSON(&req)
	if err != nil {
		ctx.JSON(400, gin.H{"message": "Invalid request data " + err.Error()})
		return
	}

	if err := req.Validate(); err != nil {
		ctx.JSON(400, gin.H{"message": "Validation failed: " + err.Error()})
		return
	}

	if err := services.SaveAdmissionLetterTemplateService(req.ToModel()); err != nil {
		ctx.JSON(500, gin.H{"message": "Failed to save admission template: " + err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"message": "Admission template saved successfully"})

}

func GetAdmissionTemplateController(ctx *gin.Context) {
	schoolNumber := ctx.Param("school_number")
	if schoolNumber == "" {
		ctx.JSON(400, gin.H{"message": "School number is required"})
		return
	}

	admissionTemplate, err := services.GetAdmissionLetterTemplateService(schoolNumber)
	if err != nil {
		ctx.JSON(500, gin.H{"message": "Failed to retrieve admission template: " + err.Error()})
		return
	}

	ctx.JSON(200, admissionTemplate)
}

func DeleteAdmissionTemplateController(ctx *gin.Context) {
	templateNumber := ctx.Param("template_number")
	if templateNumber == "" {
		ctx.JSON(400, gin.H{"message": "template number is required"})
		return
	}

	if err := services.DeleteAdmissionLetterTemplateService(templateNumber); err != nil {
		ctx.JSON(500, gin.H{"message": "Failed to delete admission template: " + err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"message": "Admission template deleted successfully"})
}
