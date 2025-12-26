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

type DeveloperRepository interface {
	Create(k *models.Developer) error
	FindByID(id uint) (*models.Developer, error)
	FindByUserAccountID(id uint) (*models.Developer, error)
	Delete(id uint) (bool, error)
	Ban(id uint) (bool, error)
	FindByIDTx(tx *gorm.DB, id uint) (*models.Developer, error)
}

type developerRepository struct {
	db *gorm.DB
}

func NewDeveloperRepository(db *gorm.DB) DeveloperRepository {
	return &developerRepository{db: db}
}

func (d *developerRepository) Create(dev *models.Developer) error {
	err := d.db.Create(dev).Error
	return err
}

func (d *developerRepository) FindByIDTx(tx *gorm.DB, id uint) (*models.Developer, error) {
	var dev models.Developer
	err := tx.First(&dev, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &dev, nil
}

func (d *developerRepository) FindByUserAccountID(id uint) (*models.Developer, error) {
	var m models.Developer
	err := d.db.Preload("UserAccount").Where("user_id = ?", id).First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (d *developerRepository) FindByID(id uint) (*models.Developer, error) {
	var dev models.Developer
	err := d.db.First(&dev, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &dev, nil
}

func (d *developerRepository) Delete(id uint) (bool, error) {
	var result models.Developer
	res := d.db.Model(&models.Developer{}).Where("id = ?", id).Delete(&result)
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

func (d *developerRepository) Ban(id uint) (bool, error) {
	err := d.db.Model(&models.Developer{}).Where("id = ?", id).Update("active", false).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
