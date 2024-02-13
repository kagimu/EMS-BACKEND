package requests

import (
	"errors"
	"html"
	"strings"
)

type ConceptInfo struct {
	Code             string `json:"code"`
	Name             string `json:"name"`
	SubTopicNumber   string `json:"sub_topic_number"`
	CurriculumNumber string `json:"curriculum_number"`
	ConceptNumber    string `json:"concept_number"`
}

type ConceptDelRequest struct {
	Number string `json:"number"`
}

func (p *ConceptInfo) Prepare() {
	p.Name = html.EscapeString(strings.TrimSpace(p.Name))
}

func (p *ConceptInfo) Validate() error {

	if p.Code == "" {
		return errors.New("code is required")
	}

	if p.Code == "0" {
		return errors.New("code must be greater than zero")
	}

	if p.Name == "" {
		return errors.New("name is required")
	}

	if p.CurriculumNumber == "" {
		return errors.New("curriculum is required")
	}

	if p.SubTopicNumber == "" {
		return errors.New("sub topic is required")
	}

	return nil
}
