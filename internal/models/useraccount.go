package models

import (
	"gorm.io/gorm"
)

type UserAccount struct {
	gorm.Model

	RoleID   uint     `json:"user_id"`
	UserRole UserRole `gorm:"foreignKey:RoleID"`
}

func NewUserAccount(role *UserRole) *UserAccount {
	return &UserAccount{RoleID: role.ID, UserRole: *role}
}
