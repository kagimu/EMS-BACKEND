package requests

import (
	"errors"
	"html"
	"strings"
)

type PromotionRequest struct {
	CurrentAcademicYearNumber string `json:"current_academic_year_number"`
	Status                    string `json:"status"`
	PromotionSettingNumber    string `json:"promotion_setting_number"`
	StudentNumber             string `json:"student_number"`
	Reason                    string `json:"reason"`
	UserNumber                string `json:"user_number"`
}

func (p *PromotionRequest) Prepare() {
	p.Status = html.EscapeString(strings.TrimSpace(p.Status))
	p.PromotionSettingNumber = html.EscapeString(strings.TrimSpace(p.PromotionSettingNumber))
	p.StudentNumber = html.EscapeString(strings.TrimSpace(p.StudentNumber))
	p.Reason = html.EscapeString(strings.TrimSpace(p.Reason))
	p.CurrentAcademicYearNumber = html.EscapeString(strings.TrimSpace(p.CurrentAcademicYearNumber))
	p.UserNumber = html.EscapeString(strings.TrimSpace(p.UserNumber))
}

func (p *PromotionRequest) Validate() error {
	if p.CurrentAcademicYearNumber == "" {
		return errors.New("current academic year is required")
	}
	if p.Status == "" {
		return errors.New("status is required")
	}
	if p.PromotionSettingNumber == "" {
		return errors.New("promotion setting number is required")
	}
	if p.StudentNumber == "" {
		return errors.New("student number is required")
	}

	return nil
}
