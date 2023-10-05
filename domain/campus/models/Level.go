package models

import "gorm.io/gorm"

type (
	Level struct {
		gorm.Model
		Number       string `json:"number" gorm:"size:36;unique;not null"`
		Name         string `json:"name" form:"name" binding:"required" gorm:"size:2000;"`
		CampusNumber string `json:"campus_number" gorm:"size:36;not null"`
	}
	LevelDataInfo struct {
		Name   string `json:"name"`
		Number string `json:"number"`
	}
)
