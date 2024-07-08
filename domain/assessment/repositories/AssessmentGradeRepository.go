package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/assessmentGrade/models"
	"gorm.io/gorm"
)

func GetAssessmentGradeByNumber(assessmentGradeNumber string, trn ...*gorm.DB) (models.AssessmentGrade, error) {

	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	var assessmentGrade models.AssessmentGrade

	if err := db.Where("number = ?", assessmentGradeNumber).First(&assessmentGrade).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return assessmentGrade, nil
		}

		return assessmentGrade, err
	}

	return assessmentGrade, nil
}
