package backend

import (
	"fmt"

	"github.com/gclkaze/evamodulerepositoryserver/internal/config"
	"github.com/gclkaze/evamodulerepositoryserver/internal/db"
	"github.com/gclkaze/evamodulerepositoryserver/internal/services"
	"github.com/gclkaze/evamodulerepositoryserver/pkg/logger"
	"github.com/magiconair/properties"
)

type EvaModuleRepositoryBackend struct {
	db         *db.EvaModuleRepositoryDatabase
	properties *properties.Properties
	logger     logger.ILogger

	moduleService          *services.ModuleService
	releaseService         *services.ReleaseService
	developerService       *services.DeveloperService
	moduleOwnershipService *services.ModuleOwnershipService
}

func NewEvaModuleRepositoryBackend() *EvaModuleRepositoryBackend {
	inst := &EvaModuleRepositoryBackend{}
	return inst
}

func (be *EvaModuleRepositoryBackend) Initialize() error {
	if config.TheConfigReader.IsOnError() {
		return fmt.Errorf("couldn't read the properties file")
	}
	be.properties = config.TheConfigReader.GetProperties()
	be.db = db.NewEvaModuleRepositoryDatabase()
	error := be.db.Initialize(be.properties)

	if error != nil {
		return error
	}
	error = be.initializeServices()
	return error
}

func (be *EvaModuleRepositoryBackend) initializeServices() error {
	be.releaseService = services.NewReleaseService(be.db.GetReleaseRepository(), be.db.GetReleaseStatusRepository())

	be.developerService = services.NewDeveloperService(be.db.GetDeveloperRepository(), be.db.GetDeveloperAccountRepository(), be.properties)
	be.moduleOwnershipService = services.NewModuleOwnershipService(be.db.GetModuleOwnerRepository(), be.db.GetModuleOwnerTypeRepository(), be.db.GetDeveloperModuleOwnerRepository())

	inst, err := services.NewModuleService(be.db.GetModuleRepository(), be.developerService, be.properties, be.moduleOwnershipService, be.db.GetKeywordRepository())
	if err != nil {
		return err
	}
	be.moduleService = inst

	/*	moduleOwnerRepo repositories.ModuleOwnerRepository,
		moduleOwnerTypeRepo repositories.ModuleOwnerTypesRepository,
		devModuleOwnerRepo repositories.DeveloperModuleOwnerRepository
	*/
	return nil
}

func (be *EvaModuleRepositoryBackend) GetModuleService() *services.ModuleService {
	return be.moduleService
}

func (be *EvaModuleRepositoryBackend) GetReleaseService() *services.ReleaseService {
	return be.releaseService
}

func (be *EvaModuleRepositoryBackend) GetDeveloperService() *services.DeveloperService {
	return be.developerService
}

func (be *EvaModuleRepositoryBackend) GetModuleOwnershipService() *services.ModuleOwnershipService {
	return be.moduleOwnershipService
}

/*

type ModuleReleaseStatus struct {
    ID   uint   `gorm:"primaryKey"`
    Name string `gorm:"unique;not null"`
}



*/
