package requests

import (
	"errors"
	"html"
	"strings"
)

type PromotionSettingValue struct {
	Type                   string  `json:"type"`
	Score                  float64 `json:"score"`
	PromotionSettingNumber string  `json:"promotion_setting_number"`
	GradingStandardNumber  string  `json:"grading_standard_number"`
}

func (p *PromotionSettingValue) Prepare() {
	p.Type = html.EscapeString(strings.TrimSpace(p.Type))
	p.PromotionSettingNumber = html.EscapeString(strings.TrimSpace(p.PromotionSettingNumber))
}

func (p *PromotionSettingValue) Validate() error {
	if p.Type == "" {
		return errors.New("type is required")
	}
	if p.PromotionSettingNumber == "" {
		return errors.New("promotion setting number is required")
	}

	return nil
}
