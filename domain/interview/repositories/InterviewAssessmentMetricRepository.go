package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/interview/models"
	"gorm.io/gorm"
)

func FindInterviewAssessmentMetricByInterviewNumberAndCampusClassNumber(interviewNumber, campusClassNumber string, trn ...*gorm.DB) ([]models.InterviewAssessmentMetric, error) {

	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	var interviewAssessmentMetrics []models.InterviewAssessmentMetric

	if err := db.
		Joins("join assessment_metrics on assessment_metrics.number = interview_assessment_metrics.assessment_metric_number AND assessment_metrics.campus_class_number = ? ", campusClassNumber).
		Where("interview_number = ? ", interviewNumber).Find(&interviewAssessmentMetrics).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return interviewAssessmentMetrics, nil
		}

		return interviewAssessmentMetrics, err
	}

	return interviewAssessmentMetrics, nil
}
