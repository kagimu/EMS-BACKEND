package responses

import "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/subtopic/api/responses"

type (
	TopicDataList struct {
		Number                string                   `json:"number"`
		Name                  string                   `json:"name"`
		Code                  int64                    `json:"code"`
		UserNumber            string                   `json:"user_number"`
		ClassNumber           string                   `json:"class_number"`
		CurriculumClassNumber string                   `json:"curriculum_class_number"`
		SubTopic              []responses.SubTopicList `json:"subtopic"`
		StrandNumber          string                   `json:"strand_number"`
		//StrandName    string `json:"strand_name"`
		//SubjectNumber string `json:"subject_number"`
		//SubjectName   string `json:"subject_name"`
	}
)
