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

package routes

import (
	"fmt"
	"net/http"

	"github.com/gclkaze/evamodulerepositoryserver/internal/backend"
	"github.com/gclkaze/evamodulerepositoryserver/internal/config"
	"github.com/gclkaze/evamodulerepositoryserver/internal/handlers"
	"github.com/gclkaze/evamodulerepositoryserver/internal/middleware"
	"github.com/gclkaze/evamodulerepositoryserver/internal/models"
	"github.com/gclkaze/evamodulerepositoryserver/pkg/utils"

	"github.com/gin-gonic/gin"
)

type EvaModuleRepositoryRouter struct {
	api  *gin.RouterGroup
	r    *gin.Engine
	port string

	moduleHandler    *handlers.ModuleHandler
	downloadHandler  *handlers.DownloadHandler
	releaseHandler   *handlers.ReleaseHandler
	authHandler      *handlers.AuthHandler
	superviseHandler *handlers.SuperviseHandler
	middleWare       *middleware.AuthMiddleWare
	uploadLimit      int64
}

func NewEvaModuleRepositoryRouter() *EvaModuleRepositoryRouter {
	inst := &EvaModuleRepositoryRouter{}
	inst.uploadLimit = 8 << 20 // 8 MB
	return inst
}

func (router *EvaModuleRepositoryRouter) SetUploadFileLimit(limit int64) int64 {
	old := router.uploadLimit
	router.uploadLimit = limit
	return old
}

func (router *EvaModuleRepositoryRouter) Initialize(r *gin.Engine, be *backend.EvaModuleRepositoryBackend) error {
	router.api = r.Group(APIGroup)
	router.moduleHandler = handlers.NewModuleHandler(be.GetModuleService())
	router.releaseHandler = handlers.NewReleaseHandler(be.GetReleaseService())
	router.authHandler = handlers.NewAuthHandler(be.GetAuthService(), be.GetUserService())
	router.downloadHandler = handlers.NewDownloadHandler(be.GetDownloadService())
	router.middleWare = middleware.NewAuthMiddleware(be.GetUserService())
	router.superviseHandler = handlers.NewSuperviseHandler(be.GetUserService())
	r.MaxMultipartMemory = router.uploadLimit

	authGroup := router.api.Group(AuthGroup)
	{
		authGroup.POST(LoginEndpoint, router.authHandler.Login)
		authGroup.POST(RefreshEndpoint, router.authHandler.Refresh)
		authGroup.POST(RegisterEndpoint, router.authHandler.Register)
	}

	modules := router.api.Group(ModulesGroup)
	{
		//find a module based on id
		modules.GET("/:id", router.moduleHandler.FindByID)                     // GET /api/modules/:id
		modules.GET(ModuleGetInfoEndpoint, router.moduleHandler.GetModuleInfo) // GET /api/modules/:id

		//modules.GET(ModuleSearchEndpoint, router.moduleHandler.SearchModulesByTags)
		modules.GET(ModuleSearchEndpoint, router.moduleHandler.SearchModulesByComponents)

		modules.GET(GetUserModulesEndpoint, router.middleWare.AuthMiddleware(be.GetJWTSecret()),
			router.middleWare.PreAuthorize(router.middleWare.HasPermissions([]models.UserPermissionTypeDef{
				models.DeleteMyModule,
				models.CreateMyModule,
			})), router.moduleHandler.GetUserModules)

		//delete a module!
		modules.POST(ModuleDeleteEndpoint,
			router.middleWare.AuthMiddleware(be.GetJWTSecret()),
			router.middleWare.PreAuthorize(router.middleWare.HasPermissions([]models.UserPermissionTypeDef{
				models.DeleteMyModule,
				models.CreateMyModule,
				//	models.DeleteModules,
			})),
			router.moduleHandler.Delete) //needs userID to be passed

		modules.POST(ModuleUploadEndpoint, router.middleWare.AuthMiddleware(be.GetJWTSecret()),

			router.middleWare.MaxBodySize(r.MaxMultipartMemory),
			func(c *gin.Context) {
				_, err := c.MultipartForm()
				if err != nil {
					c.JSON(http.StatusRequestEntityTooLarge, utils.Err(err, "file too large"))
					return
				}
			},

			router.middleWare.PreAuthorize(router.middleWare.HasPermissions([]models.UserPermissionTypeDef{
				models.DeleteMyModule,
				models.CreateMyModule,
			})), router.moduleHandler.Upload)

		modules.POST(ModuleUpdateEndpoint, router.middleWare.AuthMiddleware(be.GetJWTSecret()),

			router.middleWare.MaxBodySize(r.MaxMultipartMemory),
			func(c *gin.Context) {
				_, err := c.MultipartForm()
				if err != nil {
					c.JSON(http.StatusRequestEntityTooLarge, utils.Err(err, "file too large"))
					return
				}
			},

			router.middleWare.PreAuthorize(router.middleWare.HasPermissions([]models.UserPermissionTypeDef{
				models.DeleteMyModule,
				models.CreateMyModule,
				/*				models.UpdateModules,*/
			})), router.moduleHandler.Update)

		modules.POST(ModuleSuggestEndpoint, router.middleWare.AuthMiddleware(be.GetJWTSecret()),
			router.middleWare.PreAuthorize(router.middleWare.HasPermissions([]models.UserPermissionTypeDef{
				models.DeleteMyModule,
				models.CreateMyModule,
				models.SuggestMyModule,
			})), router.moduleHandler.SuggestRelease)
	}

	releases := router.api.Group(ReleasesGroup)
	{
		releases.GET("/:id", router.releaseHandler.GetModuleReleases)                               // GET /api/releases/:id
		releases.GET("/:id/"+ReleaseEndpoint+"/:releaseId", router.releaseHandler.GetModuleRelease) // GET /api/releases/:id/release/:releaseId
		releases.GET("/:id"+ReleaseSearchEndpoint, router.releaseHandler.SearchByKeywords)

		releases.POST("/:id"+ReleaseDeleteEndpoint+"/:releaseId", router.middleWare.AuthMiddleware(be.GetJWTSecret()),
			router.middleWare.PreAuthorize(router.middleWare.HasPermissions([]models.UserPermissionTypeDef{
				models.DeleteMyModule,
				models.CreateMyModule,
				models.SuggestMyModule,
				models.DeleteMyRelease,
			})), router.releaseHandler.DeleteModuleRelease) // GET /api/releases/:id/delete/:releaseId  ->need to check the userID is the one that initiated it

		releases.POST("/:id"+ReleaseCancelEndpoint+"/:releaseId", router.middleWare.AuthMiddleware(be.GetJWTSecret()),
			router.middleWare.PreAuthorize(router.middleWare.HasPermissions([]models.UserPermissionTypeDef{
				models.DeleteMyModule,
				models.CreateMyModule,
				models.SuggestMyModule,
			})), router.releaseHandler.CancelSuggestedRelease) // ->need to check the userID is the one that initiated it
	}

	download := router.api.Group(DownloadGroup)
	{
		download.GET("/release/:releaseId", router.downloadHandler.DownloadRelease) //the release needs to be accepted
	}

	supervision := router.api.Group(SuperviseGroup)
	{
		supervision.GET(SuperviseDownloadAnyReleaseEndpoint+"/:releaseId", router.middleWare.AuthMiddleware(be.GetJWTSecret()), router.middleWare.PreAuthorize(router.middleWare.HasPermissions([]models.UserPermissionTypeDef{
			models.UpdateReleases,
		})), router.downloadHandler.DownloadAnyRelease)

		supervision.GET(SuperviseFindReleaseEndpoint+"/:module/:version", router.middleWare.AuthMiddleware(be.GetJWTSecret()), router.middleWare.PreAuthorize(router.middleWare.HasPermissions([]models.UserPermissionTypeDef{
			models.UpdateReleases,
		})), router.releaseHandler.FindRelease)

		supervision.GET(SuperviseGetFilterReleaseEndpoint, router.middleWare.AuthMiddleware(be.GetJWTSecret()), router.middleWare.PreAuthorize(router.middleWare.HasPermissions([]models.UserPermissionTypeDef{
			models.UpdateReleases,
		})), router.releaseHandler.GetModuleReleasesByFilter)

		supervision.POST(SuperviseRejectReleaseEndpoint+"/:releaseId", router.middleWare.AuthMiddleware(be.GetJWTSecret()), router.middleWare.PreAuthorize(router.middleWare.HasPermissions([]models.UserPermissionTypeDef{
			models.RejectReleases,
		})), router.releaseHandler.RejectRelease)

		supervision.POST(SuperviseAcceptReleaseEndpoint+"/:releaseId", router.middleWare.AuthMiddleware(be.GetJWTSecret()), router.middleWare.PreAuthorize(router.middleWare.HasPermissions([]models.UserPermissionTypeDef{
			models.AcceptReleases,
		})), router.releaseHandler.AcceptRelease)

		supervision.POST(SuperviseCancelReleaseEndpoint+"/:releaseId", router.middleWare.AuthMiddleware(be.GetJWTSecret()), router.middleWare.PreAuthorize(router.middleWare.HasPermissions([]models.UserPermissionTypeDef{
			models.CancelReleases,
		})), router.releaseHandler.CancelRelease)

		supervision.POST(SupervisePendingReleaseEndpoint+"/:releaseId", router.middleWare.AuthMiddleware(be.GetJWTSecret()), router.middleWare.PreAuthorize(router.middleWare.HasPermissions([]models.UserPermissionTypeDef{
			models.CancelReleases,
		})), router.releaseHandler.ChangeToPendingRelease)

		supervision.POST(SuperviseBanUserEndpoint+"/:userId", router.middleWare.AuthMiddleware(be.GetJWTSecret()), router.middleWare.PreAuthorize(router.middleWare.HasPermissions([]models.UserPermissionTypeDef{
			models.BanUsers,
			models.UnbanUsers,
		})), router.superviseHandler.BanUser)

		supervision.POST(SuperviseUnbanUserEndpoint+"/:userId", router.middleWare.AuthMiddleware(be.GetJWTSecret()), router.middleWare.PreAuthorize(router.middleWare.HasPermissions([]models.UserPermissionTypeDef{
			models.BanUsers,
			models.UnbanUsers,
		})), router.superviseHandler.UnbanUser)

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

func (router *EvaModuleRepositoryRouter) GetRouter() *gin.Engine {
	return router.r
}

func (router *EvaModuleRepositoryRouter) GetMiddleware() *middleware.AuthMiddleWare {
	return router.middleWare
}
