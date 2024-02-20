package models

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/item/models"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	ItemNumber string      `json:"item_number" gorm:"size:36;unique"`
	Item       models.Item `gorm:"foreignKey:ItemNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}

// Add after find hook to set item
