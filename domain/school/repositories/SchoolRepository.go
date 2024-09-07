package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	schoolModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/school/models"
	"gorm.io/gorm"
)

type SchoolRepository interface {
	FindSchoolByNumber(schoolNumber string) (schoolModel.School, error)
	UpdateSchool(school schoolModel.School) error
}

type SchoolRepositoryImpl struct {
	db *gorm.DB
}

func NewSchoolRepositoryImpl(db *gorm.DB) SchoolRepository {
	return &SchoolRepositoryImpl{db}
}

func (s SchoolRepositoryImpl) FindSchoolByNumber(schoolNumber string) (schoolModel.School, error) {
	var school schoolModel.School

	if err := s.db.Where("number = ?", schoolNumber).First(&school).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return school, nil
		}

		return school, err
	}

	return school, nil
}

func (s SchoolRepositoryImpl) UpdateSchool(school schoolModel.School) error {
	return s.db.Save(&school).Error
}

func GetSchoolByNumber(schoolNumber string, tx ...*gorm.DB) (schoolModel.School, error) {

	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}

	var school schoolModel.School

	if err := db.Where("number = ?", schoolNumber).First(&school).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return school, nil
		}

		return school, err
	}

	return school, nil

}

func RemoveSchool(school schoolModel.School, tx ...*gorm.DB) error {

	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}

	return db.Unscoped().Delete(&school).Error
}
