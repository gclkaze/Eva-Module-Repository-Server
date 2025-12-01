package models

import "gorm.io/gorm"

type DeveloperModuleOwner struct {
	gorm.Model

	DeveloperID   uint `json:"developer_id"`
	ModuleOwnerID uint `json:"module_owner_id"`

	Developer   Developer   `json:"developer"`
	ModuleOwner ModuleOwner `json:"module_owner"`
}

func NewDeveloperModuleOwner(d Developer, m ModuleOwner) *DeveloperModuleOwner {
	return &DeveloperModuleOwner{DeveloperID: d.ID, ModuleOwnerID: m.ID, Developer: d, ModuleOwner: m}
}
