package reporsitory

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	reportSettings "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/models"
	"gorm.io/gorm"
)

func CountReportSettingTotalPercentageScores(reportSettingNumber string) (float64, error) {

	db := config.DB

	var totalPercentageShare float64

	if err := db.Model(&reportSettings.ReportColumn{}).
		Where("report_setting_number = ? ", reportSettingNumber).
		Select("COALESCE(sum(percentage_share),0)").
		Scan(&totalPercentageShare).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil
		}

		return totalPercentageShare, err
	}

	return totalPercentageShare, nil
}

func SaveReportColumn(reportSettingColumn reportSettings.ReportColumn, trn ...*gorm.DB) error {

	db := config.DB

	if len(trn) > 0 {
		db = trn[0]
	}

	return db.Save(&reportSettingColumn).Error
}

func FindReportColumnsByReportSettingsNumber(reportSettingsNumber string) ([]reportSettings.ReportColumn, error) {

	db := config.DB

	var reportColumns []reportSettings.ReportColumn

	if err := db.Where("report_setting_number = ?", reportSettingsNumber).Find(&reportColumns).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return reportColumns, nil
		}
		return reportColumns, err
	}

	return reportColumns, nil
}

func GetReportColumnByNameAndReportSettingsNumber(name string, reportSettingsNumber string) (reportSettings.ReportColumn, error) {
	db := config.DB
	var reportColumn reportSettings.ReportColumn

	if err := db.Where("name = ? AND report_setting_number = ?", name, reportSettingsNumber).First(&reportColumn).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return reportColumn, nil
		}
		return reportColumn, err
	}

	return reportColumn, nil
}

func GetReportColumnByCodeAndReportSettingsNumber(code int64, reportSettingsNumber string) (reportSettings.ReportColumn, error) {
	db := config.DB
	var reportColumn reportSettings.ReportColumn

	if err := db.Where("code = ? AND report_setting_number = ?", code, reportSettingsNumber).First(&reportColumn).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return reportColumn, nil
		}
		return reportColumn, err
	}

	return reportColumn, nil
}
