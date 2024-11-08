package repository

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculumdelivery/models"
	"gorm.io/gorm"
)

type CurriculumDeliveryRepository interface {
	SaveSession(delivery *models.CurriculumDelivery) error
	SaveStudentSession(delivery *models.StudentCurriculumDelivery) error
	FindCurriculumDeliveryByStreamNumberAndStepNumber(streamNumber string, stepNumber string) ([]*models.CurriculumDelivery, error)
	FindCurriculumDeliveryByStudentNumberAndStepNumber(studentNumber string, stepNumber string) ([]*models.CurriculumDelivery, error)
}

type CurriculumDeliveryRepositoryImpl struct {
	DB *gorm.DB
}

func (c CurriculumDeliveryRepositoryImpl) FindCurriculumDeliveryByStreamNumberAndStepNumber(streamNumber string, stepNumber string) ([]*models.CurriculumDelivery, error) {
	deliveries := make([]*models.CurriculumDelivery, 0)
	err := c.DB.Where("stream_number = ? AND lesson_plan_step_number = ?", streamNumber, stepNumber).Find(&deliveries).Error
	return deliveries, err
}

func (c CurriculumDeliveryRepositoryImpl) FindCurriculumDeliveryByStudentNumberAndStepNumber(studentNumber string, stepNumber string) ([]*models.CurriculumDelivery, error) {
	deliveries := make([]*models.CurriculumDelivery, 0)
	err := c.DB.Table("curriculum_deliveries").
		Joins("JOIN student_curriculum_deliveries ON curriculum_deliveries.number = student_curriculum_deliveries.curriculum_delivery_number").
		Where("student_curriculum_deliveries.student_number = ? AND curriculum_deliveries.lesson_plan_step_number = ?", studentNumber, stepNumber).
		Find(&deliveries).Error
	return deliveries, err
}

func (c CurriculumDeliveryRepositoryImpl) SaveStudentSession(delivery *models.StudentCurriculumDelivery) error {
	return c.DB.Create(delivery).Error
}

func (c CurriculumDeliveryRepositoryImpl) SaveSession(delivery *models.CurriculumDelivery) error {
	return c.DB.Create(delivery).Error
}

func NewCurriculumDeliveryRepositoryImpl(db *gorm.DB) CurriculumDeliveryRepository {
	return &CurriculumDeliveryRepositoryImpl{DB: db}
}
