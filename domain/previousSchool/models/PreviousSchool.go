package models

import (
	"errors"
	"html"
	"strings"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"gorm.io/gorm"
)

type PreviousSchool struct {
	gorm.Model
	Number            string      `json:"number" gorm:"size:36;unique;not null"`
	Name              string      `json:"name" gorm:"size:50;not null"`
	StudentNumber     string      `json:"student_number" gorm:"size:36;unique;not null"`
	StartYear         string      `json:"start_period" gorm:"size:10;not null"`
	ClassNumber       string      `json:"class_number" gorm:"size:36;index;not null"`
	ApplicationNumber string      `json:"application_number" gorm:"size:36;unique;not null"`
	EndYear           string      `json:"end_period" gorm:"size:10;not null"`
	UserNumber        string      `json:"user_number" gorm:"size:36;index;not null"`
	User              models.User `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}

func (p *PreviousSchool) Prepare() {
	p.Name = html.EscapeString(strings.TrimSpace(strings.ToLower(p.Name)))
	p.StartYear = html.EscapeString(strings.TrimSpace(strings.ToLower(p.StartYear)))
	p.EndYear = html.EscapeString(strings.TrimSpace(strings.ToLower(p.EndYear)))
}

func (p *PreviousSchool) Validate() error {
	if p.Name == "" {
		return errors.New("name is required")
	}
	if p.StartYear == "" {
		return errors.New("start year is required")
	}
	if p.EndYear == "" {
		return errors.New("end year is required")
	}
	if len(p.ClassNumber) == 0 {
		return errors.New("at least one class is required")
	}

	return nil
}
