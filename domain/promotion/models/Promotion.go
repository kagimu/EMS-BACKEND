package models

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/academicYear/models"
	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	studentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/models"
	"gorm.io/gorm"
)

type Promotion struct {
	gorm.Model
	Number                    string               `json:"number" gorm:"size:36;unique;not null"`
	CurrentAcademicYearNumber string               `json:"current_academic_year_number" gorm:"size:36;not null"`
	CurrentAcademicYear       models.AcademicYear  `gorm:"foreignKey:CurrentAcademicYearNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	Status                    string               `json:"status" gorm:"size:30;"`
	PromotionSettingNumber    string               `json:"promotion_setting_number" gorm:"size:36;not null"`
	PromotionSetting          PromotionSetting     `gorm:"foreignKey:PromotionSettingNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	StudentNumber             string               `json:"student_number" gorm:"size:36;not null"`
	Student                   studentModel.Student `gorm:"foreignKey:StudentNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	Reason                    string               `json:"reasons" gorm:"type:text"`
	UserNumber                string               `json:"user_number" gorm:"size:36;"`
	User                      authModel.User       `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
