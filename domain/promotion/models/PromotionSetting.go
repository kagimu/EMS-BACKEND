package models

import (
	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	classModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
	studyModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/studyPeriod/models"
	"gorm.io/gorm"
)

type PromotionSetting struct {
	gorm.Model
	Number                    string                  `json:"number" gorm:"size:36;unique;not null"`
	Name                      string                  `json:"name" gorm:"size:100;not null"`
	Status                    string                  `json:"status" gorm:"size:20;"`
	Type                      string                  `json:"type" gorm:"size:50"` //promotion, repeated or demotion
	NextCampusClassNumber     string                  `json:"next_campus_class_number" gorm:"size:36;not null"`
	NextCampusClass           classModels.CampusClass `json:"next_class" gorm:"foreignKey:NextCampusClassNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
	PreviousCampusClassNumber string                  `json:"previous_campus_class_number" gorm:"size:36;not null"`
	PreviousCampusClass       classModels.CampusClass `json:"previous_campus_class" gorm:"foreignKey:PreviousCampusClassNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
	UserNumber                string                  `json:"user_number" gorm:"size:36;"`
	User                      authModel.User          `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	StudyPeriodNumber         string                  `json:"study_period_number" gorm:"size:36;not null"`
	StudyPeriod               studyModels.StudyPeriod `gorm:"foreignKey:StudyPeriodNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
	CampusCurriculumNumber    string                  `json:"campus_curriculum_number" gorm:"size:36"`
}
