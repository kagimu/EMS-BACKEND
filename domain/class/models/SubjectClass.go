package models

import (
	authModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	subjectModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/subject/models"
	"gorm.io/gorm"
)

type (
	ClassLevelSubject struct {
		gorm.Model
		Number           string               `json:"number" gorm:"size:36;unique;not null"`
		ClassLevelNumber string               `json:"class_level_number" gorm:"size:36;not null"`
		SubjectNumber    string               `json:"subject_number" gorm:"size:36;not null"`
		UserNumber       string               `json:"user_number" gorm:"size:36;not null"`
		ClassLevel       ClassLevel           `gorm:"foreignKey:ClassLevelNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
		User             authModels.User      `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
		Subject          subjectModel.Subject `gorm:"foreignKey:SubjectNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	}

	SubjectInfo struct {
		SubjectNumber string `json:"subject_number"`
		Name          string `json:"name"`
		Class         string `json:"class"`
		Abbreviation  string `json:"abbreviation"`
	}
)
