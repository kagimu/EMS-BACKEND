package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/filters"
	schoolModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/school/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SchoolRepository interface {
	FindSchoolByNumber(schoolNumber string) (schoolModel.School, error)
	UpdateSchool(school schoolModel.School) error
	FindSchoolsByEntityNumber(entityNumber string, ctx *gin.Context) ([]schoolModel.School, error)
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

func (s SchoolRepositoryImpl) FindSchoolsByEntityNumber(entityNumber string, ctx *gin.Context) ([]schoolModel.School, error) {
	var schools []schoolModel.School

	if err := s.db.
		Scopes(filters.ByScopedItems("school", ctx)).
		Where("entity_number = ?", entityNumber).Find(&schools).Error; err != nil {
		return schools, err
	}

	return schools, nil
}
