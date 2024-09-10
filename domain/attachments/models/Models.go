package models

import (
	"time"
)

type Attachment struct {
	ID             uint `gorm:"primarykey"`
	CreatedAt      time.Time
	Number         string `json:"number" gorm:"unique"`
	AttachmentName string `json:"attachment_name" gorm:"not null"`
	AttachmentType string `json:"attachment_type" gorm:"not null"`
}
