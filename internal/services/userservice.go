package services

import (
	"github.com/gclkaze/evamodulerepositoryserver/internal/models"
	"github.com/gclkaze/evamodulerepositoryserver/internal/repositories"
	"github.com/gclkaze/evamodulerepositoryserver/pkg/logger"
	"github.com/gclkaze/evamodulerepositoryserver/pkg/runtime"
	"github.com/magiconair/properties"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	repo           repositories.DeveloperRepository
	accountRepo    repositories.UserAccountRepository
	permissionRepo repositories.UserPermissionRepository
	roleRepo       repositories.UserRoleRepository
	logger         logger.ILogger
}

func NewUserService(repo repositories.DeveloperRepository, accountRepo repositories.UserAccountRepository,
	permissionRepo repositories.UserPermissionRepository, roleRepo repositories.UserRoleRepository,
	p *properties.Properties) *UserService {
	l := runtime.CreateLogger(p)
	mod := &UserService{repo: repo, accountRepo: accountRepo, permissionRepo: permissionRepo, roleRepo: roleRepo}
	mod.logger = l
	return mod
}

func (s UserService) GetDevelopersUserAccount(d *models.Developer) (*models.UserAccount, error) {
	r, err := s.repo.FindByUserAccountID(d.UserID)
	if err != nil {
		return nil, err
	}
	return &r.UserAccount, nil
}

func (s *UserService) Create(handle string, firstName string, lastName string, email string, password string, active bool, role *models.UserRole) (uint, error) {
	var dev *models.Developer

	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	password = string(hash)

	account := models.NewUserAccount(role, email, password)

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

func (s *UserService) CreateUser(handle string, firstName string, lastName string, email string, password string, active bool) (uint, error) {
	var dev *models.Developer

	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	password = string(hash)

	role, err := s.roleRepo.FindByValue(models.User.String())
	if err != nil {
		return 0, err
	}

	account := models.NewUserAccount(role, email, password)

	err = s.accountRepo.Create(account)
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

func (s *UserService) FindByID(id uint) (*models.Developer, error) {
	return s.repo.FindByID(id)
}

func (s *UserService) FindByIDTx(tx *gorm.DB, id uint) (*models.Developer, error) {
	return s.repo.FindByIDTx(tx, id)
}

func (s *UserService) FindUserByID(id uint) (*models.UserAccount, error) {
	return s.accountRepo.FindByID(id)
}

func (s *UserService) GetUserPermissions(id uint) ([]models.UserPermission, error) {
	user, err := s.accountRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	role, err := s.roleRepo.FindByID(user.RoleID)
	if role != nil {
		return nil, err
	}
	return role.Permissions, nil
}

func (s UserService) UserHasPermission(id uint, perm models.UserPermissionTypeDef) bool {
	ps, err := s.GetUserPermissions(id)
	if err != nil {
		return false
	}
	for i := range ps {
		if ps[i].Value == perm.String() {
			return true
		}
	}

	return false
}

func (s *UserService) Initialize() error {
	//need to check for initialization of the user permissions, user roles, and default test user accounts. 1 admin, 2 developer
	err := s.permissionRepo.Initialize()
	if err != nil {
		return err
	}
	err = s.roleRepo.Initialize()
	if err != nil {
		return err
	}
	err = s.InitializeUserRolePermissions()
	return err
}

func (s *UserService) InitializeUserRolePermissions() error {
	for _, t := range models.GetUserRoleTypes() {
		//var res models.UserRole
		perms := s.getRolePermissions(t)
		var storedPerms []models.UserPermission
		for i := range perms {
			p, err := s.permissionRepo.FindByValue(perms[i].Value)
			if err != nil {
				return err
			}
			storedPerms = append(storedPerms, *p)
		}
		role, err := s.roleRepo.FindByValue(t.String())
		if err != nil {
			return err
		}

		role.Permissions = storedPerms
		s.roleRepo.Update(role)

	}
	return nil
}

func (s UserService) getRolePermissions(t models.UserRoleTypeDef) []models.UserPermission {
	m := map[models.UserRoleTypeDef][]models.UserPermission{
		models.Admin: {
			{Value: models.CreateMyModule.String()},
			{Value: models.SuggestMyModule.String()},
			{Value: models.DeleteModules.String()},
			{Value: models.DeleteMyModule.String()},
			{Value: models.DeleteReleases.String()},
			{Value: models.DeleteMyRelease.String()},
			{Value: models.UpdateReleases.String()},
			{Value: models.ChangeReleaseStatuses.String()},
			{Value: models.RejectReleases.String()},
			{Value: models.AcceptReleases.String()},
			{Value: models.CancelReleases.String()},
			{Value: models.BanUsers.String()},
			{Value: models.UnbanUsers.String()},
		},
		models.Maintainer: {
			{Value: models.CreateMyModule.String()},
			{Value: models.SuggestMyModule.String()},
			{Value: models.DeleteModules.String()},
			{Value: models.DeleteMyModule.String()},
			{Value: models.DeleteReleases.String()},
			{Value: models.DeleteMyRelease.String()},
			{Value: models.UpdateReleases.String()},
			{Value: models.ChangeReleaseStatuses.String()},
			{Value: models.RejectReleases.String()},
			{Value: models.AcceptReleases.String()},
			{Value: models.CancelReleases.String()},
		},
		models.User: {
			{Value: models.CreateMyModule.String()},
			{Value: models.SuggestMyModule.String()},
			{Value: models.DeleteModules.String()},
			{Value: models.DeleteMyModule.String()},
			{Value: models.DeleteReleases.String()},
			{Value: models.DeleteMyRelease.String()},
			{Value: models.UnbanUsers.String()},
		},
	}

	p, ok := m[t]
	if !ok {
		return nil
	}
	return p
}
