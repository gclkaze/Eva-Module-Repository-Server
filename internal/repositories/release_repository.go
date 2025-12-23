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
	Update(mr *models.ModuleRelease) error
	FindByID(id uint) (*models.ModuleRelease, error)
	GetRelease(id uint) (*models.ModuleRelease, error)
	DeleteModuleRelease(userID uint, id uint, releaseID uint) (bool, error)
	CancelSuggestedModuleRelease(userID uint, id uint, releaseID uint) (bool, error)
	GetModuleRelease(id uint, releaseID uint) (*models.ModuleRelease, error)
	GetModuleReleases(id uint) ([]models.ModuleRelease, error)
	GetModuleReleasesIDs(id uint) ([]uint, error)
	SearchModuleReleasesByTags(id uint, tags []string) ([]models.ModuleRelease, error)
	GetModuleReleasesWithStatus(id uint, statusID uint) ([]models.ModuleRelease, error)
	GetModuleReleaseWithStatus(modID uint, releaseID uint, statusID uint) (*models.ModuleRelease, error)
	GetMaxID() (uint, error)
	GetCount() (int64, error)
}

type releaseRepository struct {
	db *gorm.DB
}

func NewReleaseRepository(db *gorm.DB) ReleaseRepository {
	return &releaseRepository{db: db}
}

func (r releaseRepository) GetMaxID() (uint, error) {
	var maxID uint
	err := r.db.
		Model(&models.ModuleRelease{}).
		Select("COALESCE(MAX(id), 0)").
		Scan(&maxID).Error

	if err != nil {
		return 0, err
	}
	return maxID, err
}

func (r releaseRepository) GetCount() (int64, error) {
	var count int64
	err := r.db.
		Model(&models.ModuleRelease{}).
		Count(&count).Error

	if err != nil {
		return 0, err
	}
	return count, err
}

func (r *releaseRepository) Create(mod *models.ModuleRelease) error {
	return r.db.Create(mod).Error
}

func (r releaseRepository) FindByID(id uint) (*models.ModuleRelease, error) {
	var m models.ModuleRelease
	err := r.db.Preload("Status").First(&m, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &m, nil
}

func (r releaseRepository) GetRelease(id uint) (*models.ModuleRelease, error) {
	var m models.ModuleRelease
	err := r.db.Preload("Status").Preload("Statistics").First(&m, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &m, nil
}

func (r *releaseRepository) Update(mr *models.ModuleRelease) error {
	return r.db.Save(mr).Error
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

func (r releaseRepository) GetModuleReleasesIDs(id uint) ([]uint, error) {
	var ids []uint
	err := r.db.Model(&models.ModuleRelease{}).Where("module_id = ?", id).Pluck("id", &ids).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return ids, nil
}

func (r releaseRepository) GetModuleReleasesWithStatus(id uint, statusID uint) ([]models.ModuleRelease, error) {
	var results []models.ModuleRelease
	err := r.db.Where("module_id = ? AND status_id = ?", id, statusID).Find(&results).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (r releaseRepository) GetModuleReleaseWithStatus(id uint, releaseID uint, statusID uint) (*models.ModuleRelease, error) {
	var result models.ModuleRelease
	err := r.db.Preload("Creator").Where("module_id = ? AND status_id = ? AND id = ?", id, statusID, releaseID).First(&result).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r releaseRepository) GetModuleRelease(id uint, releaseID uint) (*models.ModuleRelease, error) {
	var result models.ModuleRelease
	err := r.db.Preload("Creator").Preload("Status").Where("module_id = ? AND id = ?", id, releaseID).First(&result).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r releaseRepository) DeleteModuleRelease(userID uint, id uint, releaseID uint) (bool, error) {
	var result models.ModuleRelease
	//need to check if the user issued the release
	res := r.db.Where("module_id = ? AND id = ?", id, releaseID).Delete(&result)

	if res.Error != nil {
		return false, res.Error
	}
	if res.RowsAffected == 0 {
		return false, nil
	}

	return true, nil
}

func (r releaseRepository) CancelSuggestedModuleRelease(userID uint, id uint, releaseID uint) (bool, error) {
	var result models.ModuleRelease
	//need to check if the user issued the release
	res := r.db.Where("module_id = ? AND id = ?", id, releaseID).First(&result)

	if res.Error != nil {
		return false, res.Error
	}

	return true, nil
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
