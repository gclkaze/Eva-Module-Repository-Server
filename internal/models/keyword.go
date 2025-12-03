package models

import "gorm.io/gorm"

type Keyword struct {
	gorm.Model
	Label string `gorm:"unique;not null" json:"label"`
}

func NewKeyword(label string) *Keyword {
	return &Keyword{Label: label}
}
