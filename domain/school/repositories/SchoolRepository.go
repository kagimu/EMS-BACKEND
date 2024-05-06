package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	schoolModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/school/models"
	"gorm.io/gorm"
)

func GetSchoolByNumber(schoolNumber string) (schoolModel.School, error) {

	var school schoolModel.School

	db := config.DB

	if err := db.Where("number = ?", schoolNumber).First(&school).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return school, nil
		}

		return school, err
	}

	return school, nil

}
