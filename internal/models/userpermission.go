package models

import "gorm.io/gorm"

type UserPermissionTypeDef int

const (
	CreateModule UserPermissionTypeDef = iota
	SuggestModule
	DeleteModule
	DeleteMyModule
	DeleteRelease
	DeleteMyRelease
	UpdateRelease
	ChangeReleaseStatus
	RejectRelease
	AcceptRelease
	CancelRelease
	BanUser
	UnbanUser
)

func GetUserPermissionTypes() []UserPermissionTypeDef {
	return []UserPermissionTypeDef{
		CreateModule,
		SuggestModule,
		DeleteModule,
		DeleteMyModule,
		DeleteRelease,
		DeleteMyRelease,
		UpdateRelease,
		ChangeReleaseStatus,
		RejectRelease,
		AcceptRelease,
		CancelRelease,
		BanUser,
		UnbanUser}
}

func (r UserPermissionTypeDef) String() string {
	return [...]string{
		"CreateModule",
		"SuggestModule",
		"DeleteModule",
		"DeleteMyModule",
		"DeleteRelease",
		"DeleteMyRelease",
		"UpdateRelease",
		"ChangeReleaseStatus",
		"RejectRelease",
		"AcceptRelease",
		"CancelRelease",
		"BanUser",
		"UnbanUser"}[r]
}

type UserPermission struct {
	gorm.Model
	Value string `gorm:"uniqueIndex;size:128;not null"`
}

func NewUserPermission(label string) *UserPermission {
	return &UserPermission{Value: label}
}
