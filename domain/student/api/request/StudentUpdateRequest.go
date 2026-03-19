package request

import (
	"html"
	"strings"
)

type StudentUpdateRequest struct {
	Surname                   string `json:"surname"`
	SecondName                string `json:"second_name"`
	OtherName                 string `json:"other_name"`
	DateOfBirth               string `json:"date_of_birth"`
	APNStatus                 string `json:"apn_status"`
	ReasonForChoosingSchool   string `json:"reason_for_choosing_school"`
	PreferredTypeOfEnrollment string `json:"preferred_type_of_enrollment"`
	SchoolNumber              string `json:"school_number"`
	Gender                    string `json:"gender"`
	ClassNumber               string `json:"class_number"`
	CurriculumNumber          string `json:"curriculum_number"`
	CampusCurriculumNumber    string `json:"campus_curriculum_number"`
	CampusNumber              string `json:"campus_number"`
	Lin                       string `json:"lin"`
	ModeOfStudyNumber         string `json:"mode_of_study_number"`
	ADN                       string `json:"adn"`
	Type                      string `json:"type"`
	PreferredCampusNumber     string `json:"preferred_campus_number"`
	SourceOfSchoolInfo        string `json:"source_of_school_info"`
}

func (s *StudentUpdateRequest) Prepare() {
	s.Surname = html.EscapeString(strings.TrimSpace(strings.ToLower(s.Surname)))
	s.SecondName = html.EscapeString(strings.TrimSpace(strings.ToLower(s.SecondName)))
	s.OtherName = html.EscapeString(strings.TrimSpace(strings.ToLower(s.OtherName)))
	s.ReasonForChoosingSchool = html.EscapeString(strings.TrimSpace(strings.ToLower(s.ReasonForChoosingSchool)))
	s.SecondName = html.EscapeString(strings.TrimSpace(strings.ToLower(s.SecondName)))
	s.OtherName = html.EscapeString(strings.TrimSpace(strings.ToLower(s.OtherName)))
}
