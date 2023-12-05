package models

import (
	"errors"
	"html"
	"strings"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type SourceOfSchoolInfo struct {
	gorm.Model
	ID   uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`
	Name string    `gorm:"size:100;unique;not null" json:"name" `
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
