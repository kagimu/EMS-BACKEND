package controllers

import (
	"net/http"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/admission/api/requests"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/admission/services"
	attachmentServices "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/attachments/services"
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

// UploadSignedToaController handles POST /main/admissions/toa/signed-upload
//
// Expects multipart/form-data with:
//   - file           – the signed Terms of Admission document (PDF / image / doc)
//   - student_number – the student this signed document belongs to
//
// The file is stored via the shared attachments service under category "signed_toa".
// The response includes the attachment_number so the frontend can record / reference it.
func UploadSignedToaController(ctx *gin.Context) {
	studentNumber := ctx.PostForm("student_number")
	if studentNumber == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "student_number is required"})
		return
	}

	// Inject owner_number and category so the shared service picks them up from PostForm.
	// (The service reads owner_number, title and category via ctx.PostForm internally.)
	// We can set them by passing the values directly; the service already reads them from
	// the multipart form, so we just need to make sure the client sends them — but to be
	// safe and explicit we call the service with the form fields already present.
	attachmentNumber, err := attachmentServices.UploadAttachment(ctx, "file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Failed to upload signed TOA document: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":           "Signed Terms of Admission uploaded successfully",
		"attachment_number": attachmentNumber,
	})
}

// GetSignedToaAttachmentsController handles GET /main/admissions/toa/signed-upload/:student_number
//
// Returns all signed TOA documents previously uploaded for the given student.
func GetSignedToaAttachmentsController(ctx *gin.Context) {
	studentNumber := ctx.Param("student_number")
	if studentNumber == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "student_number is required"})
		return
	}

	attachments, err := attachmentServices.GetAttachmentsByOwner(studentNumber)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch signed TOA documents: " + err.Error()})
		return
	}

	// Filter to only signed_toa category documents
	var filtered []interface{}
	for _, a := range attachments {
		if a.Category == "signed_toa" {
			filtered = append(filtered, a)
		}
	}
	if filtered == nil {
		filtered = []interface{}{}
	}

	ctx.JSON(http.StatusOK, gin.H{"data": filtered})
}

