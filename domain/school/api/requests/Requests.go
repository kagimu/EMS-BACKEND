package requests

type SchoolInfo struct {
	Name             string   `json:"name"`
	CommencementDate string   `json:"commencement_date"`
	SchoolNumber     string   `json:"school_number"`
	Curriculum       []string `json:"curriculum"`
}
