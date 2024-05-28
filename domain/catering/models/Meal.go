package models

import (
	campusModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/models"
	"time"

	authModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"gorm.io/gorm"
)

type Meal struct {
	gorm.Model
	Number               string             `json:"number" gorm:"size:36;unique;not null"`
	UserNumber           string             `json:"user_number" gorm:"size:36;not null"`
	User                 authModels.User    `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	FeedingSessionNumber string             `json:"feeding_session_number" gorm:"size:36;not null"`
	FeedingSession       FeedingSession     `gorm:"foreignKey:FeedingSessionNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	MealCourseNumber     string             `json:"meal_course_number" gorm:"size:36;not null"`
	MealCourse           MealCourse         `gorm:"foreignKey:MealCourseNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	StaffNumber          string             `json:"staff_number" gorm:"size:36;not null"`
	Staff                authModels.User    `gorm:"foreignKey:StaffNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	Date                 time.Time          `json:"date" gorm:"not null;type:date"`
	Time                 string             `json:"time" gorm:"not null"`
	Day                  string             `json:"day" gorm:"not null"`
	AddedByNumber        string             `json:"added_by_number" gorm:"size:36;not null"`
	AddedBy              authModels.User    `gorm:"foreignKey:AddedByNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	FeedingGroupNumber   string             `json:"feeding_group_number" gorm:"size:36;not null"`
	FeedingGroup         FeedingGroup       `gorm:"foreignKey:FeedingGroupNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	CampusNumber         string             `json:"campus_number" gorm:"size:36;not null"`
	Campus               campusModel.Campus `gorm:"foreignKey:CampusNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}

func (reportTemplate *Meal) TableName() string {
	return "meals"
}
