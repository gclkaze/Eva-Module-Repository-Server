package models

import (
	"gorm.io/gorm"
)

type ModuleOwnerType struct {
	gorm.Model
	Label string `json:"label"`
}
