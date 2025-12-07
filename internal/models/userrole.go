package models

import "gorm.io/gorm"

type UserRoleTypeDef int

const (
	Admin UserRoleTypeDef = iota
	Maintainer
	User
)

func GetUserRoleTypes() []UserRoleTypeDef {
	return []UserRoleTypeDef{Admin, Maintainer, User}
}

func (r UserRoleTypeDef) String() string {
	return [...]string{"Admin", "Maintainer", "User"}[r]
}

type UserRole struct {
	gorm.Model
	Name        string           `gorm:"uniqueIndex;size:64;not null"`
	Permissions []UserPermission `gorm:"many2many:user_role_permissions;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func NewUserRole(name string, permissions []UserPermission) *UserRole {
	return &UserRole{Name: name, Permissions: permissions}
}
func NewUserRoleWithoutPerms(name string) *UserRole {
	return &UserRole{Name: name}
}
