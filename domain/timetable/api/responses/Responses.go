package responses

type Subjects struct {
	Number string `json:"number"`
	Name   string `json:"name"`
}

type Classes struct {
	Number string `json:"number"`
	Name   string `json:"name"`
}

type StudyTimetables struct {
	Name   string `json:"name"`
	Number string `json:"number"`
	Status string `json:"status"`
}
