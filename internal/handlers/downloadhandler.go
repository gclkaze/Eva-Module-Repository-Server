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

}

func (h *DownloadHandler) DownloadAnyRelease(c *gin.Context) {

}
