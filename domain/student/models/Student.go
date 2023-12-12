package models

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	classModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
	sourceOfSchoolInfoModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/sourceOfSchoolInfo/models"
	"gorm.io/gorm"
)

type Student struct {
	gorm.Model
	Number                   string                                      `json:"number" gorm:"size:36;unique;not null"`
	APN                      string                                      `gorm:"unique;size:100;" json:"apn"`
	Surname                  string                                      `json:"surname" `
	SecondName               string                                      `json:"second_name"`
	OtherName                string                                      `json:"other_name"`
	Age                      uint                                        `json:"age"`
	Gender                   string                                      `json:"gender"  binding:"required" gorm:"size:8;not null"`
	Reason                   string                                      `json:"reason"`
	DateOfBirth              string                                      `json:"date_of_birth" form:"date_of_birth" binding:"required" gorm:"not null"`
	PreviousSchool           []PreviousSchools                           `gorm:"foreignKey:StudentNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	ClassNumber              string                                      `json:"class_number" gorm:"size:36;not null"`
	Class                    classModels.Class                           `gorm:"foreignKey:ClassNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	SourceOfSchoolInfoNumber string                                      `json:"source_of_school_info_number" gorm:"size:36;not null"`
	SourceOfSchoolInfo       sourceOfSchoolInfoModels.SourceOfSchoolInfo `gorm:"foreignKey:SourceOfSchoolInfoNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	UserNumber               string                                      `json:"user_number" gorm:"size:36;index;not null"`
	User                     models.User                                 `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
