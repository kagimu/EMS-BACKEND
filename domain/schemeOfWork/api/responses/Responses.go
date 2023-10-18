package responses

type (
	SchemeOfWorkResponseInfo struct {
		Number                   string                     `json:"number"`
		Name                     string                     `json:"name"`
		Status                   string                     `json:"status"`
		ClassNumber              string                     `json:"class_number"`
		UserNumber               string                     `json:"user_number"`
		SchemeOfWorkSubtopicInfo []SchemeOfWorkSubtopicInfo `json:"subtopic"`
	}
	SchemeOfWorkSubtopicInfo struct {
		Number             string `json:"number"`
		Code               int64  `json:"code"`
		Name               string `json:"name"`
		EstimatedTime      int64  `json:"estimated_time"`
		SchemeOfWorkNumber string `json:"scheme_of_work_number"`
		SubtopicNumber     string `json:"sub_topic_number"`
		UserNumber         string `json:"user_number"`
	}
)
