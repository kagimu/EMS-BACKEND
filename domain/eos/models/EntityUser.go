package models

import (
	userModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"time"
)

type EntityUser struct {
	ID            uint `gorm:"primarykey"`
	CreatedAt     time.Time
	EntityNumber  string          `json:"entity_number" gorm:"size:36"`
	Entity        Entity          `gorm:"foreignKey:EntityNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
	UserNumber    string          `json:"user_number" gorm:"size:36"`
	User          userModels.User `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
	AddedByNumber string          `json:"added_by_number" gorm:"size:36"`
	AddedBy       userModels.User `gorm:"foreignKey:AddedByNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
}
