package services

import "github.com/gclkaze/evamodulerepositoryserver/internal/repositories"

type ModuleOwnershipService struct {
	moduleOwnerRepo     repositories.ModuleOwnerRepository
	moduleOwnerTypeRepo repositories.ModuleOwnerTypesRepository
	/*	repo repositories.ModuleRepository

		moduleFolder  string
		releaseFolder string

		logger logger.ILogger

		developerService *DeveloperService*/
}

func NewModuleOwnershipService(	moduleOwnerRepo     repositories.ModuleOwnerRepository
	moduleOwnerTypeRepo repositories.ModuleOwnerTypesRepository) *ModuleOwnershipService{
  return &ModuleOwnershipService{moduleOwnerRepo: repositories.ModuleOwnerRepository}
}