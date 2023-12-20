package models

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"gorm.io/gorm"
)

type PreviousSchools struct {
	gorm.Model
	Number            string                `json:"number" gorm:"size:36;unique;not null"`
	Order             string                `json:"order" gorm:"not null"`
	Name              string                `json:"name" gorm:"size:50;not null"`
	StudentNumber     string                `json:"student_number" gorm:"size:36;unique;not null"`
	StartYear         string                `json:"start_period" gorm:"size:10;not null"`
	EndYear           string                `json:"end_period" gorm:"size:10;not null"`
	Classes           []PreviousSchoolClass `gorm:"foreignKey:ClassNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	UserNumber        string                `json:"user_number" gorm:"size:36;index;not null"`
	ApplicationNumber string                `json:"application_number" gorm:"size:36;index;not null"`
	User              models.User           `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
