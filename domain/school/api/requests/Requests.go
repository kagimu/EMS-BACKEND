package requests

type SchoolInfo struct {
	Name             string   `json:"name"`
	CommencementDate string   `json:"commencement_date"`
	Curriculum       []string `json:"curriculum"`
}
