package requests

import (
	"errors"
	"html"
	"strings"
)

type ClassInfo struct {
	Name         string `json:"name"`
	Code         string `json:"code"`
	Abbreviation string `json:"abbreviation"`
	Level        string `json:"level"`
	ClassNumber  string `json:"class_number"`
}

type StreamRequest struct {
	Name             string `json:"name"`
	ClassLevelNumber string `json:"class_level_number"`
}

type ClassDelRequest struct {
	Number string `json:"number"`
}

func (p *ClassInfo) Prepare() {
	p.Name = html.EscapeString(strings.TrimSpace(p.Name))
	p.Abbreviation = html.EscapeString(strings.TrimSpace(p.Abbreviation))
}

func (p *ClassInfo) Validate() error {

	if p.Code == "" {
		return errors.New("code is required")
	}

	if p.Code == "0" {
		return errors.New("code must be greater than zero")
	}

	if p.Name == "" {
		return errors.New("name is required")
	}

	if p.Level == "" {
		return errors.New("level is required")
	}

	return nil
}
