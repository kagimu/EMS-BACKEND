package models

import (
	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	campusModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/models"
	classModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
	schoolModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/school/models"
	studentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/models"
	"gorm.io/gorm"
)

type AdmissionOffer struct {
	gorm.Model
	Number                        string               `json:"number" gorm:"size:36;unique;not null"`
	ClassNumber                   string               `json:"class_number" gorm:"size:36;not null"`
	Class                         classModel.Class     `gorm:"foreignKey:ClassNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	SchoolNumber                  string               `json:"school_number" gorm:"size:36;not null"`
	School                        schoolModel.School   `gorm:"foreignKey:SchoolNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	CampusNumber                  string               `json:"campus_number" gorm:"size:36;not null"`
	AdmissionOfferConditionNumber string               `json:"admission_offer_condition_number" gorm:"size:36;"`
	Campus                        campusModel.Campus   `gorm:"foreignKey:CampusNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	StudentNumber                 string               `json:"student_number" gorm:"size:36;not null"`
	CustomAdmissionCondition      string               `json:"custom_admission_condition" gorm:"type:text"`
	Student                       studentModel.Student `gorm:"foreignKey:StudentNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	UserNumber                    string               `json:"user_number" gorm:"size:36;not null"`
	User                          authModel.User       `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
