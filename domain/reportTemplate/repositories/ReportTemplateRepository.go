package repositories

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/filters"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportTemplate/models"
)

func QueryReportTemplates(schoolNumber string, search string) ([]models.ReportTemplate, error) {

	var queryResults []models.ReportTemplate

	db := config.DB

	query := db.Preload("User").Scopes(filters.ByName(search), filters.BySchoolNumber(schoolNumber)).Find(queryResults)

	if query.Error != nil {
		return nil, query.Error
	}

	return queryResults, nil
}
