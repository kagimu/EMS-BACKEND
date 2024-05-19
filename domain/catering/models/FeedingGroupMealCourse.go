package models

import (
	authModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"time"
)

type FeedingGroupMealCourse struct {
	CreatedAt          time.Time       `json:"created_at"`
	Number             string          `json:"number" gorm:"size:36;primaryKey;not null"`
	FeedingGroupNumber string          `json:"feeding_group_number" gorm:"size:36;not null"`
	FeedingGroup       FeedingGroup    `gorm:"foreignKey:FeedingGroupNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
	MealCourseNumber   string          `json:"meal_course_number" gorm:"size:36;not null"`
	MealCourse         MealCourse      `gorm:"foreignKey:MealCourseNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
	Day                string          `json:"day" gorm:"size:4;not null"`
	UserNumber         string          `json:"user_number" gorm:"size:36;index;not null"`
	User               authModels.User `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
