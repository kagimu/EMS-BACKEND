package models

import (
	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	models2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/models"
	"time"
)

type Promotion struct {
	CreatedAt              time.Time
	UpdatedAt              time.Time
	Number                 string                `json:"number" gorm:"size:36;primarykey;not null"`
	ReportSummaryNumber    string                `json:"report_summary_number" gorm:"size:36;not null"`
	ReportSummary          models2.ReportSummary `gorm:"foreignKey:CurrentAcademicYearNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	Status                 string                `json:"status" gorm:"size:30;"`
	PromotionType          string                `json:"promotion_type" gorm:"size:30;"`
	PromotionSettingNumber string                `json:"promotion_setting_number" gorm:"size:36;not null"`
	PromotionSetting       PromotionSetting      `gorm:"foreignKey:PromotionSettingNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	UserNumber             string                `json:"user_number" gorm:"size:36;"`
	User                   authModel.User        `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
