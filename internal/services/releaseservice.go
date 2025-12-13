package services

import (
	"fmt"

	"github.com/gclkaze/evamodulerepositoryserver/internal/dto"
	"github.com/gclkaze/evamodulerepositoryserver/internal/models"
	"github.com/gclkaze/evamodulerepositoryserver/internal/repositories"
)

type ReleaseService struct {
	repo             repositories.ReleaseRepository
	statusRepo       repositories.ReleaseStatusRepository
	userService      *UserService
	ownershipService *ModuleOwnershipService
}

func NewReleaseService(repo repositories.ReleaseRepository, statusRepo repositories.ReleaseStatusRepository, ownershipService *ModuleOwnershipService, u *UserService) *ReleaseService {
	return &ReleaseService{repo: repo, statusRepo: statusRepo, ownershipService: ownershipService, userService: u}
}

func (s *ReleaseService) DeleteModuleRelease(userID uint, id uint, releaseID uint) (bool, error) {
	return s.repo.DeleteModuleRelease(userID, id, releaseID)
}

func (s *ReleaseService) CancelSuggestedModuleRelease(userID uint, modID uint, releaseID uint) (bool, error) {

	st, err := s.statusRepo.GetStatus(repositories.Accepted)
	if err != nil {
		return false, err
	}

	release, err := s.repo.GetModuleReleaseWithStatus(modID, releaseID, st.ID)
	if err != nil {
		return false, err
	}

	//lets check if the user is the one suggested it
	acc, err := s.userService.GetDevelopersUserAccount(&release.Creator)
	if err != nil {
		return false, err
	}
	if acc.ID != userID {
		return false, fmt.Errorf("Only initiator can cancel the release")
	}

	st, err = s.statusRepo.GetStatus(repositories.Canceled)
	if err != nil {
		return false, err
	}

	release.Status = *st
	release.StatusID = st.ID
	err = s.repo.Update(release)
	if err != nil {
		return false, err
	}
	return true, nil
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
		return 0, fmt.Errorf("there are is a pending release of that module, need to cancel, reject, accept it first to create a new release")
	}
	st, err := s.statusRepo.GetStatus(repositories.Pending)
	if err != nil {
		return 0, err
	}

	dmo, err := s.ownershipService.FindDeveloperModuleOwner(&mod.Owner)
	if err != nil {
		return 0, err
	}

	dev := dmo.Developer
	newRelease := models.NewModuleReleaseFromModule(mod, version, *st, sz, dev)
	s.repo.Create(newRelease)
	//res, err = s.userHasPendingRelease(devID)
	return newRelease.ID, nil
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
	for i := range results {
		dtos = append(dtos, *dto.NewReleaseDTO(results[i]))
	}
	return dtos, error
}

func (s *ReleaseService) Initialize() {
	s.statusRepo.Initialize()
}
