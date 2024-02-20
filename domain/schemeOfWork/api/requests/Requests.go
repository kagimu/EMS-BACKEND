package requests

type ExcelSchemeOfWorkInfo struct {
	SubTopic      string
	SortingCode   string
	EstimatedTime string
}

type SubtopicListInfo struct {
	Number        string `json:"number"`
	Name          string `json:"name"`
	Code          int64  `json:"code"`
	SortCode      int64  `json:"sort_code"`
	UserNumber    string `json:"user_number"`
	EstimatedTime int64  `json:"estimated_time"`
	TopicNumber   string `json:"topic_number"`
}

type SchemeOfWorkInfo struct {
	Name               string             `json:"name"`
	ClassNumber        string             `json:"class_number"`
	SubjectNumber      string             `json:"subject_number"`
	AcademicYearNumber string             `json:"academic_year_number"`
	SchemeNumber       string             `json:"scheme_number"`
	SubtopicList       []SubtopicListInfo `json:"subtopicList"`
}
