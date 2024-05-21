package models

type MealStock struct {
	ID             uint     `gorm:"primarykey"`
	MealNumber     string   `json:"meal_number" gorm:"size:36;not null"`
	Meal           Meal     `gorm:"foreignKey:MealNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
	StockingNumber string   `json:"stocking_number" gorm:"size:36;not null"`
	Stocking       Stocking `gorm:"foreignKey:StockingNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
}

func (reportTemplate *MealStock) TableName() string {
	return "meal_stock"
}
