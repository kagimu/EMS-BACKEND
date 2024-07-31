package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/catering/models"
	"gorm.io/gorm"
)

func GetUnitOfMeasureByNumber(unitOfMeasureNumber string) (models.UnitOfMeasure, error) {

	db := config.DB

	var unitOfMeasure models.UnitOfMeasure

	if err := db.Where("number = ? ", unitOfMeasureNumber).First(&unitOfMeasure).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return unitOfMeasure, nil
		}

		return unitOfMeasure, err

	}

	return unitOfMeasure, nil
}

type UnitOfMeasureRepository interface {
	UpdateUnitOfMeasure(unitOfMeasure models.UnitOfMeasure, trn ...*gorm.DB) error
	CreateUnitOfMeasure(unitOfMeasure models.UnitOfMeasure, trn ...*gorm.DB) error
	GetUnitOfMeasureByNumber(unitOfMeasureNumber string, trn ...*gorm.DB) (models.UnitOfMeasure, error)
	FindUnitOfMeasureByCodeAndSchoolNumber(code uint64, schoolNumber string, trn ...*gorm.DB) (models.UnitOfMeasure, error)
	FindUnitOfMeasureByNameAndSchoolNumber(name, schoolNumber string, trn ...*gorm.DB) (models.UnitOfMeasure, error)
	FindUnitOfMeasuresBySchoolNumber(schoolNumber string, trn ...*gorm.DB) ([]models.UnitOfMeasure, error)
	RemoveUnitOfMeasure(unitOfMeasure models.UnitOfMeasure, trn ...*gorm.DB) error
}

type unitOfMeasureRepo struct {
}

func (u *unitOfMeasureRepo) RemoveUnitOfMeasure(unitOfMeasure models.UnitOfMeasure, trn ...*gorm.DB) error {
	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	return db.Delete(&unitOfMeasure).Error
}

func (u *unitOfMeasureRepo) FindUnitOfMeasuresBySchoolNumber(schoolNumber string, trn ...*gorm.DB) ([]models.UnitOfMeasure, error) {
	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	var unitOfMeasure []models.UnitOfMeasure

	if err := db.Where("school_number = ?", schoolNumber).Order("code asc").Find(&unitOfMeasure).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return unitOfMeasure, nil
		}

		return unitOfMeasure, err
	}

	return unitOfMeasure, nil
}

func (u *unitOfMeasureRepo) CreateUnitOfMeasure(unitOfMeasure models.UnitOfMeasure, trn ...*gorm.DB) error {
	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	return db.Create(&unitOfMeasure).Error
}

func (u *unitOfMeasureRepo) FindUnitOfMeasureByNameAndSchoolNumber(name, schoolNumber string, trn ...*gorm.DB) (models.UnitOfMeasure, error) {

	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	var unitOfMeasure models.UnitOfMeasure

	if err := db.Where("name = ? AND school_number = ?", name, schoolNumber).First(&unitOfMeasure).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return unitOfMeasure, nil
		}

		return unitOfMeasure, err
	}

	return unitOfMeasure, nil

}

func (u *unitOfMeasureRepo) FindUnitOfMeasureByCodeAndSchoolNumber(code uint64, schoolNumber string, trn ...*gorm.DB) (models.UnitOfMeasure, error) {
	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	var unitOfMeasure models.UnitOfMeasure

	if err := db.Where("code = ? AND school_number = ?", code, schoolNumber).First(&unitOfMeasure).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return unitOfMeasure, nil
		}

		return unitOfMeasure, err
	}

	return unitOfMeasure, nil
}

func (u *unitOfMeasureRepo) UpdateUnitOfMeasure(unitOfMeasure models.UnitOfMeasure, trn ...*gorm.DB) error {
	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	return db.Model(&unitOfMeasure).
		Select("Name", "Code", "Abbreviation").
		Where("number = ? ", unitOfMeasure.Number).
		Updates(models.UnitOfMeasure{Name: unitOfMeasure.Name, Code: unitOfMeasure.Code, Abbreviation: unitOfMeasure.Abbreviation}).
		Error
}

func (u *unitOfMeasureRepo) GetUnitOfMeasureByNumber(unitOfMeasureNumber string, trn ...*gorm.DB) (models.UnitOfMeasure, error) {
	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	var unitOfMeasure models.UnitOfMeasure

	if err := db.Where("number = ?", unitOfMeasureNumber).First(&unitOfMeasure).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return unitOfMeasure, nil
		}

		return unitOfMeasure, err
	}

	return unitOfMeasure, nil
}

func NewUnitOfMeasureRepository() UnitOfMeasureRepository {
	return &unitOfMeasureRepo{}
}
