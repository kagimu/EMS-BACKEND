package repository

import (
	applicationModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/application/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/pagination"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/previousSchool/models"
	studentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ApplicationRepository interface {
	FindApplicationByNumber(applicationNumber string) (*applicationModels.Application, error)
	FindAllApplications(schoolNumber string) ([]applicationModels.Application, error)
	FindAllApplicationsPaginated(ctx *gin.Context, pagination *pagination.PaginatedResponse, schoolNumber string) ([]applicationModels.Application, error)
	FindPreviousSchoolsByApplicationNumber(applicationNumber string) ([]models.PreviousSchool, error)
	GetStudentByApplicationNumber(applicationNumber string) (studentModel.Student, error)
}

type ApplicationRepositoryImpl struct {
	DB *gorm.DB
}

func (a ApplicationRepositoryImpl) FindPreviousSchoolsByApplicationNumber(applicationNumber string) ([]models.PreviousSchool, error) {
	previousSchools := make([]models.PreviousSchool, 0)
	err := a.DB.Preload("Classes").Where("application_number = ?", applicationNumber).Find(&previousSchools).Error
	if err != nil {
		return nil, err
	}
	return previousSchools, nil
}

func (a ApplicationRepositoryImpl) GetStudentByApplicationNumber(applicationNumber string) (studentModel.Student, error) {
	var student studentModel.Student
	err := a.DB.Preload("PreferredCampus").Preload("CampusClass").Where("application_number = ?", applicationNumber).First(&student).Error
	if err != nil {
		return student, err
	}
	return student, nil
}

func (a ApplicationRepositoryImpl) FindApplicationByNumber(applicationNumber string) (*applicationModels.Application, error) {
	application := &applicationModels.Application{}
	err := a.DB.Preload("Parent").Where("number = ?", applicationNumber).First(application).Error
	if err != nil {
		return nil, err
	}
	return application, nil
}

func (a ApplicationRepositoryImpl) FindAllApplications(schoolNumber string) ([]applicationModels.Application, error) {
	applications := make([]applicationModels.Application, 0)
	err := a.DB.Where("school_number = ?", schoolNumber).Find(&applications).Error
	if err != nil {
		return nil, err

	}
	return applications, nil
}

func (a ApplicationRepositoryImpl) FindAllApplicationsPaginated(ctx *gin.Context, paginationResponse *pagination.PaginatedResponse, schoolNumber string) ([]applicationModels.Application, error) {
	applications := make([]applicationModels.Application, 0)

	err := a.DB.Scopes(pagination.Paginate(ctx, paginationResponse)).Where("school_number = ?", schoolNumber).Find(&applications).Error
	if err != nil {
		return applications, err
	}

	return applications, nil
}

func NewApplicationRepositoryImpl(db *gorm.DB) ApplicationRepository {
	return &ApplicationRepositoryImpl{DB: db}
}
