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
	context.Data(http.StatusOK, "application/octet-stream", byteFile)
}
