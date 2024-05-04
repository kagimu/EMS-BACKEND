package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/filters"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportTemplate/models"
	"gorm.io/gorm"
)

func QueryReportTemplates(schoolNumber string, search string, reportTemplateNumber string) ([]models.ReportTemplate, error) {

	var queryResults []models.ReportTemplate

	db := config.DB

	if err := db.
		Preload("User").
		Scopes(filters.ByName(search), filters.ByNumber(reportTemplateNumber), filters.BySchoolNumber(schoolNumber)).
		Find(queryResults).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return queryResults, nil
		}

		return nil, err
	}

	return queryResults, nil
}

func QueryByNameAndSchool(templateName string, schoolNumber string) (models.ReportTemplate, error) {

	var reportTemplate models.ReportTemplate

	db := config.DB

	if err := db.Where("name = ? ", templateName).Where("school_number = ? ", schoolNumber).First(reportTemplate).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return reportTemplate, nil
		}

		return reportTemplate, err
	}

	return reportTemplate, nil
}

func Save(reportTemplate models.ReportTemplate) error {

	db := config.DB

	if savingError := db.Create(&reportTemplate).Error; savingError != nil {
		return savingError
	}

	return nil
}

func QueryByReportTemplateNumber(reportTemplateNumber string) (models.ReportTemplate, error) {

	db := config.DB

	var reportTemplate models.ReportTemplate

	if err := db.Where("number = ? ", reportTemplateNumber).First(reportTemplate).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return reportTemplate, nil
		}

		return reportTemplate, err
	}

	return reportTemplate, nil
}
