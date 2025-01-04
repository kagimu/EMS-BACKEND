package requests

import (
	"errors"
	"html"
	"strings"
)

type PromotionSettingRequest struct {
	Name                   string `json:"name"`
	Status                 string `json:"status"`
	Type                   string `json:"type"`
	NextClassNumber        string `json:"next_class_number"`
	PreviousClassNumber    string `json:"previous_class_number"`
	StudyPeriodNumber      string `json:"study_period_number"`
	CampusNumber           string `json:"campus_number"`
	CurriculumNumber       string `json:"curriculum_number"`
	CampusCurriculumNumber string `json:"campus_curriculum_number"`
}

func (p *PromotionSettingRequest) Prepare() {
	p.Name = html.EscapeString(strings.TrimSpace(p.Name))
	p.Status = html.EscapeString(strings.TrimSpace(p.Status))
	p.Type = html.EscapeString(strings.TrimSpace(p.Type))
	p.NextClassNumber = html.EscapeString(strings.TrimSpace(p.NextClassNumber))
	p.PreviousClassNumber = html.EscapeString(strings.TrimSpace(p.PreviousClassNumber))
}

func (p *PromotionSettingRequest) Validate() error {
	if p.Name == "" {
		return errors.New("name is required")
	}
	if p.Status == "" {
		return errors.New("status is required")
	}
	if p.Type == "" {
		return errors.New("type is required")
	}
	if p.StudyPeriodNumber == "" {
		return errors.New("choose the promotional study period")
	}
	if p.NextClassNumber == "" {
		return errors.New("next class number is required")
	}
	if p.PreviousClassNumber == "" {
		return errors.New("previous class number is required")
	}
	return nil
}
