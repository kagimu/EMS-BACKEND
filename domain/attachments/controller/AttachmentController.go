package controller

import (
	"fmt"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/attachments/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func DownloadAttachment(context *gin.Context) {
	fileName := context.Param("attachment_number")
	fileDetails, err := services.GetAttachmentByNumber(fileName)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "File not found " + err.Error()})
		return
	}
	filePath := "/fileRepo/attachments/" + fileDetails.FileName
	// get bytes from file
	byteFile, err := os.ReadFile(filePath)
	if err != nil {
		//fmt.Println("Internal server error")
		fmt.Println(err)
		return
	}

	context.Header("Content-Disposition", "attachment; filename="+fileDetails.FileName)
	contentType := "application/octet-stream"

	switch fileDetails.Extension {
	case ".pdf":
		contentType = "application/pdf"
	case ".doc", ".docx":
		contentType = "application/msword"
	}

	context.Header("Content-Type", contentType)

	context.Data(http.StatusOK, contentType, byteFile)
}

func DeleteAttachment(context *gin.Context) {
	type DeletionRequest struct {
		AttachmentNumber string `json:"attachment_number"`
		Category         string `json:"category"`
	}

	var deletionRequest DeletionRequest
	if err := context.ShouldBindJSON(&deletionRequest); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}

	err := services.DeleteAttachment(deletionRequest.AttachmentNumber, deletionRequest.Category)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "File not found " + err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "File deleted successfully"})
}
