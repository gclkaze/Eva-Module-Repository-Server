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

type DeveloperModuleOwnerRepository interface {
	Create(d *models.Developer, mo *models.ModuleOwner) (*models.DeveloperModuleOwner, error)
	CreateTx(tx *gorm.DB, d *models.Developer, mo *models.ModuleOwner) (*models.DeveloperModuleOwner, error)
	FindByDevAndModOwner(mo *models.ModuleOwner) (*models.DeveloperModuleOwner, error)
	Delete(id uint) (bool, error)
}

type developerModuleOwnerRepository struct {
	db *gorm.DB
}

func NewDeveloperModuleOwnerRepository(db *gorm.DB) DeveloperModuleOwnerRepository {
	return &developerModuleOwnerRepository{db: db}
}

func (m *developerModuleOwnerRepository) Create(d *models.Developer, mo *models.ModuleOwner) (*models.DeveloperModuleOwner, error) {
	dmo := models.NewDeveloperModuleOwner(*d, *mo)
	m.db.Create(dmo)
	return dmo, nil
}

func (m *developerModuleOwnerRepository) CreateTx(tx *gorm.DB, d *models.Developer, mo *models.ModuleOwner) (*models.DeveloperModuleOwner, error) {
	dmo := models.NewDeveloperModuleOwner(*d, *mo)
	tx.Create(dmo)
	return dmo, nil
}

func (m developerModuleOwnerRepository) FindByDevAndModOwner(mo *models.ModuleOwner) (*models.DeveloperModuleOwner, error) {
	var dmo models.DeveloperModuleOwner
	res := m.db.Preload("Developer").Preload("ModuleOwner").Where("developer_id = ? AND module_owner_id = ?", mo.EntityID, mo.ID).First(&dmo)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &dmo, nil
}

func (m developerModuleOwnerRepository) Delete(id uint) (bool, error) {
	var r models.DeveloperModuleOwner
	res := m.db.Where(id).Delete(&r)

	if res.Error != nil {
		return false, res.Error
	}

	if res.RowsAffected == 0 {
		return false, nil
	}

	return true, nil
}
