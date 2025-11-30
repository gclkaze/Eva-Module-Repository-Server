package routes

import (
	"fmt"

	"github.com/gclkaze/evamodulerepositoryserver/internal/backend"
	"github.com/gclkaze/evamodulerepositoryserver/internal/config"
	"github.com/gclkaze/evamodulerepositoryserver/internal/handlers"

	"github.com/gin-gonic/gin"
)

type EvaModuleRepositoryRouter struct {
	api  *gin.RouterGroup
	r    *gin.Engine
	port string

	moduleHandler  *handlers.ModuleHandler
	releaseHandler *handlers.ReleaseHandler
}

func NewEvaModuleRepositoryRouter() *EvaModuleRepositoryRouter {
	inst := &EvaModuleRepositoryRouter{}
	return inst
}

func (router *EvaModuleRepositoryRouter) Initialize(r *gin.Engine, be *backend.EvaModuleRepositoryBackend) error {

	router.api = r.Group("/api")

	router.moduleHandler = handlers.NewModuleHandler(be.GetModuleService())

	router.releaseHandler = handlers.NewReleaseHandler(be.GetReleaseService())

	modules := router.api.Group("/modules")
	{
		modules.GET("/:id", router.moduleHandler.FindByID) // GET /api/modules/:id
		modules.GET("/search", router.moduleHandler.SearchModulesByTags)
	}

	releases := router.api.Group("releases")
	{
		releases.GET("/:id/releases", router.releaseHandler.GetModuleReleases)          // GET /api/releases/:id
		releases.GET("/:id/release/:releaseId", router.releaseHandler.GetModuleRelease) // GET /api/releases/:id/release/:releaseId
	}
	// /api/releases/:id/

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
