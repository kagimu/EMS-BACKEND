package models

import (
	"time"
)

type (
	User struct {
		ID          uint       `gorm:"primaryKey"`
		Number      string     `json:"number" gorm:"size:100;unique;not null"`
		OldNum      string     `json:"old_num" gorm:"size:100;"`
		Name        string     `json:"name" form:"name" binding:"required" gorm:"size:60;not null"`
		DateOfBirth string     `json:"date_of_birth" form:"date_of_birth" binding:"required" gorm:"not null"`
		Gender      string     `json:"gender" form:"gender" binding:"required" gorm:"size:8;not null"`
		Type        string     `json:"type" gorm:"size:8;not null"`
		Status      string     `json:"status" gorm:"size:100;not null"`
		Info        string     `json:"info" gorm:"size:150;"`
		Seen        time.Time  `json:"seen"`
		CreatedAt   time.Time  `json:"created_at"`
		UpdatedAt   time.Time  `json:"updated_at"`
		DeletedAt   *time.Time `sql:"index"`
		Contacts    []Contact  `json:"contacts" gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	}

	GetUserList struct {
		ID          uint          `json:"id"`
		Number      string        `json:"number"`
		Name        string        `json:"name"`
		DateOfBirth string        `json:"date_of_birth"`
		Gender      string        `json:"gender"`
		Type        string        `json:"type"`
		Status      string        `json:"status"`
		Info        string        `json:"info"`
		Seen        time.Time     `json:"seen"`
		CreatedAt   time.Time     `json:"created_at"`
		ProfilePic  string        `json:"profile_pic"`
		CoverPic    string        `json:"cover_pic"`
		ContactList []ContactList `json:"contacts"`
	}

	PostUserData struct {
		Name        string `json:"name"`
		DateOfBirth string `json:"date_of_birth"`
		Gender      string `json:"gender"`
		Email       string `json:"email"`
		Phone       string `json:"phone"`
	}

	UserAnalytics struct {
		Users     int64 `json:"users"`
		AdminType int64 `json:"admin_type"`
		UserType  int64 `json:"user_type"`
	}
)

type JSON = map[string]interface{}

func (u *User) Read(m JSON) {
	u.ID = uint(m["id"].(float64))
	u.Number = m["number"].(string)
	u.Name = m["name"].(string)
}

func (u *User) TableName() string {
	return "users"
}
