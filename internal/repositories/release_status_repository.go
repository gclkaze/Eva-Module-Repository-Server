package repositories

import (
	"github.com/gclkaze/evamodulerepositoryserver/internal/models"
	"gorm.io/gorm"
)

type ReleaseStatusRepository interface {
	Initialize() error
}

type releaseStatusRepository struct {
	db *gorm.DB
}

func NewReleaseStatusRepository(db *gorm.DB) ReleaseStatusRepository {
	return &releaseStatusRepository{db: db}
}

func (r releaseStatusRepository) Initialize() error {
	statuses := []string{"draft", "pending", "accepted", "rejected"}
	description := []string{"This is a draft of a release.", "The release is waiting to be checked by the EVA Language Team.", "The release has been accepted by the EVA Language Team.", "The release has been rejected by the EVA Language Team."}

	for i, status := range statuses {
		var count int64
		r.db.Model(&models.ModuleReleaseStatus{}).Where("label = ?", status).Count(&count)
		if count == 0 {
			if err := r.db.Create(&models.ModuleReleaseStatus{Label: status, Description: description[i]}).Error; err != nil {
				return err
			}
		}
	}
	return nil
}
