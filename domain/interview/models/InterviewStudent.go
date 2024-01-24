package models

import (
	studentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/models"
	"gorm.io/gorm"
)

type InterviewStudent struct {
	gorm.Model
	Number          string               `json:"number" gorm:"size:36;unique;not null"`
	InterviewNumber string               `json:"interview_number" gorm:"size:36;not null"`
	StudentNumber   string               `json:"student_number" gorm:"size:36;not null"`
	Student         studentModel.Student `gorm:"foreignKey:StudentNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	Interview       Interview            `gorm:"foreignKey:InterviewNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}

func (i *InterviewStudent) AfterFind(tx *gorm.DB) (err error) {
	var student studentModel.Student
	if err := tx.Preload("Application").First(&student, "number = ?", i.StudentNumber).Error; err != nil {
		return err
	}
	i.Student = student
	return nil
}
