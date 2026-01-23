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

// Package handlers contains the controllers of the Repository Server
package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gclkaze/evamodulerepositoryserver/internal/dto"
	"github.com/gclkaze/evamodulerepositoryserver/internal/models"
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

func (h *ReleaseHandler) FindRelease(c *gin.Context) {
	//userID := c.GetUint("userId")
	moduleName := c.Param("module")

	version := c.Param("version")
	if !utils.IsValidVersion(version) {
		c.JSON(http.StatusBadRequest, utils.ErrWithSimpleMessage("Invalid Version ID format"))
		return
	}

	result, err := h.service.FindModuleRelease(moduleName, version)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Err(err, "Couldn't find module release"))
		return
	}
	c.JSON(http.StatusOK, utils.OkWithMessage(result, "Release was found successfully"))
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

func (h *ReleaseHandler) GetModuleReleasesByFilter(c *gin.Context) {
	p := models.NewReleaseFilterParams()

	// string slice filters (GET query params)
	p.Status = c.QueryArray("status")
	p.Versions = c.QueryArray("versions")
	p.Tags = c.QueryArray("tags")
	p.ModuleName = c.QueryArray("module")
	p.RepoName = c.QueryArray("repo")
	p.Description = c.QueryArray("description")
	p.Creator = c.QueryArray("creator")
	p.CreatorEmail = c.QueryArray("creator-email")

	// created-after (RFC3339)
	if v := c.Query("created-after"); v != "" {
		t, err := time.Parse(time.RFC3339, v)
		if err != nil {
			c.JSON(
				http.StatusBadRequest,
				utils.Err(err, "Invalid created-after timestamp format"),
			)
			return
		}
		p.CreatedAfter = &t
	}

	// released-after (RFC3339)
	if v := c.Query("released-after"); v != "" {
		t, err := time.Parse(time.RFC3339, v)
		if err != nil {
			c.JSON(
				http.StatusBadRequest,
				utils.Err(err, "Invalid released-after timestamp format"),
			)
			return
		}
		p.ReleasedAfter = t
	}

	// normalize slices (trim, dedupe, drop empty)
	p.Status = normalizeStringSlice(p.Status)
	p.Versions = normalizeStringSlice(p.Versions)
	p.Tags = normalizeStringSlice(p.Tags)
	p.ModuleName = normalizeStringSlice(p.ModuleName)
	p.RepoName = normalizeStringSlice(p.RepoName)
	p.Description = normalizeStringSlice(p.Description)
	p.Creator = normalizeStringSlice(p.Creator)
	p.CreatorEmail = normalizeStringSlice(p.CreatorEmail)

	// service call
	result, err := h.service.GetModuleReleasesClustered(p)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			utils.Err(err, "Couldn't retrieve module releases"),
		)
		return
	}

	c.JSON(
		http.StatusOK,
		utils.OkWithMessage(result, "Module releases retrieved successfully"),
	)
}

func normalizeStringSlice(in []string) []string {
	out := make([]string, 0, len(in))
	seen := map[string]struct{}{}

	for _, v := range in {
		v = strings.TrimSpace(v)
		if v == "" {
			continue
		}
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		out = append(out, v)
	}
	return out
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

func (h *ReleaseHandler) GetModuleLatestReleaseOnAuth(c *gin.Context) {
	name := c.Param("module")

	module, err := h.service.FindModule(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Err(err, fmt.Sprintf("Couldn't find Module '%s'", name)))
		return
	}

	if module == nil {
		c.JSON(http.StatusInternalServerError, utils.ErrWithSimpleMessage(fmt.Sprintf("Couldn't find Module '%s'", name)))
		return
	}

	theRelease, err := h.service.GetLastModuleStatusIndependentRelease(module.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Err(err, fmt.Sprintf("Couldn't find Latest Release of Module: '%s'", name)))
		return
	}

	if theRelease == nil {
		c.JSON(http.StatusInternalServerError, utils.ErrWithSimpleMessage(fmt.Sprintf("Couldn't find Latest Module Release of Module:  '%s'", name)))
		return
	}

	c.JSON(http.StatusOK, utils.OkWithMessage(dto.NewReleaseDTO(*theRelease), "Latest Module Releases information retrieved successfully"))
}

func (h *ReleaseHandler) GetModuleLatestRelease(c *gin.Context) {
	name := c.Param("module")

	module, err := h.service.FindModule(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Err(err, fmt.Sprintf("Couldn't find Module '%s'", name)))
		return
	}

	if module == nil {
		c.JSON(http.StatusInternalServerError, utils.ErrWithSimpleMessage(fmt.Sprintf("Couldn't find Module '%s'", name)))
		return
	}

	theRelease, err := h.service.GetLastModuleRelease(module.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Err(err, fmt.Sprintf("Couldn't find Latest Release of Module: '%s'", name)))
		return
	}

	if theRelease == nil {
		c.JSON(http.StatusInternalServerError, utils.ErrWithSimpleMessage(fmt.Sprintf("Couldn't find Latest Module Release of Module:  '%s'", name)))
		return
	}

	c.JSON(http.StatusOK, utils.OkWithMessage(dto.NewReleaseDTO(*theRelease), "Latest Module Releases information retrieved successfully"))
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
