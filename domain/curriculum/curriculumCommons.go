package curriculum

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	learningAreaModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/learningarea/models"
	models3 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/subjectcluster/models"
	"gorm.io/gorm"
)

func FetchLearningAreaCurriculum(learningAreaNumber string) ([]learningAreaModels.LearningArea, string) {
	var learningAreas []learningAreaModels.LearningArea
	var errInfo string
	if err := config.DB.
		Scopes(filterLearningArea(learningAreaNumber)).
		Order("code asc").
		Find(&learningAreas).Error; err != nil {
		errInfo = "Learning area not found!!!"
	}
	return learningAreas, errInfo
}

func filterLearningArea(learningAreaNumber string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if learningAreaNumber != "" {
			return db.Where("number = ?", learningAreaNumber)
		}
		return db
	}
}

func FetchSubjectClusterCurriculum(learningAreaNumber string) ([]models3.SubjectCluster, string) {
	var subjectCluster []models3.SubjectCluster
	var errInfo string
	if err := config.DB.
		Where("learning_area_number = ?", learningAreaNumber).
		Order("code asc").
		Find(&subjectCluster).Error; err != nil {
		errInfo = "No Subject Cluster!"
	}
	return subjectCluster, errInfo
}
