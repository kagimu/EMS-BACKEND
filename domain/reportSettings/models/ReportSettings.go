package models

import (
	userModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"gorm.io/gorm"
)
import academicYearModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/academicYear/models"
import studyModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/studyPeriod/models"
import studentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/models"
import classModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"

type ReportSetting struct {
	gorm.Model
	Number             string                         `json:"number" gorm:"size:36;unique;not null"`
	AcademicYearNumber string                         `json:"academic_year_number" gorm:"size:36;"`
	StudyPeriodNumber  string                         `json:"study_period_number" gorm:"size:36;"`
	ClassNumber        string                         `json:"class_number" gorm:"size:36;"`
	Columns            int64                          `json:"columns" gorm:"type:bigint;"`
	Status             string                         `json:"status" gorm:"type:text;"`
	AcademicYear       academicYearModel.AcademicYear `gorm:"foreignKey:AcademicYearNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	StudyPeriod        studyModels.StudyPeriod        `gorm:"foreignKey:StudyPeriodNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	Class              classModel.Class               `gorm:"foreignKey:ClassNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	UserNumber         string                         `json:"user_number" gorm:"size:36;"`
	User               userModels.User                `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}

type ReportColumn struct {
	gorm.Model
	Number              string          `json:"number" gorm:"size:36;unique;not null"`
	ReportSettingNumber string          `json:"report_setting_number" gorm:"size:36;"`
	Name                string          `json:"name" gorm:"type:text;"`
	PercentageShare     float64         `json:"percentage_share" gorm:"type:double precision;"`
	ReportSetting       ReportSetting   `gorm:"foreignKey:ReportSettingNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	UserNumber          string          `json:"user_number" gorm:"size:36;"`
	User                userModels.User `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}

type ReportColumnAssessment struct {
	gorm.Model
	Number                 string          `json:"number" gorm:"size:36;unique;not null"`
	ReportColumnNumber     string          `json:"report_column_number" gorm:"size:36;"`
	AssessmentAreaNumber   string          `json:"assessment_number" gorm:"size:36;"`
	AssessmentMetricNumber string          `json:"assessment_metric_number" gorm:"size:36;"`
	ReportColumn           ReportColumn    `gorm:"foreignKey:ReportColumnNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	UserNumber             string          `json:"user_number" gorm:"size:36;"`
	User                   userModels.User `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}

type ReportSummary struct {
	gorm.Model
	Number              string                  `json:"number" gorm:"size:36;unique;not null"`
	ReportSettingNumber string                  `json:"report_setting_number" gorm:"size:36;"`
	StudentNumber       string                  `json:"student_number" gorm:"size:36;"`
	ClassNumber         string                  `json:"class_number" gorm:"size:36;"`
	StudyPeriodNumber   string                  `json:"study_period_number" gorm:"size:36;"`
	Grade               string                  `json:"grade" gorm:"type:text;"`
	Aggregate           float64                 `json:"aggregate" gorm:"type:double precision;"`
	Comment             string                  `json:"comment" gorm:"type:text;"`
	ReportSetting       ReportSetting           `gorm:"foreignKey:ReportSettingNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	Student             studentModel.Student    `gorm:"foreignKey:StudentNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	StudyPeriod         studyModels.StudyPeriod `gorm:"foreignKey:StudyPeriodNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	Class               classModel.Class        `gorm:"foreignKey:ClassNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	UserNumber          string                  `json:"user_number" gorm:"size:36;"`
	User                userModels.User         `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}

type ReportScore struct {
	gorm.Model
	ReportSummaryNumber  string          `json:"report_summary_number" gorm:"size:36;"`
	AssessmentAreaNumber string          `json:"assessment_area_number" gorm:"size:36;"`
	ReportColumnNumber   string          `json:"report_column_number" gorm:"size:36;"`
	Score                float64         `json:"score" gorm:"type:double precision;"`
	Aggregate            float64         `json:"aggregate" gorm:"type:double precision;"`
	Grade                string          `json:"grade" gorm:"type:text;"`
	UserNumber           string          `json:"user_number" gorm:"size:36;"`
	User                 userModels.User `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}

type ReportScoreSummary struct {
	gorm.Model
	Number               string          `json:"number" gorm:"size:36;unique;not null"`
	ReportSummaryNumber  string          `json:"report_summary_number" gorm:"size:36;"`
	AssessmentAreaNumber string          `json:"assessment_area_number" gorm:"size:36;"`
	Aggregate            float64         `json:"aggregate" gorm:"type:double precision;"`
	Grade                string          `json:"grade" gorm:"type:text;"`
	Summation            float64         `json:"summation" gorm:"type:double precision;"`
	Comment              string          `json:"comment" gorm:"type:text;"`
	UserNumber           string          `json:"user_number" gorm:"size:36;"`
	User                 userModels.User `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
