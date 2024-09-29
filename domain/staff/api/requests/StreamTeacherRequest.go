package requests

import "errors"

type StreamTeacherRequest struct {
	UserNumber           string `json:"user_number"`
	StreamNumber         string `json:"stream_number"`
	CampusClassNumber    string `json:"campus_class_number"`
	AssessmentAreaNumber string `json:"assessment_area_number"`
	SchoolNumber         string `json:"school_number"`
	CurriculumNumber     string `json:"curriculum_number"`
	CampusNumber         string `json:"campus_number"`
}

func (s *StreamTeacherRequest) Validate() error {
	if s.UserNumber == "" {
		return errors.New("user number is required")
	}
	if s.StreamNumber == "" {
		return errors.New("stream number is required")
	}
	if s.CampusClassNumber == "" {
		return errors.New("campus class number is required")
	}

	if s.AssessmentAreaNumber == "" {
		return errors.New("assessment area number is required")
	}

	if s.SchoolNumber == "" {
		return errors.New("school number is required")
	}

	if s.CurriculumNumber == "" {
		return errors.New("curriculum number is required")
	}

	return nil
}
