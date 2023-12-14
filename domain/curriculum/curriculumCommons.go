package curriculum

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/api/config"
	conceptModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/concept/models"
	learningAreaModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/learningarea/models"
	models2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/strand/models"
	models3 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/subjectcluster/models"
	subTopicModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/subtopic/models"
	topicModels "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/topic/models"
	subjectModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/subject/models"
	"gorm.io/gorm"
)

func FetchLearningAreaCurriculum(learningAreaNumber string, curriculumNumber string) ([]learningAreaModels.LearningArea, string) {
	var learningAreas []learningAreaModels.LearningArea
	var errInfo string
	if err := config.DB.
		Scopes(scopeFilter("number", learningAreaNumber)).
		Scopes(scopeFilter("curriculum_number", curriculumNumber)).
		Order("code asc").
		Find(&learningAreas).Error; err != nil {
		errInfo = "Learning area not found!!!"
	}
	return learningAreas, errInfo
}

func FetchSubjectClusterCurriculum(curriculumNumber string, learningAreaNumber string, subjectClusterNumber string) ([]models3.SubjectCluster, string) {
	var subjectCluster []models3.SubjectCluster
	var errInfo string
	if err := config.DB.
		Scopes(scopeFilter("curriculum_number", curriculumNumber)).
		Scopes(scopeFilter("number", subjectClusterNumber)).
		Scopes(scopeFilter("learning_area_number", learningAreaNumber)).
		Preload("LearningArea").
		Order("code asc").
		Find(&subjectCluster).Error; err != nil {
		errInfo = "No Subject Cluster!"
	}
	return subjectCluster, errInfo
}

func FetchSubjectCurriculum(subjectClusterNumber string, subjectNumber string) ([]subjectModel.Subject, string) {
	var subjects []subjectModel.Subject
	var errInfo string
	if err := config.DB.
		Scopes(scopeFilter("number", subjectNumber)).
		Scopes(scopeFilter("subject_cluster_number", subjectClusterNumber)).
		Order("code asc").
		Find(&subjects).Error; err != nil {
		errInfo = "No Subject!"
	}
	return subjects, errInfo
}

func FetchStrandCurriculum(subjectNumber string, strandNumber string) ([]models2.Strand, string) {
	var strand []models2.Strand
	var errInfo string
	if err := config.DB.
		Scopes(scopeFilter("number", strandNumber)).
		Scopes(scopeFilter("subject_number", subjectNumber)).
		Order("code asc").
		Find(&strand).Error; err != nil {
		errInfo = "No strand found!"
	}
	return strand, errInfo
}

func FetchTopicCurriculum(strandNumber string, topicNumber string, classNumber string) ([]topicModels.Topic, string) {
	var topic []topicModels.Topic
	var errInfo string
	if err := config.DB.
		Scopes(scopeFilter("number", topicNumber)).
		Scopes(scopeFilter("curriculum_class_number", classNumber)).
		Scopes(scopeFilter("strand_number", strandNumber)).
		//Preload("CurriculumClass.CurriculumLevel").
		//Preload("CurriculumClass").
		Order("code asc").
		Find(&topic).Error; err != nil {
		errInfo = "No topic found!"
	}
	return topic, errInfo
}

func FetchSubtopicCurriculum(topicNumber string, subtopicNumber string) ([]subTopicModels.SubTopic, string) {
	var subtopic []subTopicModels.SubTopic
	var errInfo string
	if err := config.DB.
		Scopes(scopeFilter("number", subtopicNumber)).
		Scopes(scopeFilter("topic_number", topicNumber)).
		Order("code asc").
		Find(&subtopic).Error; err != nil {
		errInfo = "No Sub topic found!"
	}
	return subtopic, errInfo
}

func FetchConceptCurriculum(subtopicNumber string) ([]conceptModels.Concept, string) {

	var concept []conceptModels.Concept
	var errInfo string
	if err := config.DB.
		Scopes(scopeFilter("sub_topic_number", subtopicNumber)).
		Order("code asc").
		Find(&concept).Error; err != nil {
		errInfo = "No Concept found!"
	}
	return concept, errInfo

}

func scopeFilter(tableType, number string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if number != "" {
			return db.Where(tableType+" = ?", number)
		}
		return db
	}
}
