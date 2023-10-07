package models

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	models2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/subject/models"
	"time"
)

type (
	Strand struct {
		ID            uint            `gorm:"primaryKey"`
		Number        string          `json:"number" gorm:"size:36;unique;not null"`
		Name          string          `json:"name" form:"name" binding:"required" gorm:"size:200;unique;not null"`
		Code          string          `json:"code" form:"code" binding:"required" gorm:"size:60;not null"`
		UserNumber    string          `json:"user_number" gorm:"size:36;not null"`
		SubjectNumber string          `json:"subject_number" gorm:"size:36;not null"`
		CreatedAt     time.Time       `json:"created_at"`
		UpdatedAt     time.Time       `json:"updated_at"`
		DeletedAt     *time.Time      `sql:"index"`
		Subject       models2.Subject `gorm:"foreignKey:SubjectNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
		User          models.User     `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	}
)

func (u *Strand) TableName() string {
	return "strands"
}
