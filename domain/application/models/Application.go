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
	ID        uuid.UUID `gorm:"type:char(36);primary_key;" json:"id"`
	Status    Status    `gorm:"type:varchar(255);not null;;default:'PENDING'" json:"status"`
	StudentID uuid.UUID `gorm:"type:uuid;index" json:"student_id"`
	ParentID  uuid.UUID `gorm:"type:uuid;index" json:"parent_id"`
}
