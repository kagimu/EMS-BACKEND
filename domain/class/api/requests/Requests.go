package requests

import (
	"errors"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/commons/constants"
	"html"
	"strings"
)

type ClassInfo struct {
	Name         string `json:"name"`
	Code         string `json:"code"`
	Abbreviation string `json:"abbreviation"`
	Level        string `json:"level"`
	ClassNumber  string `json:"class_number"`
	EntityNumber string `json:"entity_number"`
}

type StreamRequest struct {
	Name              string `json:"name"`
	CampusClassNumber string `json:"campus_class_number"`
}

type UpdateStreamRequest struct {
	StreamRequest
	StreamNumber string `json:"stream_number"`
	Status       string `json:"status"`
}

func (r UpdateStreamRequest) Validate() error {

	if r.StreamNumber == "" {
		return errors.New("choose stream to edit")
	}

	if r.Name == "" {
		return errors.New("name is required")
	}

	if r.CampusClassNumber == "" {
		return errors.New("campus class number is required")
	}

	if r.Status == "" {
		return errors.New("status is required")
	}

	if !commons.StringContains([]string{constants.StatusActive, constants.StatusInactive}, r.Status) {
		return errors.New("invalid status")
	}

	return nil
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

	if p.EntityNumber == "" {
		return errors.New("entity number is required")
	}

	return nil
}
