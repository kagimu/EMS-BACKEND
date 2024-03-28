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
	ADN                       string `json:"adn"`
	PreferredTypeOfEnrollment string `json:"preferred_type_of_enrollment"`
	ReasonForChoosingSchool   string `json:"reason_for_choosing_school"`
	ReasonForAdmissionDenial  string `json:"reason_for_admission_denial"`
	SourceOfSchoolInfo        string `json:"source_of_school_info"`
	ParentNumber              string `json:"parent_number"`
	CampusClassNumber         string `json:"campus_class_number"`
}

func (s *StudentUpdateRequest) Prepare() {
	s.Surname = html.EscapeString(strings.TrimSpace(strings.ToLower(s.Surname)))
	s.SecondName = html.EscapeString(strings.TrimSpace(strings.ToLower(s.SecondName)))
	s.OtherName = html.EscapeString(strings.TrimSpace(strings.ToLower(s.OtherName)))
	s.ReasonForChoosingSchool = html.EscapeString(strings.TrimSpace(strings.ToLower(s.ReasonForChoosingSchool)))
	s.ReasonForAdmissionDenial = html.EscapeString(strings.TrimSpace(strings.ToLower(s.ReasonForAdmissionDenial)))
	s.SourceOfSchoolInfo = html.EscapeString(strings.TrimSpace(strings.ToLower(s.SourceOfSchoolInfo)))
	s.SecondName = html.EscapeString(strings.TrimSpace(strings.ToLower(s.SecondName)))
	s.OtherName = html.EscapeString(strings.TrimSpace(strings.ToLower(s.OtherName)))
}
