package responses

type (
	TopicDataList struct {
		Number      string `json:"number"`
		Name        string `json:"name"`
		Code        int64  `json:"code"`
		UserNumber  string `json:"user_number"`
		ClassNumber string `json:"class_number"`
		//StrandNumber  string `json:"strand_number"`
		//StrandName    string `json:"strand_name"`
		//SubjectNumber string `json:"subject_number"`
		//SubjectName   string `json:"subject_name"`
	}
)
