package models

import (
	"errors"
	"html"
	"strings"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"gorm.io/gorm"
)

type Parent struct {
	gorm.Model
	Number                     string      `json:"number" gorm:"size:36;unique;not null"`
	Surname                    string      `json:"surname"`
	SecondName                 string      `json:"second_name"`
	OtherName                  string      `json:"other_name"`
	MainTelephoneNumber        string      `json:"main_telephone_number" gorm:"unique;"`
	AlternativeTelephoneNumber string      `json:"alternative_telephone_number"`
	Email                      string      `json:"email" gorm:"unique;"`
	AlternativeEmail           string      `json:"alternative_email"`
	UserNumber                 string      `json:"user_number" gorm:"size:36;not null"`
	User                       models.User `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}

func (p *Parent) TableName() string {
	return "parents"
}

func (p *Parent) Validate() error {
	if p.Surname == "" {
		return errors.New("surname is required")
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

func (p *Parent) Prepare() {
	p.Surname = html.EscapeString(strings.TrimSpace(strings.ToLower(p.Surname)))
	p.SecondName = html.EscapeString(strings.TrimSpace(strings.ToLower(p.SecondName)))
	p.OtherName = html.EscapeString(strings.TrimSpace(strings.ToLower(p.OtherName)))
	p.MainTelephoneNumber = html.EscapeString(strings.TrimSpace(strings.ToLower(p.MainTelephoneNumber)))
	p.AlternativeTelephoneNumber = html.EscapeString(strings.TrimSpace(strings.ToLower(p.AlternativeTelephoneNumber)))
	p.Email = html.EscapeString(strings.TrimSpace(strings.ToLower(p.Email)))
	p.AlternativeEmail = html.EscapeString(strings.TrimSpace(strings.ToLower(p.AlternativeEmail)))
}
