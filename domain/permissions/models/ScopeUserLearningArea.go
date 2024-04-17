package models

import (
	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	learningAreaModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/learningarea/models"
	"gorm.io/gorm"
)

type ScopeUserLearningArea struct {
	gorm.Model
	Number             string                         `json:"number" gorm:"size:36;unique;not null"`
	UserNumber         string                         `json:"user_number" gorm:"size:36;not null"`
	User               authModel.User                 `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	LearningAreaNumber string                         `json:"learning_area_number" gorm:"size:36;not null"`
	LearningArea       learningAreaModel.LearningArea `gorm:"foreignKey:LearningAreaNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	AddedByNumber      string                         `json:"added_by_number" gorm:"size:36;not null"`
	AddedBy            authModel.User                 `gorm:"foreignKey:AddedByNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
