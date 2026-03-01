package responses

type AdmissionLetterResponse struct {
	Number              string `json:"number"`
	SchoolNumber        string `json:"school_number"`
	StudentNumber       string `json:"student_number"`
	CampusNumber        string `json:"campus_number"`
	CampusClassNumber   string `json:"campus_class_number"`
	Phrase              string `json:"phrase"`
	PrincipalSignature  string `json:"principal_signature"`
	CreatedAt           string `json:"created_at"`
	UpdatedAt           string `json:"updated_at"`
}
