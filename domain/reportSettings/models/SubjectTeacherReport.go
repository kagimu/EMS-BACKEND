package models

import "time"

type SubjectTeacherReport struct {
	Number               string    `json:"number" gorm:"primaryKey;column:number"`
	SubjectTeacherName   string    `json:"subject_teacher_name" gorm:"column:subject_teacher_name"`
	Status               string    `json:"status" gorm:"column:status"`
	ClassTeacherComment  string    `json:"class_teacher_comment" gorm:"column:class_teacher_comment"`
	AssessmentAreaNumber string    `json:"assessment_area_number" gorm:"column:assessment_area_number"`
	StudyPeriodNumber    string    `json:"study_period_number" gorm:"column:study_period_number"`
	ClassNumber          string    `json:"class_number" gorm:"column:class_number"`
	Stage                string    `json:"stage" gorm:"column:stage"`
	CreatedAt            time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt            time.Time `json:"updated_at" gorm:"column:updated_at"`
}
