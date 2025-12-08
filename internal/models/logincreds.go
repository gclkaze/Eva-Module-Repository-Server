package models

type LoginRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	Handle    string `json:"handle"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type LoginResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refreshToken"`
}
