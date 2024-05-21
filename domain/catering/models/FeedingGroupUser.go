package models

import (
	authModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	campusModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/models"
	"gorm.io/gorm"
)

type FeedingGroupUser struct {
	gorm.Model
	Number             string             `json:"number" gorm:"size:36;unique;not null"`
	FeedingGroupNumber string             `json:"feeding_group_number" gorm:"size:36;not null"`
	Status             string             `json:"status"`
	FeedingGroup       FeedingGroup       `gorm:"foreignKey:FeedingGroupNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
	UserNumber         string             `json:"user_number" gorm:"size:36;index;not null"`
	User               authModels.User    `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
	CampusNumber       string             `json:"campus_number" gorm:"size:36;index;not null"`
	Campus             campusModel.Campus `gorm:"foreignKey:CampusNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
}
