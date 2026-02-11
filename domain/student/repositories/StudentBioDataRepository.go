package repositories

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/models"
	"gorm.io/gorm"
)

type StudentBioDataRepository interface {
	FindAll() ([]models.StudentBioData, error)
	GetStudentBioDataByStudentNumber(studentNumber string) ([]models.StudentBioData, error)
	GetStudentBioDataByNumber(number string) (models.StudentBioData, error)
	CreateStudentBioData(biData models.StudentBioData) error
	DeleteStudentBioData(bioData models.StudentBioData) error
	DeleteStudentBioDataByNumber(number string) error
	UpdateStudentBioData(bioData models.StudentBioData) error
}

type StudentBioDataRepositoryImpl struct {
	DB *gorm.DB
}

func (s StudentBioDataRepositoryImpl) DeleteStudentBioDataByNumber(studentNumber string) error {
	var bioData models.StudentBioData
	if err := s.DB.Where("number = ?", studentNumber).Delete(&bioData).Error; err != nil {
		return err
	}
	return nil
}

func (s StudentBioDataRepositoryImpl) UpdateStudentBioData(bioData models.StudentBioData) error {
	return s.DB.Save(&bioData).Error
}

func (s StudentBioDataRepositoryImpl) FindAll() ([]models.StudentBioData, error) {
	var bioDataList []models.StudentBioData

	if err := s.DB.Find(&bioDataList).Error; err != nil {
		return nil, err
	}

	return bioDataList, nil
}

func (s StudentBioDataRepositoryImpl) GetStudentBioDataByStudentNumber(studentNumber string) ([]models.StudentBioData, error) {
	var bioData = make([]models.StudentBioData, 0)

	if err := s.DB.Where("student_number = ?", studentNumber).Find(&bioData).Error; err != nil {
		return bioData, err
	}

	return bioData, nil
}

func (s StudentBioDataRepositoryImpl) GetStudentBioDataByNumber(number string) (models.StudentBioData, error) {
	var bioData models.StudentBioData

	if err := s.DB.Where("number = ?", number).First(&bioData).Error; err != nil {
		return models.StudentBioData{}, err
	}

	return bioData, nil
}

func (s StudentBioDataRepositoryImpl) CreateStudentBioData(bioData models.StudentBioData) error {
	return s.DB.Create(&bioData).Error
}

func (s StudentBioDataRepositoryImpl) DeleteStudentBioData(bioData models.StudentBioData) error {
	return s.DB.Delete(&bioData).Error
}

func NewStudentBioDataRepository(db *gorm.DB) StudentBioDataRepository {
	return &StudentBioDataRepositoryImpl{
		DB: db,
	}
}
