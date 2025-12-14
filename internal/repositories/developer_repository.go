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
	err := tx.First(&dev, 1).Error
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
	err := d.db.Where("user_id = ?", id).First(&m).Error
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
	err := d.db.First(&dev, 1).Error
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
