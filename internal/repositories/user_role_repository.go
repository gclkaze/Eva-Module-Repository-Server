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
	err := d.db.First(&dev, 1).Error
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
		theResult := d.db.Where("name = ?", t.String()).Find(&res)
		if errors.Is(theResult.Error, gorm.ErrRecordNotFound) {
			if err := d.db.Create(models.NewUserRoleWithoutPerms(t.String())).Error; err != nil {
				return err
			}
		}
	}
	return nil
}
