package repository

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/filters"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/enrollment/models"
	"gorm.io/gorm"
)

type EnrolmentRepository interface {
	GetStudentEnrolment(studentNumber, studyPeriodNumber, campusClassNumber, status string, txn ...*gorm.DB) (models.Enrollment, error)
	AllEnrolments(studentNumber, studyPeriodNumber, campusClassNumber, campusNumber, status string, txn ...*gorm.DB) ([]models.Enrollment, error)
	SaveEnrolment(enrollment models.Enrollment, txn ...*gorm.DB) error
}

type enrolmentRepository struct {
}

func (enrolmentRepository *enrolmentRepository) AllEnrolments(studentNumber, studyPeriodNumber, campusClassNumber, campusNumber, status string, txn ...*gorm.DB) ([]models.Enrollment, error) {
	db := config.DB
	if len(txn) > 0 {
		db = txn[0]
	}

	var enrolment []models.Enrollment
	if err := db.
		Scopes(filters.ByStatus(status), filters.ByCampusNumber(campusNumber), filters.ByStudyPeriodNumber(studyPeriodNumber), filters.ByStudentNumber(studentNumber), filters.ByCampusClassNumber(campusClassNumber)).
		Find(&enrolment).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []models.Enrollment{}, nil
		}

		return []models.Enrollment{}, err
	}

	return enrolment, nil
}

func (enrolmentRepository *enrolmentRepository) SaveEnrolment(enrollment models.Enrollment, txn ...*gorm.DB) error {
	db := config.DB
	if len(txn) > 0 {
		db = txn[0]
	}

	return db.Save(&enrollment).Error
}

func (enrolmentRepository *enrolmentRepository) GetStudentEnrolment(studentNumber, studyPeriodNumber, campusClassNumber, status string, txn ...*gorm.DB) (models.Enrollment, error) {
	db := config.DB
	if len(txn) > 0 {
		db = txn[0]
	}

	var enrolment models.Enrollment
	if err := db.
		Scopes(filters.ByStatus(status), filters.ByCampusClassNumber(campusClassNumber)).
		Where("student_number = ? AND study_period_number = ?", studentNumber, studyPeriodNumber).
		First(&enrolment).Error; err != nil {

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
