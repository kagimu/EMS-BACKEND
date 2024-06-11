package models

import userModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"

type ActiveUserEntity struct {
	EntityNumber string          `gorm:"size:36"`
	Entity       Entity          `gorm:"foreignKey:EntityNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
	UserNumber   string          `gorm:"size:36;primarykey"`
	User         userModels.User `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
}
