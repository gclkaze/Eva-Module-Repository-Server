package models

import (
	"time"

	"gorm.io/gorm"
)

type ModuleRelease struct {
	gorm.Model
	ModuleID   uint                `json:"module_id"`
	Version    string              `json:"version"`
	ReleasedAt time.Time           `json:"released_at"`
	Status     ModuleReleaseStatus `gorm:"foreignKey:StatusID"`
}
