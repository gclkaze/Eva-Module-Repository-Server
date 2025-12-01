package repositories

import (
	"github.com/gclkaze/evamodulerepositoryserver/internal/models"
	"gorm.io/gorm"
)

type DeveloperModuleOwnerRepository interface {
	Create(d models.Developer, mo models.ModuleOwner) error
}

type developerModuleOwnerRepository struct {
	db *gorm.DB
}

func NewDeveloperModuleOwnerRepository(db *gorm.DB) DeveloperModuleOwnerRepository {
	return &developerModuleOwnerRepository{db: db}
}

func (m *developerModuleOwnerRepository) Create(d models.Developer, mo models.ModuleOwner) error {
	dmo := models.NewDeveloperModuleOwner(d, mo)
	m.db.Create(dmo)
	return nil
}
