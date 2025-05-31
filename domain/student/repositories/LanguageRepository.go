package repositories

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/models"
	"gorm.io/gorm"
)

type LanguageRepository interface {
	FindAll() ([]models.StudentLanguage, error)
	FindLanguagesByStudentNumber(studentNumber string) ([]models.StudentLanguage, error)
	CreateLanguage(language models.StudentLanguage) error
	DeleteLanguage(language models.StudentLanguage) error
	DeleteLanguageByStudentNumber(studentNumber string) error
	UpdateLanguage(language models.StudentLanguage) error
}

type LanguageRepositoryImpl struct {
	DB *gorm.DB
}

func (l LanguageRepositoryImpl) DeleteLanguageByStudentNumber(studentNumber string) error {
	var language models.StudentLanguage

	if err := l.DB.Where("student_number = ?", studentNumber).Delete(&language).Error; err != nil {
		return err
	}
	return nil
}

func (l LanguageRepositoryImpl) FindAll() ([]models.StudentLanguage, error) {
	var languages []models.StudentLanguage

	if err := l.DB.Find(&languages).Error; err != nil {
		return nil, err
	}

	return languages, nil
}

func (l LanguageRepositoryImpl) FindLanguagesByStudentNumber(studentNumber string) ([]models.StudentLanguage, error) {
	var languages []models.StudentLanguage

	if err := l.DB.Where("student_number = ?", studentNumber).Find(&languages).Error; err != nil {
		return languages, err
	}

	return languages, nil
}

func (l LanguageRepositoryImpl) CreateLanguage(language models.StudentLanguage) error {
	return l.DB.Create(&language).Error
}

func (l LanguageRepositoryImpl) DeleteLanguage(language models.StudentLanguage) error {
	return l.DB.Delete(&language).Error
}

func (l LanguageRepositoryImpl) UpdateLanguage(language models.StudentLanguage) error {
	return l.DB.Save(&language).Error
}

func NewLanguageRepository(db *gorm.DB) LanguageRepository {
	return &LanguageRepositoryImpl{
		DB: db,
	}
}
