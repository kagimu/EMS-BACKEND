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

	if err := db.Where("number = ?", classNumber).First(&class).Error; err != nil {

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
		Joins("join curriculum_levels on curriculum_levels.level_number = classes.level_number").
		Joins("join curriculum_classes on curriculum_classes.curriculum_level_number = curriculum_levels.number").
		Joins("join campus_levels on campus_levels.curriculum_level_number = curriculum_levels.number").
		Joins("join curriculum_classes on curriculum_classes.class_number = classes.number").
		Joins("join campus_classes on campus_classes.curriculum_class_number = curriculum_classes.number").
		Joins("join campuses on campuses.number = campus_classes.campus_number").
		Scopes(filters.ByAnyField("campuses.school_number", schoolNumber), filters.ByAnyField("campus_levels.campus_curriculum_number", campusCurriculumNumber), filters.ByAnyField("curriculum_classes.curriculum_number", curriculumNumber), filters.ByAnyField("campus_classes.campus_number", campusNumber), filters.ByAnyField("campus_level_number", campusLevelNumber)).
		Distinct("classes.number", "classes.name", "classes.code", "classes.abbreviation", "classes.user_number", "classes.level_number", "classes.created_at", "classes.updated_at", "classes.id").
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
