package dto

import "github.com/gclkaze/evamodulerepositoryserver/internal/models"

type ModuleDTO struct {
	ID          uint   `json:"id" binding:"required"`
	Title       string `json:"title" binding:"required"`
	Repr        string `json:"repr" binding:"required"`
	Description string `json:"description" binding:"required"`
}

func NewModuleDTO(mod models.Module) *ModuleDTO {
	return &ModuleDTO{ID: mod.ID, Title: mod.Title, Repr: mod.Repr, Description: mod.Description}
}
