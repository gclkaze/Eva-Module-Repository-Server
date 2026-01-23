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

package services

import (
	"fmt"
	"path"

	"github.com/gclkaze/evamodulerepositoryserver/internal/models"
	"github.com/gclkaze/evamodulerepositoryserver/internal/repositories"
	"github.com/gclkaze/evamodulerepositoryserver/pkg/utils"
	"github.com/magiconair/properties"
	"gorm.io/gorm"
)

type DownloadService struct {
	service             *ModuleService
	defaultDistFilename string
}

func NewDownloadService(service *ModuleService, p *properties.Properties) *DownloadService {
	s := "dist.tar.gz"
	if p != nil {
		s = p.GetString("dist_name", s)
	}
	return &DownloadService{service: service, defaultDistFilename: s}
}

func (h *DownloadService) DownloadRelease(releaseID uint) (string, string, error) {
	//the release needs to be ACCEPTED
	filename := h.defaultDistFilename
	rel, err := h.service.GetReleaseService().GetRelease(releaseID)
	if err != nil {
		return "", "", fmt.Errorf("couldn't find Release %d", releaseID)
	}

	if rel.Status.Label != repositories.Accepted.String() {
		return "", "", fmt.Errorf("no release ACCEPTED found with id %d", releaseID)
	}

	mod, err := h.service.GetModule(rel.ModuleID)
	if err != nil {
		return "", "", fmt.Errorf("couldn't find Release Module related to release id %d", releaseID)
	}

	dest := h.service.GetModuleReleasePath(mod, rel)
	if !utils.FolderExists(dest) {
		return "", "", fmt.Errorf("couldn't find Release Artifact related to release id %d", releaseID)
	}

	dest = path.Join(dest, filename)
	if !utils.FileExists(dest) {
		return "", "", fmt.Errorf("couldn't find Release Artifact related to release id %d", releaseID)
	}
	err = h.IncreaseDownloadCounter(rel)
	if err != nil {
		h.service.logger.Errorf("download service", "couldn't increase the download counter for release %d, got error %s", releaseID, err.Error())
	}

	return dest, filename, nil
}

func (h DownloadService) DownloadAnyRelease(releaseID uint) (string, string, error) {
	filename := h.defaultDistFilename

	rel, err := h.service.GetReleaseService().GetRelease(releaseID)
	if err != nil {
		return "", "", fmt.Errorf("couldn't find Release %d", releaseID)
	}

	mod, err := h.service.GetModule(rel.ModuleID)
	if err != nil {
		return "", "", fmt.Errorf("couldn't find Release Module related to release id %d", releaseID)
	}

	dest := h.service.GetModuleReleasePath(mod, rel)
	if !utils.FolderExists(dest) {
		return "", "", fmt.Errorf("couldn't find Release Artifact related to release id %d", releaseID)
	}

	dest = path.Join(dest, filename)
	if !utils.FileExists(dest) {
		return "", "", fmt.Errorf("couldn't find Release Artifact related to release id %d", releaseID)
	}

	err = h.IncreaseDownloadCounter(rel)
	if err != nil {
		h.service.logger.Errorf("download service", "couldn't increase the download counter for release %d, got error %s", releaseID, err.Error())
	}
	return dest, filename, nil
}

func (h DownloadService) IncreaseDownloadCounter(release *models.ModuleRelease) error {

	if release.Statistics == nil {
		return nil
	}
	db := h.service.repo.GetDB()
	err := utils.WithGormTransaction(db, func(tx *gorm.DB) error {

		if err := tx.Model(&models.ReleaseStatistics{}).
			Where("id = ?", release.Statistics.ID).
			UpdateColumn("download_count", gorm.Expr("download_count + 1")).Error; err != nil {
			return err
		}

		return nil
	})

	return err
}

func (h *DownloadService) DownloadPublicRelease(release string) (string, string, error) {
	//the release needs to be ACCEPTED
	filename := h.defaultDistFilename
	//here we smash the release string
	moduleName, releaseVersion, err := utils.ParseModuleReleaseVersion(release)
	if err != nil {
		return "", "", err
	}
	rel, err := h.service.GetReleaseService().GetTheModuleRelease(moduleName, releaseVersion)
	if err != nil {
		return "", "", fmt.Errorf("couldn't find Module Release %s", release)
	}

	if rel.Status.Label != repositories.Accepted.String() {
		return "", "", fmt.Errorf("no ACCEPTED release found with name %s", release)
	}

	mod, err := h.service.GetModule(rel.ModuleID)
	if err != nil {
		return "", "", fmt.Errorf("couldn't find Release Module related to release %s %d", release)
	}

	dest := h.service.GetModuleReleasePath(mod, rel)
	if !utils.FolderExists(dest) {
		return "", "", fmt.Errorf("couldn't find Release Artifact related to release %s", release)
	}

	dest = path.Join(dest, filename)
	if !utils.FileExists(dest) {
		return "", "", fmt.Errorf("couldn't find Release Artifact related to release %s", release)
	}
	err = h.IncreaseDownloadCounter(rel)
	if err != nil {
		h.service.logger.Errorf("download service", "couldn't increase the download counter for release %s, got error %s", release, err.Error())
	}

	return dest, filename, nil
}

func (h *DownloadService) AuthUserDownloadSpecificRelease(userID uint, release string) (bool, string, string, error) {
	perms, err := h.service.GetReleaseService().GetUserService().GetUserPermissions(userID)
	if err != nil {
		return false, "", "", fmt.Errorf("unknown user asking for %s", release)
	}
	userCanGetAnyRelease := false
	for i := range perms {
		if perms[i].Value == models.ChangeReleaseStatuses.String() {
			userCanGetAnyRelease = true
			break
		}
	}

	//the release needs to be ACCEPTED
	filename := h.defaultDistFilename
	//here we smash the release string
	moduleName, releaseVersion, err := utils.ParseModuleReleaseVersion(release)
	if err != nil {
		return false, "", "", err
	}
	rel, err := h.service.GetReleaseService().GetTheStatusIndependentModuleRelease(moduleName, releaseVersion)
	if err != nil {
		return false, "", "", fmt.Errorf("couldn't find Module Release %s", release)
	}

	if !userCanGetAnyRelease {
		if rel.Status.Label != repositories.Accepted.String() {
			return false, "", "", fmt.Errorf("no ACCEPTED release found with name %s", release)
		}
	}

	mod, err := h.service.GetModule(rel.ModuleID)
	if err != nil {
		return false, "", "", fmt.Errorf("couldn't find Release Module related to release %s", release)
	}

	return h.RetrieveTarball(mod, rel, filename, release)
}

func (h DownloadService) RetrieveTarball(mod *models.Module, rel *models.ModuleRelease, filename string, release string) (bool, string, string, error) {
	if rel.Status.Label == repositories.Pending.String() {
		//need to check on the user's module to find it, not under releases
		dmo, err := h.service.GetModuleOwnershipService().FindDeveloperModuleOwner(&mod.Owner)
		if err != nil {
			return false, "", "", err
		}

		modPath := h.service.GetModulePath(dmo, mod)
		theDest, err := h.service.GetReleaseService().CreateTemporaryTarBall(modPath)
		if err != nil {
			return false, "", "", fmt.Errorf("couldn't create temporary Release Artifact related to release %s", release)
		}
		return true, theDest, filename, nil
	}

	dest := h.service.GetModuleReleasePath(mod, rel)
	if !utils.FolderExists(dest) {
		return false, "", "", fmt.Errorf("couldn't find Release Artifact related to release %s", release)
	}

	dest = path.Join(dest, filename)
	if !utils.FileExists(dest) {
		return false, "", "", fmt.Errorf("couldn't find Release Artifact related to release %s", release)
	}
	err := h.IncreaseDownloadCounter(rel)
	if err != nil {
		h.service.logger.Errorf("download service", "couldn't increase the download counter for release %s, got error %s", release, err.Error())
	}

	return false, dest, filename, nil
}
