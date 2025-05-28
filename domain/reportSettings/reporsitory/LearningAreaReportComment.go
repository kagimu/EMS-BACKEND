package reporsitory

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/models"
	"gorm.io/gorm"
)

func SaveLearningAreaReportComment(learningAreaReportComment models.LearningAreaReportComment, trn ...*gorm.DB) error {
	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	return db.Save(&learningAreaReportComment).Error
}

func FindLearningAreaReportComment(reportSummaryNumber, learningAreaNumber string, trn ...*gorm.DB) (models.LearningAreaReportComment, error) {
	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	var learningAreaReportComment models.LearningAreaReportComment

	if err := db.Where("report_summary_number = ? AND learning_area_number = ?", reportSummaryNumber, learningAreaNumber).First(&learningAreaReportComment).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return learningAreaReportComment, nil
		}

		return learningAreaReportComment, err
	}

	return learningAreaReportComment, nil
}
