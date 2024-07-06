package reporsitory

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	reportSettings "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/models"
	"gorm.io/gorm"
)

func GetReportSettingByNumber(reportSettingsNumber string, trn ...*gorm.DB) (reportSettings.ReportSetting, error) {

	db := config.DB

	if len(trn) > 0 {
		db = trn[0]
	}

	var reportSetting reportSettings.ReportSetting

	if err := db.Where("number = ?", reportSettingsNumber).First(&reportSetting).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return reportSetting, nil
		}

		return reportSetting, err
	}

	return reportSetting, nil
}

func SaveReportSettings(setting reportSettings.ReportSetting, trn ...*gorm.DB) error {

	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	return db.Save(&setting).Error
}

func FindDuplicateReportSettings(studyPeriodNumber, academicYearNumber, campusClassNumber, status string, trn ...*gorm.DB) (reportSettings.ReportSetting, error) {

	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	var reportSetting reportSettings.ReportSetting

	if err := db.Where("study_period_number = ? AND academic_year_number = ? AND campus_class_number = ? AND status = ?", studyPeriodNumber, academicYearNumber, campusClassNumber, status).First(&reportSetting).Error; err == nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return reportSetting, nil
		}

		return reportSetting, err
	}

	return reportSetting, nil
}
