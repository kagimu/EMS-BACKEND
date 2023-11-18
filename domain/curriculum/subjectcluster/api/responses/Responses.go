package responses

import "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/subject/api/responses"

type (
	SubjectClusterInfo struct {
		Number             string                      `json:"number"`
		Name               string                      `json:"name"`
		Code               int64                       `json:"code"`
		LearningAreaNumber string                      `json:"learning_area_number"`
		Subject            []responses.SubjectInfoList `json:"subject"`
	}
)
