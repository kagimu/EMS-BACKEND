package reporsitory

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
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
