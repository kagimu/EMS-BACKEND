package repository

import (
	"errors"
	applicationModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/application/models"
	userModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
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
	UpdateApplicationStatus(applicationNumber string, status string) error
	DeleteDraftApplication(applicationNumber string, txn *gorm.DB) error
	SaveApplication(application applicationModels.Application) error
	UpdateStudentByApplicationNumber(applicationNumber string, student studentModel.Student) error
	CreateStudent(student studentModel.Student) error
	CreateUser(user userModel.User) error
}

type applicationRepositoryImpl struct {
	DB *gorm.DB
}

func (applicationRepository *applicationRepositoryImpl) CreateUser(user userModel.User) error {
	return applicationRepository.DB.Create(&user).Error
}

func (applicationRepository *applicationRepositoryImpl) CreateStudent(student studentModel.Student) error {
	return applicationRepository.DB.Create(&student).Error
}

func (applicationRepository *applicationRepositoryImpl) SaveApplication(application applicationModels.Application) error {
	return applicationRepository.DB.Create(&application).Error
}

func (applicationRepository *applicationRepositoryImpl) UpdateStudentByApplicationNumber(applicationNumber string, student studentModel.Student) error {
	return applicationRepository.DB.Where("application_number = ?", applicationNumber).Updates(student).Error
}

func (applicationRepository *applicationRepositoryImpl) DeleteDraftApplication(applicationNumber string, txn *gorm.DB) error {
	return txn.Where("number = ? AND status = ? ", applicationNumber, constants.StatusDraft).Delete(&applicationModels.Application{}).Error
}

func (applicationRepository *applicationRepositoryImpl) UpdateApplicationStatus(applicationNumber string, status string) error {
	return applicationRepository.DB.Model(&applicationModels.Application{}).Where("number = ?", applicationNumber).Update("status", status).Error
}

func (applicationRepository *applicationRepositoryImpl) FindPreviousSchoolsByApplicationNumber(applicationNumber string) ([]models.PreviousSchool, error) {
	previousSchools := make([]models.PreviousSchool, 0)
	err := applicationRepository.DB.Preload("Classes").Where("application_number = ?", applicationNumber).Find(&previousSchools).Error
	if err != nil {
		return nil, err
	}
	return previousSchools, nil
}

func (applicationRepository *applicationRepositoryImpl) GetStudentByApplicationNumber(applicationNumber string) (studentModel.Student, error) {
	var student studentModel.Student
	err := applicationRepository.DB.Preload("PreferredCampus").Preload("CampusClass").Where("application_number = ?", applicationNumber).First(&student).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return student, nil
		}
		return student, err
	}
	return student, nil
}

func (applicationRepository *applicationRepositoryImpl) FindApplicationByNumber(applicationNumber string) (*applicationModels.Application, error) {
	application := &applicationModels.Application{}
	err := applicationRepository.DB.Preload("Parent").Where("number = ?", applicationNumber).First(application).Error
	if err != nil {
		return nil, err
	}
	return application, nil
}

func (applicationRepository *applicationRepositoryImpl) FindAllApplications(schoolNumber string) ([]applicationModels.Application, error) {
	applications := make([]applicationModels.Application, 0)
	err := applicationRepository.DB.Where("school_number = ?", schoolNumber).Find(&applications).Error
	if err != nil {
		return nil, err

	}
	return applications, nil
}

func (applicationRepository *applicationRepositoryImpl) FindAllApplicationsPaginated(ctx *gin.Context, paginationResponse *pagination.PaginatedResponse, schoolNumber string) ([]applicationModels.Application, error) {
	applications := make([]applicationModels.Application, 0)

	err := applicationRepository.DB.Scopes(pagination.Paginate(ctx, paginationResponse)).Where("school_number = ?", schoolNumber).Find(&applications).Error
	if err != nil {
		return applications, err
	}

	return applications, nil
}

func NewApplicationRepositoryImpl(db *gorm.DB) ApplicationRepository {
	return &applicationRepositoryImpl{DB: db}
}
