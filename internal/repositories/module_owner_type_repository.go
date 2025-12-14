package repositories

import (
	"fmt"

	"github.com/gclkaze/evamodulerepositoryserver/internal/models"
	"gorm.io/gorm"
)

type ModuleOwnerTypesRepository interface {
	Create(string) error
	FindByLabel(string) (*models.ModuleOwnerType, error)
	FindByLabelTx(*gorm.DB, string) (*models.ModuleOwnerType, error)
	FindById(uint) (*models.ModuleOwnerType, error)
	Initialize() error
}

type moduleOwnerTypesRepository struct {
	db *gorm.DB
}

func NewModuleOwnerTypesRepository(db *gorm.DB) ModuleOwnerTypesRepository {
	return &moduleOwnerTypesRepository{db: db}
}

func (m *moduleOwnerTypesRepository) Create(label string) error {
	r, err := m.FindByLabel(label)
	if err != nil {
		return err
	}
	if r != nil {
		return fmt.Errorf("label %s exists", label)
	}

	t := models.NewModuleOwnerType(label)
	res := m.db.Create(t)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (m *moduleOwnerTypesRepository) FindByLabel(label string) (*models.ModuleOwnerType, error) {
	var t models.ModuleOwnerType
	res := m.db.Where("label = ?", label).First(&t)
	if res.Error != nil {
		return nil, res.Error
	}
	return &t, nil
}

func (m *moduleOwnerTypesRepository) FindByLabelTx(tx *gorm.DB, label string) (*models.ModuleOwnerType, error) {
	var t models.ModuleOwnerType
	res := tx.Where("label = ?", label).First(&t)
	if res.Error != nil {
		return nil, res.Error
	}
	return &t, nil
}

func (m *moduleOwnerTypesRepository) FindById(id uint) (*models.ModuleOwnerType, error) {
	var t models.ModuleOwnerType
	res := m.db.First(&t, id)
	if res.Error != nil {
		return nil, res.Error
	}
	return &t, nil
}

func (m moduleOwnerTypesRepository) Initialize() error {
	for _, t := range models.GetModuleOwnerTypes() {
		var count int64
		m.db.Model(&models.ModuleOwnerType{}).Where("label = ?", t.String()).Count(&count)
		if count == 0 {
			if err := m.db.Create(models.NewModuleOwnerType(t.String())).Error; err != nil {
				return err
			}
		}
	}
	return nil
}
