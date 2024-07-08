package reporsitory

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/models"
	"gorm.io/gorm"
)

func GetReportSummaryByStudentNumberAndStudyPeriodNumberAndCampusClassNumber(studentNumber string, studyPeriodNumber string, campusClassNumber string, trn ...*gorm.DB) (models.ReportSummary, error) {

	db := config.DB

	if len(trn) > 0 {
		db = trn[0]
	}

	var reportSummary models.ReportSummary

	if err := db.Where("student_number = ? AND study_period_number = ? AND campus_class_number = ?", studentNumber, studyPeriodNumber, campusClassNumber).First(&reportSummary).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return reportSummary, nil
		}

		return reportSummary, err
	}

	return reportSummary, nil
}

func SaveReportSummary(summary models.ReportSummary, trn ...*gorm.DB) error {

	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	return db.Save(&summary).Error
}

func FindDuplicateReportSummary(studentNumber, studyPeriodNumber, campusClassNumber string, trn ...*gorm.DB) (models.ReportSummary, error) {

	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	var reportSummary models.ReportSummary

	if err := db.Where("student_number = ? AND study_period_number = ? AND campus_class_number = ?", studentNumber, studyPeriodNumber, campusClassNumber).First(&reportSummary).Error; err == nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return reportSummary, nil
		}

		return reportSummary, err
	}

	return reportSummary, nil
}
