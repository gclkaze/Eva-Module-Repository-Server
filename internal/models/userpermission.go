package models

import "gorm.io/gorm"

type UserPermissionTypeDef int

const (
	CreateMyModule UserPermissionTypeDef = iota
	DeleteModules
	DeleteMyModule
	SuggestMyModule
	UpdateModules
	DeleteReleases
	DeleteMyRelease
	UpdateReleases
	ChangeReleaseStatuses
	RejectReleases
	AcceptReleases
	CancelReleases
	BanUsers
	UnbanUsers
)

func GetUserPermissionTypes() []UserPermissionTypeDef {
	return []UserPermissionTypeDef{
		CreateMyModule,
		DeleteModules,
		DeleteMyModule,
		SuggestMyModule,
		UpdateModules,
		DeleteReleases,
		DeleteMyRelease,
		UpdateReleases,
		ChangeReleaseStatuses,
		RejectReleases,
		AcceptReleases,
		CancelReleases,
		BanUsers,
		UnbanUsers}
}

func (r UserPermissionTypeDef) String() string {
	return [...]string{
		"CreateMyModule",
		"DeleteModules",
		"DeleteMyModule",
		"SuggestMyModule",
		"UpdateModules",
		"DeleteReleases",
		"DeleteMyRelease",
		"UpdateReleases",
		"ChangeReleaseStatuses",
		"RejectReleases",
		"AcceptReleases",
		"CancelReleases",
		"BanUsers",
		"UnbanUsers"}[r]
}

type UserPermission struct {
	gorm.Model
	Value string `gorm:"uniqueIndex;size:128;not null"`
}

func NewUserPermission(label string) *UserPermission {
	return &UserPermission{Value: label}
}
