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
	CampusClassNumber         string `json:"campus_class_number"`
	SchoolNumber              string `json:"school_number"`
	Gender                    string `json:"gender"`
}

func (s *StudentUpdateRequest) Prepare() {
	s.Surname = html.EscapeString(strings.TrimSpace(strings.ToLower(s.Surname)))
	s.SecondName = html.EscapeString(strings.TrimSpace(strings.ToLower(s.SecondName)))
	s.OtherName = html.EscapeString(strings.TrimSpace(strings.ToLower(s.OtherName)))
	s.ReasonForChoosingSchool = html.EscapeString(strings.TrimSpace(strings.ToLower(s.ReasonForChoosingSchool)))
	s.SecondName = html.EscapeString(strings.TrimSpace(strings.ToLower(s.SecondName)))
	s.OtherName = html.EscapeString(strings.TrimSpace(strings.ToLower(s.OtherName)))
}
