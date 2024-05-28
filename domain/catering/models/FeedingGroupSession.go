package models

import (
	"gorm.io/gorm"
)

type FeedingGroupSession struct {
	gorm.Model
	Number               string         `json:"number" gorm:"size:36;unique;not null"`
	FeedingGroupNumber   string         `json:"feeding_group_number" gorm:"size:36;not null"`
	FeedingGroup         FeedingGroup   `gorm:"foreignKey:FeedingGroupNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	Status               string         `json:"status"`
	FeedingSessionNumber string         `json:"feeding_session_number" gorm:"size:36;not null"`
	FeedingSession       FeedingSession `gorm:"foreignKey:FeedingSessionNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
