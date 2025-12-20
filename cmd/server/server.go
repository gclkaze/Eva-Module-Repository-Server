package server

import (
	"fmt"

	"github.com/gclkaze/evamodulerepositoryserver/internal/backend"
	"github.com/gclkaze/evamodulerepositoryserver/internal/config"
	"github.com/gclkaze/evamodulerepositoryserver/internal/routes"
	"github.com/gclkaze/evamodulerepositoryserver/pkg/utils"
	"github.com/gin-gonic/gin"
)

type EvaModuleRepositoryServer struct {
	be     *backend.EvaModuleRepositoryBackend
	router *routes.EvaModuleRepositoryRouter
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

	moduleFolder := p.GetString("module_folder", "")
	if moduleFolder != "" && utils.FolderExists(moduleFolder) {
		utils.CleanFolder(moduleFolder)
	}
	releaseFolder := p.GetString("release_folder", "")
	if releaseFolder != "" && utils.FolderExists(releaseFolder) {
		utils.CleanFolder(releaseFolder)
	}

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
	return nil
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
	return nil
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
	return nil
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
