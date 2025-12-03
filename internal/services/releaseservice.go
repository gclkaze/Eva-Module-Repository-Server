package services

import (
	"github.com/gclkaze/evamodulerepositoryserver/internal/dto"
	"github.com/gclkaze/evamodulerepositoryserver/internal/models"
	"github.com/gclkaze/evamodulerepositoryserver/internal/repositories"
)

type ReleaseService struct {
	repo       repositories.ReleaseRepository
	statusRepo repositories.ReleaseStatusRepository

	ownershipService *ModuleOwnershipService
}

func NewReleaseService(repo repositories.ReleaseRepository, statusRepo repositories.ReleaseStatusRepository, ownershipService *ModuleOwnershipService) *ReleaseService {
	return &ReleaseService{repo: repo, statusRepo: statusRepo, ownershipService: ownershipService}
}

func (s *ReleaseService) DeleteModuleRelease(id uint, releaseID uint) (bool, error) {
	return s.repo.DeleteModuleRelease(id, releaseID)
}

func (s ReleaseService) moduleHasPendingRelease(modID uint) (bool, error) {
	releases, err := s.repo.GetModuleReleasesWithStatus(modID, 1)

	return false, nil
}

func (s ReleaseService) userHasPendingRelease(userID uint) (bool, error) {
	return false, nil
}

func (s *ReleaseService) SuggestUserModuleRelease(userID uint, mod *models.Module, version string) (uint, error) {
	dmo, err := s.ownershipService.FindDeveloperModuleOwner(&mod.Owner)
	if err != nil {
		return 0, err
	}
	devID := dmo.DeveloperID
	modID := mod.ID

	res, err := s.moduleHasPendingRelease(modID)
	res, err = s.userHasPendingRelease(devID)
	return 1, nil
}

func (s *ReleaseService) GetModuleRelease(id uint, releaseID uint) (*dto.ReleaseDTO, error) {
	result, error := s.repo.GetModuleRelease(id, releaseID)
	if error != nil {
		return nil, error
	}
	if result == nil {
		return nil, nil
	}
	res := dto.NewReleaseDTO(*result)
	return res, error
}

func (s *ReleaseService) GetModuleReleases(id uint) ([]dto.ReleaseDTO, error) {
	results, error := s.repo.GetModuleReleases(id)
	if error != nil {
		return nil, error
	}
	if results == nil {
		return nil, nil
	}
	var dtos []dto.ReleaseDTO
	for i := 0; i < len(results); i++ {
		dtos = append(dtos, *dto.NewReleaseDTO(results[i]))
	}
	return dtos, error
}

func (s *ReleaseService) SearchByKeywords(id uint, tags []string) ([]dto.ReleaseDTO, error) {
	results, error := s.repo.SearchModuleReleasesByTags(id, tags)
	if error != nil {
		return nil, error
	}
	if results == nil {
		return nil, nil
	}
	var dtos []dto.ReleaseDTO
	for i := 0; i < len(results); i++ {
		dtos = append(dtos, *dto.NewReleaseDTO(results[i]))
	}
	return dtos, error
}

func (s *ReleaseService) Initialize() {
	s.statusRepo.Initialize()
}
