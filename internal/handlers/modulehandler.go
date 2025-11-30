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
			"error": "Invalid ID format",
			/*			"details": err.Error(),*/
		})
		return

	}
	module, err := h.service.FindByID(idUint)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, module)
}

func (h *ModuleHandler) SearchModulesByTags(c *gin.Context) {
	tagsQuery := c.Query("tags")
	if tagsQuery == "" {
		c.JSON(400, gin.H{"error": "tags query parameter is required"})
		return
	}

	labels := strings.Split(tagsQuery, ",")
	modules, err := h.service.SearchByKeywords(labels)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, modules)
}
