package responses

import "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/strand/api/responses"

type (
	SubjectInfoList struct {
		Number               string                     `json:"number"`
		Name                 string                     `json:"name"`
		Code                 int64                      `json:"code"`
		UserNumber           string                     `json:"user_number"`
		SubjectClusterNumber string                     `json:"subject_cluster_number"`
		LearningAreaNumber   string                     `json:"learning_area_number"`
		Strand               []responses.StrandDataList `json:"strand"`
	}

	SubjectResponse struct {
		Name          string `json:"name"`
		SubjectNumber string `json:"subject_number"`
	}
)
