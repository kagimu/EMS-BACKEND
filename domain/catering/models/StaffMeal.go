package models

import (
	authModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	campusModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/models"
	staffMembersModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/staffMembers/models"
	"gorm.io/gorm"
	"time"
)

type StaffMeal struct {
	gorm.Model
	Number               string                             `json:"number" gorm:"size:36;unique;not null"`
	StaffNumber          string                             `json:"staff_number" gorm:"size:36;not null;index"`
	StaffDeployment      staffMembersModels.StaffDeployment `gorm:"foreignKey:StaffNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	FeedingSessionNumber string                             `json:"feeding_session_number" gorm:"size:36;not null"`
	FeedingSession       FeedingSession                     `gorm:"foreignKey:FeedingSessionNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	MealCourseNumber     string                             `json:"meal_course_number" gorm:"size:36;not null"`
	MealCourse           MealCourse                         `gorm:"foreignKey:MealCourseNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	AddedByNumber        string                             `json:"added_by_number" gorm:"size:36;not null"`
	AddedBy              authModels.User                    `gorm:"foreignKey:AddedByNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	Date                 time.Time                          `json:"date" gorm:"not null;type:date"`
	Time                 string                             `json:"time" gorm:"not null"`
	Day                  string                             `json:"day" gorm:"not null"`
	FeedingGroupNumber   string                             `json:"feeding_group_number" gorm:"size:36;not null"`
	FeedingGroup         FeedingGroup                       `gorm:"foreignKey:FeedingGroupNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	CampusNumber         string                             `json:"campus_number" gorm:"size:36;not null"`
	Campus               campusModel.Campus                 `gorm:"foreignKey:CampusNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}

func (StaffMeal) TableName() string {
	return "staff_meals"
}
