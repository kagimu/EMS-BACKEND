package models

type StaffMealStock struct {
	ID              uint      `gorm:"primarykey"`
	StaffMealNumber string    `json:"staff_meal_number" gorm:"size:36;not null"`
	StaffMeal       StaffMeal `gorm:"foreignKey:StaffMealNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
	StockingNumber  string    `json:"stocking_number" gorm:"size:36;not null"`
	Stocking        Stocking  `gorm:"foreignKey:StockingNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
}

func (StaffMealStock) TableName() string {
	return "staff_meal_stock"
}
