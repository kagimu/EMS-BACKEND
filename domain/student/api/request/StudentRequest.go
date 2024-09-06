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
	ADNStatus                 string                  `json:"adn_status"`
	PreviousSchool            []models.PreviousSchool `json:"previous_school"`
	PreferredTypeOfEnrollment string                  `json:"preferred_type_of_enrollment"`
	PreferredCampusNumber     string                  `json:"preferred_campus_number"`
	Gender                    string                  `json:"gender"`
	ReasonForChoosingSchool   string                  `json:"reason_for_choosing_school"`
	SourceOfSchoolInfo        string                  `json:"source_of_school_info"`
	ApplicationNumber         string                  `json:"application_number"`
	ParentNumber              string                  `json:"parent_number"`
	ClassNumber               string                  `json:"class_number"`
	CurriculumNumber          string                  `json:"curriculum_number"`
	CampusCurriculumNumber    string                  `json:"campus_curriculum_number"`
	CampusNumber              string                  `json:"campus_number"`
	SchoolNumber              string                  `json:"school_number"`
	Lin                       string                  `json:"lin"`
	ModeOfStudyNumber         string                  `json:"mode_of_study_number"`
}

func (s *StudentRequest) Validate() error {
	if s.Surname == "" {
		return errors.New("surname is required")
	}
	if s.CampusNumber == "" {
		return errors.New("choose campus")
	}
	if s.ClassNumber == "" {
		return errors.New("choose class")
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
	if s.ReasonForChoosingSchool == "" {
		return errors.New("reason is required")
	}

	if s.ApplicationNumber == " " {
		return errors.New("application number is required")
	}

	if s.SourceOfSchoolInfo == "" {
		return errors.New("source of school information is required")
	}
	if s.SchoolNumber == "" {
		return errors.New("school number is required to proceed")
	}

	return nil
}

func (s *StudentRequest) Prepare() {
	s.Surname = html.EscapeString(strings.TrimSpace(strings.ToLower(s.Surname)))
	s.SecondName = html.EscapeString(strings.TrimSpace(strings.ToLower(s.SecondName)))
	s.OtherName = html.EscapeString(strings.TrimSpace(strings.ToLower(s.OtherName)))
	s.ReasonForChoosingSchool = html.EscapeString(strings.TrimSpace(strings.ToLower(s.ReasonForChoosingSchool)))
	s.Gender = html.EscapeString(strings.TrimSpace(strings.ToLower(s.Gender)))
}
