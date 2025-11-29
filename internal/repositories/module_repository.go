// Package repositories contains
package repositories

import (
	"errors"

	"github.com/gclkaze/evamodulerepositoryserver/internal/dto"
	"github.com/gclkaze/evamodulerepositoryserver/internal/models"
	"github.com/gclkaze/evamodulerepositoryserver/pkg/utils"
	"gorm.io/gorm"
)

type ModuleRepository interface {
	Create(dev *models.Module) error
	FindByID(id uint) (*models.Module, error)
	SearchByKeywords(tags []string) ([]dto.ModuleDTO, error)
}

type moduleRepository struct {
	db *gorm.DB
}

func NewModuleRepository(db *gorm.DB) ModuleRepository {
	return &moduleRepository{db: db}
}

func (r *moduleRepository) Create(mod *models.Module) error {
	return r.db.Create(mod).Error
}

func (r moduleRepository) FindByID(id uint) (*models.Module, error) {
	var m models.Module
	err := r.db.First(&m, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &m, nil
}

func (r moduleRepository) SearchByKeywords(tags []string) ([]dto.ModuleDTO, error) {
	var m []dto.ModuleDTO
	var results []models.Module

	whereTitleClause := utils.BuildWhereConditionStringForUniqueAttrsContaining("title", tags)
	whereReprClause := utils.BuildWhereConditionStringForUniqueAttrsContaining("rep", tags)
	whereDescriptionClause := utils.BuildWhereConditionStringForUniqueAttrsContaining("description", tags)

	r.db.
		Where(whereTitleClause + " OR " + whereReprClause + " OR " + whereDescriptionClause).
		Find(&results)

	whereKeywordsClause := utils.BuildWhereConditionStringForUniqueAttrsContaining("label", tags)

	var taggedResults []models.Module
	q := "SELECT * FROM modules LEFT JOIN module_keywords ON modules.id = module_keywords.module_id LEFT JOIN keywords ON keywords.id = module_keywords.keyword_id WHERE %s"
	r.db.Raw(q, whereKeywordsClause).Scan(&taggedResults)

	return m, nil
}
