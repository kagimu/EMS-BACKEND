package models

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/item/models"
)

type Service struct {
	ItemNumber string      `json:"item_number" gorm:"size:36;unique;primarykey;not null"`
	Item       models.Item `gorm:"foreignKey:ItemNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
