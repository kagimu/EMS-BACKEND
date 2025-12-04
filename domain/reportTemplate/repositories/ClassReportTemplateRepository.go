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

	if err := db.
		Preload("User").
		Preload("CampusClass").
		Scopes(filters.ByReportTemplateNumber(templateNumber)).
		Find(&queryResults).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return queryResults, nil
		}

		return nil, err
	}

	return queryResults, nil
}

func SaveClassReportTemplate(classReportTemplate models.ClassReportTemplate) error {

	if err := config.DB.Create(&classReportTemplate).Error; err != nil {
		return err
	}

	return nil
}

func QueryByReportTemplateNumberAndCampusClassNumber(reportTemplateNumber string, campusClassNumber string) (models.ClassReportTemplate, error) {

	db := config.DB

	var queryResults models.ClassReportTemplate

	if err := db.
		Where("campus_class_number = ? ", campusClassNumber).
		Where("report_template_number = ? ", reportTemplateNumber).
		First(&queryResults).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return queryResults, nil
		}

		return queryResults, err
	}

	return queryResults, nil
}

func QueryByCampusClassNumber(campusClassNumber string) (models.ClassReportTemplate, error) {

	db := config.DB

	var queryResults models.ClassReportTemplate

	if err := db.
		Where("campus_class_number = ? ", campusClassNumber).
		Order("created_at DESC").
		First(&queryResults).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return queryResults, nil
		}

		return queryResults, err
	}

	return queryResults, nil
}

func DeleteCampusClassReportTemplate(classReportTemplate models.ClassReportTemplate) error {

	db := config.DB

	if err := db.
		Delete(&classReportTemplate).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}

		return err
	}

	return nil
}

func QueryByClassReportTemplateNumber(classReportTemplateNumber string) (models.ClassReportTemplate, error) {

	db := config.DB

	var queryResults models.ClassReportTemplate

	if err := db.
		Where("number = ? ", classReportTemplateNumber).
		First(&queryResults).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return queryResults, nil
		}

		return queryResults, err
	}

	return queryResults, nil
}
