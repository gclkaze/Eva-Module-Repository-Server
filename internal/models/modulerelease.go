package models

import (
	"time"

	"gorm.io/gorm"
)

type ModuleRelease struct {
	gorm.Model
	ModuleID    uint                `json:"module_id"`
	Version     string              `gorm:"not null" json:"version"`
	Description string              `json:"description"`
	ReleasedAt  time.Time           `json:"released_at"`
	StatusID    uint                `json:"status_id"`
	Status      ModuleReleaseStatus `gorm:"foreignKey:StatusID"`
	Keywords    []Keyword           `gorm:"many2many:release_keywords;" json:"keywords,omitempty"`
	DiskSize    uint                `json:"disk_size"`
}
