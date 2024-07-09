package models

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/grades/models"
)

type AggregateScale struct {
	Code         int64        `json:"code" gorm:"size:36;not null"`
	Number       string       `json:"number" gorm:"size:36;unique;not null"`
	MinValue     float64      `json:"min_value" gorm:"not null"`
	MaxValue     float64      `json:"max_value" gorm:"not null"`
	Description  string       `json:"description" gorm:"size:1000;not null"`
	ScaleNumber  string       `json:"scale_number" gorm:"size:36;not null"`
	Scale        Scale        `gorm:"foreignKey:ScaleNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
	Weight       float64      `json:"weight"`
	Illustration string       `json:"illustration" form:"illustration" gorm:"type:text;"`
	GradeNumber  string       `json:"grade_number" gorm:"size:36;not null"`
	Grade        models.Grade `gorm:"foreignKey:GradeNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
}
