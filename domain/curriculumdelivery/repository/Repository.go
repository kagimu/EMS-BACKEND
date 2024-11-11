package repository

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/filters"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/pagination"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculumdelivery/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CurriculumDeliveryRepository interface {
	SaveSession(delivery *models.CurriculumDelivery) error
	SaveStudentSession(delivery *models.StudentCurriculumDelivery) error
	FindCurriculumDeliveryByStreamNumberAndStepNumber(streamNumber string, stepNumber string) ([]*models.CurriculumDelivery, error)
	FindCurriculumDeliveryByStudentNumberAndStepNumber(studentNumber string, stepNumber string) ([]*models.CurriculumDelivery, error)
	FindCurriculumDeliveryByFilters(ctx *gin.Context, pagination *pagination.PaginatedResponse) ([]*models.CurriculumDelivery, error)
}

type CurriculumDeliveryRepositoryImpl struct {
	DB *gorm.DB
}

func (c CurriculumDeliveryRepositoryImpl) FindCurriculumDeliveryByFilters(ctx *gin.Context, pr *pagination.PaginatedResponse) ([]*models.CurriculumDelivery, error) {
	deliveries := make([]*models.CurriculumDelivery, 0)
	err := c.DB.Scopes(
		pagination.Paginate(ctx, pr),
		filters.ByAnyField("stream_number", ctx.Query("stream_number")),
		filters.ByAnyField("lesson_plan_step_number", ctx.Query("lesson_plan_step_number"))).Find(&deliveries).Error
	return deliveries, err
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
