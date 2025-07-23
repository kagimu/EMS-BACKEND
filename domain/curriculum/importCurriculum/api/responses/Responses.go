package responses

import (
	"errors"
	"html"
	"strings"
)

type (
	CurriculumResponse struct {
		Name             string `json:"name"`
		CurriculumNumber string `json:"curriculum_number"`
		Code             int64  `json:"code"`
	}

	SchoolCurriculumDataList struct {
		Number string `json:"number"`
		Name   string `json:"name"`
	}
)

func (p *CurriculumResponse) Prepare() {
	p.Name = html.EscapeString(strings.TrimSpace(p.Name))
}

func (p *CurriculumResponse) Validate() error {

	if p.Name == "" {
		return errors.New("name is required")
	}
	if p.Code <= 0 {
		return errors.New("code must be a positive integer")
	}

	return nil
}
