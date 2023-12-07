package models

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	parentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/parent/models"
	studentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/models"
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
	Number     string               `json:"number" gorm:"size:36;unique;not null"`
	Status     Status               `gorm:"type:varchar(255);not null;;default:'PENDING'" json:"status"`
	StudentID  string               `json:"student_id" gorm:"size:36;index;not null"`
	Student    studentModel.Student `gorm:"foreignKey:StudentID;references:ID;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	ParentID   string               `json:"parent_id" gorm:"size:36;index;not null"`
	Parent     parentModel.Parent   `gorm:"foreignKey:ParentID;references:ID;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	UserNumber string               `json:"user_number" gorm:"size:36;index;not null"`
	User       models.User          `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
