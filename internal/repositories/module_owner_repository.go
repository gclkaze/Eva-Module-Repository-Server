package repositories

import (
	"github.com/gclkaze/evamodulerepositoryserver/internal/models"
	"gorm.io/gorm"
)

type ModuleOwnerRepository interface {
	Create(t *models.ModuleOwnerType, entityID uint) (*models.ModuleOwner, error)
}

type moduleOwnerRepository struct {
	db *gorm.DB
}

func NewModuleOwnerRepository(db *gorm.DB) ModuleOwnerRepository {
	return &moduleOwnerRepository{db: db}
}

func (m *moduleOwnerRepository) Create(t *models.ModuleOwnerType, entityID uint) (*models.ModuleOwner, error) {

	moduleOwner := models.NewModuleOwner(*t, entityID)
	res := m.db.Create(moduleOwner)
	if res.Error != nil {
		return nil, res.Error
	}
	return moduleOwner, nil
}
