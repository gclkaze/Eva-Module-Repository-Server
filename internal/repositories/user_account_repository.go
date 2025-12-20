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
		Where("id = ?", 1).
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
		Where("id = ?", 1).
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
