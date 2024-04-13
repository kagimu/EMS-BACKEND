package models

import (
	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	lessonPlanStepModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/lessonPlan/models"
	"gorm.io/gorm"
)

type StudyPeriod struct {
	gorm.Model
	Number                  string                              `json:"number" gorm:"size:36;not null"`
	StudyingTimeTableNumber string                              `json:"studying_time_table_number" gorm:"size:36;not null"`
	StudyingTimeTable       StudyingTimeTable                   `gorm:"foreignKey:StudyingTimeTableNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	UserNumber              string                              `json:"user_number" gorm:"size:36;not null"`
	User                    authModel.User                      `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	Comment                 string                              `json:"comment" gorm:"type:text;"`
	LessonPlanStepNumber    string                              `json:"lesson_plan_step_number" gorm:"size:36;not null"`
	LessonPlanStep          lessonPlanStepModels.LessonPlanStep `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}
