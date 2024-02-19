package models

import (
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	models2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/school/models"
	"gorm.io/gorm"
)

type (
	Campus struct {
		gorm.Model
		Number           string         `json:"number" gorm:"size:36;unique;not null"`
		Name             string         `json:"name" form:"name" binding:"required" gorm:"size:200;unique;not null"`
		Phone            string         `json:"phone" form:"phone" binding:"required" gorm:"size:17;not null"`
		Email            string         `json:"email" form:"email" binding:"required" gorm:"size:50;not null"`
		Location         string         `json:"location" form:"location" binding:"required" gorm:"size:200;not null"`
		UserNumber       string         `json:"user_number" gorm:"size:36;not null"`
		SchoolNumber     string         `json:"school_number" gorm:"size:36;not null"`
		CommencementDate time.Time      `json:"commencement_date" form:"commencement_date" binding:"required" gorm:"not null"`
		User             models.User    `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
		School           models2.School `gorm:"foreignKey:SchoolNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	}
)

func (u *Campus) TableName() string {
	return "campuses"
}
