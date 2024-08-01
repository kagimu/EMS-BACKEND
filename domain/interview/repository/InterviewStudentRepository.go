package repository

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/interview/models"
	"gorm.io/gorm"
)

func FindInterviewStudentsByInterViewNumber(interviewNumber string, tx ...*gorm.DB) ([]models.InterviewStudent, error) {

	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}

	var interviewStudents = make([]models.InterviewStudent, 0)

	err := db.Preload("Student").Preload("Interview").Where("interview_number = ?", interviewNumber).Find(&interviewStudents).Error
	return interviewStudents, err
}
