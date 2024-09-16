package models

import (
	"time"
)

type Invoice struct {
	Number          string    `json:"number" gorm:"size:36;primaryKey;not null"`
	Status          string    `json:"status" gorm:"default:NOT_PAID"`
	Description     string    `json:"description" gorm:"type:text"`
	Type            string    `json:"type" gorm:"default:CLIENT"`
	DueDate         time.Time `json:"due_date" gorm:"type:date"`
	RefNumber       string    `json:"ref_number"  gorm:"size:16;not null"`
	StatementNumber string    `json:"statement_number" gorm:"size:36;not null"`
	Statement       Statement `gorm:"foreignKey:StatementNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
}
