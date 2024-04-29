package responses

type ScopeResponse struct {
	Number           string `json:"number"`
	UserNumber       string `json:"user_number"`
	UserName         string `json:"user_name"`
	CampusNumber     string `json:"campus_number"`
	CampusName       string `json:"campus_name"`
	SchoolNumber     string `json:"school_number"`
	SchoolName       string `json:"school_name"`
	CurriculumNumber string `json:"curriculum_number"`
	CurriculumName   string `json:"curriculum_name"`
	ClassNumber      string `json:"class_number"`
	ClassName        string `json:"class_name"`
	StreamNumber     string `json:"stream_number"`
	StreamName       string `json:"stream_name"`
}
