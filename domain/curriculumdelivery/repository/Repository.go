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
	FindStudentsBySession(context *gin.Context, p *pagination.PaginatedResponse) ([]*models.StudentCurriculumDelivery, error)
	DeleteStudentsBySession(sessionNumber string) error
	DeleteCurriculumDeliverySession(number string) error
	DeleteCurriculumDeliveryStep(number string) error
	UpdateCurriculumDeliverySession(delivery *models.CurriculumDelivery) error
	SaveCurriculumDeliveryStep(m *models.CurriculumDeliveryStep) error
	GetCurriculumDeliveryByNumber(number string) (*models.CurriculumDelivery, error)
}

type CurriculumDeliveryRepositoryImpl struct {
	DB *gorm.DB
}

func (c CurriculumDeliveryRepositoryImpl) GetCurriculumDeliveryByNumber(number string) (*models.CurriculumDelivery, error) {
	delivery := &models.CurriculumDelivery{}
	err := c.DB.Where("number = ?", number).First(delivery).Error
	return delivery, err
}

func (c CurriculumDeliveryRepositoryImpl) DeleteCurriculumDeliveryStep(number string) error {
	return c.DB.Where("curriculum_delivery_number = ?", number).Delete(&models.CurriculumDeliveryStep{}).Error
}

func (c CurriculumDeliveryRepositoryImpl) SaveCurriculumDeliveryStep(m *models.CurriculumDeliveryStep) error {
	return c.DB.Create(m).Error
}

func (c CurriculumDeliveryRepositoryImpl) UpdateCurriculumDeliverySession(delivery *models.CurriculumDelivery) error {
	return c.DB.Save(delivery).Error
}

func (c CurriculumDeliveryRepositoryImpl) DeleteCurriculumDeliverySession(number string) error {
	return c.DB.Where("number = ?", number).Delete(&models.CurriculumDelivery{}).Error
}

func (c CurriculumDeliveryRepositoryImpl) DeleteStudentsBySession(sessionNumber string) error {
	return c.DB.Where("curriculum_delivery_number = ?", sessionNumber).Delete(&models.StudentCurriculumDelivery{}).Error
}

func (c CurriculumDeliveryRepositoryImpl) FindStudentsBySession(context *gin.Context, p *pagination.PaginatedResponse) ([]*models.StudentCurriculumDelivery, error) {
	students := make([]*models.StudentCurriculumDelivery, 0)
	err := c.DB.Where("curriculum_delivery_number = ? ", context.Param("curriculum_delivery_number")).Scopes(pagination.Paginate(context, p)).Find(&students).Error
	return students, err
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
	err := c.DB.Joins("JOIN curriculum_delivery_steps ON curriculum_deliveries.number = curriculum_delivery_steps.curriculum_delivery_number ").Where("stream_number = ? AND curriculum_delivery_steps.lesson_plan_step_number = ?", streamNumber, stepNumber).Find(&deliveries).Error
	return deliveries, err
}

func (c CurriculumDeliveryRepositoryImpl) FindCurriculumDeliveryByStudentNumberAndStepNumber(studentNumber string, stepNumber string) ([]*models.CurriculumDelivery, error) {
	deliveries := make([]*models.CurriculumDelivery, 0)
	err := c.DB.Table("curriculum_deliveries").
		Joins("JOIN student_curriculum_deliveries ON curriculum_deliveries.number = student_curriculum_deliveries.curriculum_delivery_number").
		Joins("JOIN curriculum_delivery_steps ON curriculum_deliveries.number = curriculum_delivery_steps.curriculum_delivery_number").
		Where("student_curriculum_deliveries.student_number = ? AND curriculum_delivery_steps.lesson_plan_step_number = ?", studentNumber, stepNumber).
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
