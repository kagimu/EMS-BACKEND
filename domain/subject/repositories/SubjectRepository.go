package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	subjectModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/subject/models"
	"gorm.io/gorm"
)

func GetSubjectByNumber(subjectNumber string, trn ...*gorm.DB) (subjectModels.Subject, error) {

	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	var subject subjectModels.Subject

	if err := db.Where("number = ?", subjectNumber).First(&subject).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return subject, nil
		}

		return subject, err
	}

	return subject, nil

}
