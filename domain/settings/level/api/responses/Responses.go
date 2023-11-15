package responses

type (
	LearningLevelList struct {
		Number     string `json:"number"`
		Name       string `json:"name"`
		Code       int64  `json:"code"`
		UserNumber string `json:"user_number"`
	}
)
