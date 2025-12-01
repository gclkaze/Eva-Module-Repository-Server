// Package models contains models for the module related classes
package models

import "gorm.io/gorm"

type ModuleOwner struct {
	gorm.Model
	TypeID   uint            `json:"type_id"`
	Type     ModuleOwnerType `gorm:"foreignKey:TypeID"`
	EntityID uint            `json:"entity_id"`
}

func NewModuleOwner(t ModuleOwnerType, id uint) *ModuleOwner {
	return &ModuleOwner{TypeID: t.ID, Type: t, EntityID: id}
}
