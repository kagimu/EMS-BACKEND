package models

import "gorm.io/gorm"

type HeadTeacherComment struct {
	gorm.Model
	Comment             string        `json:"comment"`
	MinAggregate        float64       `json:"min_aggregate"`
	MaxAggregate        float64       `json:"max_aggregate"`
	UserNumber          string        `json:"user_number" gorm:"size:36"`
	ReportSettingNumber string        `json:"report_setting_number" gorm:"size:36"`
	ReportSetting       ReportSetting `gorm:"foreignKey:ReportSettingNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
