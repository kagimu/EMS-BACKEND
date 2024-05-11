package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	studentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/models"
	"gorm.io/gorm"
)

func FindStudentByUserNumber(userNumber string) (studentModel.Student, error) {

	db := config.DB
	var student studentModel.Student

	if err := db.Where("user_number = ?", userNumber).First(&student).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return student, nil
		}

		return student, err
	}

	return student, nil

}
