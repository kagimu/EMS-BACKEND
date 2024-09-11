package repository

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/attachments/models"
	"gorm.io/gorm"
)

type AttachmentRepository interface {
	CreateAttachment(attachment *models.Attachment) error
	GetAttachmentByNumber(number string) (*models.Attachment, error)
	DeleteAttachmentByNumber(number string) error
}

type AttachmentRepositoryImpl struct {
	DB *gorm.DB
}

func (a *AttachmentRepositoryImpl) DeleteAttachmentByNumber(number string) error {
	return a.DB.Where("number = ?", number).Delete(&models.Attachment{}).Error
}

func NewAttachmentRepositoryImpl(db *gorm.DB) AttachmentRepository {
	return &AttachmentRepositoryImpl{DB: db}
}

func (a *AttachmentRepositoryImpl) CreateAttachment(attachment *models.Attachment) error {
	return a.DB.Create(attachment).Error
}

func (a *AttachmentRepositoryImpl) GetAttachmentByNumber(number string) (*models.Attachment, error) {
	var attachment models.Attachment
	err := a.DB.Where("number = ?", number).First(&attachment).Error
	if err != nil {
		return nil, err
	}
	return &attachment, nil
}
