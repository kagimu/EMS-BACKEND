package requests

import (
	"errors"
	"fmt"
	"html"
	"strings"
)

type PromotionSettingValueRequest struct {
	Score                      float64                       `json:"score"`
	MinimumPoints              float64                       `json:"minimum_points"`
	PromotionSettingNumber     string                        `json:"promotion_setting_number"`
	AssessmentGradeScaleNumber string                        `json:"assessment_grade_scale_number"`
	AssessmentAreas            []AssessmentAreaRequestValues `json:"assessment_areas"`
}

type AssessmentAreaRequestValues struct {
	GradeNumber          string `json:"grade_number"`
	AssessmentAreaNumber string `json:"assessment_area_number"`
}

func (p *PromotionSettingValueRequest) Prepare() {
	p.AssessmentGradeScaleNumber = html.EscapeString(strings.TrimSpace(p.AssessmentGradeScaleNumber))
	p.PromotionSettingNumber = html.EscapeString(strings.TrimSpace(p.PromotionSettingNumber))
}

func (p *PromotionSettingValueRequest) Validate() error {
	if p.AssessmentGradeScaleNumber == "" {
		return errors.New("Choose assessment grade scale number ")
	}
	if p.PromotionSettingNumber == "" {
		return errors.New("promotion setting number is required")
	}

	if p.Score < 0 {
		return errors.New("Invalid score ")
	}

	if p.MinimumPoints < 0 {
		return errors.New("Invalid minimum aggregate ")
	}

	for index, area := range p.AssessmentAreas {

		if area.GradeNumber == "" {
			return errors.New(fmt.Sprintf("Choose grade at index %d", index+1))
		}

		if area.AssessmentAreaNumber == p.PromotionSettingNumber {
			return errors.New(fmt.Sprintf("Choose assessment area at index %d", index+1))
		}
	}

	return nil
}
