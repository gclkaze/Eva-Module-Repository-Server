package repositories

import (
	"errors"

	"github.com/gclkaze/evamodulerepositoryserver/internal/models"
	"gorm.io/gorm"
)

type ReleaseStatusTypeDef int

const (
	Draft ReleaseStatusTypeDef = iota
	Pending
	Accepted
	Rejected
)

func (t ReleaseStatusTypeDef) String() string {
	return [...]string{"draft", "pending", "accepted", "rejected"}[t]
}

type ReleaseStatusRepository interface {
	Initialize() error
	GetStatus(t ReleaseStatusTypeDef) (*models.ModuleReleaseStatus, error)
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

func (r releaseStatusRepository) GetStatus(t ReleaseStatusTypeDef) (*models.ModuleReleaseStatus, error) {
	var m models.ModuleReleaseStatus
	res := r.db.Where("label = ?", t.String())
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &m, nil
}
