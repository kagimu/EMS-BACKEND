package requests

import (
	"errors"
	"html"
	"strings"
)

type LevelInfo struct {
	Code         int64  `json:"code"`
	Name         string `json:"name"`
	LevelNumber  string `json:"level_number"`
	EntityNumber string `json:"entity_number"`
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

	if p.Code <= 0 {
		return errors.New("code must be greater than zero")
	}

	if p.EntityNumber == "" {
		return errors.New("entity number is required")
	}

	return nil
}
