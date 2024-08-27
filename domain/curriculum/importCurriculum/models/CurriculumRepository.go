package models

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"gorm.io/gorm"
)

func GetCurriculumByNumber(curriculumNumber string, txn ...*gorm.DB) (Curriculum, error) {
	db := config.DB
	if len(txn) > 0 {
		db = txn[0]
	}

	var curriculum Curriculum

	if err := db.Where("number = ?", curriculumNumber).First(&curriculum).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return curriculum, nil
		}

		return curriculum, err
	}

	return curriculum, nil
}
