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

package dto

import (
	"time"

	"github.com/gclkaze/evamodulerepositoryserver/internal/models"
)

type ReleaseDTO struct {
	ID          uint       `json:"id" binding:"required"`
	Version     string     `json:"version" binding:"required"`
	ReleasedAt  *time.Time `json:"released_at" binding:"required"`
	Description string     `json:"description" binding:"required"`
	DiskSize    int64      `json:"diskSize" binding:"required"`

	Status   string       `json:"status" binding:"required"`
	Keywords []KeywordDTO `json:"keywords" binding:"required"`
}

func NewReleaseDTO(release models.ModuleRelease) *ReleaseDTO {
	var keywordsdtos []KeywordDTO
	for i := 0; i < len(release.Keywords); i++ {
		keywordsdtos = append(keywordsdtos, *NewKeywordDTO(release.Keywords[i]))
	}

	return &ReleaseDTO{ID: release.ID, Version: release.Version, ReleasedAt: release.ReleasedAt, Description: release.Description, Keywords: keywordsdtos, DiskSize: release.DiskSize, Status: release.Status.Label}
}
