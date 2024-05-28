package requests

import "errors"

type ScopeRequest struct {
	UserNumber       string `json:"user_number"`
	CampusNumber     string `json:"campus_number"`
	SchoolNumber     string `json:"school_number"`
	CurriculumNumber string `json:"curriculum_number"`
	ClassNumber      string `json:"class_number"`
	StreamNumber     string `json:"stream_number"`
}

func (s *ScopeRequest) Validate() error {
	if s.UserNumber == "" {
		return errors.New("user number is required")
	}

	return nil
}
