package requests

type SubjectInfo struct {
	Code                 string `json:"code"`
	Name                 string `json:"name"`
	SubjectClusterNumber string `json:"subject_cluster_number"`
	CurriculumNumber     string `json:"curriculum_number"`
	SubjectNumber        string `json:"subject_number"`
}

type SubjectDelRequest struct {
	Number string `json:"number"`
}
