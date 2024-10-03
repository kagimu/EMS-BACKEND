package models

import (
	campusModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/models"
	itemModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/item/models"
	"gorm.io/gorm"
)

type CampusItem struct {
	gorm.Model
	Number       string             `json:"number" gorm:"size:36;not null"`
	CampusNumber string             `json:"campus_number" gorm:"size:36;"`
	ItemNumber   string             `json:"item_number" gorm:"size:36;"`
	Item         itemModels.Item    `gorm:"foreignKey:ItemNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	Campus       campusModel.Campus `gorm:"foreignKey:CampusNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
