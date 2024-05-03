package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/filters"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportTemplate/models"
	"gorm.io/gorm"
)

func QueryClassReportTemplates(templateNumber string) ([]models.ClassReportTemplate, error) {

	db := config.DB

	var queryResults []models.ClassReportTemplate

	err := db.
		Preload("ReportTemplate").
		Preload("User").
		Preload("CampusClass").
		Scopes(filters.ByReportTemplateNumber(templateNumber)).
		Find(queryResults).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return queryResults, nil
	}

	if err != nil {
		return nil, err
	}

	return queryResults, nil
}

func SaveClassReportTemplate(classReportTemplate models.ClassReportTemplate) error {
	return config.DB.Create(classReportTemplate).Error
}

func QueryByReportTemplateNumberAndCampusClassNumber(reportTemplateNumber string, campusClassNumber string) (models.ClassReportTemplate, error) {

	db := config.DB

	var queryResults models.ClassReportTemplate

	err := db.
		Where("campus_class_number = ? ", campusClassNumber).
		Where("report_template_number = ? ", reportTemplateNumber).
		First(queryResults).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return queryResults, nil
	}

	if err != nil {
		return queryResults, err
	}

	return queryResults, nil
}
