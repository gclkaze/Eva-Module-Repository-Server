package handlers

import (
	"net/http"

	"github.com/gclkaze/evamodulerepositoryserver/internal/services"
	"github.com/gclkaze/evamodulerepositoryserver/pkg/utils"
	"github.com/gin-gonic/gin"
)

type DownloadHandler struct {
	service *services.DownloadService
}

func NewDownloadHandler(service *services.DownloadService) *DownloadHandler {
	return &DownloadHandler{service: service}
}

func (h *DownloadHandler) DownloadRelease(c *gin.Context) {
	//the release needs to be ACCEPTED
	releaseID := c.Param("releaseId")
	releaseIDUint, err := utils.StringToUint(releaseID)

	if err != nil {
		c.JSON(http.StatusBadRequest, utils.Err(err, "Invalid Release ID format"))
		return
	}

	dest, filename, err := h.service.DownloadRelease(releaseIDUint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Err(err, err.Error()))
		return
	}

	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "application/octet-stream")
	c.File(dest)
}

func (h DownloadHandler) DownloadAnyRelease(c *gin.Context) {
	releaseID := c.Param("releaseId")
	releaseIDUint, err := utils.StringToUint(releaseID)

	if err != nil {
		c.JSON(http.StatusBadRequest, utils.Err(err, "Invalid Release ID format"))
		return
	}

	dest, filename, err := h.service.DownloadAnyRelease(releaseIDUint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Err(err, err.Error()))
		return
	}

	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "application/octet-stream")
	c.File(dest)
}
