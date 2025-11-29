package models

import "gorm.io/gorm"

/*type DeveloperModuleOwner struct {
	gorm.Model
	OwnerID     uint        `json:"owner_id"`
	Owner       ModuleOwner `gorm:"foreignKey:OwnerID"`
	DeveloperID uint        `json:"developer_id"`
	Developer   Developer   `gorm:"foreignKey:DeveloperID"`
}*/

type DeveloperModuleOwner struct {
	gorm.Model

	DeveloperID   uint `json:"developer_id"`
	ModuleOwnerID uint `json:"module_owner_id"`

	Developer   Developer   `json:"developer"`
	ModuleOwner ModuleOwner `json:"module_owner"`

	/*	Role string `json:"role"`*/
}
