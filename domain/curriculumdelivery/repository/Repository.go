package repository

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculumdelivery/models"
	"gorm.io/gorm"
)

type CurriculumDeliveryRepository interface {
	SaveSession(delivery *models.CurriculumDelivery) error
	SaveStudentSession(delivery *models.StudentCurriculumDelivery) error
}

type CurriculumDeliveryRepositoryImpl struct {
	DB *gorm.DB
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
