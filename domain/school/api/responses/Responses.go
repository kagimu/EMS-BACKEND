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

	ModeOfStudyResponse struct {
		Name              string `json:"name"`
		ModeOfStudyNumber string `json:"mode_of_study_number"`
		SchoolNumber      string `json:"school_number"`
		Status            string `json:"status"`
		CreatedAt         string `json:"created_at"`
	}
)
