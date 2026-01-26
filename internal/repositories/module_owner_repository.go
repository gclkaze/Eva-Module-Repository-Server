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

type ModuleOwnerRepository interface {
	Create(t *models.ModuleOwnerType, entityID uint) (*models.ModuleOwner, error)
	Find(t *models.ModuleOwnerType, entityID uint) (*models.ModuleOwner, error)
}

type moduleOwnerRepository struct {
	db *gorm.DB
}

func NewModuleOwnerRepository(db *gorm.DB) ModuleOwnerRepository {
	return &moduleOwnerRepository{db: db}
}

func (m *moduleOwnerRepository) Create(t *models.ModuleOwnerType, entityID uint) (*models.ModuleOwner, error) {

	moduleOwner := models.NewModuleOwner(*t, entityID)
	res := m.db.Create(moduleOwner)
	if res.Error != nil {
		return nil, res.Error
	}
	return moduleOwner, nil
}

func (m *moduleOwnerRepository) Find(t *models.ModuleOwnerType, entityID uint) (*models.ModuleOwner, error) {
	var res models.ModuleOwner
	err := m.db.Preload("Type").Where("type_id = ? AND entity_id = ?", t.ID, entityID).
		First(&res).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &res, nil
}
