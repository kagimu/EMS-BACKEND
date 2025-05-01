package request

import (
	"html"
	"strings"
)

type StudentUpdateRequest struct {
	Surname     string `json:"surname"`
	SecondName  string `json:"second_name"`
	OtherName   string `json:"other_name"`
	DateOfBirth string `json:"date_of_birth"`
	Gender      string `json:"gender"`
}

func (s *StudentUpdateRequest) Prepare() {
	s.Surname = html.EscapeString(strings.TrimSpace(strings.ToLower(s.Surname)))
	s.SecondName = html.EscapeString(strings.TrimSpace(strings.ToLower(s.SecondName)))
	s.OtherName = html.EscapeString(strings.TrimSpace(strings.ToLower(s.OtherName)))
	s.SecondName = html.EscapeString(strings.TrimSpace(strings.ToLower(s.SecondName)))
	s.OtherName = html.EscapeString(strings.TrimSpace(strings.ToLower(s.OtherName)))
}
