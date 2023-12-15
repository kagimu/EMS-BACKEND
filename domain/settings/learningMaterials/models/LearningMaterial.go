package models

import "gorm.io/gorm"
import "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
import "errors"

type LearningMaterial struct {
	gorm.Model
	Number     string      `json:"number" gorm:"size:36;unique;not null;index"`
	Code       int64       `json:"code" form:"code" binding:"required" gorm:"size:80;not null"`
	Name       string      `json:"name" form:"name" binding:"required" gorm:"size:80;not null"`
	Type       string      `json:"type" form:"type" binding:"required" gorm:"size:80;not null"`
	UserNumber string      `json:"user_number" gorm:"size:36;not null"`
	User       models.User `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}

func (l *LearningMaterial) Validate() error {
	//validate all the required fields here

	if l.Name == "" {
		return errors.New("name is required")
	}

	if l.Type == "" {
		return errors.New("type is required")
	}

	if l.UserNumber == "" {
		return errors.New("user number is required")
	}

	return nil
}
