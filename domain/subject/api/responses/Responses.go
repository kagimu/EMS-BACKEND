package responses

import "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/strand/api/responses"

type (
	SubjectInfoList struct {
		Number               string                     `json:"number"`
		Name                 string                     `json:"name"`
		Code                 int64                      `json:"code"`
		UserNumber           string                     `json:"user_number"`
		SubjectClusterNumber string                     `json:"subject_cluster_number"`
		Strand               []responses.StrandDataList `json:"strand"`
	}
)
