package responses

import "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/topic/api/responses"

type (
	StrandDataList struct {
		Number                   string                    `json:"number"`
		Name                     string                    `json:"name"`
		Code                     int64                     `json:"code"`
		UserNumber               string                    `json:"user_number"`
		SubjectNumber            string                    `json:"subject_number"`
		SubjectClusterNumber     string                    `json:"subject_cluster_number"`
		LearningAreaNumberNumber string                    `json:"learning_area_number"`
		Topic                    []responses.TopicDataList `json:"topic"`
	}
)
