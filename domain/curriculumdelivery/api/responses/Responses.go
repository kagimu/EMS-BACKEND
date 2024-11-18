package responses

type CurriculumDeliveryResponse struct {
	Number           string `json:"number"`
	Step             string `json:"step"`
	Concept          string `json:"concept"`
	Topic            string `json:"topic"`
	Subject          string `json:"subject"`
	SubTopic         string `json:"sub_topic"`
	Stream           string `json:"stream"`
	Class            string `json:"class"`
	Date             string `json:"date"`
	Time             string `json:"time"`
	Teacher          string `json:"teacher"`
	NumberOfStudents int    `json:"number_of_students"`
	Comment          string `json:"comment"`
}
