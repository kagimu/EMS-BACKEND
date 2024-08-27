package reporsitory

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/filters"
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

func FetchAllReportSettings(campusNumber, classNumber, campusClassNumber, academicYearNumber, studyPeriodNumber string, txn ...*gorm.DB) ([]reportSettingsModel.ReportSetting, error) {

	db := config.DB
	if len(txn) > 0 {
		db = txn[0]
	}

	var reportSettings []reportSettingsModel.ReportSetting

	if err := db.
		Joins("JOIN campus_classes ON campus_classes.number = report_settings.campus_class_number").
		Joins("JOIN curriculum_classes ON curriculum_classes.number = campus_classes.curriculum_class_number").
		Scopes(
			filters.ByAnyField("campus_classes.campus_number", campusNumber),
			filters.ByAnyField("curriculum_classes.class_number", classNumber),
			filters.ByAnyField("report_settings.campus_class_number", campusClassNumber),
			filters.ByAnyField("report_settings.academic_year_number", academicYearNumber),
			filters.ByAnyField("report_settings.study_period_number", studyPeriodNumber),
		).
		Distinct("report_settings.*").
		Find(&reportSettings).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return reportSettings, nil
		}

		return reportSettings, err
	}

	return reportSettings, nil
}
