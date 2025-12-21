// Package repositories contains
package repositories

import (
	"errors"
	"fmt"

	"github.com/gclkaze/evamodulerepositoryserver/internal/models"
	"github.com/gclkaze/evamodulerepositoryserver/pkg/utils"
	"gorm.io/gorm"
)

type ModuleRepository interface {
	Create(dev *models.Module) error
	CreateTx(tx *gorm.DB, dev *models.Module) error
	Update(mod *models.Module) error
	FindByID(id uint, preload bool) (*models.Module, error)
	SearchByKeywords(tags []string) ([]models.Module, error)
	Delete(id uint) (bool, error)
	GetDB() *gorm.DB
	GetMaxID() (uint, error)
}

type moduleRepository struct {
	db *gorm.DB
}

func NewModuleRepository(db *gorm.DB) ModuleRepository {
	return &moduleRepository{db: db}
}

func (r *moduleRepository) GetDB() *gorm.DB {
	return r.db
}
func (r *moduleRepository) Create(mod *models.Module) error {
	return r.db.Create(mod).Error
}

func (r *moduleRepository) CreateTx(tx *gorm.DB, mod *models.Module) error {
	return tx.Create(mod).Error
}

func (r moduleRepository) FindByID(id uint, preload bool) (*models.Module, error) {
	var m models.Module
	var err error
	if preload {
		err = r.db.Preload("Owner").Preload("Keywords").Preload("Owner.Type").Preload("Owner.Type").First(&m, id).Error
	} else {
		err = r.db.First(&m, id).Error
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &m, nil
}

func (r moduleRepository) Delete(id uint) (bool, error) {
	var m models.Module
	res := r.db.Where(id).Delete(&m)

	if res.Error != nil {
		return false, res.Error
	}

	if res.RowsAffected == 0 {
		return false, nil
	}

	return true, nil
}

func (r moduleRepository) GetMaxID() (uint, error) {
	var maxID uint
	err := r.db.
		Model(&models.Module{}).
		Select("COALESCE(MAX(id), 0)").
		Scan(&maxID).Error

	if err != nil {
		return 0, err
	}
	return maxID, err
}

func (r *moduleRepository) Update(mod *models.Module) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(mod).Error; err != nil {
			return err
		}
		if err := tx.Model(mod).Association("Keywords").Replace(mod.Keywords); err != nil {
			return err
		}
		return nil
	})
	return err
}

func (r moduleRepository) SearchByKeywords(tags []string) ([]models.Module, error) {
	var results []models.Module

	whereTitleClause := utils.BuildWhereConditionStringForUniqueAttrsContaining("title", tags)
	whereReprClause := utils.BuildWhereConditionStringForUniqueAttrsContaining("repr", tags)
	whereDescriptionClause := utils.BuildWhereConditionStringForUniqueAttrsContaining("description", tags)

	r.db.
		Where(whereTitleClause + " OR " + whereReprClause + " OR " + whereDescriptionClause).
		Find(&results)

	whereKeywordsClause := utils.BuildWhereConditionStringForUniqueAttrsContaining("keywords.label", tags)
	selection := "modules.id, modules.repr, modules.title, modules.description "
	var taggedResults []models.Module

	if results != nil {
		if len(results) != 0 {
			var ids []uint
			for i := 0; i < len(results); i++ {
				ids = append(ids, results[i].ID)
			}
			q := fmt.Sprintf("SELECT %s FROM modules LEFT JOIN module_keywords ON modules.id = module_keywords.module_id LEFT JOIN keywords ON keywords.id = module_keywords.keyword_id WHERE ( %s ) AND modules.id NOT IN ? ", selection, whereKeywordsClause)
			r.db.Raw(q, ids).Scan(&taggedResults)

		} else {
			q := fmt.Sprintf("SELECT %s FROM modules LEFT JOIN module_keywords ON modules.id = module_keywords.module_id LEFT JOIN keywords ON keywords.id = module_keywords.keyword_id WHERE ( %s )", selection, whereKeywordsClause)
			r.db.Raw(q).Scan(&taggedResults)
		}
	}

	results = append(results, taggedResults...)
	return results, nil
}
