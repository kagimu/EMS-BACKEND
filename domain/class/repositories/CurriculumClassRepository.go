package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	CurriculumClassModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
	"gorm.io/gorm"
)

func GetCurriculumClass(curriculumClassNumber string, txn ...*gorm.DB) (CurriculumClassModel.CurriculumClass, error) {
	db := config.DB
	if len(txn) > 0 {
		db = txn[0]
	}

	var curriculumClass CurriculumClassModel.CurriculumClass

	if err := db.Where("number = ?", curriculumClassNumber).First(&curriculumClass).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return CurriculumClassModel.CurriculumClass{}, nil
		}

		return CurriculumClassModel.CurriculumClass{}, err
	}

	return curriculumClass, nil
}
