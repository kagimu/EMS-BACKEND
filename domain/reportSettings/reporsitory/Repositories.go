package reporsitory

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	models3 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/aggregateScale/models"
	models2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/assessmentArea/models"
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

func GetAssessmentAreaGradingStandard(assessmentAreaNumber string) (models2.AssessmentAreaGradingStandard, error) {
	var aags models2.AssessmentAreaGradingStandard
	if err := config.DB.Where("assessment_area_number = ?", assessmentAreaNumber).First(&aags).Error; err != nil {
		return aags, err
	}
	return aags, nil

}

func GetAggregateScaleByGradingStandardNumberAndCampusClassNumber(gradingStandardNumber, campusClassNumber string) ([]models3.AggregateScale, error) {
	var aggregateScales []models3.AggregateScale
	if err := config.DB.Where("grading_standard_number = ? AND campus_class_number = ?", gradingStandardNumber, campusClassNumber).Find(&aggregateScales).Error; err != nil {
		return aggregateScales, err
	}

	//
	//if err := config.DB.Where("assessment_area_number = ?", assessmentAreaNumber).First(&aags).Error; err != nil {
	//}
	//if err := config.DB.Where("grading_standard_number = ?", aags.GradingStandardNumber).Find(&aggregateScales).Error; err != nil {
	//}

	return aggregateScales, nil
}
