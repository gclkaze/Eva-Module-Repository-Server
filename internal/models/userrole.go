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
