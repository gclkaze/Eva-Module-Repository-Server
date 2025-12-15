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

type ErroneousRequestResult struct {
	Details string `json:"details"`
	Result  bool   `json:"result"`
	Error   string `json:"error"`
}

type RequestResult[T any] struct {
	Result  bool   `json:"result"`
	Value   T      `json:"value"`
	Message string `json:"message"`
}

func NewErroneousRequestResult(d string, err string) *ErroneousRequestResult {
	return &ErroneousRequestResult{Details: d, Result: false, Error: err}
}

/*func NewErroneousRequestResult(d string, err error) *ErroneousRequestResult {
	return &ErroneousRequestResult{Details: d, Result: false, Error: err.Error()}
}*/
