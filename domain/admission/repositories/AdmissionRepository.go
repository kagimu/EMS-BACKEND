package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	admissionModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/admission/models"
	"gorm.io/gorm"
)

func FindAdmissionByStudentNumber(studentNumber string) (admissionModel.Admission, error) {
	db := config.DB
	var admission admissionModel.Admission

	if err := db.Where("student_number = ?", studentNumber).First(&admission).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return admission, nil
		}

		return admission, err
	}

	return admission, nil
}
