package models

import (
	"time"

	"gorm.io/gorm"
)

type Module struct {
	gorm.Model
	Title       string      `json:"title"`
	Repr        string      `json:"repr"`
	Description string      `json:"description"`
	OwnerID     uint        `json:"owner_id"`
	Owner       ModuleOwner `gorm:"foreignKey:OwnerID"`

	Releases []ModuleRelease `gorm:"foreignKey:ModuleID"`
}

func (m *Module) BeforeCreate(tx *gorm.DB) (err error) {
	m.CreatedAt = time.Now()
	return
}
