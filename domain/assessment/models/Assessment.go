package models

import (
	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	assessmentAreaModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/assessmentArea/models"
	levelModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/settings/level/models"
	studentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/models"
	"gorm.io/gorm"
)

type Assessment struct {
	gorm.Model
	Number                 string                             `json:"number" gorm:"size:36;unique;not null"`
	AssessmentDate         string                             `json:"assessment_date"`
	AssessmentMetricNumber string                             `json:"assessment_metric_number" gorm:"size:36;"`
	AssessmentMetric       AssessmentMetric                   `gorm:"foreignKey:AssessmentMetricNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	StudentNumber          string                             `json:"student_number" gorm:"size:36;index;not null"`
	Student                studentModel.Student               `gorm:"foreignKey:StudentNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	AssessmentAreaNumber   string                             `json:"assessment_area_number"`
	AssessmentArea         assessmentAreaModel.AssessmentArea `gorm:"foreignKey:AssessmentAreaNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	LevelNumber            string                             `json:"level_number" gorm:"size:36;not null"`
	Level                  levelModels.Level                  `gorm:"foreignKey:LevelNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	UserNumber             string                             `json:"user_number" gorm:"size:36;index;not null"`
	User                   authModel.User                     `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
