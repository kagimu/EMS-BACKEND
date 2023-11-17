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
		Scopes(scopeFilter("number", learningAreaNumber)).
		Order("code asc").
		Find(&learningAreas).Error; err != nil {
		errInfo = "Learning area not found!!!"
	}
	return learningAreas, errInfo
}

func FetchSubjectClusterCurriculum(learningAreaNumber string) ([]models3.SubjectCluster, string) {
	var subjectCluster []models3.SubjectCluster
	var errInfo string
	if err := config.DB.
		Scopes(scopeFilter("learning_area_number", learningAreaNumber)).
		Order("code asc").
		Find(&subjectCluster).Error; err != nil {
		errInfo = "No Subject Cluster!"
	}
	return subjectCluster, errInfo
}

func scopeFilter(tableType, number string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if number != "" {
			return db.Where(tableType+" = ?", number)
		}
		return db
	}
}
