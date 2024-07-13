package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/filters"
	"gorm.io/gorm"
)

func GetClassByCampusClassNumber(campusClassNumber string, trn ...*gorm.DB) (models.Class, error) {
	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	var class models.Class

	if err := db.
		Joins("JOIN curriculum_classes ON curriculum_classes.class_number = classes.number").
		Joins("JOIN campus_classes ON campus_classes.curriculum_class_number = curriculum_classes.number").
		Where("campus_classes.number = ?", campusClassNumber).First(&class).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return class, nil
		}

		return class, err
	}

	return class, nil
}

func GetClassByNumber(classNumber string, trn ...*gorm.DB) (models.Class, error) {
	db := config.DB
	if len(trn) > 0 {
		db = trn[1]
	}

	var class models.Class

	if err := db.Where("number = ?", classNumber).First(&class).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return class, nil
		}

		return class, err
	}

	return class, nil
}

func FindClassesBySchoolNumber(schoolNumber, levelNumber string, trn ...*gorm.DB) ([]models.Class, error) {
	db := config.DB
	if len(trn) > 0 {
		db = trn[1]
	}

	var classes []models.Class

	if err := db.Where(filters.BySchoolNumber(schoolNumber), filters.ByLevelNumber(levelNumber)).Find(&classes).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return classes, nil
		}

		return classes, err
	}

	return classes, nil
}
