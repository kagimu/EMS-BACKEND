package requests

import (
	"errors"
	"html"
	"strings"
)

type ParentRequest struct {
	Surname                    string `json:"surname"`
	SecondName                 string `json:"second_name"`
	OtherName                  string `json:"other_name"`
	MainTelephoneNumber        string `json:"main_telephone_number"`
	AlternativeTelephoneNumber string `json:"alternative_telephone_number"`
	Email                      string `json:"email"`
	AlternativeEmail           string `json:"alternative_email"`
	ApplicationNumber          string `json:"application_number"`
}

func (p *ParentRequest) Validate() error {
	if p.Surname == "" {
		return errors.New("surname is required")
	}
	if p.ApplicationNumber == "" {
		return errors.New("application number is required")
	}
	if p.SecondName == "" {
		return errors.New("second name is required")
	}
	if p.MainTelephoneNumber == "" {
		return errors.New("main telephone number is required")
	}
	if p.Email == "" {
		return errors.New("email is required")
	}

	return nil
}

func (p *ParentRequest) Prepare() {
	p.Surname = html.EscapeString(strings.TrimSpace(strings.ToLower(p.Surname)))
	p.SecondName = html.EscapeString(strings.TrimSpace(strings.ToLower(p.SecondName)))
	p.OtherName = html.EscapeString(strings.TrimSpace(strings.ToLower(p.OtherName)))
	p.MainTelephoneNumber = html.EscapeString(strings.TrimSpace(strings.ToLower(p.MainTelephoneNumber)))
	p.AlternativeTelephoneNumber = html.EscapeString(strings.TrimSpace(strings.ToLower(p.AlternativeTelephoneNumber)))
	p.Email = html.EscapeString(strings.TrimSpace(strings.ToLower(p.Email)))
	p.AlternativeEmail = html.EscapeString(strings.TrimSpace(strings.ToLower(p.AlternativeEmail)))
}
