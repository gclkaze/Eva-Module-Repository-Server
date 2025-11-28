package models

import "gorm.io/gorm"

type DeveloperModuleOwner struct {
	gorm.Model
	TypeID      uint            `json:"type_id"`
	Type        ModuleOwnerType `gorm:"foreignKey:TypeID"`
	OwnerID     uint            `json:"owner_id"`
	Owner       ModuleOwner     `gorm:"foreignKey:OwnerID"`
	DeveloperID uint            `json:"developer_id"`
	Developer   Developer       `gorm:"foreignKey:DeveloperID"`
}
