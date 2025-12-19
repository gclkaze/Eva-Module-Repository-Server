package server

import (
	"fmt"

	"github.com/gclkaze/evamodulerepositoryserver/internal/backend"
	"github.com/gclkaze/evamodulerepositoryserver/internal/config"
	"github.com/gclkaze/evamodulerepositoryserver/internal/routes"
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

func (inst *EvaModuleRepositoryServer) InitializeWithProperties(prop string) error {
	config.InitWithProperties(prop)

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
