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
	"net/http"
	"strings"

	"github.com/gclkaze/evamodulerepositoryserver/internal/dto"
	"github.com/gclkaze/evamodulerepositoryserver/internal/services"
	"github.com/gclkaze/evamodulerepositoryserver/pkg/utils"
	"github.com/gin-gonic/gin"
)

type ModuleHandler struct {
	service *services.ModuleService
}

func NewModuleHandler(service *services.ModuleService) *ModuleHandler {
	return &ModuleHandler{service: service}
}

func (h *ModuleHandler) FindByID(c *gin.Context) {
	id := c.Param("id")
	idUint, err := utils.StringToUint(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.Err(err, "Invalid ID format"))
		return
	}

	module, err := h.service.FindByID(idUint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Err(err, "Couldn't find module"))
		return
	}
	c.JSON(http.StatusOK, utils.OkWithMessage(module, "Module information retrieved successfully"))
}

func (h *ModuleHandler) GetModuleInfo(c *gin.Context) {
	name := c.Query("moduleName")

	module, release, err := utils.ParseModuleReleaseVersion(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Err(err, "Couldn't identify the release"))
		return
	}

	if module == "" {
		c.JSON(http.StatusInternalServerError, utils.ErrWithSimpleMessage("Couldn't identify the user's provided typed module name "+name))
		return
	}

	if release == "" {
		var module *dto.ModuleEnrichedDTO
		module, err = h.service.GetModuleInfo(name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, utils.Err(err, "Couldn't find module"))
			return
		}
		c.JSON(http.StatusOK, utils.OkWithMessage(module, "Module information retrieved successfully"))
		return
	}

	releaseDTO, err := h.service.GetModuleReleaseInfo(module, release)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Err(err, "Couldn't find module and release "+name))
		return
	}
	c.JSON(http.StatusOK, utils.OkWithMessage(releaseDTO, "Module release information retrieved successfully"))

}

func (h *ModuleHandler) Delete(c *gin.Context) {
	id := c.PostForm("id")
	userID := c.GetUint("userId")
	idUint, err := utils.StringToUint(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.Err(err, "Invalid ID format"))
		return
	}
	res, err := h.service.Delete(userID, idUint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Err(err, "Couldn't delete module"))
		return
	}

	c.JSON(http.StatusOK, utils.OkWithMessage(res, "Module deleted successfully"))
}

func (h *ModuleHandler) Upload(c *gin.Context) {
	userID := c.GetUint("userId")
	title := c.PostForm("title")
	repr := c.PostForm("repr")
	tags := c.PostForm("tags")
	description := c.PostForm("description")

	// Validate required fields
	if title == "" {
		c.JSON(http.StatusBadRequest, utils.Err(nil, "Title is required"))
		return
	}

	// Handle file upload
	_, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.Err(err, "Uploaded file is empty"))
		return
	}

	form, err := c.MultipartForm()
	if err != nil && strings.Contains(err.Error(), "http: request body too large") {
		c.JSON(http.StatusRequestEntityTooLarge, utils.Err(err, "file too large"))
		return
	}
	files, ok := form.File["file"]
	if !ok {
		c.JSON(http.StatusBadRequest, utils.Err(err, "Uploaded file is empty"))
		return
	}

	// Call service to create module
	id, err := h.service.CreateModuleTx(userID, title, description, repr, files, tags, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Err(err, "Couldn't upload module"))
		return
	}

	c.JSON(http.StatusOK, utils.OkWithMessage(id, "Module code was uploaded successfully"))
}

func (h *ModuleHandler) Update(c *gin.Context) {
	userID := c.GetUint("userId")
	modID := c.PostForm("modId")
	title := c.PostForm("title")
	repr := c.PostForm("repr")
	tags := c.PostForm("tags")
	description := c.PostForm("description")

	modIDUint, err := utils.StringToUint(modID)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.Err(err, "Invalid Module ID format"))
		return
	}

	if title == "" {
		c.JSON(http.StatusBadRequest, utils.Err(nil, "Title is required"))
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.Err(err, "File upload failed"))
		return
	}

	id, err := h.service.UpdateUserModule(userID, modIDUint, title, description, repr, file, tags, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Err(err, "Couldn't update module"))
		return
	}

	c.JSON(http.StatusOK, utils.OkWithMessage(id, "Module was updated successfully"))
}

func (h *ModuleHandler) SuggestRelease(c *gin.Context) {
	userID := c.GetUint("userId")
	modID := c.PostForm("modId")
	version := c.PostForm("version")

	modIDUint, err := utils.StringToUint(modID)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.Err(err, "Invalid Module ID format"))
		return
	}

	id, err := h.service.SuggestUserModuleRelease(userID, modIDUint, version)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Err(err, "Couldn't suggest module release"))
		return
	}

	c.JSON(http.StatusOK, utils.OkWithMessage(id, "Module suggested for release successfully"))
}

func (h *ModuleHandler) SearchModulesByTags(c *gin.Context) {
	tagsQuery := c.Query("tags")
	if tagsQuery == "" {
		c.JSON(http.StatusBadRequest, utils.Err(nil, "tags query parameter is required"))
		return
	}

	labels := strings.Split(tagsQuery, ",")
	modules, err := h.service.SearchByKeywords(labels)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Err(err, "Couldn't find the modules"))
		return
	}

	c.JSON(http.StatusOK, utils.OkWithMessage(modules, "Module search was complete"))
}

func (h *ModuleHandler) SearchModulesByComponents(c *gin.Context) {
	tagTokens := c.QueryArray("tags")

	//tagTokens := strings.Split(tagsQuery, ",")

	nameTokens := c.QueryArray("name")
	descrTokens := c.QueryArray("description")

	modules, err := h.service.SearchByComponents(nameTokens, descrTokens, tagTokens)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Err(err, "Couldn't find the modules"))
		return
	}

	c.JSON(http.StatusOK, utils.OkWithMessage(modules, "Module search was complete"))
}
