package responses

import (
	"errors"
	"html"
	"strings"
)

type (
	CurriculumInfo struct {
		Name             string `json:"name"`
		CurriculumNumber string `json:"curriculum_number"`
	}

	SchoolCurriculumDataList struct {
		Number string `json:"number"`
		Name   string `json:"name"`
	}
)

func (p *CurriculumInfo) Prepare() {
	p.Name = html.EscapeString(strings.TrimSpace(p.Name))
}

func (p *CurriculumInfo) Validate() error {

	if p.Name == "" {
		return errors.New("name is required")
	}

	return nil
}
