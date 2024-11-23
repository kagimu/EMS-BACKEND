package repository

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	conceptModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/concept/models"
	models2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/importCurriculum/models"
	subTopicModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/subtopic/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/topic/models"
	currDelivery "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculumdelivery/models"
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
	FindConceptStepByStep(step, lessonPlanConceptNumber string) (lessonPlanModel.LessonPlanStep, error)
	GetLessonPlanStepByNumber(lessonPlanStepNumber string) (lessonPlanModel.LessonPlanStep, error)
	GetConceptByNumber(conceptNumber string) (conceptModel.Concept, error)
	GetLessonPlanConceptByNumber(conceptNumber string) (lessonPlanModel.LessonPlanConcept, error)
	GetSubTopicByNumber(subTopicNumber string) (subTopicModel.SubTopic, error)
	GetTopicByNumber(topicNumber string) (models.Topic, error)
	FindCurriculumDeliveryStepsByCurriculumDeliveryNumber(number string) ([]currDelivery.CurriculumDeliveryStep, error)
	CountCurriculumsDeliveryStepsByStreamNumber(streamNumber string) (int64, error)
	CountStepsByCurriculumNumber(number string) (int64, error)
	GetCurriculumByStreamNumber(lessonPlanStepNumber string) (models2.Curriculum, error)
}

type lessonPlanRepositoryImpl struct {
	Db *gorm.DB
}

func (l *lessonPlanRepositoryImpl) CountCurriculumsDeliveryStepsByStreamNumber(streamNumber string) (int64, error) {
	var count int64
	err := l.Db.Model(&currDelivery.CurriculumDelivery{}).Where("stream_number = ?", streamNumber).Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (l *lessonPlanRepositoryImpl) CountStepsByCurriculumNumber(number string) (int64, error) {
	var count int64
	err := l.Db.Model(&lessonPlanModel.LessonPlanStep{}).
		Joins("JOIN lesson_plan_concepts ON lesson_plan_steps.lesson_plan_concept_number = lesson_plan_concepts.number ").
		Joins("JOIN concepts ON lesson_plan_concepts.concept_number = concepts.number AND concepts.curriculum_number = ?", number).Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (l *lessonPlanRepositoryImpl) GetCurriculumByStreamNumber(StreamNumber string) (models2.Curriculum, error) {
	var concept conceptModel.Concept
	var curriculum models2.Curriculum
	var LessonPlanStep lessonPlanModel.LessonPlanStep
	var cd currDelivery.CurriculumDelivery
	err := l.Db.Joins("JOIN curriculum_delivery_steps ON curriculum_delivery_steps.curriculum_delivery_number = curriculum_deliveries.number AND curriculum_deliveries.stream_number = ?", StreamNumber).First(&cd).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return curriculum, nil
		}
		return curriculum, err
	}

	err = l.Db.Joins("JOIN lesson_plan_concepts ON lesson_plan_steps.lesson_plan_concept_number = lesson_plan_concepts.number").Where("lesson_plan_steps.number = ?", LessonPlanStep.Number).Joins("JOIN concepts ON lesson_plan_concepts.concept_number = concepts.number").First(&concept).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return curriculum, nil
		}

		return curriculum, err
	}

	err = l.Db.Where("number = ?", concept.CurriculumNumber).First(&curriculum).Error
	return curriculum, nil

}

func (l *lessonPlanRepositoryImpl) FindCurriculumDeliveryStepsByCurriculumDeliveryNumber(number string) ([]currDelivery.CurriculumDeliveryStep, error) {
	steps := make([]currDelivery.CurriculumDeliveryStep, 0)
	err := l.Db.Where("curriculum_delivery_number = ?", number).Find(&steps).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {

			return steps, nil
		}

		return steps, err
	}

	return steps, err

}

func (l *lessonPlanRepositoryImpl) GetConceptByNumber(conceptNumber string) (conceptModel.Concept, error) {
	var concept conceptModel.Concept
	err := l.Db.Where("number = ?", conceptNumber).First(&concept).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return conceptModel.Concept{}, nil
		}
		return conceptModel.Concept{}, err
	}

	return concept, nil
}

func (l *lessonPlanRepositoryImpl) GetLessonPlanConceptByNumber(number string) (lessonPlanModel.LessonPlanConcept, error) {
	var lessonPlanConcept lessonPlanModel.LessonPlanConcept
	err := l.Db.Where("number = ?", number).First(&lessonPlanConcept).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return lessonPlanModel.LessonPlanConcept{}, nil
		}
		return lessonPlanModel.LessonPlanConcept{}, err
	}

	return lessonPlanConcept, nil
}

func (l *lessonPlanRepositoryImpl) GetSubTopicByNumber(subTopicNumber string) (subTopicModel.SubTopic, error) {
	var subTopic subTopicModel.SubTopic
	err := l.Db.Where("number = ?", subTopicNumber).First(&subTopic).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return subTopicModel.SubTopic{}, nil
		}
		return subTopicModel.SubTopic{}, err
	}

	return subTopic, nil
}

func (l *lessonPlanRepositoryImpl) GetTopicByNumber(topicNumber string) (models.Topic, error) {
	var topic models.Topic
	err := l.Db.Where("number = ?", topicNumber).First(&topic).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Topic{}, nil
		}
		return models.Topic{}, err
	}

	return topic, nil
}

func (l *lessonPlanRepositoryImpl) GetLessonPlanStepByNumber(lessonPlanStepNumber string) (lessonPlanModel.LessonPlanStep, error) {
	var lp lessonPlanModel.LessonPlanStep
	err := l.Db.Where("number = ?", lessonPlanStepNumber).First(&lp).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return lessonPlanModel.LessonPlanStep{}, nil
		}
		return lessonPlanModel.LessonPlanStep{}, err
	}
	return lp, nil
}

func (l *lessonPlanRepositoryImpl) FindConceptStepByStep(step, lessonPlanConceptNumber string) (lessonPlanModel.LessonPlanStep, error) {
	var lessonPlanStep lessonPlanModel.LessonPlanStep
	if err := l.Db.Where("step = ? AND lesson_plan_concept_number = ?", step, lessonPlanConceptNumber).First(&lessonPlanStep).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return lessonPlanModel.LessonPlanStep{}, nil
		}

		return lessonPlanModel.LessonPlanStep{}, err
	}

	return lessonPlanStep, nil
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
