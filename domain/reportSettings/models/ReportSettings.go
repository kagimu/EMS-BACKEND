package models

import (
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/models"
	userModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	assessmentAreaModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/assessmentArea/models"
	learningAreaModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/learningarea/models"
)
import academicYearModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/academicYear/models"
import studyModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/studyPeriod/models"
import studentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/models"
import classModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"

type ReportSetting struct {
	CreatedAt          time.Time
	UpdatedAt          time.Time
	Number             string                         `json:"number" gorm:"size:36;unique;not null;primaryKey"`
	Code               int64                          `json:"code" gorm:"type:bigint;"`
	AcademicYearNumber string                         `json:"academic_year_number" gorm:"size:36;not null;"`
	StudyPeriodNumber  string                         `json:"study_period_number" gorm:"size:36;not null;"`
	CampusClassNumber  string                         `json:"campus_class_number" gorm:"size:36;not null;"`
	Columns            int64                          `json:"columns" gorm:"type:bigint;not null;"`
	Status             string                         `json:"status" gorm:"type:text;"`
	UserNumber         string                         `json:"user_number" gorm:"size:36;not null;"`
	AcademicYear       academicYearModel.AcademicYear `gorm:"foreignKey:AcademicYearNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	StudyPeriod        studyModels.StudyPeriod        `gorm:"foreignKey:StudyPeriodNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	CampusClass        classModel.CampusClass         `gorm:"foreignKey:CampusClassNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	User               userModels.User                `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}

type ReportColumn struct {
	CreatedAt           time.Time
	UpdatedAt           time.Time
	Number              string          `json:"number" gorm:"size:36;unique;not null;primaryKey"`
	Code                int64           `json:"code" gorm:"type:bigint; not null;"`
	ReportSettingNumber string          `json:"report_setting_number" gorm:"size:36;not null;"`
	Name                string          `json:"name" gorm:"type:text;not null;"`
	PercentageShare     float64         `json:"percentage_share" gorm:"type:double precision;not null;"`
	ReportSetting       ReportSetting   `gorm:"foreignKey:ReportSettingNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
	UserNumber          string          `json:"user_number" gorm:"size:36;not null;"`
	User                userModels.User `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}

type ReportColumnAssessment struct {
	CreatedAt              time.Time
	Number                 string                              `json:"number" gorm:"size:36;unique;not null;primaryKey"`
	AssessmentAreaNumber   string                              `json:"assessment_area_number" gorm:"size:36;not null;"`
	AssessmentArea         assessmentAreaModels.AssessmentArea `gorm:"foreignKey:AssessmentAreaNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	AssessmentMetricNumber string                              `json:"assessment_metric_number" gorm:"size:36;not null;unique;"`
	AssessmentMetric       models.AssessmentMetric             `gorm:"foreignKey:AssessmentMetricNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
	ReportColumnNumber     string                              `json:"report_column_number" gorm:"size:36;not null;"`
	ReportColumn           ReportColumn                        `gorm:"foreignKey:ReportColumnNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
	UserNumber             string                              `json:"user_number" gorm:"size:36;not null;"`
	User                   userModels.User                     `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}

type ReportSummary struct {
	CreatedAt                  time.Time
	UpdatedAt                  time.Time
	Number                     string                  `json:"number" gorm:"size:36;unique;not null;primaryKey"`
	Code                       int64                   `json:"code" gorm:"type:bigint;not null;"`
	ReportSettingNumber        string                  `json:"report_setting_number" gorm:"size:36;not null;"`
	StudentNumber              string                  `json:"student_number" gorm:"size:36;not null;"`
	ReportTemplateName         string                  `json:"report_template_name" gorm:"size:250;"`
	CampusClassNumber          string                  `json:"campus_class_number" gorm:"size:36;not null;"`
	StudyPeriodNumber          string                  `json:"study_period_number" gorm:"size:36;not null;"`
	Grade                      string                  `json:"grade" gorm:"type:text;not null;"`
	Aggregate                  float64                 `json:"aggregate" gorm:"type:double precision;not null;"`
	TotalScore                 float64                 `json:"total_score" gorm:"type:double precision;not null;column:total_score"`
	Comment                    string                  `json:"comment" gorm:"type:text;not null;"`
	ReportSetting              ReportSetting           `gorm:"foreignKey:ReportSettingNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	Student                    studentModel.Student    `gorm:"foreignKey:StudentNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	StudyPeriod                studyModels.StudyPeriod `gorm:"foreignKey:StudyPeriodNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	CampusClass                classModel.CampusClass  `gorm:"foreignKey:CampusClassNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	UserNumber                 string                  `json:"user_number" gorm:"size:36;not null;"`
	User                       userModels.User         `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	StreamNumber               string                  `json:"stream_number" gorm:"size:36;"`
	Stream                     classModel.Stream       `gorm:"foreignKey:StreamNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	HeadTeacherComment         string                  `json:"head_teacher_comment" gorm:"size:100"`
	AggregateScaleUsed         string                  `json:"aggregate_scale_used" gorm:"size:100"`
	PromotionStatus            string                  `json:"promotion_status" gorm:"size:100;column:promotion_status;default:'NOT_SET'"`
	AssessmentGradeScaleNumber string                  `json:"assessment_grade_scale_number" gorm:"column:assessment_grade_scale_number;size:36;"`
}

type ReportScore struct {
	CreatedAt            time.Time
	UpdatedAt            time.Time
	Number               string                              `json:"number" gorm:"size:36;unique;not null;primaryKey"`
	Code                 int64                               `json:"code" gorm:"type:bigint;not null;"`
	ReportSummaryNumber  string                              `json:"report_summary_number" gorm:"size:36;not null;"`
	ReportSummary        ReportSummary                       `gorm:"foreignKey:ReportSummaryNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	AssessmentAreaNumber string                              `json:"assessment_area_number" gorm:"size:36;not null;"`
	AssessmentArea       assessmentAreaModels.AssessmentArea `gorm:"foreignKey:AssessmentAreaNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	ReportColumnNumber   string                              `json:"report_column_number" gorm:"size:36;not null;"`
	ReportColumn         ReportColumn                        `gorm:"foreignKey:ReportColumnNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	Score                float64                             `json:"score" gorm:"type:double precision;not null;"`
	AverageScore         float64                             `json:"average_score" gorm:"type:double precision;not null;"`
	Aggregate            float64                             `json:"aggregate" gorm:"type:double precision;not null;"`
	Grade                string                              `json:"grade" gorm:"type:text;not null;"`
	Illustration         string                              `json:"illustration" gorm:"type:text;not null;"`
	Remark               string                              `json:"remark" gorm:"type:text;"`
	UserNumber           string                              `json:"user_number" gorm:"size:36;not null;"`
	User                 userModels.User                     `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}

type ReportScoreSummary struct {
	CreatedAt            time.Time
	UpdatedAt            time.Time
	Number               string                              `json:"number" gorm:"size:36;unique;not null;primaryKey"`
	Code                 int64                               `json:"code" gorm:"type:bigint;not null;"`
	ReportSummaryNumber  string                              `json:"report_summary_number" gorm:"size:36;not null;"`
	ReportSummary        ReportSummary                       `gorm:"foreignKey:ReportSummaryNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	AssessmentAreaNumber string                              `json:"assessment_area_number" gorm:"size:36;not null;"`
	AssessmentArea       assessmentAreaModels.AssessmentArea `gorm:"foreignKey:AssessmentAreaNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	Aggregate            float64                             `json:"aggregate" gorm:"type:double precision;not null;"`
	Grade                string                              `json:"grade" gorm:"type:text;not null;"`
	Illustration         string                              `json:"illustration" gorm:"size36"`
	Summation            float64                             `json:"summation" gorm:"type:double precision;not null;"`
	Comment              string                              `json:"comment" gorm:"type:text;not null;"`
	UserNumber           string                              `json:"user_number" gorm:"size:36;not null;"`
	User                 userModels.User                     `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	TeacherInitial       string                              `json:"teacher_initial" gorm:"size:16;"`
}

type LearningAreaReportComment struct {
	CreatedAt           time.Time
	UpdatedAt           time.Time
	Number              string                          `json:"number" gorm:"size:36;not null;primaryKey"`
	ReportSummaryNumber string                          `json:"report_summary_number" gorm:"size:36;not null;"`
	ReportSummary       ReportSummary                   `gorm:"foreignKey:ReportSummaryNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
	LearningAreaNumber  string                          `json:"learning_area_number" gorm:"size:36;not null;"`
	LearningArea        learningAreaModels.LearningArea `gorm:"foreignKey:LearningAreaNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
	Comment             string                          `json:"comment" gorm:"type:text;not null;"`
	UserNumber          string                          `json:"user_number" gorm:"size:36;not null;"`
	User                userModels.User                 `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
}

type ReportSummaryComment struct {
	CreatedAt           time.Time
	UpdatedAt           time.Time
	ReportSummaryNumber string
	ReportSummary       ReportSummary `gorm:"foreignKey:ReportSummaryNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
}
