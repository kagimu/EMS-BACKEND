package models

import (
	authModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	campusModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/models"
	classModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
	studentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/models"

	"gorm.io/gorm"
)

type Admission struct {
	gorm.Model
	Number       string               `json:"number" gorm:"size:36;unique;not null"`
	ADN          string               `gorm:"unique;size:100;unique" json:"adn"`
	StudentAPN   string               `gorm:"unique;size:100;unique" json:"student_apn"`
	Status       string               `json:"status"`
	Student      studentModel.Student `gorm:"foreignKey:StudentAPN;references:APN;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	CampusNumber string               `json:"campus_number" gorm:"size:36;not null"`
	Campus       campusModels.Campus  `gorm:"foreignKey:CampusNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	ClassNumber  string               `json:"class_number" gorm:"size:36;not null"`
	Class        classModels.Class    `gorm:"foreignKey:ClassNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	StreamNumber string               `json:"stream_number" gorm:"size:36;not null"`
	Stream       classModels.Stream   `gorm:"foreignKey:StreamNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	UserNumber   string               `json:"user_number" gorm:"size:36;not null"`
	User         authModels.User      `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}

// Add an application
// Use the APN as the join foreign key
// Track application whether draft or complete
// Every form step makes an api call.
// All I do is create a parent numner, appn number
// STudent number, application number
// Query the appn, join paremt, join student, join class, join school, classes as an array, classes attendend, capture the periods attending the schools attended,
