package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	assessmenModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/assessment/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/filters"
	"gorm.io/gorm"
)

type AssessmentMetricRepository interface {
	FindAssessmentMetricByNumber(assessmentMetricNumber string, tx ...*gorm.DB) (assessmenModels.AssessmentMetric, error)
	FindAssessmentMetricsByInterviewNumber(interviewNumber string, tx ...*gorm.DB) ([]assessmenModels.AssessmentMetric, error)
	AllAssessmentMetricPaginated(assessmentAreaNumber, studyPeriodNumber, campusClassNumber, metricNumber, classNumber, campusNumber string, paginated func(db *gorm.DB) *gorm.DB) ([]assessmenModels.AssessmentMetric, error)
	AllAssessmentMetric(assessmentAreaNumber, studyPeriodNumber, campusClassNumber, metricNumber, classNumber, campusNumber string, tx ...*gorm.DB) ([]assessmenModels.AssessmentMetric, error)
	SaveAssessmentMetric(assessmentMetric assessmenModels.AssessmentMetric, tx ...*gorm.DB) error
}

type repository struct {
}

func (a *repository) AllAssessmentMetricPaginated(assessmentAreaNumber, studyPeriodNumber, campusClassNumber, metricNumber, classNumber, campusNumber string, paginated func(db *gorm.DB) *gorm.DB) ([]assessmenModels.AssessmentMetric, error) {
	db := config.DB

	var assessmentMetrics []assessmenModels.AssessmentMetric

	if err := db.
		Joins("JOIN campus_classes ON campus_classes.number = assessment_metrics.campus_class_number").
		Joins("JOIN curriculum_classes ON curriculum_classes.number = campus_classes.curriculum_class_number").
		Scopes(
			filters.ByAnyField("campus_classes.campus_number", campusNumber),
			filters.ByAnyField("curriculum_classes.class_number", classNumber),
			filters.ByStudyPeriodNumber(studyPeriodNumber),
			filters.ByAssessmentAreaNumber(assessmentAreaNumber),
			filters.ByAnyField("assessment_metrics.campus_class_number", campusClassNumber),
			filters.ByMetricNumber(metricNumber),
			paginated,
		).
		Distinct("assessment_metrics.*").
		Find(&assessmentMetrics).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return assessmentMetrics, nil
		}

		return assessmentMetrics, err

	}

	return assessmentMetrics, nil
}
func (a *repository) AllAssessmentMetric(assessmentAreaNumber, studyPeriodNumber, campusClassNumber, metricNumber, classNumber, campusNumber string, tx ...*gorm.DB) ([]assessmenModels.AssessmentMetric, error) {
	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}

	var assessmentMetrics []assessmenModels.AssessmentMetric

	if err := db.
		Joins("JOIN campus_classes ON campus_classes.number = assessment_metrics.campus_class_number").
		Joins("JOIN curriculum_classes ON curriculum_classes.number = campus_classes.curriculum_class_number").
		Scopes(
			filters.ByAnyField("campus_classes.campus_number", campusNumber),
			filters.ByAnyField("curriculum_classes.class_number", classNumber),
			filters.ByStudyPeriodNumber(studyPeriodNumber),
			filters.ByAssessmentAreaNumber(assessmentAreaNumber),
			filters.ByAnyField("assessment_metrics.campus_class_number", campusClassNumber),
			filters.ByMetricNumber(metricNumber),
		).
		Distinct("assessment_metrics.*").
		Find(&assessmentMetrics).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return assessmentMetrics, nil
		}

		return assessmentMetrics, err

	}

	return assessmentMetrics, nil
}

func (a *repository) FindAssessmentMetricByNumber(assessmentMetricNumber string, tx ...*gorm.DB) (assessmenModels.AssessmentMetric, error) {

	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}

	var assessmentMetric assessmenModels.AssessmentMetric

	if err := db.Where("number = ?", assessmentMetricNumber).First(&assessmentMetric).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return assessmentMetric, nil
		}

		return assessmentMetric, err
	}

	return assessmentMetric, nil
}

func (a *repository) FindAssessmentMetricsByInterviewNumber(interviewNumber string, tx ...*gorm.DB) ([]assessmenModels.AssessmentMetric, error) {
	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}

	var assessmentMetrics = make([]assessmenModels.AssessmentMetric, 0)

	err := db.Joins("JOIN interview_assessment_metrics ON interview_assessment_metrics.assessment_metric_number = assessment_metrics.number AND interview_assessment_metrics.interview_number = ? ", interviewNumber).Find(&assessmentMetrics).Error
	return assessmentMetrics, err
}

func (a *repository) SaveAssessmentMetric(assessmentMetric assessmenModels.AssessmentMetric, tx ...*gorm.DB) error {
	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}

	return db.Save(&assessmentMetric).Error
}

func NewAssessmentMetricRepository() AssessmentMetricRepository {
	return &repository{}
}
