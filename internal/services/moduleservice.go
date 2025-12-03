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
	repo    repositories.ModuleRepository
	keyword repositories.KeywordRepository

	moduleFolder  string
	releaseFolder string
	devFolder     string

	logger logger.ILogger

	developerService *DeveloperService
	ownershipService *ModuleOwnershipService
}

func NewModuleService(repo repositories.ModuleRepository, dev *DeveloperService, p *properties.Properties, ownershipService *ModuleOwnershipService, keywordRepo repositories.KeywordRepository) (*ModuleService, error) {
	moduleFolder := p.GetString("module_folder", "")
	releaseFolder := p.GetString("release_folder", "")
	devFolder := p.GetString("dev_folder", "")

	l := runtime.CreateLogger(p)
	if !utils.FolderExists(moduleFolder) {
		utils.CreateFolder(moduleFolder)
	}

	devPath := filepath.Join(moduleFolder, devFolder)
	if !utils.FolderExists(devPath) {
		utils.CreateFolder(devPath)
	}

	if !utils.FolderExists(releaseFolder) {
		utils.CreateFolder(releaseFolder)
	}

	mod := &ModuleService{repo: repo}
	mod.moduleFolder = moduleFolder
	mod.releaseFolder = releaseFolder
	mod.devFolder = devFolder
	mod.developerService = dev
	mod.keyword = keywordRepo

	mod.logger = l
	mod.ownershipService = ownershipService
	return mod, nil
}

func (s ModuleService) GetModulePath(dev *models.DeveloperModuleOwner, mod *models.Module) string {
	devPath := s.GetDevPath(&dev.Developer)
	modID := mod.ID
	path := fmt.Sprintf("%s/%d", devPath, modID)
	return path
}

func (s ModuleService) GetDevPath(dev *models.Developer) string {
	devID := dev.ID
	path := fmt.Sprintf("%s/%s/%d", s.moduleFolder, s.devFolder, devID)
	return path
}

func (s *ModuleService) CreateModule(userID uint, title string, descr string, repr string, file *multipart.FileHeader, tags string, c *gin.Context) (uint, error) {
	if file == nil {
		return 0, fmt.Errorf("no module file was provided")
	}
	title = strings.TrimSpace(title)
	if title == "" {
		return 0, fmt.Errorf("the module title cannot be empty")
	}

	dev, err := s.developerService.FindById(userID)
	if err != nil {
		return 0, err
	}

	//check if the folder of the user exists
	devPath := s.GetDevPath(dev)
	if !utils.FolderExists(devPath) {
		err = utils.CreateFolder(devPath)
		if err != nil {
			return 0, err
		}
	}

	//we need to make a module owner handle first to the developer
	owner, err := s.ownershipService.CreateModuleOwner(models.Dev, dev.ID)
	if err != nil {
		return 0, err
	}
	dmo, err := s.ownershipService.CreateDeveloperModuleOwner(dev, owner)
	if err != nil {
		return 0, err
	}

	//the module needs to get first an id then save it under the path of the user + id
	labels := strings.Split(tags, ",")
	keywords, err := s.CreateAndGetKeywords(labels)
	if err != nil {
		return 0, err
	}
	mod := models.NewModule(title, repr, descr, owner.ID, *owner, keywords)
	err = s.repo.Create(mod)
	if err != nil {
		return 0, err
	}

	modPath := s.GetModulePath(dmo, mod)
	if utils.FolderExists(modPath) {
		err = fmt.Errorf("module path exists %s", modPath)
		return 0, err
	}

	utils.CreateFolder(modPath)
	modPath = fmt.Sprintf("%s/%s", modPath, file.Filename)
	if err := c.SaveUploadedFile(file, modPath); err != nil {
		s.logger.Errorf("module_service", "%s", err)
		return 0, err
	}

	return mod.ID, nil
}

func (s *ModuleService) UpdateUserModule(userID uint, modID uint, title string, descr string, repr string, file *multipart.FileHeader, tags string, c *gin.Context) (uint, error) {
	mod, err := s.GetModule(modID)
	if err != nil {
		return 0, err
	}

	if mod.Owner.EntityID != userID {
		return 0, fmt.Errorf("user with ID %d didn't match with the module", userID)
	}

	if mod.Owner.Type.Label != models.Dev.String() {
		fmt.Print(models.ModuleOwnerTypeDef(mod.Owner.Type.ID))
		fmt.Print(models.Dev)

		return 0, fmt.Errorf("user with ID %d didn't match with the module type", userID)
	}
	labels := strings.Split(tags, ",")
	keywords, err := s.CreateAndGetKeywords(labels)
	if err != nil {
		return 0, err
	}

	mod.Update(title, repr, descr, keywords)
	err = s.repo.Update(mod)

	if err != nil {
		return 0, err
	}

	dmo, err := s.ownershipService.FindDeveloperModuleOwner(&mod.Owner)
	if err != nil {
		return 0, err
	}
	modPath := s.GetModulePath(dmo, mod)
	if !utils.FolderExists(modPath) {
		utils.CreateFolder(modPath)
	} else {
		utils.CleanFolder(modPath)
	}

	modPath = fmt.Sprintf("%s/%s", modPath, file.Filename)
	if err := c.SaveUploadedFile(file, modPath); err != nil {
		s.logger.Errorf("module_service", "%s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return 0, err
	}

	return mod.ID, nil
}

func (s *ModuleService) FindByID(id uint) (*dto.ModuleDTO, error) {
	result, error := s.repo.FindByID(id, false)
	if error != nil {
		return nil, error
	}
	if result == nil {
		return nil, nil
	}
	res := dto.NewModuleDTO(*result)
	return res, error
}

func (s *ModuleService) GetModule(id uint) (*models.Module, error) {
	result, error := s.repo.FindByID(id, true)
	if error != nil {
		return nil, error
	}
	if result == nil {
		return nil, nil
	}
	return result, nil
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

func (s *ModuleService) CreateAndGetKeywords(tags []string) ([]models.Keyword, error) {
	var keywords []models.Keyword
	for i := 0; i < len(tags); i++ {
		tag := tags[i]
		k, err := s.keyword.FindByLabel(tag)
		if err != nil {
			return nil, err
		}
		if k == nil {
			//need to create it
			k = models.NewKeyword(tag)
			s.keyword.Create(k)
		}
		keywords = append(keywords, *k)
	}
	return keywords, nil
}
