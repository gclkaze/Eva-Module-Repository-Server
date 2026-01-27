package dto

import "github.com/gclkaze/evamodulerepositoryserver/internal/models"

type ModuleEnrichedDTO struct {
	ID          uint   `json:"id" binding:"required"`
	Title       string `json:"title" binding:"required"`
	Repr        string `json:"repr" binding:"required"`
	Description string `json:"description" binding:"required"`
	//OwnerName   string       `json:"ownerName" binding:"required"`
	RepoName    string       `json:"repoName" binding:"required"`
	Tags        []string     `json:"tags" binding:"required"`
	ReleaseInfo []ReleaseDTO `json:"releases_info" binding:"required"`
}

func NewModuleEnrichedDTO(module *models.Module, moduleReleases []models.ModuleRelease) *ModuleEnrichedDTO {
	var tags []string
	var releaseInfo []ReleaseDTO
	for i := range module.Keywords {
		tags = append(tags, module.Keywords[i].Label)
	}
	for i := range moduleReleases {
		releaseInfo = append(releaseInfo, *NewReleaseDTO(moduleReleases[i]))
	}
	return &ModuleEnrichedDTO{ID: module.ID, Title: module.Title, Repr: module.Repr, Description: module.Description, Tags: tags, RepoName: module.RepoName, ReleaseInfo: releaseInfo}
}

func NewModuleEnrichedDTOWithReleaseDTO(module *models.Module, moduleReleases []ReleaseDTO) *ModuleEnrichedDTO {
	var tags []string
	for i := range module.Keywords {
		tags = append(tags, module.Keywords[i].Label)
	}

	return &ModuleEnrichedDTO{ID: module.ID, Title: module.Title, Repr: module.Repr, Description: module.Description, Tags: tags, RepoName: module.RepoName, ReleaseInfo: moduleReleases}
}
