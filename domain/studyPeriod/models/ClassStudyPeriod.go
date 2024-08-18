package models

import classModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"

type ClassStudyPeriod struct {
	StudyPeriodNumber string            `json:"study_period_number" gorm:"size:36;not null;primaryKey"`
	StudyPeriod       StudyPeriod       `gorm:"foreignKey:StudyPeriodNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
	ClassNumber       string            `json:"class_number" gorm:"size:36;not null;primaryKey"`
	Class             classModels.Class `gorm:"foreignKey:ClassNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
}
