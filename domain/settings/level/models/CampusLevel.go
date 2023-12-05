package models

import (
	campusModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/models"
	"gorm.io/gorm"
)

type (
	CampusLevel struct {
		gorm.Model
		Number                string              `json:"number" gorm:"size:36;unique;not null"`
		Code                  int64               `json:"code" form:"code" binding:"required" gorm:"size:60;not null"`
		CampusNumber          string              `json:"campus_number" gorm:"size:36;not null"`
		CurriculumLevelNumber string              `json:"curriculum_level_number" gorm:"size:36;not null"`
		Campus                campusModels.Campus `gorm:"foreignKey:CampusNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
		CurriculumLevel       CurriculumLevel     `gorm:"foreignKey:CurriculumLevelNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	}
)
