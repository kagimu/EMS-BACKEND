package models

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	Number         string      `json:"number" gorm:"size:36;unique;not null"`
	DatePerformed  string      `json:"date_performed"`
	Title          string      `json:"title"`
	TimeFrame      string      `json:"time_frame"`
	InputsUsed     string      `json:"inputs_used" gorm:"type:text"`
	Outputs        string      `json:"outputs" gorm:"type:text"`
	PeopleInvolved string      `json:"people_involved" gorm:"type:text"`
	StaffNumber    string      `json:"staff_number" gorm:"size:36;not null"`
	Description    string      `json:"description" gorm:"type:text"`
	Staff          Staff       `gorm:"foreignKey:StaffNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	UserNumber     string      `json:"user_number" gorm:"size:36;index;not null"`
	User           models.User `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
