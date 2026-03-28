package reporsitory

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	reportSettingsModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/models"
	"gorm.io/gorm"
)

func SaveSubjectTeacherReport(report reportSettingsModel.SubjectTeacherReport, trn ...*gorm.DB) error {
	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}
	return db.Save(&report).Error
}

func GetSubjectTeacherReportByNumber(number string, trn ...*gorm.DB) (reportSettingsModel.SubjectTeacherReport, error) {
	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}
	var report reportSettingsModel.SubjectTeacherReport
	if err := db.Where("number = ?", number).First(&report).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return reportSettingsModel.SubjectTeacherReport{}, nil
		}
		return reportSettingsModel.SubjectTeacherReport{}, err
	}
	return report, nil
}

func FetchAllSubjectTeacherReports(studyPeriodNumber, stage, status string, trn ...*gorm.DB) ([]reportSettingsModel.SubjectTeacherReport, error) {
	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	var reports []reportSettingsModel.SubjectTeacherReport
	query := db.Model(&reportSettingsModel.SubjectTeacherReport{})

	if studyPeriodNumber != "" {
		query = query.Where("study_period_number = ?", studyPeriodNumber)
	}
	if stage != "" {
		query = query.Where("stage = ?", stage)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Find(&reports).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return make([]reportSettingsModel.SubjectTeacherReport, 0), nil
		}
		return nil, err
	}

	return reports, nil
}

func DeleteSubjectTeacherReport(report reportSettingsModel.SubjectTeacherReport, trn ...*gorm.DB) error {
	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}
	return db.Delete(&report).Error
}

func SaveSubjectTeacherReportComment(comment reportSettingsModel.SubjectTeacherReportComment, trn ...*gorm.DB) error {
	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}
	return db.Save(&comment).Error
}

func FetchSubjectTeacherReportComments(reportNumber string, trn ...*gorm.DB) ([]reportSettingsModel.SubjectTeacherReportComment, error) {
	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	var comments []reportSettingsModel.SubjectTeacherReportComment
	if err := db.Where("subject_teacher_report_number = ?", reportNumber).Order("created_at desc").Find(&comments).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return make([]reportSettingsModel.SubjectTeacherReportComment, 0), nil
		}
		return nil, err
	}
	return comments, nil
}
