// Package services contains
package services

import (
	"github.com/gclkaze/evamodulerepositoryserver/internal/dto"
	"github.com/gclkaze/evamodulerepositoryserver/internal/repositories"
)

type ModuleService struct {
	repo repositories.ModuleRepository
}

func NewModuleService(repo repositories.ModuleRepository) *ModuleService {
	return &ModuleService{repo: repo}
}

func (s *ModuleService) FindByID(id uint) (*dto.ModuleDTO, error) {
	result, error := s.repo.FindByID(id)
	if error != nil {
		return nil, error
	}
	if result == nil {
		return nil, nil
	}
	res := dto.NewModuleDTO(*result)
	return res, error
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
