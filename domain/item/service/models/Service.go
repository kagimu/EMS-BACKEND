package models

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/item/models"
	"gorm.io/gorm"
)

type Service struct {
	gorm.Model
	ItemNumber     string       `json:"item_number" gorm:"size:36;unique"`
	Item           models.Item  `gorm:"foreignKey:ItemNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	CampusServices []CampusItem `gorm:"many2many:service_campus_services;" json:"campus_services"`
}
