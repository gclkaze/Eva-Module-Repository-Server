// Copyright (c) 2025 Michail Dorgiakis - gclkaze@gmail.com - https://github.com/gclkaze
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package services

import (
	"github.com/gclkaze/evamodulerepositoryserver/internal/models"
	"github.com/gclkaze/evamodulerepositoryserver/internal/repositories"
	"gorm.io/gorm"
)

type ModuleOwnershipService struct {
	moduleOwnerRepo     repositories.ModuleOwnerRepository
	moduleOwnerTypeRepo repositories.ModuleOwnerTypesRepository
	devModuleOwnerRepo  repositories.DeveloperModuleOwnerRepository
}

func NewModuleOwnershipService(moduleOwnerRepo repositories.ModuleOwnerRepository,
	moduleOwnerTypeRepo repositories.ModuleOwnerTypesRepository,
	devModuleOwnerRepo repositories.DeveloperModuleOwnerRepository) *ModuleOwnershipService {
	return &ModuleOwnershipService{moduleOwnerRepo: moduleOwnerRepo, moduleOwnerTypeRepo: moduleOwnerTypeRepo, devModuleOwnerRepo: devModuleOwnerRepo}
}

func (s *ModuleOwnershipService) CreateModuleOwner(t models.ModuleOwnerTypeDef, entityID uint) (*models.ModuleOwner, error) {
	typ, err := s.moduleOwnerTypeRepo.FindByLabel(t.String())
	if err != nil {
		return nil, err
	}
	return s.moduleOwnerRepo.Create(typ, entityID)
}

func (s *ModuleOwnershipService) GetModuleOwner(t models.ModuleOwnerTypeDef, entityID uint) (*models.ModuleOwner, error) {
	typ, err := s.moduleOwnerTypeRepo.FindByLabel(t.String())
	if err != nil {
		return nil, err
	}
	return s.moduleOwnerRepo.Find(typ, entityID)
}

func (s *ModuleOwnershipService) CreateModuleOwnerTx(tx *gorm.DB, t models.ModuleOwnerTypeDef, entityID uint) (*models.ModuleOwner, error) {
	typ, err := s.moduleOwnerTypeRepo.FindByLabelTx(tx, t.String())
	if err != nil {
		return nil, err
	}
	return s.moduleOwnerRepo.Create(typ, entityID)
}

func (s *ModuleOwnershipService) CreateDeveloperModuleOwner(d *models.Developer, mo *models.ModuleOwner) (*models.DeveloperModuleOwner, error) {
	return s.devModuleOwnerRepo.Create(d, mo)
}

func (s *ModuleOwnershipService) CreateDeveloperModuleOwnerTx(tx *gorm.DB, d *models.Developer, mo *models.ModuleOwner) (*models.DeveloperModuleOwner, error) {
	return s.devModuleOwnerRepo.CreateTx(tx, d, mo)
}

func (s ModuleOwnershipService) FindDeveloperModuleOwner(mo *models.ModuleOwner) (*models.DeveloperModuleOwner, error) {
	if mo.Type.Label != models.Dev.String() {
		return nil, nil
	}
	res, err := s.devModuleOwnerRepo.FindByDevAndModOwner(mo)
	if err != nil {
		return nil, nil
	}
	return res, nil
}

func (s ModuleOwnershipService) Delete(id uint) (bool, error) {
	return s.devModuleOwnerRepo.Delete(id)
}
