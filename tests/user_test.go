package tests

import "testing"

// register tests
func TestUserRegistrationNoPayload(t *testing.T) {
}

func TestUserRegistrationEmailNotCorrectlyFormed(t *testing.T) {
}

func TestUserRegistrationNoPassword(t *testing.T) {
}

func TestRegistedWithBannedUser(t *testing.T) {

}

func TestUnbanUserThenRegister(t *testing.T) {

}

func TestRegisterNewUser(t *testing.T) {

}

func TestKnownUserAttemptsToRegister(t *testing.T) {

}

// login tests
func TestUserLoginNoPayload(t *testing.T) {
}

func TestUserLoginEmailNotCorrectlyFormed(t *testing.T) {
}

func TestUserLoginNoPassword(t *testing.T) {
}
func TestBanUserWithConsecutiveLogin(t *testing.T) {

}

func TestUnbanUserWithLogin(t *testing.T) {

}
func TestLoginIncorrectPassword(t *testing.T) {

}

func TestUnknownUserAttemptsToLogin(t *testing.T) {

}

//refresh

func TestUserRefreshNoPayload(t *testing.T) {
}
func TestUserRefreshInvalidRefreshToken(t *testing.T) {
}
func TestUserRefreshEmptyRefreshToken(t *testing.T) {
}

func TestUserRefreshBannedUser(t *testing.T) {
}
func TestUserRefreshSuccessfull(t *testing.T) {
}
func TestUnknownUserAttemptsToRefreshToken(t *testing.T) {

}
