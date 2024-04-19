package models

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	campusModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/models"
	"gorm.io/gorm"
)

type Staff struct {
	gorm.Model
	Number            string             `json:"number" gorm:"size:36;unique;not null"`
	StaffIDNumber     string             `json:"staff_id_number"`
	JobTitleNumber    string             `json:"job_title_number" gorm:"size:36;not null"`
	JobTitle          JobTitle           `gorm:"foreignKey:JobTitleNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	CampusNumber      string             `json:"campus_number" gorm:"size:36;not null"`
	Campus            campusModel.Campus `gorm:"foreignKey:CampusNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	SalaryScaleNumber string             `json:"salary_scale_number" gorm:"size:36;not null"`
	SalaryScale       SalaryScale        `gorm:"foreignKey:SalaryScaleNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	UserNumber        string             `json:"user_number" gorm:"size:36;index;not null"`
	User              models.User        `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
