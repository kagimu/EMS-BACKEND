package models

import (
	userModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	itemModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/item/models"
	"gorm.io/gorm"
)

type EnrollmentItem struct {
	gorm.Model
	EnrollmentNumber string          `json:"enrollment_number" gorm:"size:36;not null"`
	Enrollment       Enrollment      `gorm:"foreignKey:EnrollmentNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	ItemNumber       string          `json:"item_number" gorm:"size:36; not null"`
	Item             itemModels.Item `gorm:"foreignKey:ItemNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	RequiredQuantity float64         `json:"required_quantity" gorm:"type:decimal(10,2);default:0"`
	Quantity         float64         `json:"quantity" gorm:"type:decimal(10,2);default:0"`
	Brought          bool            `json:"brought"`
	UserNumber       string          `json:"user_number" gorm:"size:36; not null"`
	User             userModel.User  `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}

type RequirementClearanceLog struct {
	gorm.Model
	EnrollmentNumber     string         `json:"enrollment_number" gorm:"size:36;index"`
	EnrollmentItemID     uint           `json:"enrollment_item_id" gorm:"index"`
	ItemNumber           string         `json:"item_number" gorm:"size:36"`
	ItemName             string         `json:"item_name" gorm:"size:255"`
	DateCleared          string         `json:"date_cleared" gorm:"size:20"`
	QuantityCleared      float64        `json:"quantity_cleared" gorm:"type:decimal(10,2)"`
	TotalQuantityBrought float64        `json:"total_quantity_brought" gorm:"type:decimal(10,2)"`
	RequiredQuantity     float64        `json:"required_quantity" gorm:"type:decimal(10,2)"`
	RemainingQuantity    float64        `json:"remaining_quantity" gorm:"type:decimal(10,2)"`
	ClearedByStaffNumber string         `json:"cleared_by_staff_number" gorm:"size:36"`
	ClearedByStaffName   string         `json:"cleared_by_staff_name" gorm:"size:255"`
	Remarks              string         `json:"remarks" gorm:"size:500"`
	User                 userModel.User `gorm:"foreignKey:ClearedByStaffNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:SET NULL;"`
}

