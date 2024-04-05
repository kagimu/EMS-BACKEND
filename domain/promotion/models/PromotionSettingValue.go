package models

import (
	assesmentGradeModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/assessmentGrade/models"
	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"gorm.io/gorm"
)

type PromotionSettingsValue struct {
	gorm.Model
	Number                 string                               `json:"number" gorm:"size:36;unique;not null"`
	PromotionSettingNumber string                               `json:"promotion_settings_number" gorm:"size:36;not null"`
	PromotionSetting       PromotionSetting                     `gorm:"foreignKey:PromotionSettingNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	AssessmentGradeNumber  string                               `json:"assessment_grade_number" gorm:"size:36;"`
	AssessmentGrade        assesmentGradeModels.AssessmentGrade `gorm:"foreignKey:AssessmentGradeNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	Type                   string                               `json:"type" gorm:"size:36;"`
	Score                  float64                              `json:"score"`
	UserNumber             string                               `json:"user_number" gorm:"size:36;"`
	User                   authModel.User                       `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
