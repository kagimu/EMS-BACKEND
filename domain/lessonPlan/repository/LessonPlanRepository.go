package repository

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	lessonPlanModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/lessonPlan/models"
	"gorm.io/gorm"
)

type LessonPlanRepository interface {
	UseTx(tx *gorm.DB) *gorm.DB
	GetLessonPlanBySchemeOfWorkSubTopicNumber(schemeOfWorkSubTopicNumber string) (lessonPlanModel.LessonPlan, error)
	SaveLessonPlan(lessonPlan lessonPlanModel.LessonPlan) error
	GetLessonPlanConcept(lessonPlanNumber, conceptNumber string) (lessonPlanModel.LessonPlanConcept, error)
	SaveLessonPlanConcept(lessonPlanConcept lessonPlanModel.LessonPlanConcept) error
	SaveConceptStep(lessonPlanStep lessonPlanModel.LessonPlanStep) error
	FindConceptStepByCode(code int64, lessonPlanConceptNumber string) (lessonPlanModel.LessonPlanStep, error)
}

type lessonPlanRepositoryImpl struct {
	Db *gorm.DB
}

func (l *lessonPlanRepositoryImpl) FindConceptStepByCode(code int64, lessonPlanConceptNumber string) (lessonPlanModel.LessonPlanStep, error) {

	var lessonPlanStep lessonPlanModel.LessonPlanStep
	if err := l.Db.Where("code = ? AND lesson_plan_concept_number = ?", code, lessonPlanConceptNumber).First(&lessonPlanStep).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return lessonPlanModel.LessonPlanStep{}, nil
		}

		return lessonPlanModel.LessonPlanStep{}, err
	}

	return lessonPlanStep, nil
}

func (l *lessonPlanRepositoryImpl) SaveLessonPlanConcept(lessonPlanConcept lessonPlanModel.LessonPlanConcept) error {
	return l.Db.Save(&lessonPlanConcept).Error
}

func (l *lessonPlanRepositoryImpl) SaveConceptStep(lessonPlanStep lessonPlanModel.LessonPlanStep) error {
	return l.Db.Save(&lessonPlanStep).Error
}

func (l *lessonPlanRepositoryImpl) GetLessonPlanConcept(lessonPlanNumber, conceptNumber string) (lessonPlanModel.LessonPlanConcept, error) {

	var lessonPlanConcept lessonPlanModel.LessonPlanConcept

	if err := l.Db.Where("lesson_plan_number = ? AND concept_number = ?", lessonPlanNumber, conceptNumber).Find(&lessonPlanConcept).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return lessonPlanModel.LessonPlanConcept{}, nil
		}
		return lessonPlanModel.LessonPlanConcept{}, err
	}

	return lessonPlanConcept, nil
}

func (l *lessonPlanRepositoryImpl) SaveLessonPlan(lessonPlan lessonPlanModel.LessonPlan) error {
	return l.Db.Save(&lessonPlan).Error
}

func (l *lessonPlanRepositoryImpl) GetLessonPlanBySchemeOfWorkSubTopicNumber(schemeSubTopicNumber string) (lessonPlanModel.LessonPlan, error) {

	var lessonPlan lessonPlanModel.LessonPlan

	if err := l.Db.Where("scheme_sub_topic_number = ?", schemeSubTopicNumber).Find(&lessonPlan).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return lessonPlanModel.LessonPlan{}, nil
		}

		return lessonPlanModel.LessonPlan{}, err
	}

	return lessonPlan, nil
}

func (l *lessonPlanRepositoryImpl) UseTx(tx *gorm.DB) *gorm.DB {
	l.Db = tx
	return l.Db
}

func NewLessonPlanRepository() LessonPlanRepository {
	return &lessonPlanRepositoryImpl{
		Db: config.DB,
	}
}
