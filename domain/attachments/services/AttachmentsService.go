package services

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/repositories"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/attachments/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/attachments/repository"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/attachments/responses"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func UploadAttachment(c *gin.Context, formKey string) (string, error) {
	attachmentFile, err := c.FormFile(formKey)
	if err != nil {
		return "", errors.New("File error: " + err.Error())
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

	ownerNumber := c.PostForm("owner_number")
	title := c.PostForm("title")
	category := c.PostForm("category")

	repo := repository.NewAttachmentRepositoryImpl(config.DB)
	attachmentNumber := uuid.NewV4().String()
	err = repo.CreateAttachment(&models.Attachment{
		AttachmentName: fileName,
		Number:         attachmentNumber,
		AttachmentType: extension,
		OwnerNumber:    ownerNumber,
		Title:          title,
		Category:       category,
	})

	if err != nil {
		return "", err
	}

	return attachmentNumber, nil

}

func GetAttachmentByNumber(number string) (responses.AttachmentResponse, error) {
	var resp responses.AttachmentResponse
	repo := repository.NewAttachmentRepositoryImpl(config.DB)
	attachment, err := repo.GetAttachmentByNumber(number)
	if err != nil {
		return resp, err
	}

	resp.Number = attachment.Number
	resp.Title = attachment.Title
	resp.Category = attachment.Category
	resp.CreatedAt = attachment.CreatedAt.Format(time.RFC3339)
	resp.FilePath = "/fileRepo/attachments/" + attachment.AttachmentName
	resp.Extension = attachment.AttachmentType
	resp.FileName = attachment.AttachmentName
	resp.URL = "/attachments/download/" + attachment.Number

	return resp, nil
}

func GetAttachmentsByOwner(ownerNumber string) ([]responses.AttachmentResponse, error) {
	repo := repository.NewAttachmentRepositoryImpl(config.DB)
	attachments, err := repo.GetAttachmentsByOwner(ownerNumber)
	if err != nil {
		return nil, err
	}

	var respList []responses.AttachmentResponse
	for _, att := range attachments {
		respList = append(respList, responses.AttachmentResponse{
			Number:    att.Number,
			Title:     att.Title,
			Category:  att.Category,
			CreatedAt: att.CreatedAt.Format(time.RFC3339),
			FilePath:  "/attachments/download/" + att.Number,
			FileName:  att.AttachmentName,
			Extension: att.AttachmentType,
			URL:       "/attachments/download/" + att.Number,
		})
	}
	return respList, nil
}

func DeleteAttachment(name string, category string) error {
	tx := config.DB.Begin()
	repo := repository.NewAttachmentRepositoryImpl(tx)

	attachment, err := repo.GetAttachmentByNumber(name)
	if err != nil {
		tx.Rollback()
		return errors.New("attachment not found")
	}

	aggregateScale, err := repositories.GetAggregateScaleByIllustration(name, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if category == "aggregate_scale" {
		aggregateScale.Illustration = ""
		err = repositories.SaveAggregateScale(aggregateScale, tx)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	// unlink file from file system
	filePath := "/fileRepo/attachments/" + attachment.AttachmentName
	//delete file
	err = os.Remove(filePath)
	if err != nil && !os.IsNotExist(err) {
		tx.Rollback()
		return err
	}

	err = repo.DeleteAttachmentByNumber(name)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil

}
