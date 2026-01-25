package dto

import "github.com/gclkaze/evamodulerepositoryserver/internal/models"

type SimpleUserAccountDTO struct {
	ID       uint   `json:"id"`
	Email    string `json:"email"`
	UserRole string `json:"user_role"`
	IsBanned bool   `json:"is_banned"`
}

func NewSimpleUserAccountDTO(u *models.UserAccount) *SimpleUserAccountDTO {
	return &SimpleUserAccountDTO{ID: u.ID, Email: u.Email, UserRole: u.UserRole.Name, IsBanned: u.IsBanned}
}
