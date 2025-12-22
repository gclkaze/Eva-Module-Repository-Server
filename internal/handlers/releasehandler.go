// Package handlers contains the controllers of the Repository Server
package handlers

import (
	"fmt"
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
		c.JSON(http.StatusBadRequest, utils.Err(err, "Invalid Release ID format"))
		return
	}

	result, err := h.service.RejectModuleRelease(userID, releaseIDUint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Err(err, "Couldn't reject module release"))
		return
	}
	c.JSON(http.StatusOK, utils.OkWithMessage(result, "Release was rejected successfully"))
}

func (h *ReleaseHandler) AcceptRelease(c *gin.Context) {
	userID := c.GetUint("userId")
	releaseID := c.Param("releaseId")
	releaseIDUint, err := utils.StringToUint(releaseID)

	if err != nil {
		c.JSON(http.StatusBadRequest, utils.Err(err, "Invalid Release ID format"))
		return
	}

	result, err := h.service.AcceptModuleRelease(userID, releaseIDUint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Err(err, "Couldn't accept module release"))
		return
	}
	c.JSON(http.StatusOK, utils.OkWithMessage(result, "Release was accepted successfully"))
}

func (h *ReleaseHandler) CancelRelease(c *gin.Context) {
	userID := c.GetUint("userId")
	releaseID := c.Param("releaseId")
	releaseIDUint, err := utils.StringToUint(releaseID)

	if err != nil {
		c.JSON(http.StatusBadRequest, utils.Err(err, "Invalid Release ID format"))
		return
	}

	result, err := h.service.CancelModuleRelease(userID, releaseIDUint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Err(err, "Couldn't cancel module release"))
		return
	}
	c.JSON(http.StatusOK, utils.OkWithMessage(result, "Release was cancelled successfully"))
}

func (h *ReleaseHandler) ChangeToPendingRelease(c *gin.Context) {
	userID := c.GetUint("userId")
	releaseID := c.Param("releaseId")
	releaseIDUint, err := utils.StringToUint(releaseID)

	if err != nil {
		c.JSON(http.StatusBadRequest, utils.Err(err, "Invalid Release ID format"))
		return
	}

	result, err := h.service.ChangeToPendingModuleRelease(userID, releaseIDUint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Err(err, "Couldn't change to pending module release"))
		return
	}
	c.JSON(http.StatusOK, utils.OkWithMessage(result, "Release Status changed to Pending successfully"))
}

func (h *ReleaseHandler) DeleteModuleRelease(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetUint("userId")

	idUint, err := utils.StringToUint(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.Err(err, "Invalid Module ID format"))
		return
	}

	releaseID := c.Param("releaseId")
	releaseIDUint, err := utils.StringToUint(releaseID)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.Err(err, "Invalid Release ID format"))
		return
	}

	result, err := h.service.DeleteModuleRelease(userID, idUint, releaseIDUint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Err(err, "Couldn't delete module release"))
		return
	}

	c.JSON(http.StatusOK, utils.OkWithMessage(result, "Release was deleted successfully"))
}

func (h *ReleaseHandler) CancelSuggestedRelease(c *gin.Context) {
	modID := c.Param("id")
	userID := c.GetUint("userId")

	modIDUint, err := utils.StringToUint(modID)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.Err(err, "Invalid Module ID format"))
		return
	}

	releaseID := c.Param("releaseId")
	releaseIDUint, err := utils.StringToUint(releaseID)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.Err(err, "Invalid Release ID format"))
		return
	}

	result, err := h.service.CancelSuggestedModuleRelease(userID, modIDUint, releaseIDUint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Err(err, "Couldn't cancel suggested release"))
		return
	}

	c.JSON(http.StatusOK, utils.OkWithMessage(result, fmt.Sprintf("Release %d was cancelled successfully", releaseIDUint)))
}

func (h *ReleaseHandler) GetModuleRelease(c *gin.Context) {
	id := c.Param("id")

	idUint, err := utils.StringToUint(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.Err(err, "Invalid Module ID format"))
		return
	}

	releaseID := c.Param("releaseId")
	releaseIDUint, err := utils.StringToUint(releaseID)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.Err(err, "Invalid Release ID format"))
		return
	}

	module, err := h.service.GetModuleRelease(idUint, releaseIDUint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Err(err, "Couldn't get Module Release"))
		return
	}

	c.JSON(http.StatusOK, utils.OkWithMessage(module, "Release information retrieved successfully"))
}

func (h *ReleaseHandler) GetModuleReleases(c *gin.Context) {
	id := c.Param("id")

	idUint, err := utils.StringToUint(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.Err(err, "Invalid Module ID format"))
		return
	}

	module, err := h.service.GetModuleReleases(idUint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Err(err, "Couldn't get Module Releases"))
		return
	}

	c.JSON(http.StatusOK, utils.OkWithMessage(module, "Module Releases information retrieved successfully"))
}

func (h *ReleaseHandler) SearchByKeywords(c *gin.Context) {
	tagsQuery := c.Query("tags")
	if tagsQuery == "" {
		c.JSON(http.StatusBadRequest, utils.Err(nil, "tags query parameter is required"))
		return
	}

	tags := strings.Split(tagsQuery, ",")
	id := c.Param("id")

	idUint, err := utils.StringToUint(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.Err(err, "Invalid Release ID format"))
		return
	}

	releases, err := h.service.SearchByKeywords(idUint, tags)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Err(err, "Couldn't search Module Releases by tags"))
		return
	}

	c.JSON(http.StatusOK, utils.OkWithMessage(releases, "Release information retrieved successfully"))
}
