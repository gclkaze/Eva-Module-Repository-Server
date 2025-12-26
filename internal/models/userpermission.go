// Copyright (c) 2025 Michail Dorgiakis - gclkaze@gmail.com - https://github.com/gclkaze
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

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
