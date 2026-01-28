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
	"strings"

	"github.com/gclkaze/evamodulerepositoryserver/internal/models"
	"gorm.io/gorm"
)

type KeywordRepository interface {
	Create(k *models.Keyword) error
	CreateTx(tx *gorm.DB, k *models.Keyword) error
	FindByID(id int64) (*models.Keyword, error)
	FindByLabel(label string) (*models.Keyword, error)
	FindByLabelTx(tx *gorm.DB, label string) (*models.Keyword, error)
	FindAll() ([]models.Keyword, error)
	Delete(id int64) error

	//utility functions
	MergeKeywords(first []models.Keyword, second []models.Keyword) []models.Keyword
}

type keywordRepository struct {
	db *gorm.DB
}

func NewKeywordRepository(db *gorm.DB) KeywordRepository {
	return &keywordRepository{db: db}
}

// Create a new keyword
func (r *keywordRepository) Create(k *models.Keyword) error {
	return r.db.Create(k).Error
}

func (r *keywordRepository) CreateTx(tx *gorm.DB, k *models.Keyword) error {
	return tx.Create(k).Error
}

// Find by ID
func (r *keywordRepository) FindByID(id int64) (*models.Keyword, error) {
	var k models.Keyword
	err := r.db.First(&k, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &k, err
}

// Find by Label
func (r *keywordRepository) FindByLabel(label string) (*models.Keyword, error) {
	var k models.Keyword
	err := r.db.Where("label = ?", label).First(&k).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &k, err
}

func (r *keywordRepository) FindByLabelTx(tx *gorm.DB, label string) (*models.Keyword, error) {
	var k models.Keyword
	err := tx.Where("label = ?", label).First(&k).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &k, err
}

// Find all keywords
func (r *keywordRepository) FindAll() ([]models.Keyword, error) {
	var keywords []models.Keyword
	err := r.db.Find(&keywords).Error
	return keywords, err
}

// Delete by ID
func (r *keywordRepository) Delete(id int64) error {
	return r.db.Delete(&models.Keyword{}, id).Error
}

func (r keywordRepository) MergeKeywords(first []models.Keyword, second []models.Keyword) []models.Keyword {
	seen := make(map[string]struct{}, len(first)+len(second))
	out := make([]models.Keyword, 0, len(first)+len(second))

	add := func(list []models.Keyword) {
		for _, k := range list {
			label := strings.ToLower(strings.TrimSpace(k.Label))
			if label == "" {
				continue // ignore empty labels defensively
			}

			if _, ok := seen[label]; ok {
				continue
			}

			seen[label] = struct{}{}

			// normalize label (optional but recommended)
			k.Label = label
			out = append(out, k)
		}
	}

	add(first)
	add(second)

	return out
}
