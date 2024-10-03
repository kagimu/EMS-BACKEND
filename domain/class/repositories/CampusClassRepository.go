package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	campusClassModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/class/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/filters"
	"gorm.io/gorm"
)

func CreateCampusClass(campusClass *campusClassModel.CampusClass, txn ...*gorm.DB) error {
	db := config.DB
	if len(txn) > 0 {
		db = txn[0]
	}
	if err := db.Create(&campusClass).Error; err != nil {
		return err
	}
	return nil

}

func GetCampusClassByNumber(campusClassNumber string) (campusClassModel.CampusClass, error) {

	var campusClass campusClassModel.CampusClass

	db := config.DB

	if err := db.Where("number = ? ", campusClassNumber).First(&campusClass).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return campusClass, nil
		}

		return campusClass, err
	}

	/*if err != nil {
		return campusClass, err
	}*/

	return campusClass, nil
}

func GetCampusClassByCurriculumClassNumberAndCampusNumberAndCampusLevelNumber(curriculumClassNumber, campusNumber, campusLevelNumber string) (campusClassModel.CampusClass, error) {
	var duplicateClass campusClassModel.CampusClass
	if err := config.DB.
		Where("curriculum_class_number = ?", curriculumClassNumber).
		Where("campus_number = ?", campusNumber).
		Where("campus_level_number = ?", campusLevelNumber).
		First(&duplicateClass).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return duplicateClass, nil
		}
		return duplicateClass, err
	}
	return duplicateClass, nil

}

func DeleteCampusClassByNumber(classNumber, campusNumber, campusLevelNumber string, txn ...*gorm.DB) error {
	db := config.DB
	if len(txn) > 0 {
		db = txn[0]
	}
	var campusClass campusClassModel.CampusClass
	//delete class if it is removed from the level
	if err := db.Unscoped().Where("curriculum_class_number = ?", classNumber).
		Where("campus_number = ?", campusNumber).
		Where("campus_level_number = ?", campusLevelNumber).
		Delete(&campusClass).Error; err != nil {
		return err
	}

	return nil
}

func GetCampusClassByClassNumberAndCampusNumber(classNumber, campusNumber, curriculumNumber, campusCurriculumNumber string, txn ...*gorm.DB) ([]campusClassModel.CampusClass, error) {
	db := config.DB
	if len(txn) > 0 {
		db = txn[0]
	}

	var campusClass []campusClassModel.CampusClass

	if err := db.
		Joins("JOIN campus_levels ON campus_levels.number = campus_classes.campus_level_number").
		Joins("JOIN curriculum_classes ON curriculum_classes.number = campus_classes.curriculum_class_number").
		Joins("JOIN curriculum_levels ON curriculum_levels.number = curriculum_classes.curriculum_level_number").
		Scopes(
			filters.ByAnyField("campus_levels.campus_curriculum_number", campusCurriculumNumber),
			filters.ByAnyField("curriculum_classes.class_number", classNumber),
			filters.ByAnyField("campus_classes.campus_number", campusNumber),
			filters.ByAnyField("curriculum_levels.curriculum_number", curriculumNumber)).
		Find(&campusClass).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return campusClass, nil
		}

		return campusClass, err
	}

	return campusClass, nil
}
