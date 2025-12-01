package models

import (
	"gorm.io/gorm"
)

type DeveloperAccount struct {
	gorm.Model
	//DeveloperID uint `json:"developer_id"`
	//Developer   Developer `gorm:"foreignKey:Developer"`
}

func NewDeveloperAccount() *DeveloperAccount {
	return &DeveloperAccount{}
}
