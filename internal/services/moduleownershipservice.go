package services

import (
	"github.com/gclkaze/evamodulerepositoryserver/internal/models"
	"github.com/gclkaze/evamodulerepositoryserver/internal/repositories"
)

type ModuleOwnershipService struct {
	moduleOwnerRepo     repositories.ModuleOwnerRepository
	moduleOwnerTypeRepo repositories.ModuleOwnerTypesRepository
	devModuleOwnerRepo  repositories.DeveloperModuleOwnerRepository
	/*	repo repositories.ModuleRepository

		moduleFolder  string
		releaseFolder string

		logger logger.ILogger

		developerService *DeveloperService*/
}

func NewModuleOwnershipService(moduleOwnerRepo repositories.ModuleOwnerRepository,
	moduleOwnerTypeRepo repositories.ModuleOwnerTypesRepository,
	devModuleOwnerRepo repositories.DeveloperModuleOwnerRepository) *ModuleOwnershipService {
	return &ModuleOwnershipService{moduleOwnerRepo: moduleOwnerRepo, moduleOwnerTypeRepo: moduleOwnerTypeRepo, devModuleOwnerRepo: devModuleOwnerRepo}
}

func (s *ModuleOwnershipService) CreateModuleOwner(t models.ModuleOwnerTypeDef, entityId uint) (*models.ModuleOwner, error) {
	typ, err := s.moduleOwnerTypeRepo.FindByLabel(t.String())
	if err != nil {
		return nil, err
	}
	return s.moduleOwnerRepo.Create(typ, entityId)
}

func (s *ModuleOwnershipService) CreateDeveloperModuleOwner(d *models.Developer, mo *models.ModuleOwner) (*models.DeveloperModuleOwner, error) {
	return s.devModuleOwnerRepo.Create(d, mo)
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
