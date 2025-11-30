package dto

import "github.com/gclkaze/evamodulerepositoryserver/internal/models"

type KeywordDTO struct {
	ID    uint
	Label string
}

func NewKeywordDTO(keyword models.Keyword) *KeywordDTO {
	return &KeywordDTO{ID: keyword.ID, Label: keyword.Label}
}
