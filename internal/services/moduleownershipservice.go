package services

import (
	"github.com/gclkaze/evamodulerepositoryserver/internal/models"
	"github.com/gclkaze/evamodulerepositoryserver/internal/repositories"
	"gorm.io/gorm"
)

type ModuleOwnershipService struct {
	moduleOwnerRepo     repositories.ModuleOwnerRepository
	moduleOwnerTypeRepo repositories.ModuleOwnerTypesRepository
	devModuleOwnerRepo  repositories.DeveloperModuleOwnerRepository
}

func NewModuleOwnershipService(moduleOwnerRepo repositories.ModuleOwnerRepository,
	moduleOwnerTypeRepo repositories.ModuleOwnerTypesRepository,
	devModuleOwnerRepo repositories.DeveloperModuleOwnerRepository) *ModuleOwnershipService {
	return &ModuleOwnershipService{moduleOwnerRepo: moduleOwnerRepo, moduleOwnerTypeRepo: moduleOwnerTypeRepo, devModuleOwnerRepo: devModuleOwnerRepo}
}

func (s *ModuleOwnershipService) CreateModuleOwner(t models.ModuleOwnerTypeDef, entityID uint) (*models.ModuleOwner, error) {
	typ, err := s.moduleOwnerTypeRepo.FindByLabel(t.String())
	if err != nil {
		return nil, err
	}
	return s.moduleOwnerRepo.Create(typ, entityID)
}

func (s *ModuleOwnershipService) CreateModuleOwnerTx(tx *gorm.DB, t models.ModuleOwnerTypeDef, entityID uint) (*models.ModuleOwner, error) {
	typ, err := s.moduleOwnerTypeRepo.FindByLabelTx(tx, t.String())
	if err != nil {
		return nil, err
	}
	return s.moduleOwnerRepo.Create(typ, entityID)
}

func (s *ModuleOwnershipService) CreateDeveloperModuleOwner(d *models.Developer, mo *models.ModuleOwner) (*models.DeveloperModuleOwner, error) {
	return s.devModuleOwnerRepo.Create(d, mo)
}

func (s *ModuleOwnershipService) CreateDeveloperModuleOwnerTx(tx *gorm.DB, d *models.Developer, mo *models.ModuleOwner) (*models.DeveloperModuleOwner, error) {
	return s.devModuleOwnerRepo.CreateTx(tx, d, mo)
}

func (s ModuleOwnershipService) FindDeveloperModuleOwner(mo *models.ModuleOwner) (*models.DeveloperModuleOwner, error) {
	if mo.Type.Label != models.Dev.String() {
		return nil, nil
	}

	res, err := s.devModuleOwnerRepo.FindByDevAndModOwner(mo)
	if err != nil {
		return nil, nil
	}
	return res, nil
}
