package reporsitory

import (
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
