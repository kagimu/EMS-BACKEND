package models

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"gorm.io/gorm"
)

type Stream struct {
	gorm.Model
	Number           string      `json:"number" gorm:"size:36;not null;unique"`
	Name             string      `json:"name" gorm:"size:200;not null"`
	UserNumber       string      `json:"user_number" gorm:"size:36;not null"`
	User             models.User `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	ClassLevelNumber string      `json:"class_level_number" gorm:"size:36;not null"`
	ClassLevel       ClassLevel  `gorm:"foreignKey:ClassLevelNumber;references:Number;constraint: onUpdate:CASCADE, onDelete:CASCADE"`
	ClassName        string      `json:"class_name" gorm:" - :migration"`
}

func (stream *Stream) AfterFind(tx *gorm.DB) (err error) {

	stream.ClassLevel = ClassLevel{}
	tx.Model(&stream.ClassLevel).Where("number = ?", stream.ClassLevelNumber).First(&stream.ClassLevel)
	stream.ClassName = stream.ClassLevel.ClassName
	return
}
