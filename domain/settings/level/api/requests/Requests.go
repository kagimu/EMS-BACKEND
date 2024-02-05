package requests

import (
	"errors"
	"html"
	"strings"
)

type LevelInfo struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	LevelNumber string `json:"level_number"`
}

type LevelDelRequest struct {
	Number string `json:"number"`
}

func (p *LevelInfo) Prepare() {
	p.Name = html.EscapeString(strings.TrimSpace(p.Name))
}

func (p *LevelInfo) Validate() error {

	if p.Name == "" {
		return errors.New("name is required")
	}

	if p.Code == "" {
		return errors.New("code is required")
	}

	if p.Code == "0" {
		return errors.New("code must be greater than zero")
	}

	if p.LevelNumber == "" {
		return errors.New("level is required")
	}

	return nil
}
