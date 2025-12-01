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
	Active           bool             `json:"is_active"`
}

func NewDeveloper(handle string, firstName string, lastName string, developerID uint, developerAccount DeveloperAccount, active bool) *Developer {
	return &Developer{Handle: handle, FirstName: firstName, LastName: lastName, DeveloperID: developerID, DeveloperAccount: developerAccount, Active: active}
}
