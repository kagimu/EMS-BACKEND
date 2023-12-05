package models

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type Status string

const (
	Pending  Status = "PENDING"
	Paid     Status = "PAID"
	Rejected Status = "REJECTED"
)

type Application struct {
	gorm.Model
	Status    Status    `gorm:"type:varchar(255);not null;;default:'PENDING'" json:"status"`
	ID        uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`
	StudentID uuid.UUID `gorm:"type:uuid;index" json:"student_id"`
	ParentID  uuid.UUID `gorm:"type:uuid;index" json:"parent_id"`
}
