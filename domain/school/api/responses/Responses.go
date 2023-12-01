package responses

type (
	SchoolResponseInfo struct {
		Number           string                         `json:"number"`
		Name             string                         `json:"name"`
		SchoolCurriculum []SchoolCurriculumResponseInfo `json:"school_curriculum"`
	}

	SchoolCurriculumResponseInfo struct {
		Number string `json:"number"`
		Name   string `json:"name"`
	}
)
