package reporsitory

import (
	"errors"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/reportSettings/models"
	"gorm.io/gorm"
)

type ReportScoreSummaryRepository interface {
	CreateReportScoreSummary(reportScoreSummary models.ReportScoreSummary) error
	SaveReportScoreSummary(reportScoreSummary models.ReportScoreSummary) error
	UpdateReportScoreSummary(reportScoreSummary models.ReportScoreSummary) error
	FindReportScoreSummaryByReportSummaryAndAssessmentArea(reportSummaryNumber, assessmentAreaNumber string) (models.ReportScoreSummary, error)
}

type ReportScoreSummaryRepositoryImpl struct {
	DB *gorm.DB
}

func (r ReportScoreSummaryRepositoryImpl) UpdateReportScoreSummary(reportScoreSummary models.ReportScoreSummary) error {
	return r.DB.Save(&reportScoreSummary).Error
}

func NewReportScoreSummaryRepository(db *gorm.DB) ReportScoreSummaryRepository {
	return ReportScoreSummaryRepositoryImpl{
		DB: db,
	}
}

func (r ReportScoreSummaryRepositoryImpl) CreateReportScoreSummary(reportScoreSummary models.ReportScoreSummary) error {

	if err := r.DB.Create(&reportScoreSummary).Error; err != nil {
		return err
	}

	return nil
}

func (r ReportScoreSummaryRepositoryImpl) SaveReportScoreSummary(reportScoreSummary models.ReportScoreSummary) error {

	if err := r.DB.Save(&reportScoreSummary).Error; err != nil {
		return err
	}

	return nil
}

func (r ReportScoreSummaryRepositoryImpl) FindReportScoreSummaryByReportSummaryAndAssessmentArea(reportSummaryNumber, assessmentAreaNumber string) (models.ReportScoreSummary, error) {

	var reportScoreSummary models.ReportScoreSummary
	if err := r.DB.Where("report_summary_number = ? AND assessment_area_number =? ", reportSummaryNumber, assessmentAreaNumber).First(&reportScoreSummary).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return reportScoreSummary, nil
		}

		return reportScoreSummary, err
	}

	return reportScoreSummary, nil
}
