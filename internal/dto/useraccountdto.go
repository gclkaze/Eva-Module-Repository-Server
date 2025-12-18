package dto

type UserAccountDTO struct {
	Handle    string `json:"handle"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	UserRole  string `json:"user_role"`
	Active    bool   `json:"is_active"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// //handle string, firstName string, lastName string, email string, password string, active bool
func NewUserAccountDTO(handle string, firstName string, lastName string, email string, password string, active bool, role string) *UserAccountDTO {
	return &UserAccountDTO{Handle: handle, Email: email, Password: password, UserRole: role, Active: active, FirstName: firstName, LastName: lastName}
}
