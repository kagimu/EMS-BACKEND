package responses

import (
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/concept/api/responses"
	responses2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/settings/learningMaterials/api/responses"
)

type (
	SubTopicList struct {
		Number               string                            `json:"number"`
		Name                 string                            `json:"name"`
		Code                 int64                             `json:"code"`
		SortCode             int64                             `json:"sort_code"`
		UserNumber           string                            `json:"user_number"`
		EstimatedTime        int64                             `json:"estimated_time"`
		TopicNumber          string                            `json:"topic_number"`
		TopicName            string                            `json:"topic_name"`
		StrandNumber         string                            `json:"strand_number"`
		StrandName           string                            `json:"strand_name"`
		SubjectNumber        string                            `json:"subject_number"`
		SubjectName          string                            `json:"subject_name"`
		SubjectClusterNumber string                            `json:"subject_cluster_number"`
		LearningAreaNumber   string                            `json:"learning_area_number"`
		ReferenceInfo        []responses2.LearningMaterialList `json:"references"`
		Concept              []responses.ConceptList           `json:"concept"`
		TeachingSequenceCode int64                             `json:"teaching_sequence_code"`
	}
)
