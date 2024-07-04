package models

import (
	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	schoolModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/school/models"
	"time"
)

type Scale struct {
	Number       string `gorm:"primarykey"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Name         string             `json:"name" gorm:"size:1000;not null"`
	SchoolNumber string             `json:"school_number" gorm:"size:36;not null"`
	School       schoolModel.School `gorm:"foreignKey:SchoolNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	Count        int                `json:"count"`
	UserNumber   string             `json:"user_number" gorm:"size:36;not null"`
	User         authModel.User     `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
