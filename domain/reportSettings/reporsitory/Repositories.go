package reporsitory

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/models"
)

func FindHeadTeacherCommentByReportSettingNumber(reportSettingNumber string) (models.HeadTeacherComment, error) {
	var headTeacherComment models.HeadTeacherComment
	err := config.DB.Where("report_setting_number = ?", reportSettingNumber).First(&headTeacherComment).Error
	if err != nil {
		return headTeacherComment, err
	}
	return headTeacherComment, nil
}
