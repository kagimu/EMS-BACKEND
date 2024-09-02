package responses

type (
	CampusClassInfo struct {
		Number                string       `json:"number"`
		ClassNumber           string       `json:"class_number"`
		Name                  string       `json:"name"`
		Code                  int64        `json:"code"`
		Abbreviation          string       `json:"abbreviation"`
		CurriculumClassNumber string       `json:"curriculum_class_number"`
		Streams               []StreamInfo `json:"streams"`
	}
	StreamInfo struct {
		StreamNumber      string `json:"stream_number"`
		Name              string `json:"name"`
		CampusClassNumber string `json:"campus_class_number"`
	}

	ClassDataList struct {
		Number       string `json:"number"`
		Name         string `json:"name"`
		Code         int64  `json:"code"`
		Abbreviation string `json:"abbreviation"`
		UserNumber   string `json:"user_number"`
		ClassNumber  string `json:"class_number"`
		LevelNumber  string `json:"level_number"`
	}

	LearningLevelDataList struct {
		Number        string          `json:"number"`
		Name          string          `json:"name"`
		Code          int64           `json:"code"`
		UserNumber    string          `json:"user_number"`
		ClassDataList []ClassDataList `json:"class"`
	}
)
