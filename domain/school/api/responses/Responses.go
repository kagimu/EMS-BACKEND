package responses

type (
	SchoolResponse struct {
		Number           string                     `json:"number"`
		Name             string                     `json:"name"`
		CommencementDate string                     `json:"commencement_date"`
		SchoolCurriculum []SchoolCurriculumResponse `json:"school_curriculum"`
		Logo             string                     `json:"logo"`
		Motto            string                     `json:"motto"`
		PrincipalName    string                     `json:"principal_name"`
	}

	SchoolCurriculumResponse struct {
		Number string `json:"number"`
		Name   string `json:"name"`
	}

	ModeOfStudyResponse struct {
		Name              string `json:"name"`
		ModeOfStudyNumber string `json:"mode_of_study_number"`
		SchoolNumber      string `json:"school_number"`
		Status            string `json:"status"`
		CreatedAt         string `json:"created_at"`
		Code              int64  `json:"code"`
	}
)
