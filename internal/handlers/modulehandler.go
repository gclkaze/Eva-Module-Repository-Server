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

func (h *ModuleHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	idUint, err := utils.StringToUint(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid ID format",
			/*			"details": err.Error(),*/
		})
		return
	}
	res, err := h.service.Delete(idUint)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, res)
}

func (h *ModuleHandler) Upload(c *gin.Context) {
	// Parse form fields
	userID := c.PostForm("userId")
	title := c.PostForm("title")
	repr := c.PostForm("repr")
	tags := c.PostForm("tags")
	description := c.PostForm("description")

	idUint, err := utils.StringToUint(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid User ID format",
			/*			"details": err.Error(),*/
		})
		return
	}

	// Validate required fields
	if title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title is required"})
		return
	}

	// Handle file upload
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File upload failed"})
		return
	}

	// Call service to create module
	id, err := h.service.CreateModule(idUint, title, description, repr, file, tags, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Save file to disk (example: ./uploads/)
	/*    uploadPath := fmt.Sprintf("./uploads/%s", file.Filename)
	      if err := c.SaveUploadedFile(file, uploadPath); err != nil {
	          c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
	          return
	      }

	      // Create record in DB
	      module := Module{
	          Title:       title,
	          Repr:        repr,
	          Description: description,
	          FilePath:    uploadPath,
	      }

	      if err := db.Create(&module).Error; err != nil {
	          c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	          return
	      }
	*/
	// Respond with created resource
	c.JSON(http.StatusCreated, gin.H{
		"message": "Module created successfully",
		"module":  id,
	})
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
