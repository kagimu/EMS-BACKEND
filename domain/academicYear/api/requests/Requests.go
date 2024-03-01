package requests

import (
	"errors"
	"html"
	"strings"
)

type AcademicYearDelRequest struct {
	Number string `json:"number"`
}

type AcademicYearInfo struct {
	AcademicYearName      string `json:"academic_year_name"`
	NoOfStudyPeriod       string `json:"no_of_study_period"`
	AcademicYearStartDate string `json:"academic_year_start_date"`
	AcademicYearEndDate   string `json:"academic_year_end_date"`
	LevelNumber           string `json:"level_number"`
	Code                  string `json:"code"`
}

func (p *AcademicYearInfo) Prepare() {
	p.AcademicYearName = html.EscapeString(strings.TrimSpace(p.AcademicYearName))
	p.NoOfStudyPeriod = html.EscapeString(strings.TrimSpace(p.NoOfStudyPeriod))
}

func (p *AcademicYearInfo) Validate() error {

	if p.Code == "" {
		return errors.New("code is required")
	}

	if p.Code == "0" {
		return errors.New("code must be greater than zero")
	}

	if p.AcademicYearName == "" {
		return errors.New("academic Year Name is required")
	}

	if p.NoOfStudyPeriod == "" {
		return errors.New("number Of Study Periods is required")
	}

	if p.AcademicYearStartDate == "" {
		return errors.New("start date is required")
	}

	if p.AcademicYearEndDate == "" {
		return errors.New("end date is required")
	}

	if p.LevelNumber == "" {
		return errors.New("level Year is required")
	}

	return nil
}
