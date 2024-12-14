package reporsitory

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/models"
	"gorm.io/gorm"
)

func CreateReportScoreSummary(reportScoreSummary models.ReportScoreSummary, trn ...*gorm.DB) error {

	db := config.DB

	if len(trn) > 0 {
		db = trn[0]
	}

	if err := db.Create(&reportScoreSummary).Error; err != nil {
		return err
	}

	return nil
}

func SaveReportScoreSummary(reportScoreSummary models.ReportScoreSummary, trn ...*gorm.DB) error {

	db := config.DB

	if len(trn) > 0 {
		db = trn[0]
	}

	if err := db.Save(&reportScoreSummary).Error; err != nil {
		return err
	}

	return nil
}

func FindReportScoreSummaryByReportSummaryAndAssessmentArea(reportSummaryNumber, assessmentAreaNumber string, trn ...*gorm.DB) (models.ReportScoreSummary, error) {

	db := config.DB

	if len(trn) > 0 {
		db = trn[0]
	}

	var reportScoreSummary models.ReportScoreSummary
	if err := db.Where("report_summary_number = ? AND assessment_area_number =? ", reportSummaryNumber, assessmentAreaNumber).First(&reportScoreSummary).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return reportScoreSummary, nil
		}

		return reportScoreSummary, err
	}

	return reportScoreSummary, nil
}
