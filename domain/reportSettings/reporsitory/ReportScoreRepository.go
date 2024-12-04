package reporsitory

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/filters"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/models"
	"gorm.io/gorm"
)

func CreateReportScore(reportScore models.ReportScore, trn ...*gorm.DB) error {

	db := config.DB

	if len(trn) > 0 {
		db = trn[0]
	}

	if err := db.Create(&reportScore).Error; err != nil {
		return err
	}

	return nil
}

func SaveReportScore(reportScore models.ReportScore, trn ...*gorm.DB) error {

	db := config.DB

	if len(trn) > 0 {
		db = trn[0]
	}

	if err := db.Save(&reportScore).Error; err != nil {
		return err
	}

	return nil
}

func FindAllReportScore(reportSummaryNumber, assessmentAreaNumber string, trn ...*gorm.DB) ([]models.ReportScore, error) {

	db := config.DB

	if len(trn) > 0 {
		db = trn[0]
	}

	var reportScores []models.ReportScore
	if err := db.
		Scopes(filters.ByAnyField("report_summary_number", reportSummaryNumber), filters.ByAnyField("assessment_area_number", assessmentAreaNumber)).
		Find(&reportScores).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []models.ReportScore{}, nil
		}
		return nil, err
	}

	return reportScores, nil
}
