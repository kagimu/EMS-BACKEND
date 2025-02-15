package models

import (
	"time"

	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"gorm.io/gorm"
)

type Attendance struct {
	gorm.Model
	Number            string         `json:"number" gorm:"size:36;unique;not null"`
	Date              time.Time      `json:"date" gorm:"size:36;"`
	CheckInTime       time.Time      `json:"check_in_time"`
	CheckOutTime      time.Time      `json:"check_out_time"`
	Present           string         `json:"present"`
	CampusClassNumber string         `json:"campus_class_number" gorm:"size:36;"`
	SchoolNumber      string         `json:"school_number" gorm:"size:36;"`
	CampusNumber      string         `json:"campus_number" gorm:"size:36;"`
	UserNumber        string         `json:"user_number" gorm:"size:36;index;not null"`
	User              authModel.User `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	AddedBy           string         `json:"added_by" gorm:"size:36;"`
	StudentNumber     string         `json:"student_number" gorm:"size:36;"`
}
