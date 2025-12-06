package models

import (
	"gorm.io/gorm"
)

type DeveloperAccount struct {
	gorm.Model
}

func NewDeveloperAccount() *DeveloperAccount {
	return &DeveloperAccount{}
}
