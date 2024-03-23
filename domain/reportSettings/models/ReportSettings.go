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
	Code               int64                          `json:"code" gorm:"type:bigint;"`
	AcademicYearNumber string                         `json:"academic_year_number" gorm:"size:36;not null;"`
	StudyPeriodNumber  string                         `json:"study_period_number" gorm:"size:36;not null;"`
	CampusClassNumber  string                         `json:"campus_class_number" gorm:"size:36;not null;"`
	Columns            int64                          `json:"columns" gorm:"type:bigint;not null;"`
	Status             string                         `json:"status" gorm:"type:text;"`
	AcademicYear       academicYearModel.AcademicYear `gorm:"foreignKey:AcademicYearNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	StudyPeriod        studyModels.StudyPeriod        `gorm:"foreignKey:StudyPeriodNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	Class              classModel.CampusClass         `gorm:"foreignKey:CampusClassNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	UserNumber         string                         `json:"user_number" gorm:"size:36;not null;"`
	User               userModels.User                `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}

type ReportColumn struct {
	gorm.Model
	Number              string          `json:"number" gorm:"size:36;unique;not null"`
	Code                int64           `json:"code" gorm:"type:bigint; not null;"`
	ReportSettingNumber string          `json:"report_setting_number" gorm:"size:36;not null;"`
	Name                string          `json:"name" gorm:"type:text;not null;"`
	PercentageShare     float64         `json:"percentage_share" gorm:"type:double precision;not null;"`
	ReportSetting       ReportSetting   `gorm:"foreignKey:ReportSettingNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	UserNumber          string          `json:"user_number" gorm:"size:36;not null;"`
	User                userModels.User `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}

type ReportColumnAssessment struct {
	gorm.Model
	Number                 string          `json:"number" gorm:"size:36;unique;not null"`
	Code                   int64           `json:"code" gorm:"type:bigint;not null;"`
	ReportColumnNumber     string          `json:"report_column_number" gorm:"size:36;not null;"`
	AssessmentAreaNumber   string          `json:"assessment_number" gorm:"size:36;not null;"`
	AssessmentMetricNumber string          `json:"assessment_metric_number" gorm:"size:36;not null;"`
	ReportColumn           ReportColumn    `gorm:"foreignKey:ReportColumnNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	UserNumber             string          `json:"user_number" gorm:"size:36;not null;"`
	User                   userModels.User `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}

type ReportSummary struct {
	gorm.Model
	Number              string                  `json:"number" gorm:"size:36;unique;not null"`
	Code                int64                   `json:"code" gorm:"type:bigint;not null;"`
	ReportSettingNumber string                  `json:"report_setting_number" gorm:"size:36;not null;"`
	StudentNumber       string                  `json:"student_number" gorm:"size:36;not null;"`
	CampusClassNumber   string                  `json:"campus_class_number" gorm:"size:36;not null;"`
	StudyPeriodNumber   string                  `json:"study_period_number" gorm:"size:36;not null;"`
	Grade               string                  `json:"grade" gorm:"type:text;not null;"`
	Aggregate           float64                 `json:"aggregate" gorm:"type:double precision;not null;"`
	Comment             string                  `json:"comment" gorm:"type:text;not null;"`
	ReportSetting       ReportSetting           `gorm:"foreignKey:ReportSettingNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	Student             studentModel.Student    `gorm:"foreignKey:StudentNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	StudyPeriod         studyModels.StudyPeriod `gorm:"foreignKey:StudyPeriodNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	Class               classModel.CampusClass  `gorm:"foreignKey:CampusClassNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	UserNumber          string                  `json:"user_number" gorm:"size:36;not null;"`
	User                userModels.User         `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}

type ReportScore struct {
	gorm.Model
	Number               string          `json:"number" gorm:"size:36;unique;not null"`
	Code                 int64           `json:"code" gorm:"type:bigint;not null;"`
	ReportSummaryNumber  string          `json:"report_summary_number" gorm:"size:36;not null;"`
	AssessmentAreaNumber string          `json:"assessment_area_number" gorm:"size:36;not null;"`
	ReportColumnNumber   string          `json:"report_column_number" gorm:"size:36;not null;"`
	Score                float64         `json:"score" gorm:"type:double precision;not null;"`
	Aggregate            float64         `json:"aggregate" gorm:"type:double precision;not null;"`
	Grade                string          `json:"grade" gorm:"type:text;not null;"`
	UserNumber           string          `json:"user_number" gorm:"size:36;not null;"`
	User                 userModels.User `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}

type ReportScoreSummary struct {
	gorm.Model
	Number               string          `json:"number" gorm:"size:36;unique;not null"`
	Code                 int64           `json:"code" gorm:"type:bigint;not null;"`
	ReportSummaryNumber  string          `json:"report_summary_number" gorm:"size:36;not null;"`
	AssessmentAreaNumber string          `json:"assessment_area_number" gorm:"size:36;not null;"`
	Aggregate            float64         `json:"aggregate" gorm:"type:double precision;not null;"`
	Grade                string          `json:"grade" gorm:"type:text;not null;"`
	Summation            float64         `json:"summation" gorm:"type:double precision;not null;"`
	Comment              string          `json:"comment" gorm:"type:text;not null;"`
	UserNumber           string          `json:"user_number" gorm:"size:36;not null;"`
	User                 userModels.User `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
