package repositories

import (
	"github.com/gclkaze/evamodulerepositoryserver/internal/models"
	"gorm.io/gorm"
)

type ReleaseStatisticsRepository interface {
	Create(t *models.ModuleRelease) (*models.ReleaseStatistics, error)
}

type releaseStatisticsRepository struct {
	db *gorm.DB
}

func NewReleaseStatisticsRepository(db *gorm.DB) ReleaseStatisticsRepository {
	return &releaseStatisticsRepository{db: db}
}

func (m *releaseStatisticsRepository) Create(rel *models.ModuleRelease) (*models.ReleaseStatistics, error) {
	stat := models.NewReleaseStatistics(&rel.ID, 0)
	res := m.db.Create(stat)
	if res.Error != nil {
		return nil, res.Error
	}
	return stat, nil
}
