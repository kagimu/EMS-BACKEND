package models

import (
	userModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	schoolModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/school/models"
	"time"
)

type StudyPeriodName struct {
	Number       string             `gorm:"primaryKey;column:number;size:36;not null"`
	CreatedAt    time.Time          `gorm:"column:created_at"`
	UpdatedAt    time.Time          `gorm:"column:created_at"`
	Status       string             `gorm:"column:status;size:16;not null"`
	Name         string             `gorm:"column:name;size:60;not null"`
	UserNumber   string             `gorm:"column:user_number;size:36;not null"`
	User         userModel.User     `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	SchoolNumber string             `gorm:"column:school_number;size:36;not null"`
	School       schoolModel.School `gorm:"foreignKey:SchoolNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
}
