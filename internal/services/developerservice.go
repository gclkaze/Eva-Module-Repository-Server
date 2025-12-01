package services

import (
	"github.com/gclkaze/evamodulerepositoryserver/internal/models"
	"github.com/gclkaze/evamodulerepositoryserver/internal/repositories"
	"github.com/gclkaze/evamodulerepositoryserver/pkg/logger"
	"github.com/gclkaze/evamodulerepositoryserver/pkg/runtime"
	"github.com/magiconair/properties"
)

type DeveloperService struct {
	repo        repositories.DeveloperRepository
	accountRepo repositories.DeveloperAccountRepository
	logger      logger.ILogger
}

func NewDeveloperService(repo repositories.DeveloperRepository, accountRepo repositories.DeveloperAccountRepository, p *properties.Properties) *DeveloperService {
	l := runtime.CreateLogger(p)
	mod := &DeveloperService{repo: repo, accountRepo: accountRepo}
	mod.logger = l
	return mod
}

/*
Handle           string           `json:"handle"`
FirstName        string           `json:"first_name"`
LastName         string           `json:"last_name"`
DeveloperID      uint             `json:"developer_id"`
DeveloperAccount DeveloperAccount `gorm:"foreignKey:DeveloperID"`
Active           bool             `json:"is_active"`
*/
func (s *DeveloperService) Create(handle string, firstName string, lastName string, active bool) (uint, error) {
	var dev *models.Developer
	var account *models.DeveloperAccount
	account = models.NewDeveloperAccount()

	err := s.accountRepo.Create(account)
	if err != nil {
		return 0, err
	}

	dev = models.NewDeveloper(handle, firstName, lastName, account.ID, *account, active)
	err = s.repo.Create(dev)
	if err != nil {
		return 0, err
	}
	return dev.ID, nil
}

func (s *DeveloperService) FindById(id uint) (*models.Developer, error) {
	return s.repo.FindByID(id)
}
