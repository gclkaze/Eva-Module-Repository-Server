// Package handlers contains the controllers of the Repository Server
package handlers

import (
	"strings"

	"github.com/gclkaze/evamodulerepositoryserver/internal/services"
	"github.com/gin-gonic/gin"
)

type ModuleHandler struct {
	service *services.ModuleService
}

func NewModuleHandler(service *services.ModuleService) *ModuleHandler {
	return &ModuleHandler{service: service}
}
func (h *ModuleHandler) GetModuleByID(c *gin.Context) {

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
