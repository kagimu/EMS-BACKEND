package responses

import "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/curriculum/subjectcluster/api/responses"

type (
	LearningAreaInfoData struct {
		Number         string                         `json:"number"`
		Name           string                         `json:"name"`
		Code           int64                          `json:"code"`
		UserNumber     string                         `json:"user_number"`
		SubjectCluster []responses.SubjectClusterInfo `json:"subject_cluster"`
	}
)
