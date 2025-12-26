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
	"time"

	"gorm.io/gorm"
)

type ModuleRelease struct {
	gorm.Model
	ModuleID    uint                `json:"module_id"`
	Version     string              `gorm:"not null" json:"version"`
	Description string              `json:"description"`
	ReleasedAt  *time.Time          `gorm:"type:datetime"`
	StatusID    uint                `json:"status_id"`
	Status      ModuleReleaseStatus `gorm:"foreignKey:StatusID"`
	Keywords    []Keyword           `gorm:"many2many:release_keywords;" json:"keywords,omitempty"`
	DiskSize    int64               `json:"disk_size"`
	CreatorID   uint                `json:"creator_id"`
	Creator     Developer           `gorm:"foreignKey:CreatorID"`
	Statistics  *ReleaseStatistics  `gorm:"foreignKey:ReleaseID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

func NewModuleReleaseFromModule(m *Module, version string, status ModuleReleaseStatus, diskSize int64, creator Developer) *ModuleRelease {
	return &ModuleRelease{ModuleID: m.ID, Version: version, Description: m.Description, StatusID: status.ID, Status: status, DiskSize: diskSize, Creator: creator, CreatorID: creator.ID}
}
