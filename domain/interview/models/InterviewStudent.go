package models

import (
	appModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/application/models"
	studentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/models"
	"gorm.io/gorm"
	"time"
)

type InterviewStudent struct {
	ID              uint `gorm:"primarykey"`
	CreatedAt       time.Time
	Number          string               `json:"number" gorm:"size:36;unique;not null"`
	InterviewNumber string               `json:"interview_number" gorm:"size:36;not null"`
	StudentNumber   string               `json:"student_number" gorm:"size:36;not null"`
	Student         studentModel.Student `gorm:"foreignKey:StudentNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	Interview       Interview            `gorm:"foreignKey:InterviewNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	Application     appModel.Application `gorm:"-"`
}

func (i *InterviewStudent) AfterFind(tx *gorm.DB) (err error) {
	var student studentModel.Student
	if err := tx.First(&student, "number = ?", i.StudentNumber).Error; err != nil {
		return err
	}
	i.Student = student
	// find application for student
	var application appModel.Application
	if err := tx.First(&application, "number = ?", student.ApplicationNumber).Error; err != nil {
		return err
	}
	i.Application = application

	return nil
}
