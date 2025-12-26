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

import (
	"gorm.io/gorm"
)

type Developer struct {
	gorm.Model
	Handle      string      `json:"handle"`
	FirstName   string      `json:"first_name"`
	LastName    string      `json:"last_name"`
	UserID      uint        `json:"user_id"`
	UserAccount UserAccount `gorm:"foreignKey:UserID"`
	Active      bool        `json:"is_active"`
}

func NewDeveloper(handle string, firstName string, lastName string, userID uint, userAccount UserAccount, active bool) *Developer {
	return &Developer{Handle: handle, FirstName: firstName, LastName: lastName, UserID: userID, UserAccount: userAccount, Active: active}
}
