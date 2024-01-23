package models

import (
	authModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/models"
	"gorm.io/gorm"
)

type Interview struct {
	gorm.Model
	Number         string         `json:"number" gorm:"size:36;unique;not null"`
	Message        string         `json:"message"`
	Venue          string         `json:"venue" gorm:"size:100;not null"`
	Name           string         `json:"name" gorm:"size:100;"`
	QuestionNumber string         `json:"question_number" gorm:"size:36;"`
	InterviewDate  string         `json:"interview_date" gorm:"size:100;not null"`
	StudentNumber  string         `json:"student_number" gorm:"size:36;"`
	Student        models.Student `gorm:"foreignKey:StudentNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	UserNumber     string         `json:"user_number" gorm:"size:36;not null"`
	User           authModel.User `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}

func (i *Interview) AfterFind(tx *gorm.DB) (err error) {
	var student models.Student
	if err := tx.First(&student, "number = ?", i.StudentNumber).Error; err != nil {
		return err
	}
	i.Student = student
	return nil
}
