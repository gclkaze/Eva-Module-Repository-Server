// Package repositories contains
package repositories

import (
	"github.com/gclkaze/evamodulerepositoryserver/internal/models"
	"gorm.io/gorm"
)

type ReleaseRepository interface {
	Create(dev *models.ModuleRelease) error
	FindByID(id uint) (*models.ModuleRelease, error)
	//SearchByKeywords(tags []string) ([]models.ModuleRelease, error)
}

type releaseRepository struct {
	db *gorm.DB
}

func NewReleaseRepository(db *gorm.DB) ReleaseRepository {
	return &releaseRepository{db: db}
}

func (r *releaseRepository) Create(mod *models.ModuleRelease) error {
	return r.db.Create(mod).Error
}

func (r releaseRepository) FindByID(id uint) (*models.ModuleRelease, error) {
	/*	var m models.Module
		err := r.db.First(&m, id).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		if err != nil {
			return nil, err
		}

		return &m, nil*/
	return nil, nil
}
