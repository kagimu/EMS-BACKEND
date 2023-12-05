package models

import (
	"errors"
	"html"
	"strings"

	applicationModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/application/models"
	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	studentModel "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/student/models"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type Parent struct {
	gorm.Model
	ID                         uuid.UUID                      `gorm:"type:uuid;primary_key;" json:"id"`
	Surname                    string                         `json:"surname"`
	SecondName                 string                         `json:"second_name"`
	OtherName                  string                         `json:"other_name"`
	MainTelephoneNumber        string                         `json:"main_telephone_number"`
	AlternativeTelephoneNumber string                         `json:"alternative_telephone_number"`
	Email                      string                         `json:"email"`
	AlternativeEmail           string                         `json:"alternative_email"`
	Learner                    []studentModel.Student         `gorm:"many2many:student_parents;" json:"students"`
	UserNumber                 string                         `json:"user_number" gorm:"size:36;not null"`
	User                       models.User                    `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	Application                []applicationModel.Application `gorm:"foreignkey:ParentID" json:"application"`
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
