package repositories

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/admission/models"
	"gorm.io/gorm"
)

type ToaRepository interface {
	SaveToa(toa models.TermsOfAdmission) error
	UpsertToa(toa models.TermsOfAdmission) error
	FindToaBySchoolNumber(schoolNumber string) (models.TermsOfAdmission, error)
}

type toaRepository struct {
	DB *gorm.DB
}

func (r *toaRepository) SaveToa(toa models.TermsOfAdmission) error {
	return r.DB.Create(&toa).Error
}

// UpsertToa updates if a record for the school exists, otherwise creates a new one.
func (r *toaRepository) UpsertToa(toa models.TermsOfAdmission) error {
	var existing models.TermsOfAdmission
	err := r.DB.Where("school_number = ?", toa.SchoolNumber).First(&existing).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return r.DB.Create(&toa).Error
		}
		return err
	}
	existing.Content = toa.Content
	return r.DB.Save(&existing).Error
}

func (r *toaRepository) FindToaBySchoolNumber(schoolNumber string) (models.TermsOfAdmission, error) {
	var toa models.TermsOfAdmission
	err := r.DB.Where("school_number = ?", schoolNumber).Order("created_at desc").First(&toa).Error
	if err != nil {
		return toa, err
	}
	return toa, nil
}

func NewToaRepository(db *gorm.DB) ToaRepository {
	return &toaRepository{DB: db}
}
