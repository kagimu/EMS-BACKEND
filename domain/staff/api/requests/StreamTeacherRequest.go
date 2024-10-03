package requests

import "errors"

type StreamTeacherRequest struct {
	TeacherIds           []string `json:"teacher_ids"`
	Streams              []string `json:"streams"`
	ClassNumber          string   `json:"class_number"`
	AssessmentAreaNumber string   `json:"assessment_area_number"`
	SchoolNumber         string   `json:"school_number"`
	CurriculumNumber     string   `json:"curriculum_number"`
	CampusNumber         string   `json:"campus_number"`
}

func (s *StreamTeacherRequest) Validate() error {
	if len(s.TeacherIds) == 0 {
		return errors.New("please select at-least teacher/staff member to attach to stream")
	}
	if len(s.Streams) == 0 {
		return errors.New("stream number is required")
	}
	if s.ClassNumber == "" {
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
