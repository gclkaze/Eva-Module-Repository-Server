package db

import (
	"fmt"

	"github.com/gclkaze/evamodulerepositoryserver/logger"
	"github.com/gclkaze/evamodulerepositoryserver/models"
	"github.com/gclkaze/evamodulerepositoryserver/runtime"
	"github.com/magiconair/properties"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type EvaModuleRepositoryDatabase struct {
	logger logger.ILogger
	c      *EvaModuleRepositoryDatabaseConfig
	db     *gorm.DB
}

func NewEvaModuleRepositoryDatabase() *EvaModuleRepositoryDatabase {
	inst := &EvaModuleRepositoryDatabase{}
	return inst
}

func (db *EvaModuleRepositoryDatabase) Initialize(p *properties.Properties) error {
	db.logger = runtime.CreateLogger(p)

	if p == nil {
		return fmt.Errorf("couln't read application properties")
	}

	db.c = NewEvaModuleRepositoryDatabaseConfig()
	error := db.c.LoadFromProperties(p)
	if error != nil {
		return error
	}
	connectionString := db.c.GetConnectionString()
	db.db, error = gorm.Open(mysql.Open(connectionString), &gorm.Config{})
	if error != nil {
		return fmt.Errorf("failed to connect to database: ", error)
	}

	// Auto migrate your models
	db.db.AutoMigrate(&models.Module{})
	return nil
}
