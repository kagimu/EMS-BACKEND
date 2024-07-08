package reporsitory

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	reportSettingsModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/models"
	"gorm.io/gorm"
)

func GetReportSettingByNumber(reportSettingsNumber string, trn ...*gorm.DB) (reportSettingsModel.ReportSetting, error) {

	db := config.DB

	if len(trn) > 0 {
		db = trn[0]
	}

	var reportSetting reportSettingsModel.ReportSetting

	if err := db.Where("number = ?", reportSettingsNumber).First(&reportSetting).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return reportSetting, nil
		}

		return reportSetting, err
	}

	return reportSetting, nil
}

func SaveReportSettings(setting reportSettingsModel.ReportSetting, trn ...*gorm.DB) error {

	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	return db.Save(&setting).Error
}

func FindDuplicateReportSettings(studyPeriodNumber, academicYearNumber, campusClassNumber, status string, trn ...*gorm.DB) (reportSettingsModel.ReportSetting, error) {

	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	var reportSetting reportSettingsModel.ReportSetting

	if err := db.Where("study_period_number = ? AND academic_year_number = ? AND campus_class_number = ? AND status = ?", studyPeriodNumber, academicYearNumber, campusClassNumber, status).First(&reportSetting).Error; err == nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return reportSetting, nil
		}

		return reportSetting, err
	}

	return reportSetting, nil
}

func GetReportSettingByCampusClassNumberAndStatus(campusClassNumber, status string) (reportSettingsModel.ReportSetting, error) {
	var reportSetting reportSettingsModel.ReportSetting

	if err := config.DB.Where("campus_class_number = ? AND status = ?", campusClassNumber, status).First(&reportSetting).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return reportSettingsModel.ReportSetting{}, nil
		}
		return reportSettingsModel.ReportSetting{}, err

	}
	return reportSetting, nil
}

func FindReportColumnsByReportSettingNumber(reportSettingNumber string) ([]reportSettingsModel.ReportColumn, error) {
	var reportColumns []reportSettingsModel.ReportColumn

	if err := config.DB.Where("report_setting_number = ?", reportSettingNumber).Find(&reportColumns).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return make([]reportSettingsModel.ReportColumn, 0), errors.New("report columns for this class not yet set up")

		}
		return reportColumns, err
	}

	return reportColumns, nil
}
