package models

import (
	"time"

	"gorm.io/gorm"
)

type Developer struct {
	gorm.Model
	Handle    string           `json:"handle"`
	FirstName string           `json:"first_name"`
	LastName  string           `json:"last_name"`
	Account   DeveloperAccount `gorm:"foreignKey:DeveloperID"`
}

func (m *Developer) BeforeCreate(tx *gorm.DB) (err error) {
	m.CreatedAt = time.Now()
	return
}
