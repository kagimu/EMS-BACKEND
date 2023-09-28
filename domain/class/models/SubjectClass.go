package models

import (
	models2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/models"
	models3 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/subject/models"
	"time"
)

type (
	SubjectCLass struct {
		ID            uint            `gorm:"primaryKey"`
		Number        string          `json:"number" gorm:"size:100;unique;not null"`
		ClassNumber   string          `json:"class_number" gorm:"size:40;not null"`
		SubjectNumber string          `json:"subject_number" gorm:"size:40;not null"`
		CampusNumber  string          `json:"campus_number" gorm:"size:40;not null"`
		LevelNumber   string          `json:"level_number" gorm:"size:40;not null"`
		UserNumber    string          `json:"user_number" gorm:"size:40;not null"`
		CreatedAt     time.Time       `json:"created_at"`
		UpdatedAt     time.Time       `json:"updated_at"`
		DeletedAt     *time.Time      `sql:"index"`
		User          models2.User    `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
		Subject       models3.Subject `gorm:"foreignKey:SubjectNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
		Class         Class           `gorm:"foreignKey:ClassNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
		Campus        models.Campus   `gorm:"foreignKey:CampusNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
		Level         models.Level    `gorm:"foreignKey:LevelNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	}
)

func (u *SubjectCLass) TableName() string {
	return "subject_class"
}
