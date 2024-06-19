package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	assessmentMetricModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/models"
	"gorm.io/gorm"
)

func FindAssessmentMetricByNumber(assessmentMetricNumber string) (assessmentMetricModel.AssessmentMetric, error) {

	db := config.DB

	var assessmentMetric assessmentMetricModel.AssessmentMetric

	if err := db.Where("number = ?", assessmentMetricNumber).First(&assessmentMetric).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return assessmentMetric, nil
		}

		return assessmentMetric, err
	}

	return assessmentMetric, nil
}
