package models

import "gorm.io/gorm"

// Module is the persisted registry state for a single EMS addon/module.
//
// The descriptive catalogue (code, name, type, version, dependencies) lives in
// Go code in package addons. This table stores the per-addon state so that
// enable/disable and, later, licensing can reference it without re-reading code.
//
// EMS Core Foundation A: this is the only table introduced by this release.
type Module struct {
	gorm.Model
	Code      string `json:"code" gorm:"size:64;unique;not null"`
	Name      string `json:"name" gorm:"size:128"`
	Version   string `json:"version" gorm:"size:32"`
	DependsOn string `json:"depends_on" gorm:"size:255"`
	Enabled   bool   `json:"enabled" gorm:"not null;default:false"`
	Status    string `json:"status" gorm:"size:20"`
}

// TableName pins the table to ems_modules. GORM would otherwise derive
// "modules", which could collide with an Education concept.
func (Module) TableName() string {
	return "ems_modules"
}
