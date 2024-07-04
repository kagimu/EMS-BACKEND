package repositories

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/assessmentGrade/models"
	"gorm.io/gorm"
)

func SaveGradingStandardClassScale(gradingStandardClassScale models.GradingStandardClassScale, trn ...*gorm.DB) error {

	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	return db.Save(&gradingStandardClassScale).Error
}
