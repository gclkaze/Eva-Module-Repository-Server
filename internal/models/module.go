package models

import (
	"gorm.io/gorm"
)

type Module struct {
	gorm.Model
	Title       string      `json:"title"`
	Repr        string      `json:"repr"`
	Description string      `json:"description"`
	OwnerID     uint        `json:"owner_id"`
	Owner       ModuleOwner `gorm:"foreignKey:OwnerID"`

	//Releases []ModuleRelease `gorm:"foreignKey:ModuleID"`

	Keywords []Keyword `gorm:"many2many:module_keywords;" json:"keywords,omitempty"`
}
