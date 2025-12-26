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

type UserPermissionRepository interface {
	Create(k *models.UserPermission) error
	FindByID(id uint) (*models.UserPermission, error)
	FindByValue(v string) (*models.UserPermission, error)
	Delete(id uint) (bool, error)
	Initialize() error
}

type userPermissionRepository struct {
	db *gorm.DB
}

func NewUserPermissionRepository(db *gorm.DB) UserPermissionRepository {
	return &userPermissionRepository{db: db}
}
func (d *userPermissionRepository) Create(dev *models.UserPermission) error {
	err := d.db.Create(dev).Error
	return err
}

func (d *userPermissionRepository) FindByValue(v string) (*models.UserPermission, error) {
	var dev models.UserPermission
	err := d.db.Where("value = ?", v).First(&dev).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &dev, nil
}

func (d *userPermissionRepository) FindByID(id uint) (*models.UserPermission, error) {
	var dev models.UserPermission
	err := d.db.First(&dev, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &dev, nil
}

func (d *userPermissionRepository) Delete(id uint) (bool, error) {
	var result models.UserPermission
	res := d.db.Model(&models.UserPermission{}).Where("id = ?", id).Delete(&result)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if res.Error != nil {
		return false, res.Error
	}
	if res.RowsAffected == 0 {
		return false, nil
	}
	return true, nil
}

func (d *userPermissionRepository) Initialize() error {
	for _, t := range models.GetUserPermissionTypes() {
		var count int64
		d.db.Model(&models.UserPermission{}).Where("value = ?", t.String()).Count(&count)
		if count == 0 {
			if err := d.db.Create(models.NewUserPermission(t.String())).Error; err != nil {
				return err
			}
		}
	}
	return nil
}
