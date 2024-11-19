package responses

type CurriculumDeliveryResponse struct {
	Number                string   `json:"number"`
	Steps                 []Step   `json:"steps"`
	Stream                string   `json:"stream"`
	Class                 string   `json:"class"`
	Date                  string   `json:"date"`
	Time                  string   `json:"time"`
	Teacher               string   `json:"teacher"`
	NumberOfStudents      int      `json:"number_of_students"`
	Comment               string   `json:"comment"`
	StreamNumber          string   `json:"stream_number"`
	LessonPlanStepNumbers []string `json:"lesson_plan_step_number"`
	TeacherNumber         string   `json:"teacher_number"`
}

type Step struct {
	Name     string `json:"name"`
	Concept  string `json:"concept"`
	Topic    string `json:"topic"`
	Subject  string `json:"subject"`
	SubTopic string `json:"sub_topic"`
	Date     string `json:"date"`
	Time     string `json:"time"`
	Comment  string `json:"comment"`
}
