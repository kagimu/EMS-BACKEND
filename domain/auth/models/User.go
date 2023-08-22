package models

import (
	"time"
)

type (
	User struct {
		ID         uint       `gorm:"primaryKey"`
		Number     string     `json:"number" gorm:"size:100;unique;not null"`
		FirstName  string     `json:"first_name" form:"name" binding:"required" gorm:"size:60;not null"`
		SecondName string     `json:"second_name" form:"name" binding:"required" gorm:"size:60;not null"`
		LastName   string     `json:"last_name" form:"name" binding:"required" gorm:"size:60;"`
		Gender     string     `json:"gender" form:"gender" binding:"required" gorm:"size:8;not null"`
		Type       string     `json:"type" gorm:"size:8;not null"`
		Status     string     `json:"status" gorm:"size:100;not null"`
		CreatedAt  time.Time  `json:"created_at"`
		UpdatedAt  time.Time  `json:"updated_at"`
		DeletedAt  *time.Time `sql:"index"`
		Contacts   []Contact  `json:"contacts" gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	}

	GetUserList struct {
		ID         uint          `json:"id"`
		Number     string        `json:"number"`
		FirstName  string        `json:"first_name"`
		SecondName string        `json:"second_name"`
		LastName   string        `json:"last_name"`
		Gender     string        `json:"gender"`
		Type       string        `json:"type"`
		Status     string        `json:"status"`
		CreatedAt  time.Time     `json:"created_at"`
		UpdatedAt  time.Time     `json:"updated_at"`
		DeletedAt  *time.Time    `sql:"index"`
		Contacts   []ContactList `json:"contacts"`
	}

	UserDataRequest struct {
		FirstName  string `json:"first_name"`
		SecondName string `json:"second_name"`
		LastName   string `json:"last_name"`
		Gender     string `json:"gender"`
		Email      string `json:"email"`
		Phone      string `json:"phone"`
	}
)

func (u *User) TableName() string {
	return "users"
}
