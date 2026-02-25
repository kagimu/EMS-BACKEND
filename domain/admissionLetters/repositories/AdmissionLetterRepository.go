package repositories

import (
	"errors"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/admissionLetters/models"
	"gorm.io/gorm"
)

func SaveAdmissionLetter(letter *models.AdmissionLetter, tx ...*gorm.DB) error {
	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}
	return db.Save(letter).Error
}

func FindAdmissionLetterByNumber(number string, tx ...*gorm.DB) (models.AdmissionLetter, error) {
	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}
	var letter models.AdmissionLetter
	if err := db.Where("number = ?", number).First(&letter).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return letter, nil
		}
		return letter, err
	}
	return letter, nil
}

func FindAdmissionLettersByStudentNumber(studentNumber string, tx ...*gorm.DB) ([]models.AdmissionLetter, error) {
	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}
	var letters []models.AdmissionLetter
	if err := db.Where("student_number = ?", studentNumber).Order("created_at desc").Find(&letters).Error; err != nil {
		return nil, err
	}
	return letters, nil
}

func FindAdmissionLettersBySchoolNumber(schoolNumber string, tx ...*gorm.DB) ([]models.AdmissionLetter, error) {
	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}
	var letters []models.AdmissionLetter
	if err := db.Where("school_number = ?", schoolNumber).Order("created_at desc").Find(&letters).Error; err != nil {
		return nil, err
	}
	return letters, nil
}

func DeleteAdmissionLetter(letter *models.AdmissionLetter, tx ...*gorm.DB) error {
	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}
	return db.Delete(letter).Error
}
