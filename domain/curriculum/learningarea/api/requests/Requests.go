package requests

import (
	"errors"
	"html"
	"strings"
)

type LearningAreaInfo struct {
	Name               string `json:"name"`
	Code               string `json:"code"`
	CurriculumNumber   string `json:"curriculum_number"`
	LearningAreaNumber string `json:"learning_area_number"`
}

type LearningAreaDelRequest struct {
	Number string `json:"number"`
}

func (p *LearningAreaInfo) Prepare() {
	p.Name = html.EscapeString(strings.TrimSpace(p.Name))
	p.Code = html.EscapeString(strings.TrimSpace(p.Code))
	p.CurriculumNumber = html.EscapeString(strings.TrimSpace(p.CurriculumNumber))
}

func (p *LearningAreaInfo) Validate() error {

	if p.Name == "" {
		return errors.New("name is required")
	}

	if p.Code == "" {
		return errors.New("code is required")
	}

	if p.Code == "0" {
		return errors.New("code must be greater than zero")
	}

	if p.CurriculumNumber == "" {
		return errors.New("curriculum is required")
	}

	return nil
}
