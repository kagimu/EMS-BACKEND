package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	assessmenModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/interview/models"
	"gorm.io/gorm"
)

func FindAssessmentMetricByNumber(assessmentMetricNumber string) (assessmenModels.AssessmentMetric, error) {

	db := config.DB

	var assessmentMetric assessmenModels.AssessmentMetric

	if err := db.Where("number = ?", assessmentMetricNumber).First(&assessmentMetric).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return assessmentMetric, nil
		}

		return assessmentMetric, err
	}

	return assessmentMetric, nil
}

func FetchAllAssessmentMetrics() ([]assessmenModels.AssessmentMetric, error) {

	db := config.DB

	var assessmentMetrics []assessmenModels.AssessmentMetric

	if err := db.Find(&assessmentMetrics).Error; err != nil {
		return assessmentMetrics, err
	}

	return assessmentMetrics, nil
}

func FindAssessmentMetricsByInterviewNumber(interviewNumber string) ([]assessmenModels.AssessmentMetric, error) {

	db := config.DB

	var assessmentMetrics = make([]assessmenModels.AssessmentMetric, 0)

	err := db.Joins("JOIN interview_assessment_metrics ON interview_assessment_metrics.assessment_metric_number = assessment_metrics.number AND interview_assessment_metrics.interview_number = ? ", interviewNumber).Find(&assessmentMetrics).Error
	return assessmentMetrics, err
}

func FindInterviewStudentsByInterViewNumber(interviewNumber string) ([]models.InterviewStudent, error) {

	var interviewStudents = make([]models.InterviewStudent, 0)
	err := config.DB.Preload("Student").Preload("Interview").Where("interview_number = ?", interviewNumber).Find(&interviewStudents).Error
	return interviewStudents, err

}
