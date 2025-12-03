package requests

import (
	"errors"
	"html"
	"strings"
)

type PromotionSettingRequest struct {
	Name                      string `json:"name"`
	Status                    string `json:"status"`
	Type                      string `json:"type"`
	NextClassNumber           string `json:"next_class_number"`
	PreviousClassNumber       string `json:"previous_class_number"`
	StudyPeriodNumber         string `json:"study_period_number"`
	CurrentStudyPeriodNumber  string `json:"current_study_period_number"`
	CampusNumber              string `json:"campus_number"`
	CurriculumNumber          string `json:"curriculum_number"`
	CampusCurriculumNumber    string `json:"campus_curriculum_number"`
	CampusLevelNumber         string `json:"campus_level_number"`
	CurrentAcademicYearNumber string `json:"current_academic_year_number"`
	CurrentLevelNumber        string `json:"current_level_number"`
	NextLevelNumber           string `json:"next_level_number"`
	NextAcademicYearNumber    string `json:"next_academic_year_number"`
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

	if len(p.Status) > 90 {
		return errors.New("The status should be 900 characters maximum ")
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
