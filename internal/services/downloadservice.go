package services

import (
	"fmt"
	"path"

	"github.com/gclkaze/evamodulerepositoryserver/internal/repositories"
	"github.com/gclkaze/evamodulerepositoryserver/pkg/utils"
	"github.com/magiconair/properties"
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

	return dest, filename, nil
}
