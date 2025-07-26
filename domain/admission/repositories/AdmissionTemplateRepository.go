package repositories

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/admission/models"
	"gorm.io/gorm"
)

type AdmissionTemplateRepository interface {
	SaveAdmissionTemplate(admissionTemplate models.AdmissionLetterTemplate) error
	FindAdmissionTemplateBySchoolNumber(schoolNumber string) (models.AdmissionLetterTemplate, error)
	FindAdmissionTemplateByNumber(templateNumber string) (models.AdmissionLetterTemplate, error)
	DeleteAdmissionTemplateBySchoolNumber(schoolNumber string) error
	DeleteAdmissionTemplateByNumber(templateNumber string) error
	UpdateAdmissionTemplate(admissionTemplate models.AdmissionLetterTemplate) error
}

type admissionTemplateRepository struct {
	DB *gorm.DB
}

func (a admissionTemplateRepository) UpdateAdmissionTemplate(admissionTemplate models.AdmissionLetterTemplate) error {
	if err := a.DB.Save(&admissionTemplate).Error; err != nil {
		return err
	}
	return nil
}

func (a admissionTemplateRepository) SaveAdmissionTemplate(admissionTemplate models.AdmissionLetterTemplate) error {
	if err := a.DB.Create(&admissionTemplate).Error; err != nil {
		return err
	}
	return nil
}

func (a admissionTemplateRepository) FindAdmissionTemplateBySchoolNumber(schoolNumber string) (models.AdmissionLetterTemplate, error) {
	var admissionLetterTemplate models.AdmissionLetterTemplate
	if err := a.DB.Where("school_number = ?", schoolNumber).Order("created_at desc").First(&admissionLetterTemplate).Error; err != nil {
		return admissionLetterTemplate, err
	}
	return admissionLetterTemplate, nil
}

func (a admissionTemplateRepository) FindAdmissionTemplateByNumber(templateNumber string) (models.AdmissionLetterTemplate, error) {
	var admissionLetterTemplate models.AdmissionLetterTemplate
	if err := a.DB.Where("number = ?", templateNumber).Order("created_at desc").First(&admissionLetterTemplate).Error; err != nil {
		return admissionLetterTemplate, err
	}
	return admissionLetterTemplate, nil
}

func (a admissionTemplateRepository) DeleteAdmissionTemplateBySchoolNumber(schoolNumber string) error {
	var AdmissionLetterTemplate models.AdmissionLetterTemplate
	if err := a.DB.Where("school_number = ?", schoolNumber).Delete(&AdmissionLetterTemplate).Error; err != nil {
		return err
	}
	return nil
}

func (a admissionTemplateRepository) DeleteAdmissionTemplateByNumber(templateNumber string) error {
	var AdmissionLetterTemplate models.AdmissionLetterTemplate
	if err := a.DB.Where("number = ?", templateNumber).Delete(&AdmissionLetterTemplate).Error; err != nil {
		return err
	}
	return nil
}

func NewAdmissionTemplateRepository(db *gorm.DB) AdmissionTemplateRepository {
	return &admissionTemplateRepository{
		DB: db,
	}
}
