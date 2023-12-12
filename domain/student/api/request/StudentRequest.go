package request

import (
	"errors"
	"html"
	"strings"
)

type StudentRequest struct {
	Surname         string   `json:"surname"`
	SecondName      string   `json:"second_name"`
	OtherName       string   `json:"other_name"`
	Age             uint     `json:"age"`
	APN             string   `json:"apn"`
	PreviousSchools []string `json:"previous_schools"`
	Gender          string   `json:"gender"`
	Reason          string   `json:"reason"`
}

func (s *StudentRequest) Validate() error {
	if s.Surname == "" {
		return errors.New("surname is required")
	}

	if s.SecondName == "" {
		return errors.New("second name is required")
	}
	if s.Age == 0 {
		return errors.New("age must be greater than 0")
	}
	if s.Gender == "" {
		return errors.New("gender is required")
	}
	if s.Reason == "" {
		return errors.New("reason is required")
	}

	return nil
}

func (s *StudentRequest) Prepare() {
	s.Surname = html.EscapeString(strings.TrimSpace(strings.ToLower(s.Surname)))
	s.SecondName = html.EscapeString(strings.TrimSpace(strings.ToLower(s.SecondName)))
	s.OtherName = html.EscapeString(strings.TrimSpace(strings.ToLower(s.OtherName)))
	s.Reason = html.EscapeString(strings.TrimSpace(strings.ToLower(s.Reason)))
	s.Gender = html.EscapeString(strings.TrimSpace(strings.ToLower(s.Gender)))
}
