package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	campusModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/campus/models"
	"gorm.io/gorm"
)

func GetCampusByNumber(campusNumber string, tx ...*gorm.DB) (campusModel.Campus, error) {

	db := config.DB

	if len(tx) > 0 {
		db = tx[0]
	}

	var campus campusModel.Campus
	if err := db.Where("number = ?", campusNumber).First(&campus).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return campus, nil
		}

		return campus, err
	}
	return campus, nil
}

func FindCampuses(schoolNumber string, tx ...*gorm.DB) ([]campusModel.Campus, error) {
	db := config.DB

	if len(tx) > 0 {
		db = tx[0]
	}

	var campus []campusModel.Campus

	if err := db.Where("school_number = ?", schoolNumber).Find(&campus).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return campus, nil
		}

		return campus, err
	}
	return campus, nil
}

func GetCampusByCampusClassNumber(campusClassNumber string, tx ...*gorm.DB) (campusModel.Campus, error) {

	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}

	var campus campusModel.Campus

	if err := db.
		Joins("JOIN campus_classes ON campus.number = campus_classes.campus_number").
		Where("campus_classes.number = ?", campusClassNumber).
		First(&campus).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return campus, nil
		}

		return campus, err
	}

	return campus, nil
}
