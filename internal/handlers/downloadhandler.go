package handlers

import (
	"github.com/gclkaze/evamodulerepositoryserver/internal/services"
	"github.com/gin-gonic/gin"
)

type DownloadHandler struct {
	service *services.ModuleService
}

func NewDownloadHandler(service *services.ModuleService) *DownloadHandler {
	return &DownloadHandler{service: service}
}

func (h *DownloadHandler) DownloadRelease(c *gin.Context) {
	/*	userID := c.GetUint("userId")
		releaseID := c.Param("releaseId")
		releaseIDUint, err := utils.StringToUint(releaseID)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid Release ID format",
			})
			return
		}

		result, err := h.service.AcceptModuleRelease(userID, releaseIDUint)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, result)*/
}

func (h *DownloadHandler) DownloadAnyRelease(c *gin.Context) {

}
