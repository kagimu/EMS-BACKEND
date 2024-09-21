package models

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"gorm.io/gorm"
)

type Stream struct {
	gorm.Model
	Number            string      `json:"number" gorm:"size:36;not null;unique"`
	Name              string      `json:"name" gorm:"size:200;not null"`
	UserNumber        string      `json:"user_number" gorm:"size:36;not null"`
	User              models.User `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	CampusClassNumber string      `json:"campus_class_number" gorm:"size:36;not null"`
	CampusClass       CampusClass `gorm:"foreignKey:CampusClassNumber;references:Number;constraint: onUpdate:CASCADE, onDelete:CASCADE"`
	ClassName         string      `json:"class_name" gorm:"-"`
}

func (stream *Stream) AfterFind(tx *gorm.DB) (err error) {
	stream.CampusClass = CampusClass{}
	tx.Model(&stream.CampusClass).Where("number = ?", stream.CampusClassNumber).First(&stream.CampusClass)
	stream.ClassName = stream.CampusClass.ClassName
	return
}
