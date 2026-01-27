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
	userService            *services.UserService
	moduleOwnershipService *services.ModuleOwnershipService
	authService            *services.AuthService
	downloadService        *services.DownloadService
}

func NewEvaModuleRepositoryBackend() *EvaModuleRepositoryBackend {
	inst := &EvaModuleRepositoryBackend{}
	return inst
}

func (be EvaModuleRepositoryBackend) GetJWTSecret() string {
	return be.authService.GetJWTSecret()
}

func (be *EvaModuleRepositoryBackend) GetProperties() *properties.Properties {
	return be.properties
}

func (be *EvaModuleRepositoryBackend) CleanDB() {
	be.db.CleanDB()
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
	be.moduleOwnershipService = services.NewModuleOwnershipService(be.db.GetModuleOwnerRepository(), be.db.GetModuleOwnerTypeRepository(), be.db.GetDeveloperModuleOwnerRepository())

	be.userService = services.NewUserService(be.db.GetDeveloperRepository(), be.db.GetUserAccountRepository(),
		be.db.GetUserPermissionRepository(), be.db.GetUserRoleRepository(),
		be.properties)

	be.releaseService = services.NewReleaseService(be.db.GetReleaseRepository(), be.db.GetReleaseStatusRepository(), be.moduleOwnershipService, be.GetUserService(), be.properties, be.db.GetReleaseStatisticsRepository())

	inst, err := services.NewModuleService(be.db.GetModuleRepository(), be.userService, be.properties, be.moduleOwnershipService, be.db.GetKeywordRepository(), be.releaseService)
	if err != nil {
		return err
	}
	be.moduleService = inst
	err = be.userService.Initialize()

	be.downloadService = services.NewDownloadService(be.moduleService, be.properties)
	be.authService = services.NewAuthService(be.db.GetUserAccountRepository(), be.db.GetDevAccountRepository(), be.userService, be.properties)
	return err
}

func (be *EvaModuleRepositoryBackend) GetModuleService() *services.ModuleService {
	return be.moduleService
}

func (be *EvaModuleRepositoryBackend) GetDownloadService() *services.DownloadService {
	return be.downloadService
}

func (be *EvaModuleRepositoryBackend) GetReleaseService() *services.ReleaseService {
	return be.releaseService
}

func (be *EvaModuleRepositoryBackend) GetUserService() *services.UserService {
	return be.userService
}

func (be *EvaModuleRepositoryBackend) GetModuleOwnershipService() *services.ModuleOwnershipService {
	return be.moduleOwnershipService
}
func (be *EvaModuleRepositoryBackend) GetAuthService() *services.AuthService {
	return be.authService
}
