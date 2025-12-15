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

func (h *ReleaseHandler) RejectRelease(c *gin.Context) {
	userID := c.GetUint("userId")
	releaseID := c.Param("releaseId")
	releaseIDUint, err := utils.StringToUint(releaseID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid Release ID format",
			"result":  false,
			"details": err.Error(),
		})
		return
	}

	result, err := h.service.RejectModuleRelease(userID, releaseIDUint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"details": err.Error(),
			"result":  false,
			"error":   "Coulnd't reject module release",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":   "Release was rejected successfully",
		"releaseId": result,
		"result":    true,
	})
}

func (h *ReleaseHandler) AcceptRelease(c *gin.Context) {
	userID := c.GetUint("userId")
	releaseID := c.Param("releaseId")
	releaseIDUint, err := utils.StringToUint(releaseID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid Release ID format",
			"result":  false,
			"details": err.Error(),
		})
		return
	}

	result, err := h.service.AcceptModuleRelease(userID, releaseIDUint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":   "Release was accepted successfully",
		"releaseId": result,
	})
}

func (h *ReleaseHandler) CancelRelease(c *gin.Context) {
	userID := c.GetUint("userId")
	releaseID := c.Param("releaseId")
	releaseIDUint, err := utils.StringToUint(releaseID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid Release ID format",
			"result":  false,
			"details": err.Error(),
		})
		return
	}

	result, err := h.service.CancelModuleRelease(userID, releaseIDUint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":   "Release was cencelled successfully",
		"releaseId": result,
	})
}

func (h *ReleaseHandler) ChangeToPendingRelease(c *gin.Context) {
	userID := c.GetUint("userId")
	releaseID := c.Param("releaseId")
	releaseIDUint, err := utils.StringToUint(releaseID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid Release ID format",
			"result":  false,
			"details": err.Error(),
		})
		return
	}

	result, err := h.service.ChangeToPendingModuleRelease(userID, releaseIDUint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":   "Release Status changed to Pending successfully",
		"releaseId": result,
	})
}

func (h *ReleaseHandler) DeleteModuleRelease(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetUint("userId")

	idUint, err := utils.StringToUint(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid Release ID format",
			"result":  false,
			"details": err.Error(),
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

	result, err := h.service.DeleteModuleRelease(userID, idUint, releaseIDUint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Release was deleted successfully",
		"result":  result,
	})
}

func (h *ReleaseHandler) CancelSuggestedRelease(c *gin.Context) {
	modID := c.Param("id")
	userID := c.GetUint("userId")

	modIDUint, err := utils.StringToUint(modID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid Module ID format",
			"result":  false,
			"details": err.Error(),
		})
		return
	}

	releaseID := c.Param("releaseId")
	releaseIDUint, err := utils.StringToUint(releaseID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid Release ID format",
			"result":  false,
			"details": err.Error(),
		})
		return
	}

	result, err := h.service.CancelSuggestedModuleRelease(userID, modIDUint, releaseIDUint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Release was cancelled successfully",
		"result":  result,
	})
}

func (h *ReleaseHandler) GetModuleRelease(c *gin.Context) {
	id := c.Param("id")

	idUint, err := utils.StringToUint(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid Module Release ID format",
			"result":  false,
			"details": err.Error(),
		})
		return
	}

	releaseID := c.Param("releaseId")
	releaseIDUint, err := utils.StringToUint(releaseID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid Release ID format",
			"result":  false,
			"details": err.Error(),
		})
		return
	}

	module, err := h.service.GetModuleRelease(idUint, releaseIDUint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Couldn't get Module Release",
			"result":  false,
			"details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Release information retrieved successfully",
		"release": module,
		"result":  true,
	})
}

func (h *ReleaseHandler) GetModuleReleases(c *gin.Context) {
	id := c.Param("id")

	idUint, err := utils.StringToUint(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid Module ID format",
			"result":  false,
			"details": err.Error(),
		})
		return
	}

	module, err := h.service.GetModuleReleases(idUint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Couldn't get Module Releases",
			"result":  false,
			"details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Module Releases information retrieved successfully",
		"releases": module,
		"result":   true,
	})
}

func (h *ReleaseHandler) SearchByKeywords(c *gin.Context) {
	tagsQuery := c.Query("tags")
	if tagsQuery == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "tags query parameter is required",
			"result":  false,
			"details": "tags query parameter is required",
		})
		return
	}

	tags := strings.Split(tagsQuery, ",")
	id := c.Param("id")

	idUint, err := utils.StringToUint(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Invalid Module Release ID format",
			"result":  false,
			"details": err.Error()})
		return
	}

	releases, err := h.service.SearchByKeywords(idUint, tags)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Couldn't search Module Releases by tags",
			"result":  false,
			"details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Release information retrieved successfully",
		"releases": releases,
		"result":   true,
	})
}
