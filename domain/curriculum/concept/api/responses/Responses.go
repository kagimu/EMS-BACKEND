package responses

type (
	ConceptList struct {
		Number         string `json:"number"`
		Name           string `json:"name"`
		Code           int64  `json:"code"`
		UserNumber     string `json:"user_number"`
		SubTopicNumber string `json:"sub_topic_number"`
		SubTopicName   string `json:"sub_topic_name"`
		TopicNumber    string `json:"topic_number"`
		TopicName      string `json:"topic_name"`
		StrandNumber   string `json:"strand_number"`
		StrandName     string `json:"strand_name"`
		SubjectNumber  string `json:"subject_number"`
		SubjectName    string `json:"subject_name"`
	}
)
