package models

import (
	"errors"
	"html"
	"strings"

	"gorm.io/gorm"
)

type SourceOfSchoolInfo struct {
	gorm.Model
	Name string `json:"name" gorm:"size:100;unique;not null"`
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
