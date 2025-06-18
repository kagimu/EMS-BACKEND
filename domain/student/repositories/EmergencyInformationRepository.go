package repositories

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/models"
	"gorm.io/gorm"
)

type EmergencyInformationRepository interface {
	FindAll() ([]models.EmergencyInformation, error)
	GetEmergencyInformationByStudentNumber(studentNumber string) (models.EmergencyInformation, error)
	CreateEmergencyInformation(emergencyInformation models.EmergencyInformation) error
	DeleteEmergencyInformation(emergencyInformation models.EmergencyInformation) error
	DeleteEmergencyInformationByStudentNumber(studentNumber string) error
}

type EmergencyInformationRepositoryImpl struct {
	DB *gorm.DB
}

func (e EmergencyInformationRepositoryImpl) DeleteEmergencyInformationByStudentNumber(studentNumber string) error {
	var emergencyInformation models.EmergencyInformation
	err := e.DB.Where("student_number = ?", studentNumber).Delete(&emergencyInformation).Error
	if err != nil {
		return err
	}
	return nil
}

func (e EmergencyInformationRepositoryImpl) DeleteEmergencyInformation(emergencyInformation models.EmergencyInformation) error {
	return e.DB.Delete(&emergencyInformation).Error
}

func (e EmergencyInformationRepositoryImpl) FindAll() ([]models.EmergencyInformation, error) {
	var emergencyInformation []models.EmergencyInformation

	if err := e.DB.Find(&emergencyInformation).Error; err != nil {
		return nil, err
	}

	return emergencyInformation, nil
}

func (e EmergencyInformationRepositoryImpl) GetEmergencyInformationByStudentNumber(studentNumber string) (models.EmergencyInformation, error) {
	var emergencyInformation models.EmergencyInformation

	if err := e.DB.Where("student_number = ?", studentNumber).First(&emergencyInformation).Error; err != nil {
		return models.EmergencyInformation{}, err
	}

	return emergencyInformation, nil
}

func (e EmergencyInformationRepositoryImpl) CreateEmergencyInformation(emergencyInformation models.EmergencyInformation) error {
	return e.DB.Create(&emergencyInformation).Error
}

func NewEmergencyInformationRepository(db *gorm.DB) EmergencyInformationRepository {
	return &EmergencyInformationRepositoryImpl{
		DB: db,
	}
}
