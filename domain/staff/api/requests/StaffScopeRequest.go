package requests

import "errors"

type StaffScopeRequest struct {
	StaffNumber  string             `json:"staff_number"`
	SchoolNumber string             `json:"school_number"`
	Campuses     []StaffScopeCampus `json:"campuses"`
}

type StaffScopeCampus struct {
	CampusNumber string                 `json:"campus_number"`
	Curriculums  []StaffScopeCurriculum `json:"curriculums"`
}

type StaffScopeCurriculum struct {
	CurriculumNumber string            `json:"curriculum_number"`
	Levels           []StaffScopeLevel `json:"levels"`
}

type StaffScopeLevel struct {
	CampusLevelNumber string            `json:"campus_level_number"`
	Classes           []StaffScopeClass `json:"classes"`
}

type StaffScopeClass struct {
	ClassNumber string             `json:"class_number"`
	Streams     []StaffScopeStream `json:"streams"`
}

type StaffScopeStream struct {
	StreamNumber    string   `json:"stream_number"`
	AssessmentAreas []string `json:"assessment_areas"`
}

func (s *StaffScopeRequest) Validate() error {
	if s.StaffNumber == "" {
		return errors.New("staff number is required")
	}
	if s.SchoolNumber == "" {
		return errors.New("school number is required")
	}
	return nil
}
