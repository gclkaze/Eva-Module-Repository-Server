package models

import (
	"gorm.io/gorm"
)

type UserAccount struct {
	gorm.Model

	Email    string `gorm:"type:varchar(255);uniqueIndex" json:"email"`
	Password string
	RoleID   uint     `json:"user_id"`
	UserRole UserRole `gorm:"foreignKey:RoleID"`
	IsBanned bool     `gorm:"is_banned"`
}

func NewUserAccount(role *UserRole, email string, password string) *UserAccount {
	return &UserAccount{RoleID: role.ID, UserRole: *role, Email: email, Password: password}
}
