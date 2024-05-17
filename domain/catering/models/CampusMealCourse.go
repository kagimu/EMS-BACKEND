package models

import (
	userModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	campusModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/models"
	"gorm.io/gorm"
)

type CampusMealCourse struct {
	gorm.Model
	Number                         string                   `json:"number" gorm:"size:36;unique;not null"`
	UserNumber                     string                   `json:"user_number" gorm:"size:36;not null"`
	User                           userModel.User           `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	CampusNumber                   string                   `json:"campus_number" gorm:"size:36;not null"`
	Campus                         campusModel.Campus       `gorm:"foreignKey:CampusNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	Day                            string                   `json:"day" gorm:"size:4;not null"`
	MealCourseFeedingSessionNumber string                   `json:"meal_course_feeding_session_number" gorm:"size:36;not null"`
	MealCourseFeedingSession       MealCourseFeedingSession `gorm:"foreignKey:MealCourseFeedingSessionNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
}
