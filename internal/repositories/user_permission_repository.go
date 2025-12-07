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
	err := d.db.First(&dev, 1).Error
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
