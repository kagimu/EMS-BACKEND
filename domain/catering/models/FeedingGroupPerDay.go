package models

import (
	authModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"gorm.io/gorm"
)

type FeedingGroupPerDay struct {
	gorm.Model
	Number             string          `json:"number" gorm:"size:36;unique;not null"`
	FeedingGroupNumber string          `json:"feeding_group_number" gorm:"size:36;"`
	Day                string          `json:"day" gorm:"not null"`
	FeedingGroup       FeedingGroup    `gorm:"foreignKey:FeedingGroupNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	UserNumber         string          `json:"user_number" gorm:"size:36;index;not null"`
	User               authModels.User `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}

// Check if the breakfast is meant to be served on the given day.
//  Respond with an error
// Add a filter
//  get the list of users
// append the menu items - append to each user
