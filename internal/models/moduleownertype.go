package models

import (
	"gorm.io/gorm"
)

type ModuleOwnerTypeDef int

const (
	Dev ModuleOwnerTypeDef = iota
	Group
	Organization
)

func GetModuleOwnerTypes() []ModuleOwnerTypeDef {
	return []ModuleOwnerTypeDef{Dev, Group, Organization}
}

func (r ModuleOwnerTypeDef) String() string {
	return [...]string{"Dev", "Group", "Organization"}[r]
}

type ModuleOwnerType struct {
	gorm.Model
	Label string `gorm:"unique;not null" json:"label"`
}

func NewModuleOwnerType(label string) *ModuleOwnerType {
	return &ModuleOwnerType{Label: label}
}
