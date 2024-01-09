package models

import (
	"errors"
	"html"
	"strings"

	"gorm.io/gorm"
)

type PreviousSchool struct {
	gorm.Model
	Number            string                `json:"number" gorm:"size:36;unique;not null"`
	Name              string                `json:"name" gorm:"size:50;not null"`
	ApplicationNumber string                `json:"application_number" gorm:"size:36;not null"`
	Classes           []PreviousSchoolClass `gorm:"foreignKey:PreviousSchoolNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}

func (p *PreviousSchool) Prepare() {
	p.Name = html.EscapeString(strings.TrimSpace(strings.ToLower(p.Name)))

}

func (p *PreviousSchool) Validate() error {
	if p.Name == "" {
		return errors.New("name of school  is required")
	}

	return nil
}

func (p *PreviousSchool) AfterFind(tx *gorm.DB) (err error) {
	var classes []PreviousSchoolClass
	tx.Model(&PreviousSchoolClass{}).Where("previous_school_number = ?", p.Number).Find(&classes)
	p.Classes = classes

	return
}
