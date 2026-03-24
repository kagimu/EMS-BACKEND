package repositories

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/admission/models"
	"gorm.io/gorm"
)

type AdmissionRequirementsRepository interface {
	SaveAdmissionRequirements(admissionRequirement models.AdmissionRequirement) error
	SaveAdmissionRequirementsBatch(requirements []models.AdmissionRequirement) error
	FindAdmissionRequirements(schoolNumber, campusClassNumber, studyPeriodNumber, academicYearNumber, itemNumber string) ([]models.AdmissionRequirement, error)
	DeleteAdmissionRequirementByNumber(number string) error
	UpdateAdmissionRequirementStatus(number, status string) error
}

type admissionRequirementsRepository struct {
	DB *gorm.DB
}

func (a admissionRequirementsRepository) SaveAdmissionRequirements(admissionRequirement models.AdmissionRequirement) error {
	return a.DB.Create(&admissionRequirement).Error
}

func (a admissionRequirementsRepository) SaveAdmissionRequirementsBatch(requirements []models.AdmissionRequirement) error {
	if len(requirements) == 0 {
		return nil
	}
	return a.DB.Create(&requirements).Error
}

func (a admissionRequirementsRepository) DeleteAdmissionRequirementByNumber(number string) error {
	return a.DB.Where("number = ?", number).Delete(&models.AdmissionRequirement{}).Error
}

func (a admissionRequirementsRepository) UpdateAdmissionRequirementStatus(number, status string) error {
	return a.DB.Model(&models.AdmissionRequirement{}).
		Where("number = ?", number).
		Update("status", status).Error
}

// FindAdmissionRequirements returns admission requirements filtered by non-empty query params.
func (a admissionRequirementsRepository) FindAdmissionRequirements(schoolNumber, campusClassNumber, studyPeriodNumber, academicYearNumber, itemNumber, modeOfStudyNumber string) ([]models.AdmissionRequirement, error) {
	var list []models.AdmissionRequirement
	query := a.DB.Model(&models.AdmissionRequirement{})
	if schoolNumber != "" {
		query = query.Where("school_number = ?", schoolNumber)
	}
	if campusClassNumber != "" {
		query = query.Where("campus_class_number = ?", campusClassNumber)
	}
	if studyPeriodNumber != "" {
		query = query.Where("study_period_number = ?", studyPeriodNumber)
	}
	if academicYearNumber != "" {
		query = query.Where("academic_year_number = ?", academicYearNumber)
	}
	if itemNumber != "" {
		query = query.Where("item_number = ?", itemNumber)
	}
	if modeOfStudyNumber != "" {
		query = query.Where("mode_of_study_number = ?", modeOfStudyNumber)
	}
	if err := query.Order("created_at DESC").Find(&list).Error; err != nil {
		return nil, err
	}
	
	return list, nil
}

func NewAdmissionRequirementsRepository(db *gorm.DB) AdmissionRequirementsRepository {
	return &admissionRequirementsRepository{DB: db}
}
