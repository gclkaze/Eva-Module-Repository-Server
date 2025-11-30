// Package handlers contains the controllers of the Repository Server
package handlers

import (
	"github.com/gclkaze/evamodulerepositoryserver/internal/services"
	"github.com/gin-gonic/gin"
)

type ReleaseHandler struct {
	service *services.ReleaseService
}

func NewReleaseHandler(service *services.ReleaseService) *ReleaseHandler {
	return &ReleaseHandler{service: service}
}

/*
	func (h *ReleaseHandler) FindByID(c *gin.Context) {
		id := c.Param("id")

		idUint, err := utils.StringToUint(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid ID format",
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

	func (h *ReleaseHandler) SearchModulesByTags(c *gin.Context) {
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
*/
func (h *ReleaseHandler) GetModuleRelease(c *gin.Context) {

}

func (h *ReleaseHandler) GetModuleReleases(c *gin.Context) {

}
