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

func NewModule(title string, repr string, description string, ownerID uint, owner ModuleOwner, keywords []Keyword) *Module {
	return &Module{Title: title, Repr: repr, Description: description, OwnerID: ownerID, Owner: owner, Keywords: keywords}
}

func (m *Module) Update(title string, repr string, description string, keywords []Keyword) {
	m.Title = title
	m.Repr = repr
	m.Description = description
	m.Keywords = keywords
}
