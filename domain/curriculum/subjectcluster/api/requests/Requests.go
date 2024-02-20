package requests

type SubjectClusterData struct {
	Code                 string `json:"code"`
	Name                 string `json:"name"`
	LearningAreaNumber   string `json:"learning_area_number"`
	CurriculumNumber     string `json:"curriculum_number"`
	SubjectClusterNumber string `json:"subject_cluster_number"`
}

type SubjectClusterDelRequest struct {
	Number string `json:"number"`
}
