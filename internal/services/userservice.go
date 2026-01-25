// Copyright (c) 2025 Michail Dorgiakis - gclkaze@gmail.com - https://github.com/gclkaze
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package services

import (
	"fmt"

	"github.com/gclkaze/evamodulerepositoryserver/internal/dto"
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
	p              *properties.Properties
}

func NewUserService(repo repositories.DeveloperRepository, accountRepo repositories.UserAccountRepository,
	permissionRepo repositories.UserPermissionRepository, roleRepo repositories.UserRoleRepository,
	p *properties.Properties) *UserService {
	l := runtime.CreateLogger(p)
	mod := &UserService{repo: repo, accountRepo: accountRepo, permissionRepo: permissionRepo, roleRepo: roleRepo, p: p}
	mod.logger = l
	return mod
}

func (s UserService) GetDeveloperRepository() repositories.DeveloperRepository {
	return s.repo
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
	//the email and the handle need to be unique
	hFound, err := s.repo.FindByHandle(handle)
	if err != nil {
		return 0, err
	}
	if hFound != nil {
		return 0, fmt.Errorf("the handle is used")
	}

	eFound, err := s.accountRepo.FindByEmail(email)
	if err != nil {
		return 0, err
	}
	if eFound != nil {
		return 0, fmt.Errorf("the email is used")
	}

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

func (s *UserService) GetFirstWithRole(roleString string) (*models.UserAccount, error) {
	role, err := s.roleRepo.FindByValue(roleString)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, nil
	}
	ua, err := s.accountRepo.GetFirstWithRole(role)
	return ua, err
}

func (s *UserService) CreateUserWithRole(handle string, firstName string, lastName string, email string, password string, active bool, roleString string) (uint, error) {
	var dev *models.Developer

	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	password = string(hash)

	role, err := s.roleRepo.FindByValue(roleString)
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

func (s *UserService) CreateUserFromDTO(dto *dto.UserAccountDTO) (uint, error) {
	return s.CreateUserWithRole(dto.Handle, dto.FirstName, dto.LastName, dto.Email, dto.Password, dto.Active, dto.UserRole)
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
	if user != nil && user.IsBanned {
		s.logger.Errorf("user service", "banned user tries to access with id %d", id)
		return nil, fmt.Errorf("unknown user with id %d", id)
	}
	if user == nil {
		s.logger.Errorf("user service", "unknown user with id %d", id)
		return nil, fmt.Errorf("unknown user with id %d", id)
	}
	role, err := s.roleRepo.FindByID(user.RoleID)
	if err != nil {
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
	if err != nil {
		return err
	}

	return s.initializeDefaultUsers()
}

func (s *UserService) initializeDefaultUsers() error {
	var users []*dto.UserAccountDTO

	defaultPassword := s.p.GetString("default_password", "thisisapass")
	users = append(users, dto.NewUserAccountDTO(1, "gclkaze", "gcl", "kaze", "gclkaze@gmail.com", defaultPassword, true, models.Admin.String()))
	users = append(users, dto.NewUserAccountDTO(2, "mdor", "michail", "dorgiakis", "michail.dorgiakis@gmail.com", defaultPassword, true, models.User.String()))

	for i := range users {
		email := users[i].Email
		user, err := s.accountRepo.FindByEmail(email)
		if err != nil {
			s.logger.Errorf("user service", "An error has occurred while searcing for user %s : %s", email, err.Error())
			return err
		}

		if user == nil {
			//lets add the user
			_, err = s.CreateUserFromDTO(users[i])
			if err != nil {
				s.logger.Errorf("user service", "An error has occurred while adding user %s : %s", email, err.Error())
				return err
			}
		}
	}
	return nil
}

func (s UserService) FindByEmail(email string) (*models.UserAccount, error) {
	return s.accountRepo.FindByEmail(email)
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

func (s *UserService) BanUser(userID uint) error {
	dev, err := s.accountRepo.FindByID(userID)
	if err != nil {
		return err
	}
	if dev == nil {
		return fmt.Errorf("unknown user with id %d", userID)
	}
	if dev.IsBanned {
		return fmt.Errorf("user with id %d is already banned", userID)
	}
	return s.accountRepo.BanUser(userID)
}

func (s *UserService) UnbanUser(userID uint) error {
	dev, err := s.accountRepo.FindByID(userID)
	if err != nil {
		return err
	}
	if dev == nil {
		return fmt.Errorf("unknown user with id %d", userID)
	}
	if !dev.IsBanned {
		return fmt.Errorf("user with id %d is not banned", userID)
	}
	return s.accountRepo.UnbanUser(userID)
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
			{Value: models.DeleteMyModule.String()},
			{Value: models.DeleteMyRelease.String()},
		},
	}

	p, ok := m[t]
	if !ok {
		return nil
	}
	return p
}
