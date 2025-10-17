package backend

import (
	"fmt"

	"github.com/gclkaze/evamodulerepositoryserver/config"
	"github.com/gclkaze/evamodulerepositoryserver/db"
	"github.com/gclkaze/evamodulerepositoryserver/logger"
	"github.com/magiconair/properties"
)

type EvaModuleRepositoryBackend struct {
	db         *db.EvaModuleRepositoryDatabase
	properties *properties.Properties
	logger     logger.ILogger
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
	return error
}
