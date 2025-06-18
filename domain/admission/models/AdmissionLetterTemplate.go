package models

import "time"

type AdmissionLetterTemplate struct {
	Number       string    `json:"number" gorm:"size:36;primaryKey"`
	SchoolNumber string    `json:"school_number" gorm:"size:36;not null"`
	Phrase       string    `json:"phase" gorm:"not null:type:text"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
