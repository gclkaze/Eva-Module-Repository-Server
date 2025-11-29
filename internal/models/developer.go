package models

import (
	"gorm.io/gorm"
)

type Developer struct {
	gorm.Model
	Handle           string           `json:"handle"`
	FirstName        string           `json:"first_name"`
	LastName         string           `json:"last_name"`
	DeveloperID      uint             `json:"developer_id"`
	DeveloperAccount DeveloperAccount `gorm:"foreignKey:DeveloperID"`
}
