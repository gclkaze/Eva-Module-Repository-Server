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

type Module struct {
	gorm.Model
	Title       string      `json:"title"`
	Repr        string      `gorm:"type:varchar(50);uniqueIndex:idx_module_repr" json:"repr"`
	Description string      `json:"description"`
	OwnerID     uint        `json:"owner_id"`
	Owner       ModuleOwner `gorm:"foreignKey:OwnerID"`
	RepoName    string      `json:"repoName"`

	//Releases []ModuleRelease `gorm:"foreignKey:ModuleID"`

	Keywords []Keyword `gorm:"many2many:module_keywords;" json:"keywords,omitempty"`
}

func NewModule(title string, repr string, description string, ownerID uint, owner ModuleOwner, keywords []Keyword, repoName string) *Module {
	return &Module{Title: title, Repr: repr, Description: description, OwnerID: ownerID, Owner: owner, Keywords: keywords, RepoName: repoName}
}

func (m *Module) Update(title string, repr string, description string, keywords []Keyword /*, repoName string*/) {
	m.Title = title
	m.Repr = repr
	m.Description = description
	m.Keywords = keywords
	//	m.RepoName = repoName
}
