package models

import "gorm.io/gorm"

type LandingPage struct {
	gorm.Model
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Image       string `json:"image"`
}
