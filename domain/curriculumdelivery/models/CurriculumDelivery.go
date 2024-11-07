package models

import "time"

type CurriculumDelivery struct {
	Number               string    `json:"number" gorm:"size:36;unique;not null"`
	StreamNumber         string    `json:"stream_number" gorm:"size:36;not null"`
	LessonPlanStepNumber string    `json:"lesson_plan_step_number" gorm:"size:36;not null"`
	Date                 time.Time `json:"date" gorm:"not null"`
	Time                 string    `json:"time" gorm:"size:10;not null"`
	TeacherNumber        string    `json:"teacher_number" gorm:"size:36;not null"`
	Comment              string    `json:"comment" gorm:"size:250"`
	UserNumber           string    `json:"user_number" gorm:"size:36;not null"`
	ID                   uint      `json:"id" gorm:"primaryKey"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

type StudentCurriculumDelivery struct {
	ID                       uint      `json:"id" gorm:"primaryKey"`
	CreatedAt                time.Time `json:"created_at"`
	StudentNumber            string    `json:"student_number" gorm:"size:36;not null"`
	CurriculumDeliveryNumber string    `json:"curriculum_delivery_number" gorm:"size:36;not null"`
	Comment                  string    `json:"comment" gorm:"size:250"`
}
