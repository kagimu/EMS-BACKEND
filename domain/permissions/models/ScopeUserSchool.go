package models

import (
	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	schoolModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/school/models"
	"gorm.io/gorm"
)

type ScopeUserSchool struct {
	gorm.Model
	Number        string             `json:"number" gorm:"size:36;unique;not null"`
	UserNumber    string             `json:"user_number" gorm:"size:36;not null"`
	User          authModel.User     `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	SchoolNumber  string             `json:"school_number" gorm:"size:36;not null"`
	School        schoolModel.School `gorm:"foreignKey:SchoolNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	AddedByNumber string             `json:"added_by_number" gorm:"size:36;not null"`
	AddedBy       authModel.User     `gorm:"foreignKey:AddedByNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
