package repository

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/enrollment/models"
	"gorm.io/gorm"
)

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
