// Package handlers contains the controllers of the Repository Server
package handlers

import (
	"net/http"
	"strings"

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
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid ID format",
			"result":  false,
			"details": err.Error(),
		})
		return
	}

	module, err := h.service.FindByID(idUint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Couldn't find module",
			"result":  false,
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Module information retrieved successfully",
		"releases": module,
		"result":   true,
	})
}

func (h *ModuleHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetUint("userId")
	idUint, err := utils.StringToUint(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid ID format",
			"result":  false,
			"details": err.Error(),
		})
		return
	}
	res, err := h.service.Delete(userID, idUint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Couldn't delete module",
			"result":  false,
			"details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Module deleted successfully",
		"result":  res,
	})
}

func (h *ModuleHandler) Upload(c *gin.Context) {
	userID := c.GetUint("userId")
	title := c.PostForm("title")
	repr := c.PostForm("repr")
	tags := c.PostForm("tags")
	description := c.PostForm("description")

	// Validate required fields
	if title == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Title is required",
			"result":  false,
			"details": "",
		})
		return
	}

	// Handle file upload
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "File upload failed",
			"result":  false,
			"details": err.Error()})
		return
	}

	// Call service to create module
	id, err := h.service.CreateModuleTx(userID, title, description, repr, file, tags, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Couldn't upload module",
			"result":  false,
			"details": err.Error()})
		return
	}

	// Respond with created resource
	c.JSON(http.StatusOK, gin.H{
		"message": "Module code was uploaded successfully",
		"module":  id,
		"result":  true,
	})
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
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid Module ID format",
			"result":  false,
			"details": err.Error(),
		})
		return
	}

	if title == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Title is required",
			"result":  false,
			"details": "",
		})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "File upload failed",
			"result":  false,
			"details": err.Error(),
		})
		return
	}

	id, err := h.service.UpdateUserModule(userID, modIDUint, title, description, repr, file, tags, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Couldn't update module",
			"result":  false,
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Module was updated successfully",
		"module":  id,
		"result":  true,
	})
}

func (h *ModuleHandler) SuggestRelease(c *gin.Context) {
	userID := c.GetUint("userId")
	modID := c.PostForm("modId")
	version := c.PostForm("version")

	modIDUint, err := utils.StringToUint(modID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid Module ID format",
			"result":  false,
			"details": err.Error(),
		})
		return
	}

	id, err := h.service.SuggestUserModuleRelease(userID, modIDUint, version)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Couldn't suggest module release",
			"result":  false,
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Module suggested for release successfully",
		"release": id,
		"result":  true,
	})
}

func (h *ModuleHandler) SearchModulesByTags(c *gin.Context) {
	tagsQuery := c.Query("tags")
	if tagsQuery == "" {
		c.JSON(400, gin.H{
			"error":  "tags query parameter is required",
			"result": false,
		})
		return
	}

	labels := strings.Split(tagsQuery, ",")
	modules, err := h.service.SearchByKeywords(labels)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Couldn't find the modules",
			"result":  false,
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Module search was complete",
		"modules": modules,
		"result":  true,
	})
}
