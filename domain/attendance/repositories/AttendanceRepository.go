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
	var accountType = context.Query("type")
	if err := a.DB.Scopes(pagination.Paginate(context, p),
		filters.ByAnyField("attendances.user_number", context.Query("user_number")),
		filters.ByAnyField("attendances.campus_class_number", context.Query("campus_class_number")),
		filters.ByAnyField("attendances.campus_number", context.Query("campus_number")),
		filters.ByAnyField("attendances.school_number", context.Query("school_number")),
		filters.ByAnyField("attendances.stream_number", context.Query("stream_number")),
		filters.ByAnyField("attendances.campus_number", context.Query("campus_number")),
		filters.ByAnyField("attendances.student_number", context.Query("student_number")),
		filters.ByAnyField("attendances.study_period_number", context.Query("study_period_number")),
		filters.ByAnyField("attendances.date", context.Query("date")),
		filters.ByUserNamesLike(context.Query("search")),
	).Joins("JOIN users ON users.number = attendances.user_number").Where("users.type =?", accountType).Find(&attendances).Error; err != nil {
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
