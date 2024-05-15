package models

import "time"

type MealCourseFeedingSession struct {
	ID                   uint           `gorm:"primarykey"`
	CreatedAt            time.Time      `json:"created_at"`
	Number               string         `json:"number" gorm:"size:36;unique;not null"`
	FeedingSessionNumber string         `json:"feeding_session_number" gorm:"size:36;"`
	FeedingSession       FeedingSession `gorm:"foreignKey:FeedingSessionNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	MealCourseNumber     string         `json:"meal_course_number" gorm:"size:36;"`
	MealCourse           MealCourse     `gorm:"foreignKey:MealCourseNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
