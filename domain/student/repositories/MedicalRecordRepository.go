package repositories

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/models"
	"gorm.io/gorm"
)

type MedicalRecordRepository interface {
	FindAll() ([]models.MedicalRecord, error)
	GetMedicalRecordByStudentNumber(studentNumber string) (models.MedicalRecord, error)
	CreateMedicalRecord(medicalRecord models.MedicalRecord) error
	DeleteMedicalRecord(medicalRecord models.MedicalRecord) error
	DeleteMedicalRecordByStudentNumber(studentNumber string) error
	UpdateMedicalRecord(medicalRecord models.MedicalRecord) error
}

type MedicalRecordRepositoryImpl struct {
	DB *gorm.DB
}

func (m MedicalRecordRepositoryImpl) DeleteMedicalRecordByStudentNumber(studentNumber string) error {
	var medicalRecord models.MedicalRecord
	if err := m.DB.Where("student_number = ?", studentNumber).Delete(&medicalRecord).Error; err != nil {
		return err
	}
	return nil
}

func (m MedicalRecordRepositoryImpl) UpdateMedicalRecord(medicalRecord models.MedicalRecord) error {
	return m.DB.Save(&medicalRecord).Error
}

func (m MedicalRecordRepositoryImpl) FindAll() ([]models.MedicalRecord, error) {
	var medicalRecords []models.MedicalRecord

	if err := m.DB.Find(&medicalRecords).Error; err != nil {
		return nil, err
	}

	return medicalRecords, nil
}

func (m MedicalRecordRepositoryImpl) GetMedicalRecordByStudentNumber(studentNumber string) (models.MedicalRecord, error) {
	var medicalRecord models.MedicalRecord

	if err := m.DB.Where("student_number = ?", studentNumber).First(&medicalRecord).Error; err != nil {
		return models.MedicalRecord{}, err
	}

	return medicalRecord, nil
}

func (m MedicalRecordRepositoryImpl) CreateMedicalRecord(medicalRecord models.MedicalRecord) error {
	return m.DB.Create(&medicalRecord).Error
}

func (m MedicalRecordRepositoryImpl) DeleteMedicalRecord(medicalRecord models.MedicalRecord) error {
	return m.DB.Delete(&medicalRecord).Error
}

func NewMedicalRecordRepository(db *gorm.DB) MedicalRecordRepository {
	return &MedicalRecordRepositoryImpl{
		DB: db,
	}
}
