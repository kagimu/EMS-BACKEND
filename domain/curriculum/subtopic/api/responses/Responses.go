package responses

type (
	SubTopicList struct {
		Number        string `json:"number"`
		Name          string `json:"name"`
		Code          string `json:"code"`
		SortCode      string `json:"sort_code"`
		UserNumber    string `json:"user_number"`
		EstimatedTime int64  `json:"estimated_time"`
		TopicNumber   string `json:"topic_number"`
	}
)
