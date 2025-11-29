package repositories

import (
	"errors"

	"github.com/gclkaze/evamodulerepositoryserver/internal/models"
	"gorm.io/gorm"
)

type KeywordRepository interface {
	Create(k *models.Keyword) error
	FindByID(id int64) (*models.Keyword, error)
	FindByLabel(label string) (*models.Keyword, error)
	FindAll() ([]models.Keyword, error)
	Delete(id int64) error
}

type keywordRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) KeywordRepository {
	return &keywordRepository{db: db}
}

// Create a new keyword
func (r *keywordRepository) Create(k *models.Keyword) error {
	return r.db.Create(k).Error
}

// Find by ID
func (r *keywordRepository) FindByID(id int64) (*models.Keyword, error) {
	var k models.Keyword
	err := r.db.First(&k, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &k, err
}

// Find by Label
func (r *keywordRepository) FindByLabel(label string) (*models.Keyword, error) {
	var k models.Keyword
	err := r.db.Where("label = ?", label).First(&k).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &k, err
}

// Find all keywords
func (r *keywordRepository) FindAll() ([]models.Keyword, error) {
	var keywords []models.Keyword
	err := r.db.Find(&keywords).Error
	return keywords, err
}

// Delete by ID
func (r *keywordRepository) Delete(id int64) error {
	return r.db.Delete(&models.Keyword{}, id).Error
}
