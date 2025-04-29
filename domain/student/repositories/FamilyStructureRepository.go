package repositories

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/models"
	"gorm.io/gorm"
)

type FamilyStructureRepository interface {
	FindAll() ([]models.FamilyStructure, error)
	GetFamilyStructureByStudentNumber(studentNumber string) (models.FamilyStructure, error)
	CreateFamilyStructure(familyStructure models.FamilyStructure) error
	DeleteFamilyStructure(familyStructure models.FamilyStructure) error
	UpdateFamilyStructure(familyStructure models.FamilyStructure) error
}

type FamilyStructureRepositoryImpl struct {
	DB *gorm.DB
}

func (f FamilyStructureRepositoryImpl) FindAll() ([]models.FamilyStructure, error) {
	var familyStructures []models.FamilyStructure

	if err := f.DB.Find(&familyStructures).Error; err != nil {
		return nil, err
	}

	return familyStructures, nil
}

func (f FamilyStructureRepositoryImpl) GetFamilyStructureByStudentNumber(studentNumber string) (models.FamilyStructure, error) {
	var familyStructure models.FamilyStructure

	if err := f.DB.Where("student_number = ?", studentNumber).First(&familyStructure).Error; err != nil {
		return models.FamilyStructure{}, err
	}

	return familyStructure, nil
}

func (f FamilyStructureRepositoryImpl) CreateFamilyStructure(familyStructure models.FamilyStructure) error {
	return f.DB.Create(&familyStructure).Error
}

func (f FamilyStructureRepositoryImpl) DeleteFamilyStructure(familyStructure models.FamilyStructure) error {
	return f.DB.Delete(&familyStructure).Error
}

func (f FamilyStructureRepositoryImpl) UpdateFamilyStructure(familyStructure models.FamilyStructure) error {
	return f.DB.Save(&familyStructure).Error
}

func NewFamilyStructureRepository(db *gorm.DB) FamilyStructureRepository {
	return &FamilyStructureRepositoryImpl{
		DB: db,
	}
}
