package models

import authModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"

type MealCourse struct {
	Number                string          `json:"number" gorm:"size:36;unique;not null;primarykey"`
	MealCourseCountNumber string          `json:"meal_course_count_number" gorm:"size:36;not null"`
	MealCourseCount       MealCourseCount `gorm:"foreignKey:MealCourseCountNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	Name                  string          `json:"name" gorm:"size:120;not null"`
	UserNumber            string          `json:"user_number" gorm:"size:36;not null"`
	User                  authModels.User `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	Status                string          `json:"status" gorm:"size:36;not null"`
	Code                  int             `json:"code" gorm:"size:10;not null"`
}
