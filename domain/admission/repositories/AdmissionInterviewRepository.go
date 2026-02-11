package repositories

import (
	"errors"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	admissionModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/admission/models"
	"gorm.io/gorm"
)

type AdmissionInterviewRepository interface {
	SaveAdmissionInterview(admissionInterview admissionModel.AdmissionInterview, txn ...*gorm.DB) error
	FindAdmissionInterviewByAdmissionNumber(admissionNumber string, txn ...*gorm.DB) (admissionModel.AdmissionInterview, error)
	DeleteAdmissionInterview(admissionNumber string, txn ...*gorm.DB) error
}

type admissionInterviewRepository struct {
}

func (aIRepository *admissionInterviewRepository) SaveAdmissionInterview(admissionInterview admissionModel.AdmissionInterview, txn ...*gorm.DB) error {
	db := config.DB
	if len(txn) > 0 {
		db = txn[0]
	}

	return db.Save(&admissionInterview).Error
}

func (aIRepository *admissionInterviewRepository) FindAdmissionInterviewByAdmissionNumber(admissionNumber string, txn ...*gorm.DB) (admissionModel.AdmissionInterview, error) {
	db := config.DB
	if len(txn) > 0 {
		db = txn[0]
	}

	var admissionInterview admissionModel.AdmissionInterview

	if err := db.Where("admission_number = ?", admissionNumber).First(&admissionInterview).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return admissionInterview, nil
		}
		return admissionInterview, err
	}

	return admissionInterview, nil
}

func (aIRepository *admissionInterviewRepository) DeleteAdmissionInterview(admissionNumber string, txn ...*gorm.DB) error {
	db := config.DB
	if len(txn) > 0 {
		db = txn[0]
	}

	return db.Where("admission_number = ?", admissionNumber).Delete(&admissionModel.AdmissionInterview{}).Error
}

func NewAdmissionInterviewRepository() AdmissionInterviewRepository {
	return &admissionInterviewRepository{}
}
