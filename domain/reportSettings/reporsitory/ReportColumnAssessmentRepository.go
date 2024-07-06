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

func FindReportColumnAssessmentsByReportColumnNumber(reportColumnNumber string) ([]models.ReportColumnAssessment, error) {

	db := config.DB
	var reportColumnAssessments []models.ReportColumnAssessment

	if err := db.Where("report_column_number = ?", reportColumnNumber).Find(&reportColumnAssessments).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return reportColumnAssessments, nil
		}

		return reportColumnAssessments, err
	}

	return reportColumnAssessments, nil
}

func FindReportColumnAssessmentByReportColumnNumberAndAssessmentAreaNumber(reportColumnNumber string, assessmentAreaNumber string) ([]models.ReportColumnAssessment, error) {

	db := config.DB
	var reportColumnAssessment []models.ReportColumnAssessment

	if err := db.Where("report_column_number = ? AND assessment_area_number = ?", reportColumnNumber, assessmentAreaNumber).Find(&reportColumnAssessment).Error; err != nil {
		return reportColumnAssessment, err
	}

	return reportColumnAssessment, nil
	//tx.Where("assessment_area_number = ? AND report_column_number = ? ", assessmentArea.Number, reportColum.Number).Find(&rcas).Error;
}
