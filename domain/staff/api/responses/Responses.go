package responses

import "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/responses"

type SingleStaffResponse struct {
	Number string `json:"number"`
	responses.UserResponse
	Classes []string `json:"classes"`
	Streams []string `json:"streams"`
}

type TeacherResponse struct {
	Number               string `json:"number"`
	FirstName            string `json:"first_name"`
	SecondName           string `json:"second_name"`
	LastName             string `json:"last_name"`
	Gender               string `json:"gender"`
	Type                 string `json:"type"`
	Status               string `json:"status"`
	CurriculumNumber     string `json:"curriculum_number"`
	CurriculumName       string `json:"curriculum_name"`
	CampusNumber         string `json:"campus_number"`
	CampusName           string `json:"campus_name"`
	SchoolNumber         string `json:"school_number"`
	SchoolName           string `json:"school_name"`
	StreamNumber         string `json:"stream_number"`
	StreamName           string `json:"stream_name"`
	CampusClassNumber    string `json:"campus_class_number"`
	CampusClassName      string `json:"campus_class_name"`
	AssessmentAreaNumber string `json:"assessment_area_number"`
	AssessmentAreaName   string `json:"assessment_area_name"`
	AddedBy              string `json:"added_by"`
	CreatedAt            string `json:"created_at"`
}
