package routes

import (
	"fmt"

	"github.com/gclkaze/evamodulerepositoryserver/config"
	"github.com/gclkaze/evamodulerepositoryserver/controllers"

	"github.com/gin-gonic/gin"
)

type EvaModuleRepositoryRouter struct {
	api  *gin.RouterGroup
	r    *gin.Engine
	port string
}

func NewEvaModuleRepositoryRouter() *EvaModuleRepositoryRouter {
	inst := &EvaModuleRepositoryRouter{}
	return inst
}

func (router *EvaModuleRepositoryRouter) Initialize(r *gin.Engine) error {
	router.api = r.Group("/api")
	{
		router.api.GET("/moduleSearch", controllers.SearchModules)
		//api.POST("/albums", controllers.CreateAlbum)
	}

	if config.TheConfigReader.IsOnError() {
		return fmt.Errorf("couldn't read the properties file")
	}

	p := config.TheConfigReader.GetProperties()
	thePort := p.GetString("server_port", "")
	if thePort == "" {
		return fmt.Errorf("no server_port declared")
	}
	router.port = thePort
	router.r = r
	return nil
}

func (router *EvaModuleRepositoryRouter) Run() error {
	if router.r == nil {
		return fmt.Errorf("the engine is uninitialized..call Initialize() first")
	}
	router.r.Run(router.port)
	return nil
}
