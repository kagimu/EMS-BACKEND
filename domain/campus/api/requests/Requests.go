package requests

type CampusInfo struct {
	Name             string   `json:"name"`
	Phone            string   `json:"phone"`
	Email            string   `json:"email"`
	Location         string   `json:"location"`
	CommencementDate string   `json:"commencement_date"`
	SchoolNumber     string   `json:"school_number"`
	CampusCurriculum []string `json:"campus_curriculum"`
	CampusNumber     string   `json:"campus_number"`
}
type CampusLevelInfo struct {
	LevelNumber            string `json:"level_number"`
	CampusCurriculumNumber string `json:"campus_curriculum_number"`
}
