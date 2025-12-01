// Package handlers contains the controllers of the Repository Server
package handlers

import (
	"net/http"
	"strings"

	"github.com/gclkaze/evamodulerepositoryserver/internal/services"
	"github.com/gclkaze/evamodulerepositoryserver/pkg/utils"
	"github.com/gin-gonic/gin"
)

type ReleaseHandler struct {
	service *services.ReleaseService
}

func NewReleaseHandler(service *services.ReleaseService) *ReleaseHandler {
	return &ReleaseHandler{service: service}
}

func (h *ReleaseHandler) DeleteModuleRelease(c *gin.Context) {
	id := c.Param("id")

	idUint, err := utils.StringToUint(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid ID format",
		})
		return
	}

	releaseID := c.Param("releaseId")
	releaseIDUint, err := utils.StringToUint(releaseID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Release ID format",
		})
		return
	}

	result, err := h.service.DeleteModuleRelease(idUint, releaseIDUint)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, result)
}

func (h *ReleaseHandler) GetModuleRelease(c *gin.Context) {
	id := c.Param("id")

	idUint, err := utils.StringToUint(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid ID format",
		})
		return
	}

	releaseID := c.Param("releaseId")
	releaseIDUint, err := utils.StringToUint(releaseID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Release ID format",
		})
		return
	}

	module, err := h.service.GetModuleRelease(idUint, releaseIDUint)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, module)
}

func (h *ReleaseHandler) GetModuleReleases(c *gin.Context) {
	id := c.Param("id")

	idUint, err := utils.StringToUint(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid ID format",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid ID format",
		})
		return
	}

	module, err := h.service.GetModuleReleases(idUint)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, module)
}

func (h *ReleaseHandler) SearchByKeywords(c *gin.Context) {
	tagsQuery := c.Query("tags")
	if tagsQuery == "" {
		c.JSON(400, gin.H{"error": "tags query parameter is required"})
		return
	}

	tags := strings.Split(tagsQuery, ",")

	id := c.Param("id")

	idUint, err := utils.StringToUint(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid ID format",
		})
		return
	}

	releases, err := h.service.SearchByKeywords(idUint, tags)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, releases)
}
