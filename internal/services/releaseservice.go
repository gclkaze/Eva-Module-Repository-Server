package services

import (
	"github.com/gclkaze/evamodulerepositoryserver/internal/dto"
	"github.com/gclkaze/evamodulerepositoryserver/internal/repositories"
)

type ReleaseService struct {
	repo repositories.ReleaseRepository
}

func NewReleaseService(repo repositories.ReleaseRepository) *ReleaseService {
	return &ReleaseService{repo: repo}
}

func (s *ReleaseService) FindByID(id uint) (*dto.ReleaseDTO, error) {
	/*result, error := s.repo.FindByID(id)
	if error != nil {
		return nil, error
	}
	if result == nil {
		return nil, nil
	}
	res := dto.NewModuleDTO(*result)
	return res, error*/
	return nil, nil
}

func (s *ReleaseService) GetModuleRelease(id uint) (*dto.ReleaseDTO, error) {
	/*result, error := s.repo.FindByID(id)
	if error != nil {
		return nil, error
	}
	if result == nil {
		return nil, nil
	}
	res := dto.NewModuleDTO(*result)
	return res, error*/
	return nil, nil
}

func (s *ReleaseService) GetModuleReleases(id uint) (*dto.ReleaseDTO, error) {
	/*result, error := s.repo.FindByID(id)
	if error != nil {
		return nil, error
	}
	if result == nil {
		return nil, nil
	}
	res := dto.NewModuleDTO(*result)
	return res, error*/
	return nil, nil
}

func (s *ReleaseService) SearchByKeywords(tags []string) ([]dto.ReleaseDTO, error) {
	/*	results, error := s.repo.SearchByKeywords(tags)
		if error != nil {
			return nil, error
		}
		var dtos []dto.ModuleDTO
		for i := 0; i < len(results); i++ {
			dtos = append(dtos, *dto.NewModuleDTO(results[i]))
		}
		return dtos, error*/
	return nil, nil
}
