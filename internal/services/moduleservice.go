// Package services contains
package services

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gclkaze/evamodulerepositoryserver/internal/dto"
	"github.com/gclkaze/evamodulerepositoryserver/internal/models"
	"github.com/gclkaze/evamodulerepositoryserver/internal/repositories"
	"github.com/gclkaze/evamodulerepositoryserver/pkg/logger"
	"github.com/gclkaze/evamodulerepositoryserver/pkg/runtime"
	"github.com/gclkaze/evamodulerepositoryserver/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/magiconair/properties"
)

type ModuleService struct {
	repo repositories.ModuleRepository

	moduleFolder  string
	releaseFolder string

	logger logger.ILogger

	developerService *DeveloperService
}

func NewModuleService(repo repositories.ModuleRepository, dev *DeveloperService, p *properties.Properties) (*ModuleService, error) {
	moduleFolder := p.GetString("module_folder", "")
	releaseFolder := p.GetString("release_folder", "")

	l := runtime.CreateLogger(p)
	if !utils.FolderExists(moduleFolder) {
		err := fmt.Errorf("module folder %s does not exist", moduleFolder)
		l.Errorf("module service", "%s", err)
		return nil, err
	}

	if !utils.FolderExists(releaseFolder) {
		err := fmt.Errorf("release folder %s does not exist", releaseFolder)
		l.Errorf("module service", "%s", err)
		return nil, err
	}

	mod := &ModuleService{repo: repo}
	mod.moduleFolder = moduleFolder
	mod.releaseFolder = releaseFolder
	mod.developerService = dev

	mod.logger = l
	return mod, nil
}

func (s *ModuleService) Create(userID uint, title string, descr string, repr string, file *multipart.FileHeader, tags string, c *gin.Context) (uint, error) {
	if file == nil {
		return 0, fmt.Errorf("no module file was provided")
	}
	title = strings.TrimSpace(title)
	if title == "" {
		return 0, fmt.Errorf("the module title cannot be empty")
	}

	//check if the folder of the user exists
	fullPath := filepath.Join(s.moduleFolder, utils.UintToString(userID))
	if !utils.FolderExists(fullPath) {
		err := utils.CreateFolder(fullPath)
		if err != nil {
			return 0, fmt.Errorf("couldn't create user folder %s", fullPath)
		}
	}

	dev, err := s.developerService.FindById(userID)
	if err != nil {
		return 0, err
	}
	//the module needs to get first an id then save it under the path of the user + id
	mod := models.NewModule(title, repr, descr, ownerID, owner, keywords)
	err := s.repo.Create(mod)
	if err != nil {
		return 0, err
	}
	if err := c.SaveUploadedFile(file, uploadPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return 0, err
	}

	return 0, nil
}

func (s *ModuleService) FindByID(id uint) (*dto.ModuleDTO, error) {
	result, error := s.repo.FindByID(id)
	if error != nil {
		return nil, error
	}
	if result == nil {
		return nil, nil
	}
	res := dto.NewModuleDTO(*result)
	return res, error
}

func (s *ModuleService) Delete(id uint) (bool, error) {
	return s.repo.Delete(id)
}

func (s *ModuleService) SearchByKeywords(tags []string) ([]dto.ModuleDTO, error) {
	results, error := s.repo.SearchByKeywords(tags)
	if error != nil {
		return nil, error
	}
	var dtos []dto.ModuleDTO
	for i := 0; i < len(results); i++ {
		dtos = append(dtos, *dto.NewModuleDTO(results[i]))
	}
	return dtos, error
}
