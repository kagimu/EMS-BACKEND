package repository

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/filters"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/enrollment/models"
	"gorm.io/gorm"
)

type EnrolmentRepository interface {
	GetStudentEnrolment(studentNumber, studyPeriodNumber, status string, txn ...*gorm.DB) (models.Enrollment, error)
}

type enrolmentRepository struct {
}

func (enrolmentRepository *enrolmentRepository) GetStudentEnrolment(studentNumber, studyPeriodNumber, status string, txn ...*gorm.DB) (models.Enrollment, error) {
	db := config.DB
	if len(txn) > 0 {
		db = txn[0]
	}

	var enrolment models.Enrollment
	if err := db.
		Scopes(filters.ByStatus(status)).
		Where("student_number = ? AND study_period_number = ?", studentNumber, studyPeriodNumber).
		Find(&enrolment).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Enrollment{}, nil
		}

		return models.Enrollment{}, err
	}

	return enrolment, nil
}

func NewEnrolmentRepository() EnrolmentRepository {
	return &enrolmentRepository{}
}

func GetEnrollmentByStudentNumber(studentNumber string, txn ...*gorm.DB) (models.Enrollment, error) {
	db := config.DB

	if len(txn) > 0 {
		db = txn[0]

	}

	var enrollment models.Enrollment

	if err := db.Where("student_number = ?", studentNumber).First(&enrollment).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return enrollment, nil
		}

		return enrollment, err
	}

	return enrollment, nil
}
