package handlers

import (
	"net/http"
	"path"

	"github.com/gclkaze/evamodulerepositoryserver/internal/repositories"
	"github.com/gclkaze/evamodulerepositoryserver/internal/services"
	"github.com/gclkaze/evamodulerepositoryserver/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/magiconair/properties"
)

type DownloadHandler struct {
	service             *services.ModuleService
	defaultDistFilename string
}

func NewDownloadHandler(service *services.ModuleService, p *properties.Properties) *DownloadHandler {
	s := "dist.tar.gz"
	if p != nil {
		s = p.GetString("dist_name", s)
	}
	return &DownloadHandler{service: service, defaultDistFilename: s}
}

func (h *DownloadHandler) DownloadRelease(c *gin.Context) {
	//the release needs to be ACCEPTED
	releaseID := c.Param("releaseId")
	releaseIDUint, err := utils.StringToUint(releaseID)
	filename := h.defaultDistFilename

	if err != nil {
		c.JSON(http.StatusBadRequest, utils.Err(err, "Invalid Release ID format"))
		return
	}

	rel, err := h.service.GetReleaseService().GetRelease(releaseIDUint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Err(err, "Couldn't find Release"))
		return
	}

	if rel.Status.Label != repositories.Accepted.String() {
		c.JSON(http.StatusInternalServerError, utils.Err(nil, "no release found"))
		return
	}

	mod, err := h.service.GetModule(rel.ModuleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Err(err, "Couldn't find Release Module"))
		return
	}

	dest := h.service.GetModuleReleasePath(mod, rel)
	if !utils.FolderExists(dest) {
		c.JSON(http.StatusInternalServerError, utils.Err(err, "Couldn't find Release artifact."))
		return
	}

	dest = path.Join(dest, filename)
	if !utils.FileExists(dest) {
		c.JSON(http.StatusInternalServerError, utils.Err(err, "Couldn't find Release artifact."))
		return
	}
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "application/octet-stream")
	c.File(dest)
}

func (h DownloadHandler) DownloadAnyRelease(c *gin.Context) {
	releaseID := c.Param("releaseId")
	releaseIDUint, err := utils.StringToUint(releaseID)
	filename := h.defaultDistFilename

	if err != nil {
		c.JSON(http.StatusBadRequest, utils.Err(err, "Invalid Release ID format"))
		return
	}

	rel, err := h.service.GetReleaseService().GetRelease(releaseIDUint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Err(err, "Couldn't find Release"))
		return
	}

	mod, err := h.service.GetModule(rel.ModuleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Err(err, "Couldn't find Release Module"))
		return
	}

	dest := h.service.GetModuleReleasePath(mod, rel)
	if !utils.FolderExists(dest) {
		c.JSON(http.StatusInternalServerError, utils.Err(err, "Couldn't find Release artifact."))
		return
	}

	dest = path.Join(dest, filename)
	if !utils.FileExists(dest) {
		c.JSON(http.StatusInternalServerError, utils.Err(err, "Couldn't find Release artifact."))
		return
	}

	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "application/octet-stream")
	c.File(dest)
}
