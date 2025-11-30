package models

import (
	"gorm.io/gorm"
)

type ModuleReleaseStatus struct {
	gorm.Model
	Label       string `json:"label"` // e.g., "draft", "published", "deprecated"
	Description string `json:"description"`
}
