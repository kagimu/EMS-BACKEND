package reporsitory

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/filters"
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
}

func AllReportColumnAssessments(reportColumnNumber, assessmentAreaNumber, assessmentMetricNumber, academicYearNumber, studyPeriodNumber, campusClassNumber, reportSettingNumber string) ([]models.ReportColumnAssessment, error) {

	db := config.DB
	var reportColumnAssessments []models.ReportColumnAssessment

	if err := db.
		Joins("JOIN report_columns ON report_columns.number = report_column_assessments.report_column_number").
		Joins("JOIN report_settings ON report_settings.number = report_columns.report_setting_number").
		Scopes(filters.ByAnyField("report_column_number", reportColumnNumber), filters.ByAcademicYearNumber(academicYearNumber), filters.ByAnyField("report_setting_number", reportSettingNumber), filters.ByCampusClassNumber(campusClassNumber), filters.ByStudyPeriodNumber(studyPeriodNumber), filters.ByAnyField("assessment_metric_number", assessmentMetricNumber), filters.ByAnyField("assessment_area_number", assessmentAreaNumber)).
		Find(&reportColumnAssessments).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return reportColumnAssessments, nil
		}

		return reportColumnAssessments, err
	}

	return reportColumnAssessments, nil
}

func DeleteReportColumnAssessmentsByReportColumnNumber(number string, tx *gorm.DB) error {

	db := config.DB

	if tx != nil {
		db = tx
	}

	if err := db.Where("report_column_number = ?", number).Delete(&models.ReportColumnAssessment{}).Error; err != nil {
		return err
	}

	return nil
}

func DeleteReportColumnAssessment(reportColumnAssessment models.ReportColumnAssessment, tx ...*gorm.DB) error {

	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}

	return db.Delete(&reportColumnAssessment).Error
}

func GetReportColumnAssessmentByNumber(number string) (models.ReportColumnAssessment, error) {
	db := config.DB

	var reportColumnAssessment models.ReportColumnAssessment
	if err := db.Where("number = ?", number).First(&reportColumnAssessment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return reportColumnAssessment, nil
		}
		return reportColumnAssessment, err
	}
	return reportColumnAssessment, nil
}
