package models

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	models2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
	models3 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/strand/models"
	"time"
)

type (
	Topic struct {
		ID           uint           `gorm:"primaryKey"`
		Number       string         `json:"number" gorm:"size:100;unique;not null"`
		Name         string         `json:"name" form:"name" binding:"required" gorm:"size:60;unique;not null"`
		UserNumber   string         `json:"user_number" gorm:"size:40;not null"`
		StrandNumber string         `json:"strand_number" gorm:"size:40;not null"`
		ClassNumber  string         `json:"class_number" gorm:"size:40;not null"`
		CreatedAt    time.Time      `json:"created_at"`
		UpdatedAt    time.Time      `json:"updated_at"`
		DeletedAt    *time.Time     `sql:"index"`
		Strand       models3.Strand `gorm:"foreignKey:StrandNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
		Class        models2.Class  `gorm:"foreignKey:ClassNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
		User         models.User    `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	}
)

func (u *Topic) TableName() string {
	return "topics"
}
