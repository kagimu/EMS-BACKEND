package reporsitory

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/models"
	"gorm.io/gorm"
)

func GetHeadTeacherCommentByReportSettingNumber(reportSettingNumber string) (models.HeadTeacherComment, error) {
	db := config.DB
	var reportComment models.HeadTeacherComment
	if err := db.Where("report_setting_number = ?", reportSettingNumber).First(&reportComment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return reportComment, nil
		}
		return reportComment, err
	}
	return reportComment, nil
}
