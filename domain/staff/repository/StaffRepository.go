package repository

import (
	"errors"
	userModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/filters"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/pagination"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/staff/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type StaffRepository interface {
	AttachStaffToStream(streamTeacher *models.StreamTeacher) error
	DetachStaffFromStream(streamTeacher *models.StreamTeacher) error
	GetTeacherStreamDeployments(ctx *gin.Context, pagination *pagination.PaginatedResponse) ([]models.StreamTeacher, error)
	GetTeacherByNumberAndStreamNumber(userNumber string, streamNumber string) (*models.StreamTeacher, error)
	GetDefaultAssessmentAreaStreamTeacher(assessmentAreaNumber, streamNumber string) (userModel.User, error)
}

type StaffRepo struct {
	DB *gorm.DB
}

func (s StaffRepo) GetDefaultAssessmentAreaStreamTeacher(assessmentAreaNumber, streamNumber string) (userModel.User, error) {
	var user userModel.User

	err := s.DB.
		Joins("JOIN stream_teachers ON users.number = stream_teachers.user_number").
		Where("stream_teachers.assessment_area_number = ? AND stream_teachers.stream_number = ?", assessmentAreaNumber, streamNumber).
		First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user, nil
		}
		return user, err
	}

	return user, nil
}

func (s StaffRepo) GetTeacherByNumberAndStreamNumber(userNumber string, streamNumber string) (*models.StreamTeacher, error) {
	var streamTeacher models.StreamTeacher
	err := s.DB.Where("user_number = ? AND stream_number = ?", userNumber, streamNumber).First(&streamTeacher).Error
	if err != nil {
		return nil, err
	}

	return &streamTeacher, nil
}

func (s StaffRepo) GetTeacherStreamDeployments(ctx *gin.Context, pr *pagination.PaginatedResponse) ([]models.StreamTeacher, error) {
	var streamTeachers []models.StreamTeacher
	err := s.DB.Scopes(
		pagination.Paginate(ctx, pr),
		filters.ByAnyField("school_number", ctx.Query("school_number")),
		filters.ByAnyField("campus_number", ctx.Query("campus_number")),
		filters.ByAnyField("curriculum_number", ctx.Query("curriculum_number")),
		filters.ByAnyField("campus_class_number", ctx.Query("campus_class_number")),
		filters.ByAnyField("stream_number", ctx.Query("stream_number")),
		filters.ByAnyField("user_number", ctx.Query("user_number")),
		filters.ByAnyField("added_by", ctx.Query("added_by")),
		filters.ByAnyField("assessment_area_number", ctx.Query("assessment_area_number"))).Find(&streamTeachers).Error
	if err != nil {
		return nil, err
	}

	return streamTeachers, nil
}

func (s StaffRepo) AttachStaffToStream(streamTeacher *models.StreamTeacher) error {
	return s.DB.Create(&streamTeacher).Error
}

func (s StaffRepo) DetachStaffFromStream(streamTeacher *models.StreamTeacher) error {
	return s.DB.Delete(streamTeacher).Error
}

func NewStaffRepository(db *gorm.DB) StaffRepository {
	return &StaffRepo{DB: db}
}
