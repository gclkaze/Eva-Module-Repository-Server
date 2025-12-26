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

type UserRoleRepository interface {
	Create(k *models.UserRole) error
	FindByID(id uint) (*models.UserRole, error)
	FindByValue(name string) (*models.UserRole, error)
	Delete(id uint) (bool, error)
	Update(*models.UserRole) error
	Initialize() error
}

type userRoleRepository struct {
	db *gorm.DB
}

func NewUserRoleRepository(db *gorm.DB) UserRoleRepository {
	return &userRoleRepository{db: db}
}
func (d *userRoleRepository) Create(dev *models.UserRole) error {
	err := d.db.Create(dev).Error
	return err
}

func (d *userRoleRepository) FindByID(id uint) (*models.UserRole, error) {
	var dev models.UserRole
	err := d.db.Preload("Permissions").First(&dev, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &dev, nil
}

func (d *userRoleRepository) FindByValue(name string) (*models.UserRole, error) {
	var dev models.UserRole
	err := d.db.Where("name = ?", name).First(&dev).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &dev, nil
}

func (d *userRoleRepository) Delete(id uint) (bool, error) {
	var result models.UserRole
	res := d.db.Model(&models.UserRole{}).Where("id = ?", id).Delete(&result)
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

func (d *userRoleRepository) Update(role *models.UserRole) error {
	res := d.db.Save(role)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (d *userRoleRepository) Initialize() error {
	for _, t := range models.GetUserRoleTypes() {
		var res models.UserRole
		result := d.db.Where("name = ?", t.String()).Find(&res)
		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			if err := d.db.Create(models.NewUserRoleWithoutPerms(t.String())).Error; err != nil {
				return err
			}
		}
	}
	return nil
}
