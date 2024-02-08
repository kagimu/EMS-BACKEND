package requests

import (
	"errors"
	"html"
	"strings"
)

type SubjectClusterInfo struct {
	Code             string `json:"code"`
	Name             string `json:"name"`
	SubjectNumber    string `json:"subject_number"`
	CurriculumNumber string `json:"curriculum_number"`
	StrandNumber     string `json:"strand_number"`
}

type StrandDelRequest struct {
	Number string `json:"number"`
}

func (p *SubjectClusterInfo) Prepare() {
	p.Name = html.EscapeString(strings.TrimSpace(p.Name))
}

func (p *SubjectClusterInfo) Validate() error {

	if p.Code == "" {
		return errors.New("code is required")
	}

	if p.Code == "0" {
		return errors.New("code must be greater than zero")
	}

	if p.Name == "" {
		return errors.New("name is required")
	}

	if p.SubjectNumber == "" {
		return errors.New("subjectNumber is required")
	}

	return nil
}
