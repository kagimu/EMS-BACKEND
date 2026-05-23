package models

import (
	"time"

	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	campusModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/models"
	schoolModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/school/models"
)

type Credit struct {
	Number          string              `json:"number" gorm:"size:36;primaryKey;not null"`
	CreatedAt       time.Time           `json:"created_at"`
	UpdatedAt       time.Time           `json:"updated_at"`
	Date            time.Time           `json:"date" gorm:"not null"`
	Amount          float64             `json:"amount"`
	Balance         float64             `json:"balance"`
	Status          string              `json:"status" gorm:"size:36;not null;default:ACTIVE"`
	ClientReference string              `json:"client_reference" gorm:"size:62"`
	Description     string              `json:"description" gorm:"size:500"`
	SchoolNumber    string              `json:"school_number" gorm:"size:36;not null"`
	School          schoolModel.School  `gorm:"foreignKey:SchoolNumber;references:Number;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	CampusNumber    string              `json:"campus_number" gorm:"size:36;not null"`
	Campus          campusModels.Campus `gorm:"foreignKey:CampusNumber;references:Number;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	UserNumber      string              `json:"user_number" gorm:"size:36;not null"`
	User            authModel.User      `gorm:"foreignKey:UserNumber;references:Number;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	AddedByNumber   string              `json:"added_by_number" gorm:"size:36;not null"`
	AddedBy         authModel.User      `gorm:"foreignKey:AddedByNumber;references:Number;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
}

func (Credit) TableName() string {
	return "credits"
}
