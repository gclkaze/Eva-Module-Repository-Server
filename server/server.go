package server

import (
	"fmt"

	"github.com/gclkaze/evamodulerepositoryserver/backend"
	"github.com/gclkaze/evamodulerepositoryserver/config"
	"github.com/gclkaze/evamodulerepositoryserver/routes"
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

func (inst *EvaModuleRepositoryServer) Initialize() error {
	config.Init()

	inst.be = backend.NewEvaModuleRepositoryBackend()
	inst.router = routes.NewEvaModuleRepositoryRouter()
	error := inst.be.Initialize()
	if error != nil {
		return error
	}

	r := gin.Default()
	error = inst.router.Initialize(r)
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
