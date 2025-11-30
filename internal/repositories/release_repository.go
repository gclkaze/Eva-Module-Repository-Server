// Package repositories contains
package repositories

import (
	"errors"
	"fmt"

	"github.com/gclkaze/evamodulerepositoryserver/internal/models"
	"github.com/gclkaze/evamodulerepositoryserver/pkg/utils"
	"gorm.io/gorm"
)

type ReleaseRepository interface {
	Create(dev *models.ModuleRelease) error
	FindByID(id uint) (*models.ModuleRelease, error)
	GetModuleRelease(id uint, releaseID uint) (*models.ModuleRelease, error)
	GetModuleReleases(id uint) ([]models.ModuleRelease, error)
	SearchModuleReleasesByTags(id uint, tags []string) ([]models.ModuleRelease, error)
}

type releaseRepository struct {
	db *gorm.DB
}

func NewReleaseRepository(db *gorm.DB) ReleaseRepository {
	return &releaseRepository{db: db}
}

func (r *releaseRepository) Create(mod *models.ModuleRelease) error {
	return r.db.Create(mod).Error
}

func (r releaseRepository) FindByID(id uint) (*models.ModuleRelease, error) {
	var m models.ModuleRelease
	err := r.db.First(&m, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &m, nil
}

func (r releaseRepository) GetModuleReleases(id uint) ([]models.ModuleRelease, error) {
	var results []models.ModuleRelease
	err := r.db.Where("module_id = ?", id).Find(&results).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (r releaseRepository) GetModuleRelease(id uint, releaseID uint) (*models.ModuleRelease, error) {
	var result models.ModuleRelease
	err := r.db.Where("module_id = ? AND id = ?", id, releaseID).First(&result).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r releaseRepository) SearchModuleReleasesByTags(id uint, tags []string) ([]models.ModuleRelease, error) {
	var results []models.ModuleRelease

	whereKeywordsClause := utils.BuildWhereConditionStringForUniqueAttrsContaining("keywords.label", tags)
	selection := "module_releases.id, module_releases.version, module_releases.released_at, module_releases.description,module_releases.disk_size "

	q := fmt.Sprintf("SELECT %s FROM module_releases LEFT JOIN release_keywords ON module_releases.id = release_keywords.module_release_id LEFT JOIN keywords ON keywords.id = release_keywords.keyword_id WHERE ( %s ) AND module_releases.module_id = ? ", selection, whereKeywordsClause)

	err := r.db.Raw(q, id).Scan(&results).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return results, nil
}
