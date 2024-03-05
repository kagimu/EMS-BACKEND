package models

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	subjectModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/subject/models"
	"gorm.io/gorm"
)

type AssessmentArea struct {
	gorm.Model
	Number        string               `json:"number" gorm:"size:36;unique;not null"`
	Name          string               `json:"name" gorm:"not null;comment:Other name"`
	ClassNumber   string               `json:"class_number" gorm:"size:36;not null"`
	Class         string               `gorm:"foreignKey:ClassNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	UserNumber    string               `json:"user_number" gorm:"size:36;not null"`
	User          models.User          `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	SubjectNumber string               `json:"subject_number"`
	Subject       subjectModel.Subject `gorm:"foreignKey:SubjectNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}

//

func (as *AssessmentArea) AfterFind(tx *gorm.DB) (err error) {
	if as.SubjectNumber != "" {
		var subject subjectModel.Subject
		tx.Model(&as).Where("number = ?", as.SubjectNumber).First(&subject)
		as.Name = subject.Name
	}
	return nil
}
