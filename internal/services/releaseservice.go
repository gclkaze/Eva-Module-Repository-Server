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

func (s *ReleaseService) GetUserService() *UserService {
	return s.userService
}

func (s *ReleaseService) DeleteModuleRelease(userID uint, modID uint, releaseID uint) (bool, error) {
	if s.userService.UserHasPermission(userID, models.DeleteModules) {
		return s.repo.DeleteModuleRelease(userID, modID, releaseID)
	}

	release, err := s.repo.GetModuleRelease(modID, releaseID)
	if err != nil {
		return false, err
	}
	acc, err := s.userService.GetDevelopersUserAccount(&release.Creator)
	if err != nil {
		return false, err
	}
	if acc.ID != userID {
		return false, fmt.Errorf("only initiator can delete the release")
	}

	return s.repo.DeleteModuleRelease(userID, modID, releaseID)
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
	if !s.userService.UserHasPermission(userID, models.CancelReleases) {
		acc, er := s.userService.GetDevelopersUserAccount(&release.Creator)
		if er != nil {
			return false, er
		}
		if acc.ID != userID {
			return false, fmt.Errorf("only initiator can cancel the release")
		}
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

func (s *ReleaseService) AcceptModuleRelease(userID uint, releaseID uint) (uint, error) {
	rel, err := s.repo.FindByID(releaseID)
	if err != nil {
		return 0, err
	}
	if rel.Status.Label == repositories.Accepted.String() {
		return 0, fmt.Errorf("the release is already accepted")
	}

	st, err := s.statusRepo.GetStatus(repositories.Accepted)
	if err != nil {
		return 0, err
	}

	rel.Status = *st
	rel.StatusID = st.ID
	err = s.repo.Update(rel)
	if err != nil {
		return 0, err
	}
	return releaseID, nil
}

func (s *ReleaseService) RejectModuleRelease(userID uint, releaseID uint) (uint, error) {
	rel, err := s.repo.FindByID(releaseID)
	if err != nil {
		return 0, err
	}
	if rel.Status.Label == repositories.Rejected.String() {
		return 0, fmt.Errorf("the release is already rejected")
	}

	st, err := s.statusRepo.GetStatus(repositories.Rejected)
	if err != nil {
		return 0, err
	}

	rel.Status = *st
	rel.StatusID = st.ID
	err = s.repo.Update(rel)
	if err != nil {
		return 0, err
	}
	return releaseID, nil
}

func (s *ReleaseService) CancelModuleRelease(userID uint, releaseID uint) (uint, error) {
	rel, err := s.repo.FindByID(releaseID)
	if err != nil {
		return 0, err
	}
	if rel.Status.Label == repositories.Canceled.String() {
		return 0, fmt.Errorf("the release is already canceled")
	}
	st, err := s.statusRepo.GetStatus(repositories.Canceled)
	if err != nil {
		return 0, err
	}

	rel.Status = *st
	rel.StatusID = st.ID
	err = s.repo.Update(rel)
	if err != nil {
		return 0, err
	}
	return releaseID, nil
}

func (s *ReleaseService) ChangeToPendingModuleRelease(userID uint, releaseID uint) (uint, error) {
	rel, err := s.repo.FindByID(releaseID)
	if err != nil {
		return 0, err
	}
	if rel.Status.Label == repositories.Pending.String() {
		return 0, fmt.Errorf("the release is already pending")
	}
	st, err := s.statusRepo.GetStatus(repositories.Pending)
	if err != nil {
		return 0, err
	}

	rel.Status = *st
	rel.StatusID = st.ID
	err = s.repo.Update(rel)
	if err != nil {
		return 0, err
	}
	return releaseID, nil
}

func (s *ReleaseService) Initialize() {
	s.statusRepo.Initialize()
}
