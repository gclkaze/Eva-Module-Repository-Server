package models

import (
	"gorm.io/gorm"
)

type ReleaseStatistics struct {
	gorm.Model
	ReleaseID     *uint  `gorm:"uniqueIndex"`
	DownloadCount uint64 `gorm:"not null;default:0"`
}

func NewReleaseStatistics(releaseID *uint, cnt uint64) *ReleaseStatistics {
	return &ReleaseStatistics{ReleaseID: releaseID, DownloadCount: cnt}
}

/*

db.Transaction(func(tx *gorm.DB) error {
	if err := tx.Model(&File{}).
		Where("id = ?", fileID).
		UpdateColumn("download_count", gorm.Expr("download_count + 1")).Error; err != nil {
		return err
	}

	// other DB changes here

	return nil
})
*/
