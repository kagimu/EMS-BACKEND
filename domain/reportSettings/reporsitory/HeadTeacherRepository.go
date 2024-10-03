package reporsitory

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/models"
)

func SaveHeadTeacherComment(headTeacherComment models.HeadTeacherComment) error {
	if err := config.DB.Save(&headTeacherComment).Error; err != nil {
		return err
	}
	return nil
}

func CreateHeadTeacherComment(headTeacherComment models.HeadTeacherComment) error {
	if err := config.DB.Create(&headTeacherComment).Error; err != nil {
		return err
	}
	return nil
}

func FindHeadTeacherComments() ([]models.HeadTeacherComment, error) {
	var headTeacherComment []models.HeadTeacherComment
	if err := config.DB.Preload("ReportSetting").
		Preload("ReportSetting.AcademicYear").
		Preload("ReportSetting.StudyPeriod").
		Preload("ReportSetting.CampusClass").
		Find(&headTeacherComment).Error; err != nil {
		return headTeacherComment, err
	}
	return headTeacherComment, nil
}

func FindHeadTeacherCommentsByReportSettingNumber(reportSettingNumber string) ([]models.HeadTeacherComment, error) {
	var headTeacherComment []models.HeadTeacherComment
	if err := config.DB.Preload("ReportSetting").
		Preload("ReportSetting.AcademicYear").
		Preload("ReportSetting.StudyPeriod").
		Preload("ReportSetting.CampusClass").
		Where("report_setting_number = ?", reportSettingNumber).
		Find(&headTeacherComment).Error; err != nil {
		return headTeacherComment, err
	}
	return headTeacherComment, nil
}
