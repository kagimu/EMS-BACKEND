package models

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"gorm.io/gorm"
)

type (
	ImportSchemeOfWork struct {
		gorm.Model
		Number   string `json:"number" gorm:"size:36;unique;not null"`
		Name     string `json:"name" form:"name" binding:"required" gorm:"size:200;not null"`
		FileName string `json:"file_name" form:"file_name" binding:"required" gorm:"size:200;not null"`
		//CurriculumNumber string      `json:"curriculum_number" gorm:"size:36;not null"`
		UserNumber string `json:"user_number" gorm:"size:36;not null"`
		//Curriculum       Curriculum  `gorm:"foreignKey:CurriculumNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
		User models.User `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	}
)
