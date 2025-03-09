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
}

type AttendanceRepositoryImpl struct {
	DB *gorm.DB
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
	if err := a.DB.Scopes(pagination.Paginate(context, p), filters.ByAnyField("user_number", context.Query("user_number"))).Find(&attendances).Error; err != nil {
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
