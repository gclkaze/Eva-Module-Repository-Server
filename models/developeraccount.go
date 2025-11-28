package models

import (
	"time"

	"gorm.io/gorm"
)

type DeveloperAccount struct {
	gorm.Model
	DeveloperID uint `json:"developer_id"`
}

func (m *DeveloperAccount) BeforeCreate(tx *gorm.DB) (err error) {
	m.CreatedAt = time.Now()
	return
}
