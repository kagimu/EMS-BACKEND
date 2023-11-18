package responses

import "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/concept/api/responses"

type (
	SubTopicList struct {
		Number        string                  `json:"number"`
		Name          string                  `json:"name"`
		Code          int64                   `json:"code"`
		SortCode      int64                   `json:"sort_code"`
		UserNumber    string                  `json:"user_number"`
		EstimatedTime int64                   `json:"estimated_time"`
		TopicNumber   string                  `json:"topic_number"`
		TopicName     string                  `json:"topic_name"`
		StrandNumber  string                  `json:"strand_number"`
		StrandName    string                  `json:"strand_name"`
		SubjectNumber string                  `json:"subject_number"`
		SubjectName   string                  `json:"subject_name"`
		Concept       []responses.ConceptList `json:"concept"`
	}
)
