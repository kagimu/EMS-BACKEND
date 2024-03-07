package models

import "gorm.io/gorm"

type Task struct {
	gorm.Model
	Number         string `json:"number" gorm:"size:36;unique;not null"`
	DatePerformed  string `json:"date_performed"`
	TimeFrame      string `json:"time_frame"`
	InputsUsed     string `json:"inputs_used" gorm:"type:text"`
	Outputs        string `json:"outputs" gorm:"type:text"`
	PeopleInvolved string `json:"people_involved" gorm:"type:text"`
	StaffNumber    string `json:"staff_number" gorm:"size:36;not null"`
	Staff          Staff  `gorm:"foreignKey:StaffNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
