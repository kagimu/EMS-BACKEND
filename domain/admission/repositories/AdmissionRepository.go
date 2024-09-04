package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	admissionModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/admission/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/filters"
	"gorm.io/gorm"
)

type AdmissionRepository interface {
	FindAdmissionByStudentNumber(studentNumber string, txn ...*gorm.DB) (admissionModel.Admission, error)
	SaveAdmission(admission admissionModel.Admission, txn ...*gorm.DB) error
	AllAdmissions(campusNumber, studyPeriodNumber, streamNumber, interviewNumber, schoolNumber, campusClassNumber string, paginateFunc func(db *gorm.DB) *gorm.DB, txn ...*gorm.DB) ([]admissionModel.Admission, error)
}

type admissionRepository struct {
}

func (admissionRepository *admissionRepository) AllAdmissions(campusNumber, studyPeriodNumber, streamNumber, interviewNumber, schoolNumber, campusClassNumber string, paginateFunc func(db *gorm.DB) *gorm.DB, txn ...*gorm.DB) ([]admissionModel.Admission, error) {

	db := config.DB
	if len(txn) > 0 {
		db = txn[0]
	}

	var admissions []admissionModel.Admission

	if err := db.Scopes(
		filters.ByStudyPeriodNumber(studyPeriodNumber),
		filters.ByCampusNumber(campusNumber),
		filters.ByStreamNumber(streamNumber),
		filters.ByInterviewNumber(interviewNumber),
		filters.BySchoolNumber(schoolNumber),
		filters.ByCampusClassNumber(campusClassNumber),
		paginateFunc,
	).Find(&admissions).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return admissions, nil
		}

		return admissions, err
	}

	return admissions, nil
}

func (admissionRepository *admissionRepository) SaveAdmission(admission admissionModel.Admission, txn ...*gorm.DB) error {
	db := config.DB
	if len(txn) > 0 {
		db = txn[0]
	}

	return db.Save(&admission).Error
}

func (admissionRepository *admissionRepository) FindAdmissionByStudentNumber(studentNumber string, txn ...*gorm.DB) (admissionModel.Admission, error) {
	db := config.DB
	if len(txn) > 0 {
		db = txn[0]
	}

	var admission admissionModel.Admission

	if err := db.Where("student_number = ?", studentNumber).First(&admission).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return admission, nil
		}

		return admission, err
	}

	return admission, nil
}

func NewAdmissionRepository() AdmissionRepository {
	return &admissionRepository{}
}
