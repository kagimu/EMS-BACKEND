package models

import (
	"time"
)

type StudentBioData struct {
	Number        string    `json:"number" gorm:"size:36;primaryKey"`
	StudentNumber string    `json:"student_number" gorm:"size:36;not null"`
	FileName      string    `json:"file_name" gorm:"size:255;not null"`
	FilePath      string    `json:"file_path" gorm:"type:text;not null"`
	FileType      string    `json:"file_type" gorm:"size:50"`
	FileSize      int64     `json:"file_size"`
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	Student       Student   `gorm:"foreignKey:StudentNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
}
