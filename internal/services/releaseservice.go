package services

import (
	"fmt"

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
	st, err := s.statusRepo.GetStatus(repositories.Pending)
	if err != nil {
		return false, err
	}

	releases, err := s.repo.GetModuleReleasesWithStatus(modID, st.ID)
	if err != nil {
		return false, err
	}
	return len(releases) != 0, nil
}

func (s ReleaseService) userHasPendingRelease(userID uint) (bool, error) {
	return false, nil
}

func (s *ReleaseService) SuggestUserModuleRelease(userID uint, mod *models.Module, version string, sz int64) (uint, error) {
	/*dmo, err := s.ownershipService.FindDeveloperModuleOwner(&mod.Owner)
	if err != nil {
		return 0, err
	}*/
	//devID := dmo.DeveloperID
	modID := mod.ID

	res, err := s.moduleHasPendingRelease(modID)
	if err != nil {
		return 0, err
	}

	if res {
		return 0, fmt.Errorf("therer are is a pending release of that module, need to cancel, reject, accept it first to create a new release.")
	}
	st, err := s.statusRepo.GetStatus(repositories.Pending)
	if err != nil {
		return 0, err
	}
	newRelease := models.NewModuleReleaseFromModule(mod, version, *st, sz)
	s.repo.Create(newRelease)
	//res, err = s.userHasPendingRelease(devID)
	return newRelease.ID, nil
}

func (s ReleaseService) GetFolderSize()

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
