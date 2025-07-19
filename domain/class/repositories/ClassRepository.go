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
		db = trn[0]
	}

	var class models.Class

	if err := db.Unscoped().Where("number = ?", classNumber).First(&class).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return class, nil
		}

		return class, err
	}

	return class, nil
}

func GetClassByCode(code int64, entityNumber string, trn ...*gorm.DB) (models.Class, error) {
	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	var class models.Class

	if err := db.Unscoped().Where("code = ? AND entity_number = ?", code, entityNumber).First(&class).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return class, nil
		}

		return class, err
	}

	return class, nil
}

func GetClassByName(name, entityNumber string, trn ...*gorm.DB) (models.Class, error) {
	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	var class models.Class

	if err := db.Unscoped().Where("name = ? AND entity_number = ?", name, entityNumber).First(&class).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return class, nil
		}

		return class, err
	}

	return class, nil
}

func GetClassByAbbreviation(abbreviation, entityNumber string, trn ...*gorm.DB) (models.Class, error) {
	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	var class models.Class

	if err := db.Unscoped().Where("abbreviation = ? AND entity_number = ?", abbreviation, entityNumber).First(&class).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return class, nil
		}

		return class, err
	}

	return class, nil
}

func FindClassesBySchoolNumber(schoolNumber, campusLevelNumber, campusNumber, curriculumNumber, campusCurriculumNumber string, trn ...*gorm.DB) ([]models.Class, error) {
	db := config.DB
	if len(trn) > 0 {
		db = trn[1]
	}

	var classes []models.Class

	if err := db.
		Joins("JOIN curriculum_classes ON curriculum_classes.class_number = classes.number").
		Joins("JOIN campus_classes ON campus_classes.curriculum_class_number = curriculum_classes.number").
		Joins("JOIN campus_levels ON campus_levels.number = campus_classes.campus_level_number").
		Joins("JOIN curriculum_levels ON curriculum_levels.number = curriculum_classes.curriculum_level_number").
		Joins("JOIN campuses ON campuses.number = campus_classes.campus_number").
		Scopes(
			filters.ByAnyField("campuses.school_number", schoolNumber),
			filters.ByAnyField("campus_levels.campus_curriculum_number", campusCurriculumNumber),
			filters.ByAnyField("campus_levels.number", campusLevelNumber),
			filters.ByAnyField("campus_classes.campus_number", campusNumber),
			filters.ByAnyField("curriculum_levels.curriculum_number", curriculumNumber)).
		Distinct("classes.*").
		Order("classes.code asc").
		Find(&classes).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return classes, nil
		}

		return classes, err
	}

	return classes, nil
}

func FindClassesByAssessmentGradeNumber(assessmentGradeNumber string, trn ...*gorm.DB) ([]models.Class, error) {

	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	var classes []models.Class

	if err := db.
		Joins("join assessment_grade_classes on assessment_grade_classes.class_number = classes.number").
		Where("assessment_grade_number = ?", assessmentGradeNumber).Find(&classes).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return classes, nil
		}

		return classes, err
	}

	return classes, nil
}

func SaveClass(class models.Class, tx ...*gorm.DB) error {

	db := config.DB
	if len(tx) > 0 {
		db = tx[0]
	}

	return db.Save(&class).Error
}

func FindClassAfterLevelCode(entityNumber string, currentLevelCode int64, trn ...*gorm.DB) (models.Class, error) {
	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	var class models.Class

	err := db.
		Joins("JOIN levels ON levels.number = classes.level_number").
		Where("classes.entity_number = ? AND levels.code > ?", entityNumber, currentLevelCode).
		Order("levels.code ASC, classes.code ASC").
		First(&class).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return class, nil
	}
	return class, err
}

func FindClassBeforeLevelCode(entityNumber string, currentLevelCode int64, trn ...*gorm.DB) (models.Class, error) {
	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	var class models.Class

	err := db.
		Joins("JOIN levels ON levels.number = classes.level_number").
		Where("classes.entity_number = ? AND levels.code < ?", entityNumber, currentLevelCode).
		Order("levels.code DESC, classes.code DESC").
		First(&class).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return class, nil
	}
	return class, err
}
