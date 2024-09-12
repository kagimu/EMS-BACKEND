package services

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/attachments/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/attachments/repository"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/attachments/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"path/filepath"
	"strconv"
	"time"
)

func UploadAttachment(c *gin.Context, formKey string) (string, error) {
	attachmentFile, err := c.FormFile(formKey)
	if err != nil {
		return "", err
	}

	if attachmentFile.Size > 20971520 {
		return "", errors.New("file size should be less than 20MB")
	}

	t := time.Now()
	name := commons.GeneratedFormattedUUID() + "-" + strconv.Itoa(t.Year()) + "-" + strconv.Itoa(int(t.Month())) + "-" + strconv.Itoa(t.Day()) + "-" + strconv.Itoa(t.Minute())
	extension := filepath.Ext(attachmentFile.Filename)

	fileName := name + extension
	filePath := "/fileRepo/attachments/" + fileName

	if err := c.SaveUploadedFile(attachmentFile, filePath); err != nil {
		return "", err
	}

	repo := repository.NewAttachmentRepositoryImpl(config.DB)
	err = repo.CreateAttachment(&models.Attachment{
		AttachmentName: fileName,
		Number:         uuid.NewV4().String(),
		AttachmentType: extension,
	})

	if err != nil {
		return "", err
	}

	return fileName, nil

}

func GetAttachmentByNumber(number string) (responses.AttachmentResponse, error) {
	var resp responses.AttachmentResponse
	repo := repository.NewAttachmentRepositoryImpl(config.DB)
	attachment, err := repo.GetAttachmentByNumber(number)
	if err != nil {
		return resp, err
	}

	resp.FilePath = "/fileRepo/attachments/" + attachment.AttachmentName
	resp.Extension = attachment.AttachmentType
	resp.FileName = attachment.AttachmentName

	return resp, nil
}
