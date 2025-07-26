package models

import (
	"errors"
	models2 "github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/eos/models"
	"html"
	"strings"
	"time"

	"github.com/Noregrets-Systems/wholesomeEdBackEnd/domain/auth/models"
)

type SourceOfSchoolInfo struct {
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Number       string         `json:"number" gorm:"size:36;primaryKey"`
	Name         string         `gorm:"size:100;not null" json:"name"`
	UserNumber   string         `json:"user_number" gorm:"size:36;index;not null"`
	User         models.User    `gorm:"foreignKey:UserNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:RESTRICT;"`
	EntityNumber string         `json:"entity_number" gorm:"size:36;index;not null"`
	Entity       models2.Entity `gorm:"foreignKey:EntityNumber;references:Number;constraint: OnUpdate:CASCADE, OnDelete:CASCADE;"`
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
