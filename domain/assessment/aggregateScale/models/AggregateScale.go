package models

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/grades/models"
	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	campusModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/models"
	classModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
	"gorm.io/gorm"
)

type AggregateScale struct {
	gorm.Model
	Code                  int64                  `json:"code" gorm:"size:36;not null"`
	Number                string                 `json:"number" gorm:"size:36;unique;not null"`
	MinValue              float64                `json:"min_value" gorm:"not null"`
	MaxValue              float64                `json:"max_value" gorm:"not null"`
	Name                  string                 `json:"name" gorm:"-"`
	Weight                float64                `json:"weight" gorm:"not null"`
	Illustration          string                 `json:"illustration" form:"illustration" gorm:"type:text;"`
	Description           string                 `json:"description" gorm:"type:text"`
	CampusClassNumber     string                 `json:"class_number" gorm:"size:36;not null"`
	CampusClass           classModel.CampusClass `gorm:"foreignKey:CampusClassNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	CampusNumber          string                 `json:"campus_number" gorm:"size:36;not null"`
	Campus                campusModel.Campus     `gorm:"foreignKey:CampusNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	UserNumber            string                 `json:"user_number" gorm:"size:36;not null"`
	User                  authModel.User         `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	GradingStandardNumber string                 `json:"grading_standard_number" gorm:"size:36;not null"`
	GradingStandard       GradingStandard        `gorm:"foreignKey:GradingStandardNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	GradeNumber           string                 `json:"grade_number" gorm:"size:36;not null"`
	Grade                 models.Grade           `gorm:"foreignKey:GradeNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}

// AfterFind hook to set the name of the aggregate scale

func (a *AggregateScale) AfterFind(tx *gorm.DB) (err error) {
	if a.GradeNumber != "" {
		var grade models.Grade
		if err := tx.Where("number = ?", a.GradeNumber).First(&grade).Error; err != nil {
			return err
		}
		a.Name = grade.Name
	}
	return
}
