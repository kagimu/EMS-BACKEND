package repositories

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/models"
	"gorm.io/gorm"
)

type FamilyMemberRepository interface {
	FindAll() ([]models.FamilyMember, error)
	FindFamilyMembersByStudentNumber(studentNumber string) ([]models.FamilyMember, error)
	CreateFamilyMember(FamilyMember models.FamilyMember) error
	DeleteFamilyMember(FamilyMember models.FamilyMember) error
	DeleteFamilyMemberByStudentNumber(studentNumber string) error
	UpdateFamilyMember(FamilyMember models.FamilyMember) error
}

type FamilyMemberRepositoryImpl struct {
	DB *gorm.DB
}

func (f FamilyMemberRepositoryImpl) DeleteFamilyMemberByStudentNumber(studentNumber string) error {
	var familyMember models.FamilyMember

	if err := f.DB.Where("student_number = ?", studentNumber).Delete(&familyMember).Error; err != nil {
		return err
	}
	return nil
}

func (f FamilyMemberRepositoryImpl) FindAll() ([]models.FamilyMember, error) {
	var FamilyMembers []models.FamilyMember

	if err := f.DB.Find(&FamilyMembers).Error; err != nil {
		return nil, err
	}

	return FamilyMembers, nil
}

func (f FamilyMemberRepositoryImpl) FindFamilyMembersByStudentNumber(studentNumber string) ([]models.FamilyMember, error) {
	var FamilyMembers []models.FamilyMember

	if err := f.DB.Where("student_number = ?", studentNumber).Find(&FamilyMembers).Error; err != nil {
		return FamilyMembers, err
	}

	return FamilyMembers, nil
}

func (f FamilyMemberRepositoryImpl) CreateFamilyMember(FamilyMember models.FamilyMember) error {
	return f.DB.Create(&FamilyMember).Error
}

func (f FamilyMemberRepositoryImpl) DeleteFamilyMember(FamilyMember models.FamilyMember) error {
	return f.DB.Delete(&FamilyMember).Error
}

func (f FamilyMemberRepositoryImpl) UpdateFamilyMember(FamilyMember models.FamilyMember) error {
	return f.DB.Save(&FamilyMember).Error
}

func NewFamilyMemberRepository(db *gorm.DB) FamilyMemberRepository {
	return &FamilyMemberRepositoryImpl{
		DB: db,
	}
}
