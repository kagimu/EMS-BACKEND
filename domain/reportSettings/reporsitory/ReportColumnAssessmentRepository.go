package reporsitory

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/models"
	"gorm.io/gorm"
)

func GetReportColumnAssessmentByAssessmentMetricNumber(assessmentMetricNumber string) (models.ReportColumnAssessment, error) {

	db := config.DB
	var reportColumnAssessment models.ReportColumnAssessment

	if err := db.Where("assessment_metric_number = ? ", assessmentMetricNumber).First(&reportColumnAssessment).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return reportColumnAssessment, nil
		}

		return reportColumnAssessment, err
	}

	return reportColumnAssessment, nil
}

func SaveReportColumnAssessment(reportColumnAssessment models.ReportColumnAssessment, trn ...*gorm.DB) error {

	db := config.DB

	if len(trn) > 0 {
		db = trn[0]
	}

	if err := db.Save(&reportColumnAssessment).Error; err != nil {
		return err
	}

	return nil
}
