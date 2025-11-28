// Package models contains models for the module related classes
package models

import "gorm.io/gorm"

type ModuleOwner struct {
	gorm.Model
	Type    ModuleOwnerType `gorm:"foreignKey:ModuleTypeID"`
	OwnerID uint            `json:"owner_id"`
}
