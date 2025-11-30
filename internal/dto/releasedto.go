package dto

import (
	"time"

	"github.com/gclkaze/evamodulerepositoryserver/internal/models"
)

type ReleaseDTO struct {
	ID          uint      `json:"id" binding:"required"`
	Version     string    `json:"title" binding:"required"`
	ReleasedAt  time.Time `json:"released_at" binding:"required"`
	Description string    `json:"description" binding:"required"`
	DiskSize    uint      `json:"diskSize" binding:"required"`

	Keywords []KeywordDTO `json:"keywords" binding:"required"`
}

func NewReleaseDTO(release models.ModuleRelease) *ReleaseDTO {
	var keywordsdtos []KeywordDTO
	for i := 0; i < len(release.Keywords); i++ {
		keywordsdtos = append(keywordsdtos, *NewKeywordDTO(release.Keywords[i]))
	}

	return &ReleaseDTO{ID: release.ID, Version: release.Version, ReleasedAt: release.ReleasedAt, Description: release.Description, Keywords: keywordsdtos, DiskSize: release.DiskSize}
}
