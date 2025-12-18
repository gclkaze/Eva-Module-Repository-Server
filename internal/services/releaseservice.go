package services

import (
	"fmt"
	"path"

	"github.com/gclkaze/evamodulerepositoryserver/internal/dto"
	"github.com/gclkaze/evamodulerepositoryserver/internal/models"
	"github.com/gclkaze/evamodulerepositoryserver/internal/repositories"
	"github.com/gclkaze/evamodulerepositoryserver/pkg/logger"
	"github.com/gclkaze/evamodulerepositoryserver/pkg/runtime"
	"github.com/gclkaze/evamodulerepositoryserver/pkg/utils"
	"github.com/magiconair/properties"
)

type ReleaseService struct {
	repo                repositories.ReleaseRepository
	statusRepo          repositories.ReleaseStatusRepository
	statRepo            repositories.ReleaseStatisticsRepository
	userService         *UserService
	ownershipService    *ModuleOwnershipService
	logger              logger.ILogger
	moduleService       *ModuleService
	defaultDistFilename string
}

func NewReleaseService(repo repositories.ReleaseRepository, statusRepo repositories.ReleaseStatusRepository, ownershipService *ModuleOwnershipService, u *UserService, p *properties.Properties, statRepo repositories.ReleaseStatisticsRepository) *ReleaseService {
	l := runtime.CreateLogger(p)
	s := "dist.tar.gz"
	if p != nil {
		s = p.GetString("dist_name", s)
	}
	return &ReleaseService{repo: repo, statusRepo: statusRepo, ownershipService: ownershipService, userService: u, logger: l, defaultDistFilename: s, statRepo: statRepo}
}

func (s *ReleaseService) GetUserService() *UserService {
	return s.userService
}

func (s *ReleaseService) SetModuleService(mod *ModuleService) {
	s.moduleService = mod
}
func (s ReleaseService) FindByID(id uint) (*models.ModuleRelease, error) {
	return s.repo.FindByID(id)
}

func (s ReleaseService) GetRelease(id uint) (*models.ModuleRelease, error) {
	return s.repo.GetRelease(id)
}

func (s *ReleaseService) removeModuleReleaseFolder(modID uint, releaseID uint) (bool, error) {
	release, err := s.repo.GetModuleRelease(modID, releaseID)
	if err != nil {
		return false, err
	}

	err = s.cleanReleaseFolder(release)
	if err != nil {
		return false, err
	}
	return true, nil

}

func (s *ReleaseService) DeleteModuleRelease(userID uint, modID uint, releaseID uint) (bool, error) {
	if s.userService.UserHasPermission(userID, models.DeleteModules) {
		res, err := s.repo.DeleteModuleRelease(userID, modID, releaseID)
		if res && err != nil {
			//lets remove the release folder
			return s.removeModuleReleaseFolder(modID, releaseID)
		}
		s.logger.Errorf("release service", "couldn't delete Module Release folder for mod: %d and release: %d", modID, releaseID)
		return res, err
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

	res, err := s.repo.DeleteModuleRelease(userID, modID, releaseID)
	if res && err != nil {
		//lets remove the release folder
		return s.removeModuleReleaseFolder(modID, releaseID)
	}
	s.logger.Errorf("release service", "couldn't delete user's %d Module Release folder for mod: %d and release: %d", userID, modID, releaseID)
	return res, err
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
	err = s.repo.Create(newRelease)
	if err != nil {
		return 0, err
	}

	stat, err := s.statRepo.Create(newRelease)
	if err != nil {
		return 0, err
	}

	newRelease.Statistics = stat
	err = s.repo.Update(newRelease)
	if err != nil {
		return 0, err
	}

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

func (s *ReleaseService) GetModuleReleaseIds(id uint) ([]uint, error) {
	results, error := s.repo.GetModuleReleasesIDs(id)
	if error != nil {
		return nil, error
	}
	if results == nil {
		return nil, nil
	}
	return results, nil
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
	//need to create the release in the file system
	err = s.createReleaseFolder(rel)
	if err != nil {
		return 0, err
	}
	return releaseID, err
}

func (s *ReleaseService) createReleaseFolder(rel *models.ModuleRelease) error {
	mod, err := s.moduleService.GetModule(rel.ModuleID)
	if err != nil {
		s.logger.Errorf("release service", "couldnt find module with module ID %d", rel.ModuleID)
		return err
	}
	dmo, err := s.ownershipService.FindDeveloperModuleOwner(&mod.Owner)
	if err != nil {
		s.logger.Errorf("release service", "couldnt find module owner with module owner ID %d", mod.OwnerID)
		return err
	}

	dest := s.moduleService.GetModuleReleasePath(mod, rel)
	err = utils.CreateFolderPath(dest)
	if err != nil {
		s.logger.Errorf("release service", "couldnt create the module release path %s", dest)
		return err
	}

	modPath := s.moduleService.GetModulePath(dmo, mod)
	err = utils.CopyDir(modPath, dest)
	if err != nil {
		s.logger.Errorf("release service", "couldnt copy the module directory to the release directory %s", dest)
		return err
	}

	dest = path.Join(dest, s.defaultDistFilename)
	err = utils.CreateTarGz(modPath, dest)
	if err != nil {
		s.logger.Errorf("release service", "couldnt create tar ball of the module directory and copy to the release directory %s", dest)
		return err
	}

	return nil
}

func (s *ReleaseService) cleanReleaseFolder(rel *models.ModuleRelease) error {
	mod, err := s.moduleService.GetModule(rel.ModuleID)
	if err != nil {
		s.logger.Errorf("release service", "couldnt find module with module ID %d", rel.ModuleID)
		return err
	}

	dest := s.moduleService.GetModuleReleasePath(mod, rel)
	if utils.FolderExists(dest) {
		return utils.CleanFolder(dest)
	}

	return nil
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
	//need to clean the release from the file system
	err = s.cleanReleaseFolder(rel)
	if err != nil {
		s.logger.Errorf("release service", "couldnt clean the release folder with release ID %d", rel.ID)
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

	//need to clean the release from the file system if there
	err = s.cleanReleaseFolder(rel)
	if err != nil {
		s.logger.Errorf("release service", "couldnt clean the release folder with release ID %d", rel.ID)
		return 0, err
	}

	return releaseID, nil
}

func (s *ReleaseService) Initialize() {
	s.statusRepo.Initialize()
}
