package models

import (
	"time"
)

type FamilyMember struct {
	Number                 string    `json:"number" gorm:"size:36;primaryKey"`
	StudentNumber          string    `json:"student_number" gorm:"size:36;not null"`
	Surname                string    `json:"surname" gorm:"size:50;not null"`
	IsLiving               string    `json:"is_living" gorm:"size:50"`
	OtherName              string    `json:"other_name" gorm:"size:50;not null"`
	PlaceOfEmployment      string    `json:"place_of_employment" gorm:"size:50;not null"`
	PhysicalContactAddress string    `json:"physical_contact_address" gorm:"size:100;not null"`
	Relationship           string    `json:"relationship" gorm:"size:100;not null"`
	CreatedAt              time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt              time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	Student                Student   `gorm:"foreignKey:StudentNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
}

type FamilyMemberContact struct {
	Number             string       `json:"number" gorm:"size:36;primaryKey"`
	FamilyMemberNumber string       `json:"family_member_number" gorm:"size:36;not null"`
	PhoneNumber        string       `json:"phone_number" gorm:"size:50"`
	Email              string       `json:"email" gorm:"size:100;"`
	CreatedAt          time.Time    `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt          time.Time    `json:"updated_at" gorm:"autoUpdateTime"`
	FamilyMember       FamilyMember `gorm:"foreignKey:FamilyMemberNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
}

type StudentLanguage struct {
	Number        string    `json:"number" gorm:"size:36;primaryKey"`
	StudentNumber string    `json:"student_number" gorm:"size:36;not null"`
	Code          int64     `json:"code"`
	Name          string    `json:"name" gorm:"size:50;not null"`
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	Student       Student   `gorm:"foreignKey:StudentNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
}
