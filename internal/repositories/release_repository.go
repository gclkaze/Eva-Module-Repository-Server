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

// Package repositories contains
package repositories

import (
	"errors"
	"fmt"
	"strings"

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
	GetAllModuleReleases(id uint) ([]models.ModuleRelease, error)
	GetModuleReleasesIDs(id uint) ([]uint, error)
	SearchModuleReleasesByTags(id uint, tags []string) ([]models.ModuleRelease, error)
	GetModuleReleasesWithStatus(id uint, statusID uint) ([]models.ModuleRelease, error)
	GetModuleReleaseWithStatus(modID uint, releaseID uint, statusID uint) (*models.ModuleRelease, error)
	GetMaxID() (uint, error)
	GetCount() (int64, error)
	FindByModuleIDAndVersionExceptOne(ID uint, modID uint, version string) (*models.ModuleRelease, error)
	GetModuleReleaseByVersion(id uint, version string) (*models.ModuleRelease, error)
	GetModuleReleaseByVersionAndStatus(id uint, version string, stID uint) (*models.ModuleRelease, error)
	GetReleaseCountForModule(moduleID uint, statusID uint) (int64, error)
	GetLastModuleRelease(id uint, stID uint) (*models.ModuleRelease, error)
	GetLastModuleStatusIndependentRelease(id uint) (*models.ModuleRelease, error)

	GetModuleReleasesByFilter(p *models.ReleaseFilterParams) ([]models.ModuleRelease, error)
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

func (r releaseRepository) GetReleaseCountForModule(moduleID uint, statusID uint) (int64, error) {
	var count int64
	err := r.db.
		Model(&models.ModuleRelease{}).
		Where("module_id = ? AND status_id = ?", moduleID, statusID).
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
	err := r.db.Preload("Status").Where("module_id = ? AND released_at IS NOT NULL", id).Order("released_at DESC").Find(&results).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (r releaseRepository) GetAllModuleReleases(id uint) ([]models.ModuleRelease, error) {
	var results []models.ModuleRelease
	err := r.db.Preload("Status").Preload("Keywords").Where("module_id = ?", id).Order("created_at DESC").Find(&results).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (r releaseRepository) GetModuleReleaseByVersion(id uint, version string) (*models.ModuleRelease, error) {
	var result models.ModuleRelease
	version = strings.TrimSpace(version)
	versionWithoutV := version
	if version[0] == 'v' {
		versionWithoutV = version[1:]
	}
	err := r.db.Preload("Status").Where("module_id = ? AND ( version = ? OR version = ? )", id, version, versionWithoutV).First(&result).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r releaseRepository) GetModuleReleaseByVersionAndStatus(id uint, version string, stID uint) (*models.ModuleRelease, error) {
	var result models.ModuleRelease
	version = strings.TrimSpace(version)
	versionWithoutV := version
	if version[0] == 'v' {
		versionWithoutV = version[1:]
	}
	err := r.db.Where("module_id = ? AND ( version = ? OR version = ? ) AND status_id = ?", id, version, versionWithoutV, stID).First(&result).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r releaseRepository) GetLastModuleRelease(id uint, stID uint) (*models.ModuleRelease, error) {
	var result models.ModuleRelease
	err := r.db.Preload("Status").Where("module_id = ? AND status_id = ? AND released_at IS NOT NULL", id, stID).
		Order("released_at DESC").
		First(&result).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r releaseRepository) GetLastModuleStatusIndependentRelease(id uint) (*models.ModuleRelease, error) {
	var result models.ModuleRelease
	err := r.db.Preload("Status").Where("module_id = ? AND released_at IS NOT NULL", id).
		Order("released_at DESC").
		First(&result).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &result, nil
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

func (r releaseRepository) FindByModuleIDAndVersionExceptOne(ID uint, modID uint, version string) (*models.ModuleRelease, error) {
	var result models.ModuleRelease
	err := r.db.Where("id != ? AND module_id = ? AND version = ?", ID, modID, version).First(&result).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r releaseRepository) GetModuleReleasesWithStatus(id uint, statusID uint) ([]models.ModuleRelease, error) {
	var results []models.ModuleRelease
	err := r.db.Where("module_id = ? AND status_id = ?", id, statusID).Order("released_at DESC").Find(&results).Error
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

	q := fmt.Sprintf("SELECT %s FROM module_releases LEFT JOIN release_keywords ON module_releases.id = release_keywords.module_release_id LEFT JOIN keywords ON keywords.id = release_keywords.keyword_id WHERE ( %s ) AND module_releases.module_id = ? AND module_releases.released_at IS NOT NULL ORDER BY released_at DESC", selection, whereKeywordsClause)

	err := r.db.Raw(q, id).Scan(&results).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (r *releaseRepository) GetModuleReleasesByFilter(p *models.ReleaseFilterParams) ([]models.ModuleRelease, error) {
	var results []models.ModuleRelease

	q := r.db.Model(&models.ModuleRelease{}).
		Preload("Status").
		Preload("Keywords").
		Preload("Creator.UserAccount").
		Preload("Statistics").
		Joins("JOIN modules m ON m.id = module_releases.module_id").
		Joins("JOIN developers d ON d.id = module_releases.creator_id").
		Joins("JOIN user_accounts ua ON ua.id = d.user_id")

	if p == nil || p.IsEmpty() {
		err := q.Order("module_releases.created_at DESC").Find(&results).Error
		return results, err
	}

	// Status (by label)
	if len(p.Status) > 0 {
		q = q.Joins("JOIN module_release_statuses s ON s.id = module_releases.status_id").
			Where("s.label IN ?", p.Status)
	}

	// Versions (exact)
	if len(p.Versions) > 0 {
		q = q.Where("module_releases.version IN ?", p.Versions)
	}

	// Tags (keywords) - matches ANY of the tags
	/*	if len(p.Tags) > 0 {
		q = q.Joins("JOIN release_keywords rk ON rk.module_release_id = module_releases.id").
			Joins("JOIN keywords k ON k.id = rk.keyword_id").
			Where("k.label IN ?", p.Tags).
			Distinct("module_releases.*")
	}*/

	if len(p.Tags) > 0 {
		// Match if release has at least one tag OR module has at least one tag
		q = q.Where(`
		EXISTS (
			SELECT 1
			FROM release_keywords rk
			JOIN keywords k ON k.id = rk.keyword_id
			WHERE rk.module_release_id = module_releases.id
			  AND k.label IN ?
		)
		OR EXISTS (
			SELECT 1
			FROM module_keywords mk
			JOIN keywords k2 ON k2.id = mk.keyword_id
			WHERE mk.module_id = m.id
			  AND k2.label IN ?
		)
	`, p.Tags, p.Tags)
	}

	// ModuleName: match module repr OR title (LIKE, OR across terms)
	if len(p.ModuleName) > 0 {
		sub := r.db
		for i, term := range p.ModuleName {
			like := "%" + term + "%"
			if i == 0 {
				sub = sub.Where("(m.repr LIKE ? OR m.title LIKE ?)", like, like)
			} else {
				sub = sub.Or("(m.repr LIKE ? OR m.title LIKE ?)", like, like)
			}
		}
		q = q.Where(sub)
	}

	// RepoName: match module repo_name (LIKE, OR across terms)
	if len(p.RepoName) > 0 {
		sub := r.db
		for i, term := range p.RepoName {
			like := "%" + term + "%"
			if i == 0 {
				sub = sub.Where("m.repo_name LIKE ?", like)
			} else {
				sub = sub.Or("m.repo_name LIKE ?", like)
			}
		}
		q = q.Where(sub)
	}

	// CreatedAfter (gorm.Model.CreatedAt)
	if p.CreatedAfter != nil {
		q = q.Where("module_releases.created_at > ?", *p.CreatedAfter)
	}

	// ReleasedAfter (ReleasedAt is nullable)
	if !p.ReleasedAfter.IsZero() {
		q = q.Where("module_releases.released_at IS NOT NULL").
			Where("module_releases.released_at > ?", p.ReleasedAfter)
	}

	// Description (LIKE, OR across terms)
	if len(p.Description) > 0 {
		sub := r.db
		for i, term := range p.Description {
			like := "%" + term + "%"
			if i == 0 {
				sub = sub.Where("module_releases.description LIKE ?", like)
			} else {
				sub = sub.Or("module_releases.description LIKE ?", like)
			}
		}
		q = q.Where(sub)
	}

	// Creator (LIKE against handle/first/last, OR across terms)
	if len(p.Creator) > 0 {
		sub := r.db
		for i, term := range p.Creator {
			like := "%" + term + "%"
			if i == 0 {
				sub = sub.Where(
					"(d.handle LIKE ? OR d.first_name LIKE ? OR d.last_name LIKE ?)",
					like, like, like,
				)
			} else {
				sub = sub.Or(
					"(d.handle LIKE ? OR d.first_name LIKE ? OR d.last_name LIKE ?)",
					like, like, like,
				)
			}
		}
		q = q.Where(sub)
	}

	// CreatorEmail (LIKE, OR across terms)
	if len(p.CreatorEmail) > 0 {
		sub := r.db
		for i, term := range p.CreatorEmail {
			like := "%" + term + "%"
			if i == 0 {
				sub = sub.Where("ua.email LIKE ?", like)
			} else {
				sub = sub.Or("ua.email LIKE ?", like)
			}
		}
		q = q.Where(sub)
	}

	err := q.Order("module_releases.created_at DESC").Find(&results).Error

	/*	query := q.Statement.SQL.String()
		fmt.Print(query)
		fmt.Println(q.Statement.Vars)*/

	return results, err
}
