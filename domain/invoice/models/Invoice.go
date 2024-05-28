package models

import (
	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	campusModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/models"
	schoolModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/school/models"
	"gorm.io/gorm"
	"time"
)

type Invoice struct {
	gorm.Model
	Number        string              `json:"number" gorm:"size:36;unique;not null"`
	Status        string              `json:"status" gorm:"default:NOT_PAID"`
	Description   string              `json:"description" gorm:"type:text"`
	Type          string              `json:"type" gorm:"default:CLIENT"`
	DueDate       time.Time           `json:"due_date" gorm:"type:date"`
	RefNumber     string              `json:"ref_number"  gorm:"size:16;unique;not null"`
	UserNumber    string              `json:"user_number" gorm:"size:36;not null"`
	User          authModel.User      `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	AddedByNumber string              `json:"added_by_number" gorm:"size:36;not null"`
	AddedBy       authModel.User      `gorm:"foreignKey:AddedByNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	SchoolNumber  string              `json:"school_number" gorm:"size:36;not null"`
	School        schoolModel.School  `gorm:"foreignKey:SchoolNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	CampusNumber  string              `json:"campus_number" gorm:"size:36;not null"`
	Campus        campusModels.Campus `gorm:"foreignKey:CampusNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
