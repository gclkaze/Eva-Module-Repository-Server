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

// Package services contains
package services

import (
	"fmt"
	"mime/multipart"
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
	"gorm.io/gorm"
)

type ModuleService struct {
	repo    repositories.ModuleRepository
	keyword repositories.KeywordRepository

	moduleFolder  string
	releaseFolder string
	devFolder     string

	logger logger.ILogger

	developerService *UserService
	ownershipService *ModuleOwnershipService
	releaseService   *ReleaseService
}

func (s *ModuleService) GetReleaseService() *ReleaseService {
	return s.releaseService
}

func (s *ModuleService) GetModuleOwnershipService() *ModuleOwnershipService {
	return s.ownershipService
}

func NewModuleService(repo repositories.ModuleRepository, dev *UserService, p *properties.Properties, ownershipService *ModuleOwnershipService, keywordRepo repositories.KeywordRepository, releaseService *ReleaseService) (*ModuleService, error) {
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
	mod.releaseService = releaseService

	mod.releaseService.SetModuleService(mod)
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

func (s ModuleService) GetReleasePath(mod *models.Module, release *models.ModuleRelease) string {
	mName := utils.GetRepoName(mod.Repr)
	path := fmt.Sprintf("%s/%s", s.releaseFolder, mName)
	return path
}

func (s ModuleService) GetReleaseBasePath(mod *models.Module) string {
	mName := utils.GetRepoName(mod.Repr)
	path := fmt.Sprintf("%s/%s", s.releaseFolder, mName)
	return path
}

func (s ModuleService) GetReleaseBasePathWithName(name string) string {
	path := fmt.Sprintf("%s/%s", s.releaseFolder, name)
	return path
}

func (s ModuleService) GetReleaseBasePathWithNameAndVersion(name string, version string) string {
	path := fmt.Sprintf("%s/%s/%s", s.releaseFolder, name, version)
	return path
}

func (s ModuleService) GetModuleReleasePath(mod *models.Module, release *models.ModuleRelease) string {
	devPath := s.GetReleasePath(mod, release)
	version := release.Version
	path := fmt.Sprintf("%s/%s", devPath, version)
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

	dev, err := s.developerService.FindByID(userID)
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

	repr = strings.TrimSpace(repr)
	if !utils.IsValidModuleName(repr) {
		err = fmt.Errorf("the module handle %s should not contain strange characters and its length should be between 3 and 20 characters long", repr)
		return 0, err
	}
	reprExists, err := s.ReprExists(repr)
	if err != nil {
		return 0, err
	}

	if reprExists {
		err = fmt.Errorf("the module handle %s is already taken..use a different one", repr)
		return 0, err
	}
	mod := models.NewModule(title, repr, descr, owner.ID, *owner, keywords, utils.GetRepoName(repr))
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

func (s ModuleService) GetTheModulePath(mod *models.Module) {

}

func (s ModuleService) ReprExists(repr string) (bool, error) {
	return s.repo.ReprExists(repr)
}

func (s ModuleService) ReprExistsExcept(id uint, repr string) (bool, error) {
	return s.repo.ReprExistsExcept(id, repr)
}

func (s *ModuleService) CreateModuleTx(
	userID uint,
	title string,
	descr string,
	repr string,
	files []*multipart.FileHeader,
	tags string,
	c *gin.Context,
) (uint, error) {

	if files == nil {
		return 0, fmt.Errorf("no module file was provided")
	}

	if len(files) == 0 {
		return 0, fmt.Errorf("no module file was provided")
	}

	title = strings.TrimSpace(title)
	if title == "" {
		return 0, fmt.Errorf("the module title cannot be empty")
	}

	repr = strings.TrimSpace(repr)

	if !utils.IsValidModuleName(repr) {
		err := fmt.Errorf("the module handle %s should not contain strange characters and its length should be between 3 and 50 characters long", repr)
		return 0, err
	}
	reprExists, errFileExists := s.ReprExists(repr)
	if errFileExists != nil {
		return 0, errFileExists
	}
	if reprExists {
		return 0, fmt.Errorf("the module handle %s is already taken..use a different one", repr)
	}

	var (
		mod     *models.Module
		dev     *models.Developer
		owner   *models.ModuleOwner
		dmo     *models.DeveloperModuleOwner
		modPath string
	)

	db := s.repo.GetDB()
	err := utils.WithGormTransaction(db, func(tx *gorm.DB) error {

		var err error

		dev, err = s.developerService.FindByIDTx(tx, userID)
		if err != nil {
			return err
		}

		owner, err = s.ownershipService.GetModuleOwner(models.Dev, dev.ID)
		if err != nil {
			return err
		}
		if owner == nil {
			owner, err = s.ownershipService.CreateModuleOwnerTx(tx, models.Dev, dev.ID)
			if err != nil {
				return err
			}

			dmo, err = s.ownershipService.CreateDeveloperModuleOwnerTx(tx, dev, owner)
			if err != nil {
				return err
			}
		} else {
			dmo, err = s.ownershipService.FindDeveloperModuleOwner(owner)
			if err != nil {
				return err
			}
		}

		labels := strings.Split(tags, ",")
		keywords, err := s.CreateAndGetKeywordsTx(tx, labels)
		if err != nil {
			return err
		}

		mod = models.NewModule(title, repr, descr, owner.ID, *owner, keywords, utils.GetRepoName(repr))

		if err := s.repo.CreateTx(tx, mod); err != nil {
			return err
		}

		// prepare path, but DO NOT create folders yet
		modPath = s.GetModulePath(dmo, mod)

		if utils.FolderExists(modPath) {
			return fmt.Errorf("module path exists %s", modPath)
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	devPath := s.GetDevPath(dev)
	if !utils.FolderExists(devPath) {
		if err := utils.CreateFolder(devPath); err != nil {
			return 0, err
		}
	}

	if err := utils.CreateFolder(modPath); err != nil {
		return 0, err
	}

	for i := range files {
		file := files[i]
		filePath := fmt.Sprintf("%s/%s", modPath, file.Filename)
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			s.logger.Errorf("module_service", "%s", err)
			return 0, err
		}

	}

	return mod.ID, nil
}

func (s *ModuleService) UpdateUserModule(userID uint, modID uint, title string, descr string, repr string, files []*multipart.FileHeader, tags string, c *gin.Context) (uint, error) {
	mod, err := s.GetModule(modID)
	if err != nil {
		return 0, err
	}

	repr = strings.TrimSpace(repr)
	if !utils.IsValidModuleName(repr) {
		return 0, fmt.Errorf("the module handle %s should not contain strange characters and its length should be between 3 and 20 characters long", repr)
	}
	reprExists, errReprExists := s.ReprExistsExcept(modID, repr)
	if errReprExists != nil {
		return 0, errReprExists
	}
	if reprExists {
		err = fmt.Errorf("the module handle %s is already taken..use a different one", repr)
		return 0, err
	}

	if !s.releaseService.userService.UserHasPermission(userID, models.UpdateModules) {

		if mod.Owner.EntityID != userID {
			return 0, fmt.Errorf("user with ID %d didn't match with the module", userID)
		}

		if mod.Owner.Type.Label != models.Dev.String() {
			return 0, fmt.Errorf("user with ID %d didn't match with the module type", userID)
		}
	}

	labels := strings.Split(tags, ",")
	keywords, err := s.CreateAndGetKeywords(labels)
	if err != nil {
		return 0, err
	}

	mod.Update(title, repr, descr, keywords /*, utils.GetRepoName(repr)*/)
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

	for i := range files {
		file := files[i]
		filePath := fmt.Sprintf("%s/%s", modPath, file.Filename)
		//filePath := fmt.Sprintf("%s/%s", modPath, file.Filename)
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			s.logger.Errorf("module_service", "%s", err)
			return 0, err
		}

	}

	return mod.ID, nil
}

func (s ModuleService) GetFolderSize(mod *models.Module) (int64, error) {

	dmo, err := s.ownershipService.FindDeveloperModuleOwner(&mod.Owner)
	if err != nil {
		return 0, err
	}

	modPath := s.GetModulePath(dmo, mod)
	if !utils.FolderExists(modPath) {
		err = fmt.Errorf("module path does not exist %s", modPath)
		return 0, err
	}

	sz, err := utils.ComputeFolderSizeBytes(modPath)
	if err != nil {
		return 0, err
	}
	return sz, nil
}

func (s *ModuleService) SuggestUserModuleRelease(userID uint, modID uint, version string) (uint, error) {
	mod, err := s.GetModule(modID)
	if err != nil {
		return 0, err
	}
	if mod == nil {
		return 0, fmt.Errorf("couldn't find module with id %d that was suggested by user %d", modID, userID)
	}

	if mod.Owner.EntityID != userID {
		return 0, fmt.Errorf("user with ID %d didn't match with the module", userID)
	}

	if mod.Owner.Type.Label != models.Dev.String() {
		return 0, fmt.Errorf("user with ID %d didn't match with the module type", userID)
	}

	diskSize, err := s.GetFolderSize(mod)
	if err != nil {
		return 0, err
	}
	if diskSize == 0 {
		return 0, fmt.Errorf("the module size is 0. Cannot proceed with the release")
	}

	if !utils.IsValidVersion(version) {
		return 0, fmt.Errorf("invalid version suggested version string %s", version)
	}

	return s.releaseService.SuggestUserModuleRelease(userID, mod, version, diskSize)
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

func (s *ModuleService) GetModuleReleaseInfo(moduleName string, releaseName string) (*dto.ModuleEnrichedDTO, error) {
	moduleName = strings.TrimSpace(moduleName)
	if moduleName == "" {
		return nil, fmt.Errorf("empty module name provided")
	}

	releaseName = strings.TrimSpace(releaseName)
	if releaseName == "" {
		return nil, fmt.Errorf("empty module version provided")
	}

	//lets see if we know this version token, is it equal to "latest"?
	moduleName = strings.ToLower(moduleName)
	releaseName = strings.ToLower(releaseName)

	theModule, err := s.repo.FindByNameOrReprNameOrRepoName(moduleName)
	if err != nil {
		return nil, err
	}
	if theModule == nil {
		return nil, fmt.Errorf("no module found with name %s", moduleName)
	}

	if releaseName == "latest" {
		var theRelease *models.ModuleRelease
		theRelease, err = s.releaseService.GetLastModuleRelease(theModule.ID)
		if err != nil {
			return nil, err
		}
		if theRelease == nil {
			return nil, fmt.Errorf("no release found with version %s@%s", moduleName, releaseName)
		}
		return dto.NewModuleEnrichedDTO(theModule, []models.ModuleRelease{*theRelease}), nil
	}

	//if not, we need to retrieve the release from the db,the release needs to be accepted ofc!
	theRelease, err := s.releaseService.GetAcceptedModuleReleaseByVersion(theModule.ID, releaseName)
	if err != nil {
		return nil, err
	}
	if theRelease == nil {
		return nil, fmt.Errorf("no release found with version %s@%s", moduleName, releaseName)
	}

	return dto.NewModuleEnrichedDTO(theModule, []models.ModuleRelease{*theRelease}), nil
}

func (s *ModuleService) GetModuleInfo(moduleName string) (*dto.ModuleEnrichedDTO, error) {
	moduleName = strings.TrimSpace(moduleName)
	if moduleName == "" {
		return nil, fmt.Errorf("empty module name provided")
	}

	theModule, err := s.repo.FindByNameOrReprNameOrRepoName(moduleName)
	if err != nil {
		return nil, err
	}
	if theModule == nil {
		return nil, fmt.Errorf("no module found with name %s", moduleName)
	}

	//the module should have at least one ACCEPTED release so thus it is available in the REPOSITORY List
	id := theModule.ID
	moduleReleases, error := s.releaseService.GetModuleAcceptedReleases(id)
	if error != nil {
		return nil, error
	}
	if len(moduleReleases) == 0 {
		return nil, fmt.Errorf("the module with id %d has no accepted releases", id)
	}
	res := dto.NewModuleEnrichedDTO(theModule, moduleReleases)
	return res, error
}

func (s *ModuleService) GetMaxID() (uint, error) {
	return s.repo.GetMaxID()
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

func (s *ModuleService) deleteModuleReleases(userID uint, modID uint) (bool, error) {
	releaseIDS, err := s.releaseService.GetModuleReleaseIds(modID)
	if err != nil {
		return false, err
	}

	for i := range releaseIDS {
		res, err := s.releaseService.DeleteModuleRelease(userID, modID, releaseIDS[i])
		if err != nil {
			s.logger.Errorf("module service", "couldn't delete Module Release folder for mod: %d and release: %d", modID, releaseIDS[i])
			return false, err
		}
		if !res {
			s.logger.Errorf("module service", "couldn't delete Module Release folder for mod: %d and release: %d", modID, releaseIDS[i])
		}
	}
	return true, nil
}

func (s *ModuleService) deleteModule(userID uint, modID uint) (bool, error) {
	mod, err := s.GetModule(modID)
	if err != nil {
		s.logger.Errorf("module service", "couldn't find the module %d of user %d", modID, userID)
		return false, err
	}

	res, err := s.deleteModuleReleases(userID, modID)
	if res && err == nil {
		//let's delete also the module's folder
		dmo, dmoErr := s.ownershipService.FindDeveloperModuleOwner(&mod.Owner)
		if dmoErr != nil {
			return false, dmoErr
		}
		modPath := s.GetModulePath(dmo, mod)
		modDeletionRes, modDeletionErr := s.repo.Delete(modID)
		if modDeletionErr != nil {
			return false, modDeletionErr
		}
		if !modDeletionRes {
			s.logger.Errorf("module service", "the module has already been deleted %d with user %d", userID, modID)
		}
		//let's remove the folder
		if utils.FolderExists(modPath) {
			utils.CleanFolder(modPath)
			utils.RemoveFolder(modPath)
		}

		//does this module has any releases ?
		//we need to remove also these
		modpath := s.GetReleaseBasePath(mod)
		if utils.FolderExists(modpath) {
			err = utils.RemoveFolder(modpath)
			if err != nil {
				s.logger.Errorf("module service", "couldn't remove module folder %s while removing module with ID %d ", modpath, modID)
			}
		}
		//lets remove the dmo
		return s.ownershipService.Delete(dmo.ID)

	}
	return res, err
}

func (s *ModuleService) DeleteTx(userID uint, modID uint) (bool, error) {
	if s.releaseService.userService.UserHasPermission(userID, models.DeleteModules) {
		return s.deleteModule(userID, modID)
	}
	mod, err := s.GetModule(modID)
	if err != nil {
		return false, err
	}

	if mod.Owner.EntityID != userID {
		return false, fmt.Errorf("user with ID %d didn't match with the module", userID)
	}

	if mod.Owner.Type.Label != models.Dev.String() {
		return false, fmt.Errorf("user with ID %d didn't match with the module type", userID)
	}
	return s.deleteModule(userID, modID)
}

func (s *ModuleService) Delete(userID uint, modID uint) (bool, error) {
	if s.releaseService.userService.UserHasPermission(userID, models.DeleteModules) {
		return s.deleteModule(userID, modID)
	}
	mod, err := s.GetModule(modID)
	if err != nil {
		return false, err
	}

	if mod.Owner.EntityID != userID {
		return false, fmt.Errorf("user with ID %d didn't match with the module", userID)
	}

	if mod.Owner.Type.Label != models.Dev.String() {
		return false, fmt.Errorf("user with ID %d didn't match with the module type", userID)
	}
	return s.deleteModule(userID, modID)
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

func (s *ModuleService) CreateAndGetKeywordsTx(tx *gorm.DB, tags []string) ([]models.Keyword, error) {
	var keywords []models.Keyword
	for i := 0; i < len(tags); i++ {
		tag := tags[i]
		k, err := s.keyword.FindByLabelTx(tx, tag)
		if err != nil {
			return nil, err
		}
		if k == nil {
			//need to create it
			k = models.NewKeyword(tag)
			s.keyword.CreateTx(tx, k)
		}
		keywords = append(keywords, *k)
	}
	return keywords, nil
}

func (s ModuleService) SearchByComponents(nameTokens []string, descrTokens []string, tags []string) ([]dto.ModuleDTO, error) {
	if len(nameTokens) == 0 && len(descrTokens) == 0 && len(tags) == 0 {
		return nil, fmt.Errorf("no search criteria provided")
	}

	results, err := s.repo.SearchByComponents(nameTokens, descrTokens, tags)
	if err != nil {
		return nil, err
	}

	var dtos []dto.ModuleDTO
	for i := range results {
		releaseCNT, countErr := s.releaseService.GetReleaseAcceptedCountForModule(results[i].ID)
		if countErr != nil {
			s.logger.Errorf("module service", "error while counting accepted releases for the modules %s", countErr.Error())
			continue
		}
		dtos = append(dtos, *dto.NewModuleDTOWithReleaseInformation(results[i], releaseCNT))
	}
	return dtos, err
}

func (s ModuleService) GetUserModules(userID uint) ([]dto.ModuleDTO, error) {

	userAccount, err := s.developerService.FindUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("couldn't find user with ID %d", userID)
	}

	owner, err := s.ownershipService.GetModuleOwner(models.Dev, userAccount.ID)
	if err != nil {
		return nil, err
	}
	if owner == nil {
		return nil, nil
	}
	results, err := s.repo.FindByOwner(owner)
	if err != nil {
		return nil, err
	}
	var dtos []dto.ModuleDTO
	for i := range results {
		releaseCNT, countErr := s.releaseService.GetReleaseAcceptedCountForModule(results[i].ID)
		if countErr != nil {
			s.logger.Errorf("module service", "error while counting accepted releases for the modules %s", countErr.Error())
			continue
		}
		dtos = append(dtos, *dto.NewModuleDTOWithReleaseInformation(results[i], releaseCNT))
	}
	return dtos, err
}

func (s ModuleService) GetModulesByOrder(order []uint) ([]models.Module, error) {
	return s.repo.GetModules(order)
}

/*
func (s ModuleService) getColumnTokenConditionString(columnName string, tokens []string) (string, []any) {
	if len(tokens) == 0 {
		return "", nil
	}
	var conditions []string
	var args []any

	for _, t := range tokens {
		conditions = append(conditions, columnName+" LIKE ?")
		args = append(args, "%"+t+"%")
	}
	return strings.Join(conditions, " OR "), args
}*/
