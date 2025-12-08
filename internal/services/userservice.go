package services

import (
	"github.com/gclkaze/evamodulerepositoryserver/internal/models"
	"github.com/gclkaze/evamodulerepositoryserver/internal/repositories"
	"github.com/gclkaze/evamodulerepositoryserver/pkg/logger"
	"github.com/gclkaze/evamodulerepositoryserver/pkg/runtime"
	"github.com/magiconair/properties"
	"golang.org/x/crypto/bcrypt"
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

func (s *UserService) FindByID(id uint) (*models.Developer, error) {
	return s.repo.FindByID(id)
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
			{Value: models.CreateModule.String()},
			{Value: models.SuggestModule.String()},
			{Value: models.DeleteModule.String()},
			{Value: models.DeleteMyModule.String()},
			{Value: models.DeleteRelease.String()},
			{Value: models.DeleteMyRelease.String()},
			{Value: models.UpdateRelease.String()},
			{Value: models.ChangeReleaseStatus.String()},
			{Value: models.RejectRelease.String()},
			{Value: models.AcceptRelease.String()},
			{Value: models.CancelRelease.String()},
			{Value: models.BanUser.String()},
			{Value: models.UnbanUser.String()},
		},
		models.Maintainer: {
			{Value: models.CreateModule.String()},
			{Value: models.SuggestModule.String()},
			{Value: models.DeleteModule.String()},
			{Value: models.DeleteMyModule.String()},
			{Value: models.DeleteRelease.String()},
			{Value: models.DeleteMyRelease.String()},
			{Value: models.UpdateRelease.String()},
			{Value: models.ChangeReleaseStatus.String()},
			{Value: models.RejectRelease.String()},
			{Value: models.AcceptRelease.String()},
			{Value: models.CancelRelease.String()},
		},
		models.User: {
			{Value: models.CreateModule.String()},
			{Value: models.SuggestModule.String()},
			{Value: models.DeleteModule.String()},
			{Value: models.DeleteMyModule.String()},
			{Value: models.DeleteRelease.String()},
			{Value: models.DeleteMyRelease.String()},
			{Value: models.UnbanUser.String()},
		},
	}

	p, ok := m[t]
	if !ok {
		return nil
	}
	return p
}
