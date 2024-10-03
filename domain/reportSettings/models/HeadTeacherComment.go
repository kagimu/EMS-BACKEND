package models

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"time"
)

type HeadTeacherComment struct {
	ID                  uint `gorm:"primarykey"`
	CreatedAt           time.Time
	UpdatedAt           time.Time
	Comment             string        `json:"comment"`
	MinAggregate        float64       `json:"min_aggregate"`
	MaxAggregate        float64       `json:"max_aggregate"`
	UserNumber          string        `json:"user_number" gorm:"size:36"`
	User                models.User   `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
	ReportSettingNumber string        `json:"report_setting_number" gorm:"size:36"`
	ReportSetting       ReportSetting `gorm:"foreignKey:ReportSettingNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
}
