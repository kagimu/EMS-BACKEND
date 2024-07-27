package models

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/aggregateScale/models"
	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"time"
)

type AssessmentGrade struct {
	Number              string `gorm:"primarykey;size:36;not null"`
	UpdatedAt           time.Time
	Status              string         `json:"status"`
	AssessmentAreaCount int            `json:"assessment_area_count"`
	GradePoints         int            `json:"grade_points"`
	UserNumber          string         `json:"user_number" gorm:"size:36;index;not null"`
	User                authModel.User `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	ScaleNumber         string         `json:"scale_number" gorm:"size:36;index;not null"`
	Scale               models.Scale   `gorm:"foreignKey:ScaleNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
}

type AssessmentGradeScale struct {
	Code                  int64           `json:"code"`
	Number                string          `json:"number" gorm:"size:36;primarykey;not null"`
	Name                  string          `json:"name"`
	Description           string          `json:"description" gorm:"type:text"`
	MinimumPoints         float64         `json:"minimum_points"`
	MaximumPoints         float64         `json:"maximum_points"`
	Abbreviation          string          `json:"abbreviation" gorm:"size:20;"`
	Remark                string          `json:"remark" gorm:"type:text"`
	Illustration          string          `json:"illustration" form:"illustration" gorm:"type:text"`
	AssessmentGradeNumber string          `json:"assessment_grade_number" gorm:"size:36;not null"`
	AssessmentGrade       AssessmentGrade `gorm:"foreignKey:AssessmentGradeNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
}
