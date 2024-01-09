package models

import (
	"errors"
	"html"
	"strings"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
	"gorm.io/gorm"
)

type SourceOfSchoolInfo struct {
	gorm.Model
	Number     string      `json:"number" gorm:"size:36;unique;not null"`
	Name       string      `gorm:"size:100;unique;not null" json:"name"`
	UserNumber string      `json:"user_number" gorm:"size:36;index;not null"`
	User       models.User `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
}

func (s *SourceOfSchoolInfo) TableName() string {
	return "sources_of_school_info"
}

func (s *SourceOfSchoolInfo) Prepare() {
	s.Name = html.EscapeString(strings.TrimSpace(strings.ToLower(s.Name)))
}

func (s *SourceOfSchoolInfo) Validate() error {
	if s.Name == "" {
		return errors.New("name is required")
	}

	return nil
}
