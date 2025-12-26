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

package repositories

import (
	"errors"

	"github.com/gclkaze/evamodulerepositoryserver/internal/models"
	"gorm.io/gorm"
)

type ReleaseStatusTypeDef int

const (
	Draft ReleaseStatusTypeDef = iota
	Pending
	Accepted
	Rejected
	Canceled
)

func (t ReleaseStatusTypeDef) String() string {
	return [...]string{"draft", "pending", "accepted", "rejected", "canceled"}[t]
}

type ReleaseStatusRepository interface {
	Initialize() error
	GetStatus(t ReleaseStatusTypeDef) (*models.ModuleReleaseStatus, error)
}

type releaseStatusRepository struct {
	db *gorm.DB
}

func NewReleaseStatusRepository(db *gorm.DB) ReleaseStatusRepository {
	return &releaseStatusRepository{db: db}
}

func (r releaseStatusRepository) Initialize() error {
	statuses := []string{"draft", "pending", "accepted", "rejected", "canceled"}
	description := []string{"This is a draft of a release.", "The release is waiting to be checked by the EVA Language Team.",
		"The release has been accepted by the EVA Language Team.", "The release has been rejected by the EVA Language Team.",
		"The release has been canceled by the EVA Language Team."}

	for i, status := range statuses {
		var count int64
		r.db.Model(&models.ModuleReleaseStatus{}).Where("label = ?", status).Count(&count)
		if count == 0 {
			if err := r.db.Create(&models.ModuleReleaseStatus{Label: status, Description: description[i]}).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

func (r releaseStatusRepository) GetStatus(t ReleaseStatusTypeDef) (*models.ModuleReleaseStatus, error) {
	var m models.ModuleReleaseStatus
	res := r.db.Where("label = ?", t.String()).First(&m)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &m, nil
}
