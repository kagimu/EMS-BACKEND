package models

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/aggregateScale/models"
	classModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
	"gorm.io/gorm"
)

type GradingStandardClassScale struct {
	gorm.Model
	GradingStandardNumber string                 `json:"grading_standard_number" gorm:"size:36;not null"`
	GradingStandard       models.GradingStandard `gorm:"foreignKey:GradingStandardNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
	ClassNumber           string                 `json:"class_number" gorm:"size:36;not null"`
	Class                 classModels.Class      `gorm:"foreignKey:ClassNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
	ScaleNumber           string                 `json:"scale_number" gorm:"size:36;not null"`
	Scale                 models.Scale           `gorm:"foreignKey:ScaleNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
}
