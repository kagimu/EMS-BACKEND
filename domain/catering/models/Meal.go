package models

import (
	"time"

	authModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	campusModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/models"
	"gorm.io/gorm"
)

type Meal struct {
	gorm.Model
	Number               string              `json:"number" gorm:"size:36;unique;not null"`
	UserNumber           string              `json:"user_number" gorm:"size:36;not null"`
	User                 authModels.User     `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	FeedingSessionNumber string              `json:"feeding_session_number" gorm:"size:36;not null"`
	FeedingSession       FeedingSession      `gorm:"foreignKey:FeedingSessionNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	StockingNumber       string              `json:"stocking_number" gorm:"size:36;not null"`
	Stocking             Stocking            `gorm:"foreignKey:StockingNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	MenuItemNumber       string              `json:"menu_item_number" gorm:"size:36;not null"`
	Quantity             float64             `json:"quantity"`
	MenuItem             MenuItem            `gorm:"foreignKey:MenuItemNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	StaffNumber          string              `json:"staff_number" gorm:"size:36;not null"`
	Staff                authModels.User     `gorm:"foreignKey:StaffNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	Date                 time.Time           `json:"date" gorm:"not null"`
	Time                 string              `json:"time" gorm:"not null"`
	CampusNumber         string              `json:"campus_number" gorm:"size:36;not null"`
	Campus               campusModels.Campus `gorm:"foreignKey:CampusNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	AddedByNumber        string              `json:"added_by_number" gorm:"size:36;not null"`
	AddedBy              authModels.User     `gorm:"foreignKey:AddedByNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
