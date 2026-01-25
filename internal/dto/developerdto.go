package dto

import "github.com/gclkaze/evamodulerepositoryserver/internal/models"

type DeveloperDTO struct {
	UserID    uint   `json:"user_id"`
	Handle    string `json:"handle"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `gorm:"email"`
	UserRole  string `gorm:"string"`
	IsBanned  bool   `gorm:"is_banned"`
}

func NewDeveloperDTO(developer *models.Developer) *DeveloperDTO {
	return &DeveloperDTO{UserID: developer.UserID, Handle: developer.Handle, FirstName: developer.FirstName, LastName: developer.LastName, Email: developer.UserAccount.Email, UserRole: developer.UserAccount.UserRole.Name, IsBanned: developer.UserAccount.IsBanned}
}
