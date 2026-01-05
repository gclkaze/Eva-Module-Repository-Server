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

package server

import (
	"fmt"
	"path"

	"github.com/gclkaze/evamodulerepositoryserver/internal/backend"
	"github.com/gclkaze/evamodulerepositoryserver/internal/config"
	"github.com/gclkaze/evamodulerepositoryserver/internal/routes"
	"github.com/gclkaze/evamodulerepositoryserver/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/magiconair/properties"
)

type EvaModuleRepositoryServer struct {
	be     *backend.EvaModuleRepositoryBackend
	router *routes.EvaModuleRepositoryRouter

	moduleFolder    string
	releaseFolder   string
	developerFolder string
}

func NewEvaModuleRepositoryServer() *EvaModuleRepositoryServer {
	inst := &EvaModuleRepositoryServer{}
	return inst
}

func (inst *EvaModuleRepositoryServer) GetBackend() *backend.EvaModuleRepositoryBackend {
	return inst.be
}

func (inst *EvaModuleRepositoryServer) ClearModuleFolders() {
	p := inst.be.GetProperties()

	inst.moduleFolder = p.GetString("module_folder", "")
	if inst.moduleFolder != "" && utils.FolderExists(inst.moduleFolder) {
		utils.CleanFolder(inst.moduleFolder)
		utils.RemoveFolder(inst.moduleFolder)
	}
	inst.releaseFolder = p.GetString("release_folder", "")
	if inst.releaseFolder != "" && utils.FolderExists(inst.releaseFolder) {
		utils.CleanFolder(inst.releaseFolder)
		utils.RemoveFolder(inst.releaseFolder)
	}

	inst.developerFolder = p.GetString("dev_folder", "")
	if inst.developerFolder != "" {
		inst.developerFolder = path.Join(inst.moduleFolder, inst.developerFolder)
	}
	if inst.developerFolder != "" && utils.FolderExists(inst.developerFolder) {
		utils.CleanFolder(inst.developerFolder)
		utils.RemoveFolder(inst.developerFolder)
	}
}

func (inst *EvaModuleRepositoryServer) setupRepositoryFolders(p *properties.Properties) error {
	inst.moduleFolder = p.GetString("module_folder", "")
	inst.releaseFolder = p.GetString("release_folder", "")
	inst.developerFolder = p.GetString("dev_folder", "")
	if inst.developerFolder != "" {
		inst.developerFolder = path.Join(inst.moduleFolder, inst.developerFolder)
	}

	err := utils.CreateFolder(inst.moduleFolder)
	if err != nil {
		return err
	}
	err = utils.CreateFolder(inst.developerFolder)
	if err != nil {
		return err
	}

	err = utils.CreateFolder(inst.releaseFolder)
	return err
}

func (inst EvaModuleRepositoryServer) GetDeveloperModuleFolder(devID uint, modID uint) string {
	p := path.Join(inst.developerFolder, utils.UintToString(devID), utils.UintToString(modID))
	return p
}

func (inst EvaModuleRepositoryServer) GetDeveloperFolder(devID uint) string {
	p := path.Join(inst.developerFolder, utils.UintToString(devID))
	return p
}

func (inst EvaModuleRepositoryServer) GetReleaseBasePathWithName(name string) string {
	return inst.be.GetModuleService().GetReleaseBasePathWithName(utils.GetRepoName(name))
}

func (inst *EvaModuleRepositoryServer) SetUploadFileLimit(limit int64) int64 {
	old := inst.router.GetRouter().MaxMultipartMemory
	inst.router.SetUploadFileLimit(limit)
	inst.router.GetRouter().Use(inst.router.GetMiddleware().MaxBodySize(limit))
	return old
}

func (inst EvaModuleRepositoryServer) GetReleaseBasePathWithNameAndVersion(name string, version string) string {
	return inst.be.GetModuleService().GetReleaseBasePathWithNameAndVersion(utils.GetRepoName(name), version)
}
func (inst *EvaModuleRepositoryServer) InitializeWithPropertiesPath(prop string) error {
	config.InitWithPropertiesPath(prop)

	inst.be = backend.NewEvaModuleRepositoryBackend()
	inst.router = routes.NewEvaModuleRepositoryRouter()
	error := inst.be.Initialize()
	if error != nil {
		return error
	}

	r := gin.Default()
	error = inst.router.Initialize(r, inst.be)
	if error != nil {
		return error
	}

	error = inst.setupRepositoryFolders(inst.be.GetProperties())
	return error
}

func (inst *EvaModuleRepositoryServer) ResetRouterWithUploadLimit(prop string, limit int64) (int64, error) {
	config.InitWithPropertiesPath(prop)

	inst.be = backend.NewEvaModuleRepositoryBackend()
	inst.router = routes.NewEvaModuleRepositoryRouter()
	old := inst.router.SetUploadFileLimit(limit)
	error := inst.be.Initialize()
	if error != nil {
		return old, error
	}

	r := gin.Default()
	error = inst.router.Initialize(r, inst.be)
	if error != nil {
		return old, error
	}

	error = inst.setupRepositoryFolders(inst.be.GetProperties())
	return old, error
}

func (inst *EvaModuleRepositoryServer) CleanDB() {
	inst.be.CleanDB()
}

func (inst *EvaModuleRepositoryServer) InitializeWithPropertiesMap(m *map[string]string) error {
	config.InitWithPropertiesMap(m)

	inst.be = backend.NewEvaModuleRepositoryBackend()
	inst.router = routes.NewEvaModuleRepositoryRouter()
	error := inst.be.Initialize()
	if error != nil {
		return error
	}

	r := gin.Default()
	error = inst.router.Initialize(r, inst.be)
	if error != nil {
		return error
	}
	error = inst.setupRepositoryFolders(inst.be.GetProperties())
	return error
}

func (inst *EvaModuleRepositoryServer) Initialize() error {
	config.Init()

	inst.be = backend.NewEvaModuleRepositoryBackend()
	inst.router = routes.NewEvaModuleRepositoryRouter()
	error := inst.be.Initialize()
	if error != nil {
		return error
	}

	r := gin.Default()
	error = inst.router.Initialize(r, inst.be)
	if error != nil {
		return error
	}
	error = inst.setupRepositoryFolders(inst.be.GetProperties())
	return error
}

func (inst *EvaModuleRepositoryServer) Run() error {
	if inst.router == nil {
		return fmt.Errorf("the router is uninitialized..call Initialize() first")
	}
	inst.router.Run()
	return nil
}

func (inst *EvaModuleRepositoryServer) GetRouter() *gin.Engine {
	if inst.router == nil {
		return nil
	}
	return inst.router.GetRouter()
}
