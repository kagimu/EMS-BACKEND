package models

import (
	campusModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/models"
	campusClass "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
	"gorm.io/gorm"
)
import userModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"

type GradingStandard struct {
	gorm.Model
	Code              int64                   `json:"code" gorm:"size:36;unique;not null"`
	Number            string                  `json:"number" gorm:"size:36;unique;not null"`
	Name              string                  `json:"name" gorm:"size:255;not null"`
	Description       string                  `json:"description" gorm:"type:text;"`
	Type              string                  `json:"type" gorm:"size:255;not null"`
	CampusNumber      string                  `json:"campus_number" gorm:"size:36;not null"`
	Campus            campusModel.Campus      `gorm:"foreignKey:CampusNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	CampusClassNumber string                  `json:"campus_class_number" gorm:"size:36;not null"`
	CampusClass       campusClass.CampusClass `gorm:"foreignKey:CampusClassNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	UserNumber        string                  `json:"user_number" gorm:"size:36;not null"`
	User              userModel.User          `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
