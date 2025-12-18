package routes

import (
	"fmt"

	"github.com/gclkaze/evamodulerepositoryserver/internal/backend"
	"github.com/gclkaze/evamodulerepositoryserver/internal/config"
	"github.com/gclkaze/evamodulerepositoryserver/internal/handlers"
	"github.com/gclkaze/evamodulerepositoryserver/internal/middleware"
	"github.com/gclkaze/evamodulerepositoryserver/internal/models"

	"github.com/gin-gonic/gin"
)

type EvaModuleRepositoryRouter struct {
	api  *gin.RouterGroup
	r    *gin.Engine
	port string

	moduleHandler   *handlers.ModuleHandler
	downloadHandler *handlers.DownloadHandler
	releaseHandler  *handlers.ReleaseHandler
	authHandler     *handlers.AuthHandler

	middleWare *middleware.AuthMiddleWare
}

func NewEvaModuleRepositoryRouter() *EvaModuleRepositoryRouter {
	inst := &EvaModuleRepositoryRouter{}
	return inst
}

func (router *EvaModuleRepositoryRouter) Initialize(r *gin.Engine, be *backend.EvaModuleRepositoryBackend) error {
	router.api = r.Group("/api")
	router.moduleHandler = handlers.NewModuleHandler(be.GetModuleService())
	router.releaseHandler = handlers.NewReleaseHandler(be.GetReleaseService())
	router.authHandler = handlers.NewAuthHandler(be.GetAuthService(), be.GetDeveloperService())
	router.downloadHandler = handlers.NewDownloadHandler(be.GetDownloadService())
	router.middleWare = middleware.NewAuthMiddleware(be.GetDeveloperService())

	r.MaxMultipartMemory = 8 << 20 // 8 MB

	authGroup := router.api.Group("/auth")
	{
		authGroup.POST("/login", router.authHandler.Login)
		authGroup.POST("/refresh", router.authHandler.Refresh)
		authGroup.POST("/register", router.authHandler.Register)
	}

	modules := router.api.Group("/modules")
	{
		//find a module based on id
		modules.GET("/:id", router.moduleHandler.FindByID) // GET /api/modules/:id
		modules.GET("/search", router.moduleHandler.SearchModulesByTags)

		//delete a module!
		modules.GET("/:id/delete",
			router.middleWare.AuthMiddleware(be.GetJWTSecret()),
			router.middleWare.PreAuthorize(router.middleWare.HasPermissions([]models.UserPermissionTypeDef{
				models.DeleteMyModule,
				models.CreateMyModule,
				models.DeleteModules,
			})),
			router.moduleHandler.Delete) //needs userID to be passed

		modules.POST("/upload", router.middleWare.AuthMiddleware(be.GetJWTSecret()),
			router.middleWare.PreAuthorize(router.middleWare.HasPermissions([]models.UserPermissionTypeDef{
				models.DeleteMyModule,
				models.CreateMyModule,
			})), router.moduleHandler.Upload)

		modules.POST("/update", router.middleWare.AuthMiddleware(be.GetJWTSecret()),
			router.middleWare.PreAuthorize(router.middleWare.HasPermissions([]models.UserPermissionTypeDef{
				models.DeleteMyModule,
				models.CreateMyModule,
				models.UpdateModules,
			})), router.moduleHandler.Update)

		modules.POST("/suggest", router.middleWare.AuthMiddleware(be.GetJWTSecret()),
			router.middleWare.PreAuthorize(router.middleWare.HasPermissions([]models.UserPermissionTypeDef{
				models.DeleteMyModule,
				models.CreateMyModule,
				models.SuggestMyModule,
			})), router.moduleHandler.SuggestRelease)
	}

	releases := router.api.Group("releases")
	{
		releases.GET("/:id", router.releaseHandler.GetModuleReleases)                   // GET /api/releases/:id
		releases.GET("/:id/release/:releaseId", router.releaseHandler.GetModuleRelease) // GET /api/releases/:id/release/:releaseId
		releases.GET("/:id/search", router.releaseHandler.SearchByKeywords)

		releases.POST("/:id/delete/:releaseId", router.middleWare.AuthMiddleware(be.GetJWTSecret()),
			router.middleWare.PreAuthorize(router.middleWare.HasPermissions([]models.UserPermissionTypeDef{
				models.DeleteMyModule,
				models.CreateMyModule,
				models.SuggestMyModule,
				models.DeleteReleases,
			})), router.releaseHandler.DeleteModuleRelease) // GET /api/releases/:id/delete/:releaseId  ->need to check the userID is the one that initiated it

		releases.POST("/:id/cancel/:releaseId", router.middleWare.AuthMiddleware(be.GetJWTSecret()),
			router.middleWare.PreAuthorize(router.middleWare.HasPermissions([]models.UserPermissionTypeDef{
				models.DeleteMyModule,
				models.CreateMyModule,
				models.SuggestMyModule,
				models.CancelReleases,
			})), router.releaseHandler.CancelSuggestedRelease) // ->need to check the userID is the one that initiated it
	}

	download := router.api.Group("download")
	{
		download.GET("/release/:releaseId", router.downloadHandler.DownloadRelease) //the release needs to be accepted
	}

	supervision := router.api.Group("supervise")
	{
		supervision.GET("/download/release/:releaseId", router.middleWare.AuthMiddleware(be.GetJWTSecret()), router.middleWare.PreAuthorize(router.middleWare.HasPermissions([]models.UserPermissionTypeDef{
			models.UpdateReleases,
		})), router.downloadHandler.DownloadAnyRelease)

		supervision.POST("/reject/release/:releaseId", router.middleWare.AuthMiddleware(be.GetJWTSecret()), router.middleWare.PreAuthorize(router.middleWare.HasPermissions([]models.UserPermissionTypeDef{
			models.RejectReleases,
		})), router.releaseHandler.RejectRelease)

		supervision.POST("/accept/release/:releaseId", router.middleWare.AuthMiddleware(be.GetJWTSecret()), router.middleWare.PreAuthorize(router.middleWare.HasPermissions([]models.UserPermissionTypeDef{
			models.AcceptReleases,
		})), router.releaseHandler.AcceptRelease)

		supervision.POST("/cancel/release/:releaseId", router.middleWare.AuthMiddleware(be.GetJWTSecret()), router.middleWare.PreAuthorize(router.middleWare.HasPermissions([]models.UserPermissionTypeDef{
			models.CancelReleases,
		})), router.releaseHandler.CancelRelease)

		supervision.POST("/pending/release/:releaseId", router.middleWare.AuthMiddleware(be.GetJWTSecret()), router.middleWare.PreAuthorize(router.middleWare.HasPermissions([]models.UserPermissionTypeDef{
			models.CancelReleases,
		})), router.releaseHandler.ChangeToPendingRelease)

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
