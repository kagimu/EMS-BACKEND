package models

import (
	"gorm.io/gorm"
)

type StaffBioDataEducationItem struct {
	gorm.Model
	Number          string `json:"number" gorm:"size:36;unique;index;not null"`
	EducationNumber string `json:"education_number" gorm:"size:36;not null"`
	Level           string `json:"level" gorm:"size:50;not null"` // Primary, Secondary, HighSchool, PostHighSchool
	Institution     string `json:"institution" gorm:"size:255"`
	StartYear       string `json:"startYear" gorm:"size:50"`
	EndYear         string `json:"endYear" gorm:"size:50"`
	StartDate       string `json:"startDate" gorm:"size:50"`
	EndDate         string `json:"endDate" gorm:"size:50"`
	Award           string `json:"award" gorm:"size:255"`
	AwardDate       string `json:"awardDate" gorm:"size:50"`
}
