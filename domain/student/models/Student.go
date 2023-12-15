package models

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	sourceOfSchoolInfoModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/sourceOfSchoolInfo/models"
	"gorm.io/gorm"
)

type Student struct {
	gorm.Model
	Number                    string                                      `json:"number" gorm:"size:36;unique;not null"`
	APN                       string                                      `gorm:"unique;size:100;" json:"apn"`
	APNStatus                 string                                      `gorm:"size:20;default:Inactive" json:"apn_status"`
	Surname                   string                                      `json:"surname" `
	SecondName                string                                      `json:"second_name"`
	OtherName                 string                                      `json:"other_name"`
	Age                       uint                                        `json:"age"`
	Gender                    string                                      `json:"gender"  binding:"required" gorm:"size:8;not null"`
	Reason                    string                                      `json:"reason"`
	DateOfBirth               string                                      `json:"date_of_birth" form:"date_of_birth" binding:"required" gorm:"not null"`
	PreferredTypeOfEnrollment string                                      `json:"preferred_type_of_enrollment" gorm:"size:100;not null"`
	PreviousSchoolNumber      string                                      `json:"previous_school_number" gorm:"size:36;not null"`
	PreviousSchool            PreviousSchools                             `gorm:"foreignKey:PreviousSchoolNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	SourceOfSchoolInfoNumber  string                                      `json:"source_of_school_info_number" gorm:"size:36;not null"`
	SourceOfSchoolInfo        sourceOfSchoolInfoModels.SourceOfSchoolInfo `gorm:"foreignKey:SourceOfSchoolInfoNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	UserNumber                string                                      `json:"user_number" gorm:"size:36;index;not null"`
	User                      models.User                                 `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
