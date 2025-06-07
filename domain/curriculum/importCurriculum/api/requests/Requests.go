package requests

import (
	"errors"
	"html"
	"strings"
)

type CurriculumDelRequest struct {
	Number string `json:"number"`
}

type CurriculumCreateRequest struct {
	Name string `json:"name"`
}

func (p *CurriculumCreateRequest) Prepare() {
	p.Name = html.EscapeString(strings.TrimSpace(p.Name))
}

func (p *CurriculumCreateRequest) Validate() error {

	if p.Name == "" {
		return errors.New("name is required")
	}

	return nil
}

func (p *CurriculumDelRequest) ValidateDelete() error {

	if p.Number == "" {
		return errors.New("curriculum number is required")
	}

	return nil
}
