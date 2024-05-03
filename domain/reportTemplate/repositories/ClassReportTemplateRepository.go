package repositories

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/filters"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportTemplate/models"
)

func QueryClassReportTemplates(templateNumber string) ([]models.ClassReportTemplate, error) {

	db := config.DB

	var queryResults []models.ClassReportTemplate

	err := db.
		Preload("ReportTemplate").
		Preload("User").
		Preload("CampusClass").
		Scopes(filters.ByReportTemplateNumber(templateNumber)).
		Find(queryResults)

	if err.Error != nil {
		return nil, err.Error
	}

	return queryResults, nil
}
