package models

import (
	assesmentGradeModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/assessmentGrade/models"
	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"gorm.io/gorm"
	"time"
)

type PromotionSettingValue struct {
	CreatedAt                  time.Time
	UpdatedAt                  time.Time
	Number                     string                                   `json:"number" gorm:"size:36;primarykey;not null"`
	PromotionSettingNumber     string                                   `json:"promotion_setting_number" gorm:"column:promotion_setting_number;size:36;not null"`
	PromotionSetting           PromotionSetting                         `gorm:"foreignKey:PromotionSettingNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
	AssessmentGradeScaleNumber string                                   `json:"assessment_grade_scale_number" gorm:"column:assessment_grade_scale_number;size:36;"`
	AssessmentGradeScale       assesmentGradeModel.AssessmentGradeScale `gorm:"foreignKey:AssessmentGradeScaleNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
	Score                      float64                                  `gorm:"column:total_score;comment:'Summation of total percentage score e.g 300'"`
	MinimumPoints              float64                                  `gorm:"column:minimum_points"`
	UserNumber                 string                                   `json:"user_number" gorm:"size:36;"`
	User                       authModel.User                           `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}

type PromotionAssessmentAreaValue struct {
	gorm.Model
	Number string `json:"number" gorm:"size:36;unique;not null"`
}
