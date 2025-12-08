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

	moduleRepo          repositories.ModuleRepository
	releaseRepo         repositories.ReleaseRepository
	releaseStatusRepo   repositories.ReleaseStatusRepository
	developerRepo       repositories.DeveloperRepository
	userAccountRepo     repositories.UserAccountRepository
	moduleOwnerRepo     repositories.ModuleOwnerRepository
	devModuleOwnerRepo  repositories.DeveloperModuleOwnerRepository
	moduleOwnerTypeRepo repositories.ModuleOwnerTypesRepository
	keywordRepo         repositories.KeywordRepository
	userPermissionRepo  repositories.UserPermissionRepository
	userRoleRepo        repositories.UserRoleRepository
	devAccountRepo      repositories.DeveloperRepository
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
		&models.UserPermission{},
		&models.UserRole{},
		&models.UserAccount{},
		&models.Developer{},
		&models.Module{},
		&models.ModuleRelease{},
		&models.DeveloperModuleOwner{},
		&models.Keyword{},
	)

	db.initializeRepositories()
	error = db.initializeData()
	return error
}

func (db *EvaModuleRepositoryDatabase) initializeData() error {
	err := db.releaseStatusRepo.Initialize()
	if err != nil {
		return err
	}
	err = db.moduleOwnerTypeRepo.Initialize()
	return err
}

func (db *EvaModuleRepositoryDatabase) initializeRepositories() error {
	db.moduleRepo = repositories.NewModuleRepository(db.db)
	db.releaseRepo = repositories.NewReleaseRepository(db.db)
	db.releaseStatusRepo = repositories.NewReleaseStatusRepository(db.db)
	db.developerRepo = repositories.NewDeveloperRepository(db.db)
	db.userAccountRepo = repositories.NewUserAccountRepository(db.db)
	db.devAccountRepo = repositories.NewDeveloperRepository(db.db)
	db.moduleOwnerRepo = repositories.NewModuleOwnerRepository(db.db)
	db.devModuleOwnerRepo = repositories.NewDeveloperModuleOwnerRepository(db.db)
	db.moduleOwnerTypeRepo = repositories.NewModuleOwnerTypesRepository(db.db)
	db.keywordRepo = repositories.NewKeywordRepository(db.db)
	db.userPermissionRepo = repositories.NewUserPermissionRepository(db.db) //repositories.UserPermissionRepository
	db.userRoleRepo = repositories.NewUserRoleRepository(db.db)             //repositories.UserRoleRepository
	return nil
}

func (db EvaModuleRepositoryDatabase) GetModuleRepository() repositories.ModuleRepository {
	return db.moduleRepo
}

func (db EvaModuleRepositoryDatabase) GetReleaseRepository() repositories.ReleaseRepository {
	return db.releaseRepo
}

func (db EvaModuleRepositoryDatabase) GetReleaseStatusRepository() repositories.ReleaseStatusRepository {
	return db.releaseStatusRepo
}

func (db EvaModuleRepositoryDatabase) GetDeveloperRepository() repositories.DeveloperRepository {
	return db.developerRepo
}

func (db EvaModuleRepositoryDatabase) GetUserAccountRepository() repositories.UserAccountRepository {
	return db.userAccountRepo
}

func (db EvaModuleRepositoryDatabase) GetDevAccountRepository() repositories.DeveloperRepository {
	return db.devAccountRepo
}

func (db EvaModuleRepositoryDatabase) GetModuleOwnerRepository() repositories.ModuleOwnerRepository {
	return db.moduleOwnerRepo
}

func (db EvaModuleRepositoryDatabase) GetDeveloperModuleOwnerRepository() repositories.DeveloperModuleOwnerRepository {
	return db.devModuleOwnerRepo
}

func (db EvaModuleRepositoryDatabase) GetModuleOwnerTypeRepository() repositories.ModuleOwnerTypesRepository {
	return db.moduleOwnerTypeRepo
}

func (db EvaModuleRepositoryDatabase) GetKeywordRepository() repositories.KeywordRepository {
	return db.keywordRepo
}

func (db EvaModuleRepositoryDatabase) GetUserPermissionRepository() repositories.UserPermissionRepository {
	return db.userPermissionRepo
}

func (db EvaModuleRepositoryDatabase) GetUserRoleRepository() repositories.UserRoleRepository {
	return db.userRoleRepo
}
