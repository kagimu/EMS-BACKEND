package reporsitory

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	reportSettings "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/models"
	"gorm.io/gorm"
)

func SaveReportColumn(reportSettingColumn reportSettings.ReportColumn, trn ...*gorm.DB) error {

	db := config.DB

	if len(trn) > 0 {
		db = trn[0]
	}

	return db.Save(&reportSettingColumn).Error
}

func FindReportColumnsByReportSettingsNumber(reportSettingsNumber string, tx ...*gorm.DB) ([]reportSettings.ReportColumn, error) {

	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}

	var reportColumns = make([]reportSettings.ReportColumn, 0)

	if err := db.Where("report_setting_number = ?", reportSettingsNumber).
		Order("code asc").
		Find(&reportColumns).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return reportColumns, nil
		}
		return reportColumns, err
	}

	return reportColumns, nil
}

func GetReportColumnByNumber(reportColumnNumber string, tx ...*gorm.DB) (reportSettings.ReportColumn, error) {
	db := config.DB

	if len(tx) > 0 {
		db = tx[0]
	}

	var reportColumn reportSettings.ReportColumn
	if err := db.Where("number = ?", reportColumnNumber).First(&reportColumn).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return reportColumn, nil
		}
		return reportColumn, err
	}
	return reportColumn, nil
}

func DeleteReportColumn(reportColumn reportSettings.ReportColumn, tx ...*gorm.DB) error {

	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}

	return db.Delete(&reportColumn).Error
}

func DeleteReportColumnsByReportSettingNumber(number string, tx *gorm.DB) error {
	return tx.Where("report_setting_number = ?", number).Delete(&reportSettings.ReportColumn{}).Error
}
