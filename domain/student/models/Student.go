package models

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	campusModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/models"
	classModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
	"gorm.io/gorm"
)

type Student struct {
	gorm.Model
	Number    string `json:"number" gorm:"size:36;unique;not null"`
	APN       string `gorm:"unique;size:100;" json:"apn"`
	ADN       string `gorm:"unique;size:100;default:null" json:"adn"`
	APNStatus string `gorm:"size:20;default:Inactive" json:"apn_status"`
	ADNStatus string `gorm:"size:20;" json:"adn_status"`
	Surname   string `json:"surname" `
	// UserNumber
	SecondName                string                  `json:"second_name"`
	OtherName                 string                  `json:"other_name"`
	Photo                     string                  `json:"photo"`
	AcademicReport            string                  `json:"academic_report"`
	Gender                    string                  `json:"gender"  binding:"required" gorm:"size:8;not null"`
	ReasonForChoosingSchool   string                  `json:"reason_for_choosing_school" gorm:"type:text"`
	ApplicationNumber         string                  `json:"application_number" gorm:"size:36;not null"`
	ReasonForAdmissionDenial  string                  `json:"reason_for_admission_denial"`
	AttentionNeeded           string                  `json:"attention_needed"`
	DateOfBirth               string                  `json:"date_of_birth" form:"date_of_birth" binding:"required" gorm:"not null"`
	PreferredTypeOfEnrollment string                  `json:"preferred_type_of_enrollment" gorm:"size:100;not null"`
	PreferedCampusNumber      string                  `json:"preferred_campus_number" gorm:"size:36; not null"`
	PreferredCampus           campusModels.Campus     `gorm:"foreignKey:PreferedCampusNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	CampusClassNumber         string                  `json:"campus_class_number" gorm:"size:36; not null"`
	CampusClass               classModels.CampusClass `gorm:"foreignKey:CampusClassNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	SourceOfSchoolInfo        string                  `json:"source_of_school_info" gorm:"size:100;not null"`
	UserNumber                string                  `json:"user_number" gorm:"size:36;index;not null"`
	User                      models.User             `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
