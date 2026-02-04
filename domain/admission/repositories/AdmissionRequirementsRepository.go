package repositories

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/admission/models"
	"gorm.io/gorm"
)

type AdmissionRequirementsRepository interface {
	SaveAdmissionRequirements(admissionTemplate models.AdmissionRequirement) error
}

type admissionRequirementsRepository struct {
	DB *gorm.DB
}

func (a admissionRequirementsRepository) SaveAdmissionRequirements(admissionRequirements models.AdmissionRequirement) error {
	if err := a.DB.Create(&admissionRequirements).Error; err != nil {
		return err
	}
	return nil
}

func NewAdmissionRequirementsRepository(db *gorm.DB) AdmissionRequirementsRepository {
	return &admissionRequirementsRepository{DB: db}
}
