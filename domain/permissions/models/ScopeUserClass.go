package models

import (
	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	classModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
	"gorm.io/gorm"
)

type ScopeUserClass struct {
	gorm.Model
	Number        string           `json:"number" gorm:"size:36;unique;not null"`
	UserNumber    string           `json:"user_number" gorm:"size:36;not null"`
	User          authModel.User   `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	ClassNumber   string           `json:"class_number" gorm:"size:36;not null"`
	Class         classModel.Class `gorm:"foreignKey:ClassNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	AddedByNumber string           `json:"added_by_number" gorm:"size:36;not null"`
	AddedBy       authModel.User   `gorm:"foreignKey:AddedByNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
