package responses

type (
	SchoolResponseInfo struct {
		Number           string                         `json:"number"`
		Name             string                         `json:"name"`
		CommencementDate string                         `json:"commencement_date"`
		SchoolCurriculum []SchoolCurriculumResponseInfo `json:"school_curriculum"`
	}

	SchoolCurriculumResponseInfo struct {
		Number string `json:"number"`
		Name   string `json:"name"`
	}
)
