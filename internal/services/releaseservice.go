package services

import (
	"github.com/gclkaze/evamodulerepositoryserver/internal/dto"
	"github.com/gclkaze/evamodulerepositoryserver/internal/repositories"
)

type ReleaseService struct {
	repo       repositories.ReleaseRepository
	statusRepo repositories.ReleaseStatusRepository
}

func NewReleaseService(repo repositories.ReleaseRepository, statusRepo repositories.ReleaseStatusRepository) *ReleaseService {
	return &ReleaseService{repo: repo, statusRepo: statusRepo}
}

func (s *ReleaseService) DeleteModuleRelease(id uint, releaseID uint) (bool, error) {
	return s.repo.DeleteModuleRelease(id, releaseID)
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
