package repositories

import (
	"errors"

	"github.com/gclkaze/evamodulerepositoryserver/internal/models"
	"gorm.io/gorm"
)

type DeveloperAccountRepository interface {
	Create(k *models.DeveloperAccount) error
	FindByID(id uint) (*models.DeveloperAccount, error)
	Delete(id uint) (bool, error)
}

type developerAccountRepository struct {
	db *gorm.DB
}

func NewDeveloperAccountRepository(db *gorm.DB) DeveloperAccountRepository {
	return &developerAccountRepository{db: db}
}
func (d *developerAccountRepository) Create(dev *models.DeveloperAccount) error {
	err := d.db.Create(dev).Error
	return err
}

func (d *developerAccountRepository) FindByID(id uint) (*models.DeveloperAccount, error) {
	var dev models.DeveloperAccount
	err := d.db.First(&dev, 1).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &dev, nil
}

func (d *developerAccountRepository) Delete(id uint) (bool, error) {
	var result models.DeveloperAccount
	res := d.db.Model(&models.DeveloperAccount{}).Where("id = ?", id).Delete(&result)
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
