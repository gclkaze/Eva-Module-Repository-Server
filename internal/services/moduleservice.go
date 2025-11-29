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

func (s *ModuleService) SearchByKeywords(tags []string) ([]dto.ModuleDTO, error) {
	results, error := s.repo.SearchByKeywords(tags)
	return results, error
}
