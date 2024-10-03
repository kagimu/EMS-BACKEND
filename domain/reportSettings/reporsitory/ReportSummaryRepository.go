package reporsitory

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	models4 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/aggregateScale/models"
	models3 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/assessmentGrade/models"
	models2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/assessmentArea/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/api/responses"
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

func CreateReportSummary(summary models.ReportSummary, trn ...*gorm.DB) error {

	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	return db.Create(&summary).Error
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

func GetGradingStandardNumberByAssessmentAreaNumber(assessmentAreaNumber string) (gs models4.GradingStandard, err error) {
	db := config.DB
	var assessmentAreaGradingStandard models2.AssessmentAreaGradingStandard

	if err = db.Preload("GradingStandard").Where("assessment_area_number = ?", assessmentAreaNumber).First(&assessmentAreaGradingStandard).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return gs, nil
		}

		return gs, err
	}

	return assessmentAreaGradingStandard.GradingStandard, nil
}

func GetScaleNumberByGradingStandardNumberAndClassNumber(gradingStandardNumber, classNumber string) (string, error) {
	db := config.DB
	var gradingStandardScale models3.GradingStandardClassScale

	if err := db.Where("grading_standard_number = ? AND class_number = ?", gradingStandardNumber, classNumber).First(&gradingStandardScale).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil
		}

		return "", err
	}

	return gradingStandardScale.ScaleNumber, nil
}

func GetAppliedScaleByScaleNumber(scaleNumber string) ([]models4.AggregateScale, error) {
	var ags models4.AggregateScale
	if err := config.DB.Where("number = ?", scaleNumber).First(&ags).Error; err != nil {
		return nil, err
	}

	scales := make([]models4.AggregateScale, 0)
	if err := config.DB.Preload("Grade").Where("scale_number = ?", ags.ScaleNumber).Order("code asc").Find(&scales).Error; err != nil {
		return nil, err
	}

	return scales, nil

}

func BuildAppliedScaleResponse(scales []models4.AggregateScale) []responses.AppliedScale {

	response := make([]responses.AppliedScale, 0)
	for _, scale := range scales {
		response = append(response, responses.AppliedScale{
			Name:      scale.Grade.Name,
			ShortName: scale.Grade.ShortName,
			Min:       scale.MinValue,
			Max:       scale.MaxValue,
		})
	}

	return response
}

func GetReportSummaryByReportSettingNumber(number string) (models.ReportSummary, error) {
	var rs models.ReportSummary
	if err := config.DB.Where("report_setting_number = ?", number).First(&rs).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.ReportSummary{}, nil
		}
		return rs, err
	}

	return rs, nil

}
