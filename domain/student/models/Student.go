package models

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	campusModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/models"
	classModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
	parentModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/parent/models"
	serviceModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/service/models"
	sourceOfSchoolInfoModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/sourceOfSchoolInfo/models"
	"gorm.io/gorm"
)

type Student struct {
	gorm.Model
	Number                    string                                      `json:"number" gorm:"size:36;unique;not null"`
	APN                       string                                      `gorm:"unique;size:100;unique" json:"apn"`
	APNStatus                 string                                      `gorm:"size:20;default:Inactive" json:"apn_status"`
	Surname                   string                                      `json:"surname" `
	SecondName                string                                      `json:"second_name"`
	OtherName                 string                                      `json:"other_name"`
	AcademicReport            string                                      `json:"academic_report"`
	Gender                    string                                      `json:"gender"  binding:"required" gorm:"size:8;not null"`
	Reason                    string                                      `json:"reason" gorm:"type:text"`
	ReasonForAdmissionDenial  string                                      `json:"reason_for_admission_denial"`
	AttentionNeeded           string                                      `json:"attention_needed"`
	DateOfBirth               string                                      `json:"date_of_birth" form:"date_of_birth" binding:"required" gorm:"not null"`
	PreferredTypeOfEnrollment string                                      `json:"preferred_type_of_enrollment" gorm:"size:100;not null"`
	PreferedCampusNumber      string                                      `json:"preferred_campus_number" gorm:"size:36; not null"`
	PreferredCampus           campusModels.Campus                         `gorm:"foreignKey:PreferedCampusNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	ClassNumber               string                                      `json:"class_number" gorm:"size:36; not null"`
	Class                     classModels.Class                           `gorm:"foreignKey:ClassNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	Services                  []serviceModel.Service                      `gorm:"foreignKey:StudentNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	SourceOfSchoolInfoNumber  string                                      `json:"source_of_school_info_number" gorm:"size:36;not null"`
	SourceOfSchoolInfo        sourceOfSchoolInfoModels.SourceOfSchoolInfo `gorm:"foreignKey:SourceOfSchoolInfoNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	ParentNumber              string                                      `json:"parent_number" gorm:"size:36;not null"`
	Parent                    parentModels.Parent                         `gorm:"foreignKey:ParentNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	UserNumber                string                                      `json:"user_number" gorm:"size:36;index;not null"`
	User                      models.User                                 `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
