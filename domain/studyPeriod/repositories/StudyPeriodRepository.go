package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	studyPeriodModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/studyPeriod/models"
	"gorm.io/gorm"
)

func GetStudyPeriodByNumber(studyPeriodNumber string) (studyPeriodModel.StudyPeriod, error) {
	db := config.DB
	var studyPeriod studyPeriodModel.StudyPeriod
	if err := db.Where("number = ?", studyPeriodNumber).First(&studyPeriod).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return studyPeriod, nil
		}

		return studyPeriod, err
	}
	return studyPeriod, nil
}
