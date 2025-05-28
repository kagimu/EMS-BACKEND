package repositories

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/filters"
	assessmentAreaModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/assessmentArea/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SaveAssessmentArea(assessmentArea assessmentAreaModels.AssessmentArea, trn ...*gorm.DB) (assessmentAreaModels.AssessmentArea, error) {

	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	if err := db.Save(&assessmentArea).Error; err != nil {
		return assessmentArea, err
	}

	return assessmentArea, nil
}

func FindAssessmentAreaBySubjectNumberAndCurriculumClassNumber(subjectNumber string, curriculumClassNUmber string, trn ...*gorm.DB) (assessmentAreaModels.AssessmentArea, error) {

	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	var assessmentArea assessmentAreaModels.AssessmentArea

	if err := db.Where("curriculum_class_number = ? AND subject_number = ? ", curriculumClassNUmber, subjectNumber).First(&assessmentArea).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return assessmentArea, nil
		}

		return assessmentArea, err
	}

	return assessmentArea, nil
}

func FindAssessmentAreaByNameAndCurriculumClassNumber(otherName string, curriculumClassNUmber string, trn ...*gorm.DB) (assessmentAreaModels.AssessmentArea, error) {

	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	var assessmentArea assessmentAreaModels.AssessmentArea

	if err := db.Where("curriculum_class_number = ? AND name = ? ", curriculumClassNUmber, otherName).First(&assessmentArea).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return assessmentArea, nil
		}

		return assessmentArea, err
	}

	return assessmentArea, nil
}

func GetAssessmentAreaByNumber(assessmentAreaNumber string, trn ...*gorm.DB) (assessmentAreaModels.AssessmentArea, error) {

	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	var assessmentArea assessmentAreaModels.AssessmentArea

	if err := db.Where("number = ?", assessmentAreaNumber).First(&assessmentArea).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return assessmentArea, nil
		}

		return assessmentArea, err
	}

	return assessmentArea, nil

}

func FindAssessmentAreasByCurriculumClassNumber(curriculumClassNumber string, trn ...*gorm.DB) ([]assessmentAreaModels.AssessmentArea, error) {

	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	var assessmentAreas []assessmentAreaModels.AssessmentArea

	if err := db.Where("curriculum_class_number = ?", curriculumClassNumber).Find(&assessmentAreas).Error; err != nil {

		return assessmentAreas, err
	}

	return assessmentAreas, nil
}

func FindAllAssessmentAreas(curriculumClassNumber, classNumber, curriculumNumber, levelNumber, curriculumLevelNumber string, ctx *gin.Context, trn ...*gorm.DB) ([]assessmentAreaModels.AssessmentArea, error) {

	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	var assessmentAreas []assessmentAreaModels.AssessmentArea

	if err := db.
		Scopes(filters.ByScopedItems("assessmentArea", ctx, "assessment_areas.number")).
		Joins("JOIN curriculum_classes ON curriculum_classes.number = assessment_areas.curriculum_class_number").
		Joins("JOIN curriculum_levels ON curriculum_levels.number = curriculum_classes.curriculum_level_number").
		Scopes(filters.ByAnyField("curriculum_class_number", curriculumClassNumber), filters.ByAnyField("curriculum_classes.curriculum_level_number", curriculumLevelNumber), filters.ByAnyField("curriculum_levels.level_number", levelNumber), filters.ByAnyField("curriculum_classes.class_number", classNumber), filters.ByAnyField("curriculum_levels.curriculum_number", curriculumNumber)).
		Distinct("assessment_areas.*").
		Find(&assessmentAreas).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return assessmentAreas, nil
		}

		return assessmentAreas, err
	}

	return assessmentAreas, nil
}

func FindAllAssessmentAreasBySubjectCluster(curriculumClassNumber, curriculumNumber, learningAreaNumber string, trn ...*gorm.DB) ([]assessmentAreaModels.AssessmentArea, error) {
	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	var assessmentAreas []assessmentAreaModels.AssessmentArea

	if err := db.Where("assessment_areas.curriculum_class_number = ?", curriculumClassNumber).
		Joins("JOIN subjects on subjects.number = assessment_areas.subject_number AND subjects.curriculum_number=?", curriculumNumber).
		Joins("JOIN subject_clusters on subject_clusters.number = subjects.subject_cluster_number AND subject_clusters.learning_area_number =?", learningAreaNumber).
		Find(&assessmentAreas).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return assessmentAreas, nil
		}
		return assessmentAreas, err
	}

	return assessmentAreas, nil
}

func DeleteAssessmentArea(assessmentAreaNumber string, trn ...*gorm.DB) error {

	db := config.DB
	if len(trn) > 0 {
		db = trn[0]
	}

	if err := db.Where("number = ?", assessmentAreaNumber).Delete(&assessmentAreaModels.AssessmentArea{}).Error; err != nil {
		return err
	}

	return nil
}
