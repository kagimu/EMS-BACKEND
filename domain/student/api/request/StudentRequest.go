package request

import (
	"errors"
	"html"
	"strings"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/previousSchool/models"
)

type StudentRequest struct {
	Surname                   string                  `json:"surname"`
	SecondName                string                  `json:"second_name"`
	OtherName                 string                  `json:"other_name"`
	DateOfBirth               string                  `json:"date_of_birth"`
	APN                       string                  `json:"apn"`
	PreviousSchool            []models.PreviousSchool `json:"previous_school"`
	PreferredTypeOfEnrollment string                  `json:"preferred_type_of_enrollment"`
	PreferredCampusNumber     string                  `json:"preferred_campus_number"`
	Gender                    string                  `json:"gender"`
	Reason                    string                  `json:"reason"`
	SourceOfSchoolInfoNumber  string                  `json:"source_of_school_info_number"`
	ParentNumber              string                  `json:"parent_number"`
	ClassNumber               string                  `json:"class_number"`
}

func (s *StudentRequest) Validate() error {
	if s.Surname == "" {
		return errors.New("surname is required")
	}
	if s.SecondName == "" {
		return errors.New("second name is required")
	}
	if s.DateOfBirth == "" {
		return errors.New("date of birth is required")
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
