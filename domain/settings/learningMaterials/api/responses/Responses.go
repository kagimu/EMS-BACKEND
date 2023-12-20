package responses

type (
	LearningMaterialList struct {
		Number                 string `json:"number"`
		Code                   int64  `json:"code"`
		Description            string `json:"description"`
		LearningMaterialNumber string `json:"learning_material_number"`
		SubTopicNumber         string `json:"sub_topic_number"`
		UserNumber             string `json:"user_number"`
	}
)
