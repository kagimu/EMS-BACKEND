package reporsitory

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	models4 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/aggregateScale/models"
	models3 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/assessmentGrade/models"
	models2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/assessmentArea/models"
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
