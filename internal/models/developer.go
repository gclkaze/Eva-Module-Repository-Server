package models

import (
	"gorm.io/gorm"
)

type Developer struct {
	gorm.Model
	Handle      string      `json:"handle"`
	FirstName   string      `json:"first_name"`
	LastName    string      `json:"last_name"`
	UserID      uint        `json:"user_id"`
	UserAccount UserAccount `gorm:"foreignKey:UserID"`
	Active      bool        `json:"is_active"`
}

func NewDeveloper(handle string, firstName string, lastName string, userID uint, userAccount UserAccount, active bool) *Developer {
	return &Developer{Handle: handle, FirstName: firstName, LastName: lastName, UserID: userID, UserAccount: userAccount, Active: active}
}
