package models

type EntityUser struct {
	EntityNumber string `json:"entity_number" gorm:"size:36"`
	UserNumber   string `json:"user_number" gorm:"size:36"`
	AddedBy      string `json:"added_by" gorm:"size:36"`
}
