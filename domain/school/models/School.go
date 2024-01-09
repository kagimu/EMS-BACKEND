package models

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"gorm.io/gorm"
	"time"
)

type (
	School struct {
		gorm.Model
		Number           string             `json:"number" gorm:"size:36;unique;not null"`
		Name             string             `json:"name" form:"name" binding:"required" gorm:"size:200;unique;not null"`
		CommencementDate time.Time          `json:"commencement_date" form:"name" binding:"required" gorm:"not null"`
		SchoolCurriculum []SchoolCurriculum `gorm:"foreignKey:SchoolNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
		UserNumber       string             `json:"user_number" gorm:"size:36;index;not null"`
		User             models.User        `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	}
)
