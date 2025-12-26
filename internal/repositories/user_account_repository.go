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

type UserAccountRepository interface {
	Create(k *models.UserAccount) error
	FindByID(id uint) (*models.UserAccount, error)
	FindByEmail(email string) (*models.UserAccount, error)
	GetByID(id uint) (*models.UserAccount, error)
	Delete(id uint) (bool, error)
	BanUser(userID uint) error
	UnbanUser(userID uint) error
	GetFirstWithRole(r *models.UserRole) (*models.UserAccount, error)
}

type userAccountRepository struct {
	db *gorm.DB
}

func NewUserAccountRepository(db *gorm.DB) UserAccountRepository {
	return &userAccountRepository{db: db}
}
func (d *userAccountRepository) Create(dev *models.UserAccount) error {
	err := d.db.Create(dev).Error
	return err
}

func (d *userAccountRepository) BanUser(userID uint) error {
	return d.db.Model(&models.UserAccount{}).
		Where("id = ?", userID).
		Update("is_banned", "true").Error
}

func (d *userAccountRepository) GetFirstWithRole(r *models.UserRole) (*models.UserAccount, error) {
	var dev models.UserAccount
	err := d.db.Where("role_id = ?", r.ID).First(&dev).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &dev, nil
}

func (d *userAccountRepository) UnbanUser(userID uint) error {
	return d.db.Model(&models.UserAccount{}).
		Where("id = ?", userID).
		Update("is_banned", "false").Error
}

func (d *userAccountRepository) FindByID(id uint) (*models.UserAccount, error) {
	var dev models.UserAccount
	err := d.db.First(&dev, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &dev, nil
}

func (d *userAccountRepository) GetByID(id uint) (*models.UserAccount, error) {
	var dev models.UserAccount
	err := d.db.Preload("UserRole").First(&dev, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &dev, nil
}

func (d *userAccountRepository) FindByEmail(email string) (*models.UserAccount, error) {
	var dev models.UserAccount
	err := d.db.Where("email = ?", email).First(&dev).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &dev, nil
}

func (d *userAccountRepository) Delete(id uint) (bool, error) {
	var result models.UserAccount
	res := d.db.Model(&models.UserAccount{}).Where("id = ?", id).Delete(&result)
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
