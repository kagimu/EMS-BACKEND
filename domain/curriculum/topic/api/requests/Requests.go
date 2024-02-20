package requests

import (
	"errors"
	"html"
	"strings"
)

type TopicInfo struct {
	Code             string `json:"code"`
	Name             string `json:"name"`
	StrandNumber     string `json:"strand_number"`
	ClassNumber      string `json:"Class_number"`
	CurriculumNumber string `json:"curriculum_number"`
	TopicNumber      string `json:"topic_number"`
}

type TopicDelRequest struct {
	Number string `json:"number"`
}

func (p *TopicInfo) Prepare() {
	p.Name = html.EscapeString(strings.TrimSpace(p.Name))
}

func (p *TopicInfo) Validate() error {

	if p.Code == "" {
		return errors.New("code is required")
	}

	if p.Code == "0" {
		return errors.New("code must be greater than zero")
	}

	if p.Name == "" {
		return errors.New("name is required")
	}

	if p.StrandNumber == "" {
		return errors.New("subject is required")
	}

	if p.ClassNumber == "" {
		return errors.New("class is required")
	}

	if p.CurriculumNumber == "" {
		return errors.New("curriculum is required")
	}

	return nil
}
