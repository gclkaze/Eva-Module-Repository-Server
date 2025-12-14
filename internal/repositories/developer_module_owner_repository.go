package repositories

import (
	"errors"

	"github.com/gclkaze/evamodulerepositoryserver/internal/models"
	"gorm.io/gorm"
)

type DeveloperModuleOwnerRepository interface {
	Create(d *models.Developer, mo *models.ModuleOwner) (*models.DeveloperModuleOwner, error)
	CreateTx(tx *gorm.DB, d *models.Developer, mo *models.ModuleOwner) (*models.DeveloperModuleOwner, error)
	FindByDevAndModOwner(mo *models.ModuleOwner) (*models.DeveloperModuleOwner, error)
}

type developerModuleOwnerRepository struct {
	db *gorm.DB
}

func NewDeveloperModuleOwnerRepository(db *gorm.DB) DeveloperModuleOwnerRepository {
	return &developerModuleOwnerRepository{db: db}
}

func (m *developerModuleOwnerRepository) Create(d *models.Developer, mo *models.ModuleOwner) (*models.DeveloperModuleOwner, error) {
	dmo := models.NewDeveloperModuleOwner(*d, *mo)
	m.db.Create(dmo)
	return dmo, nil
}

func (m *developerModuleOwnerRepository) CreateTx(tx *gorm.DB, d *models.Developer, mo *models.ModuleOwner) (*models.DeveloperModuleOwner, error) {
	dmo := models.NewDeveloperModuleOwner(*d, *mo)
	tx.Create(dmo)
	return dmo, nil
}

func (m developerModuleOwnerRepository) FindByDevAndModOwner(mo *models.ModuleOwner) (*models.DeveloperModuleOwner, error) {
	var dmo models.DeveloperModuleOwner
	res := m.db.Preload("Developer").Preload("ModuleOwner").Where("developer_id = ? AND module_owner_id = ?", mo.EntityID, mo.ID).First(&dmo)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &dmo, nil
}
