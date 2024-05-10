package models

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"gorm.io/gorm"
)

type Item struct {
	gorm.Model
	Number             string      `json:"number" gorm:"size:36;unique;not null"`
	Name               string      `json:"name" gorm:"not null"`
	Description        string      `json:"description" gorm:"type:text"`
	Type               string      `json:"type" gorm:"not null;comment:'Service or Product'"`
	Cost               float64     `json:"cost" gorm:"not null"`
	Prefix             string      `json:"prefix"`
	PaymentMode        string      `json:"payment_mode" gorm:"default:'CASH'"`
	Mandatory          bool        `json:"mandatory" gorm:"default:false"`
	SerialCode         string      `json:"serial_code"`
	DiscountPercentage uint        `json:"discount_percentage"`
	DiscountAmount     uint        `json:"discount_amount"`
	OneTimePurchase    bool        `json:"one_time_purchase" gorm:"default:false"`
	RepeatCycle        string      `json:"repeat_cycle" gorm:"size:20"`
	UserNumber         string      `json:"user_number" gorm:"size:36;not null"`
	User               models.User `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
