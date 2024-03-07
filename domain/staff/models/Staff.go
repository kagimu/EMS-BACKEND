package models

import (
	campusModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/models"
	"gorm.io/gorm"
)

type Staff struct {
	gorm.Model
	Number                     string             `json:"number" gorm:"size:36;unique;not null"`
	UniqueStaffNumber          string             `json:"unique_staff_number" gorm:"unique;"`
	Surname                    string             `json:"surname" gorm:"not null"`
	Firstname                  string             `json:"firstname" gorm:"not null"`
	OtherName                  string             `json:"other_name"`
	MainTelephoneNumber        string             `json:"main_telephone_number" gorm:"size:20; not null"`
	AlternativeTelephoneNumber string             `json:"alternative_telephone_number" gorm:"size:20"`
	MainEmailAddress           string             `json:"main_email_address" gorm:"size:100; not null"`
	AlternativeEmailAddress    string             `json:"alternative_email_address" gorm:"size:100"`
	JobTitleNumber             string             `json:"job_title_number" gorm:"size:36;not null"`
	JobTitle                   JobTitle           `gorm:"foreignKey:JobTitleNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	CampusNumber               string             `json:"campus_number" gorm:"size:36;not null"`
	Campus                     campusModel.Campus `gorm:"foreignKey:CampusNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	SalaryScaleNumber          string             `json:"salary_scale_number" gorm:"size:36;not null"`
	SalaryScale                SalaryScale        `gorm:"foreignKey:SalaryScaleNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
