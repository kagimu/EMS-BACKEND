package repositories

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/attendance/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/filters"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/pagination"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AttendanceRepository interface {
	Save(attendance *models.Attendance) (*models.Attendance, error)
	FindByUserNumber(userNumber string) ([]models.Attendance, error)
	AllAttendance(*gin.Context, *pagination.PaginatedResponse) ([]models.Attendance, error)
	Update(attendance *models.Attendance) (*models.Attendance, error)
	GetAttendanceByNumber(number string) (models.Attendance, error)
}

type AttendanceRepositoryImpl struct {
	DB *gorm.DB
}

func (a *AttendanceRepositoryImpl) GetAttendanceByNumber(number string) (models.Attendance, error) {
	var attendance models.Attendance
	err := a.DB.Where("number = ?", number).First(&attendance).Error
	if err != nil {
		return attendance, err
	}

	return attendance, nil
}

func (a *AttendanceRepositoryImpl) Update(attendance *models.Attendance) (*models.Attendance, error) {
	if err := a.DB.Save(&attendance).Error; err != nil {
		return nil, err
	}
	return attendance, nil
}

func (a *AttendanceRepositoryImpl) FindByUserNumber(userNumber string) ([]models.Attendance, error) {
	var attendances []models.Attendance = make([]models.Attendance, 0)
	if err := a.DB.Where("user_number = ?", userNumber).Find(&attendances).Error; err != nil {
		return attendances, err
	}
	return attendances, nil
}

func (a *AttendanceRepositoryImpl) AllAttendance(context *gin.Context, p *pagination.PaginatedResponse) ([]models.Attendance, error) {
	var attendances = make([]models.Attendance, 0)
	if err := a.DB.Scopes(pagination.Paginate(context, p),
		filters.ByAnyField("user_number", context.Query("user_number")),
		filters.ByAnyField("campus_class_number", context.Query("campus_class_number")),
		filters.ByAnyField("campus_number", context.Query("campus_number")),
		filters.ByAnyField("school_number", context.Query("school_number")),
		filters.ByAnyField("stream_number", context.Query("stream_number")),
		filters.ByAnyField("campus_number", context.Query("campus_number")),
		filters.ByAnyField("student_number", context.Query("student_number")),
		filters.ByAnyField("study_period_number", context.Query("study_period_number")),
	).Find(&attendances).Error; err != nil {
		return attendances, err
	}
	return attendances, nil
}

func NewAttendanceRepository(db *gorm.DB) AttendanceRepository {
	return &AttendanceRepositoryImpl{DB: db}
}

func (a *AttendanceRepositoryImpl) Save(attendance *models.Attendance) (*models.Attendance, error) {
	if err := a.DB.Create(&attendance).Error; err != nil {
		return nil, err
	}
	return attendance, nil
}
