package models

import (
	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	campusModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/models"
	schoolModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/school/models"
	"time"
)

type Statement struct {
	Number       string `json:"number" gorm:"size:36;primaryKey;not null"`
	CreatedAt    time.Time
	RefNumber    string              `json:"ref_number"  gorm:"size:16;not null"`
	Amount       float64             `json:"amount"`
	Balance      float64             `json:"balance"`
	Type         string              `json:"type"`
	Description  string              `json:"description" gorm:"size:750;not null"`
	SchoolNumber string              `json:"school_number" gorm:"size:36;not null"`
	School       schoolModel.School  `gorm:"foreignKey:SchoolNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
	CampusNumber string              `json:"campus_number" gorm:"size:36;not null"`
	Campus       campusModels.Campus `gorm:"foreignKey:CampusNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
	//StreamNumber  string              `json:"stream_number" gorm:"size:36;not null"`
	//Stream        models.Stream       `gorm:"foreignKey:StreamNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
	UserNumber    string         `json:"user_number" gorm:"size:36;not null"`
	User          authModel.User `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
	AddedByNumber string         `json:"added_by_number" gorm:"size:36;not null"`
	AddedBy       authModel.User `gorm:"foreignKey:AddedByNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
}
