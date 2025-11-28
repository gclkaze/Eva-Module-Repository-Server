package models

import "gorm.io/gorm"

type DeveloperModuleOwner struct {
	gorm.Model
	Type        ModuleOwnerType `gorm:"foreignKey:ModuleTypeID"`
	OwnerID     uint            `json:"owner_id"`
	Owner       ModuleOwner     `gorm:"foreignKey:OwnerID"`
	DeveloperID uint            `json:"developer_id"`
}
