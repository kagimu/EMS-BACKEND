package requests

import "errors"

type StreamTeacherRequest struct {
	UserNumber           string `json:"user_number"`
	StreamNumber         string `json:"stream_number"`
	CampusClassNumber    string `json:"campus_class_number"`
	AssessmentAreaNumber string `json:"assessment_area_number"`
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

	return nil
}
