package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gclkaze/evamodulerepositoryserver/internal/models"
	"github.com/gclkaze/evamodulerepositoryserver/internal/repositories"
	"github.com/gclkaze/evamodulerepositoryserver/pkg/logger"
	"github.com/gclkaze/evamodulerepositoryserver/pkg/runtime"
	"github.com/magiconair/properties"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
)

type EvaModuleRepositoryDatabase struct {
	logger logger.ILogger
	c      *EvaModuleRepositoryDatabaseConfig
	db     *gorm.DB

	moduleRepo  repositories.ModuleRepository
	releaseRepo repositories.ReleaseRepository
}

func NewEvaModuleRepositoryDatabase() *EvaModuleRepositoryDatabase {
	inst := &EvaModuleRepositoryDatabase{}
	return inst
}

func (db EvaModuleRepositoryDatabase) getConfig(p *properties.Properties) *gorm.Config {
	sqlDebug := p.GetBool("sql_debug", false)
	if !sqlDebug {
		return &gorm.Config{}
	}
	newLogger := glogger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // Writer
		glogger.Config{
			SlowThreshold: time.Second,  // Log queries slower than this
			LogLevel:      glogger.Info, // Log level: Silent, Error, Warn, Info
			Colorful:      true,         // Enable color
		},
	)
	return &gorm.Config{
		Logger: newLogger,
	}

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

	db.db, error = gorm.Open(mysql.Open(connectionString), db.getConfig(p))
	if error != nil {
		return fmt.Errorf("failed to connect to database: ", error)
	}

	// Auto migrate your models
	db.db.AutoMigrate(
		&models.ModuleReleaseStatus{},
		&models.ModuleOwnerType{},
		&models.ModuleOwner{},
		&models.DeveloperAccount{},
		&models.Developer{},
		&models.Module{},
		&models.ModuleRelease{},
		&models.DeveloperModuleOwner{},
		&models.Keyword{},
	)

	db.initializeRepositories()
	return nil
}

func (db *EvaModuleRepositoryDatabase) initializeRepositories() error {
	db.moduleRepo = repositories.NewModuleRepository(db.db)
	db.releaseRepo = repositories.NewReleaseRepository(db.db)
	return nil
}

func (db EvaModuleRepositoryDatabase) GetModuleRepository() repositories.ModuleRepository {
	return db.moduleRepo
}

func (db EvaModuleRepositoryDatabase) GetReleaseRepository() repositories.ReleaseRepository {
	return db.releaseRepo
}
