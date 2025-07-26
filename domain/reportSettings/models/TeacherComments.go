package models

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	campusClassModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
	assessmentAreaModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/assessmentArea/models"
	models2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/studyPeriod/models"
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

type ClassTeacherComment struct {
	ID                uint
	CreatedAt         time.Time
	UpdatedAt         time.Time
	Date              time.Time
	Number            string                 `json:"number" gorm:"size:36;unique;not null;primaryKey"`
	Comment           string                 `json:"comment"`
	StudyPeriodNumber string                 `json:"study_period_number" gorm:"size:36"`
	StudyPeriod       models2.StudyPeriod    `gorm:"foreignKey:StudyPeriodNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
	ClassNumber       string                 `json:"class_number" gorm:"size:36"`
	Class             campusClassModel.Class `gorm:"foreignKey:ClassNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
	UserNumber        string                 `json:"user_number" gorm:"size:36"`
	User              models.User            `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
}

type SubjectTeacherComment struct {
	ID                   uint
	CreatedAt            time.Time
	UpdatedAt            time.Time
	Number               string                             `json:"number" gorm:"size:36;unique;not null;primaryKey"`
	Comment              string                             `json:"comment"`
	AssessmentAreaNumber string                             `json:"assessment_area_number" gorm:"size:36"`
	AssessmentArea       assessmentAreaModel.AssessmentArea `gorm:"foreignKey:AssessmentAreaNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
	StudyPeriodNumber    string                             `json:"study_period_number" gorm:"size:36"`
	StudyPeriod          models2.StudyPeriod                `gorm:"foreignKey:StudyPeriodNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
	UserNumber           string                             `json:"user_number" gorm:"size:36"`
	User                 models.User                        `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
}
