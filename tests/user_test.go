// Copyright (c) 2025 Michail Dorgiakis - gclkaze@gmail.com - https://github.com/gclkaze
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package tests

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gclkaze/evamodulerepositoryserver/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/magiconair/properties/assert"
)

// register tests
func TestUserRegistrationNoPayload(t *testing.T) {
	w := UserRegister(nil, t, TheTestServer.GetRouter())
	assert.Equal(t, w.Code, http.StatusBadRequest)
}

func TestRegisterNewUser(t *testing.T) {

	body := models.LoginRequest{
		Email:     "newuser+" + t.Name() + "@example.com",
		Password:  "testpass",
		Handle:    "hnnn",
		FirstName: "Fn",
		LastName:  "Ln",
	}
	w := UserRegister(&body, t, TheTestServer.GetRouter())
	assert.Equal(t, w.Code, http.StatusOK)
	var rr models.RequestResult[models.LoginResponse]
	if err := json.Unmarshal(w.Body.Bytes(), &rr); err != nil {
		t.Fatalf("invalid response JSON: %v", err)
	}
	assert.Equal(t, (rr.Value.AccessToken != "" && rr.Value.RefreshToken != ""), true)
}

func TestKnownUserAttemptsToRegister(t *testing.T) {

	// default user seeded by Initialize
	body := models.LoginRequest{
		Email:     "gclkaze@gmail.com",
		Password:  "thisisapass",
		FirstName: "A first name",
		LastName:  "A Last name",
		Handle:    "hnnn22",
	}
	w := UserRegister(&body, t, TheTestServer.GetRouter())
	assert.Equal(t, w.Code, http.StatusInternalServerError)
	var rr models.ErrorResult
	if err := json.Unmarshal(w.Body.Bytes(), &rr); err != nil {
		t.Fatalf("invalid response JSON: %v", err)
	}
	assert.Equal(t, rr.Error == "Couldn't register user, the email is used", true)
}

// check for erroneous first name and last name
func TestRegisterNewUserWithErroneousFirstNameLastName(t *testing.T) {

	body := models.LoginRequest{
		Email:     "newuser+" + t.Name() + "@example.com",
		Password:  "testpass",
		Handle:    "hn",
		FirstName: "asd ?? asd '!2",
		LastName:  "alastname",
	}
	w := UserRegister(&body, t, TheTestServer.GetRouter())
	assert.Equal(t, w.Code, http.StatusBadRequest)
	var rr models.ErrorResult
	if err := json.Unmarshal(w.Body.Bytes(), &rr); err != nil {
		t.Fatalf("invalid response JSON: %v", err)
	}
	assert.Equal(t, rr.Details == "Invalid first name provided", true)
}

func TestRegisterNewUserWithErroneousLastName(t *testing.T) {

	body := models.LoginRequest{
		Email:     "newuser+" + t.Name() + "@example.com",
		Password:  "testpass",
		Handle:    "hn",
		FirstName: "george",
		LastName:  "Ln asd ?? asd '!2",
	}
	w := UserRegister(&body, t, TheTestServer.GetRouter())
	assert.Equal(t, w.Code, http.StatusBadRequest)
	var rr models.ErrorResult
	if err := json.Unmarshal(w.Body.Bytes(), &rr); err != nil {
		t.Fatalf("invalid response JSON: %v", err)
	}
	assert.Equal(t, rr.Details == "Invalid last name provided", true)
}

func TestRegisterNewUserWitEmptyFirstName(t *testing.T) {

	body := models.LoginRequest{
		Email:     "newuser+" + t.Name() + "@example.com",
		Password:  "testpass",
		Handle:    "hn",
		FirstName: " ",
		LastName:  "test",
	}
	w := UserRegister(&body, t, TheTestServer.GetRouter())
	assert.Equal(t, w.Code, http.StatusBadRequest)
	var rr models.ErrorResult
	if err := json.Unmarshal(w.Body.Bytes(), &rr); err != nil {
		t.Fatalf("invalid response JSON: %v", err)
	}
	assert.Equal(t, rr.Details == "No first name provided", true)
}

func TestRegisterNewUserWithEmptyLastName(t *testing.T) {
	body := models.LoginRequest{
		Email:     "newuser+" + t.Name() + "@example.com",
		Password:  "testpass",
		Handle:    "hn",
		FirstName: "george",
		LastName:  "  ",
	}
	w := UserRegister(&body, t, TheTestServer.GetRouter())
	assert.Equal(t, w.Code, http.StatusBadRequest)
	var rr models.ErrorResult
	if err := json.Unmarshal(w.Body.Bytes(), &rr); err != nil {
		t.Fatalf("invalid response JSON: %v", err)
	}
	assert.Equal(t, rr.Details == "No last name provided", true)
}

// check for handle uniqueness
func TestRegisterNewUsersWithWithDuplicateHandle(t *testing.T) {
	body := models.LoginRequest{
		Email:     "newuser1+" + t.Name() + "@example.com",
		Password:  "testpass",
		Handle:    "newuser",
		FirstName: "george",
		LastName:  "jenkins",
	}
	w1 := UserRegister(&body, t, TheTestServer.GetRouter())
	assert.Equal(t, w1.Code, http.StatusOK)

	body = models.LoginRequest{
		Email:     "newuser2+" + t.Name() + "@example.com",
		Password:  "testpass",
		Handle:    "newuser",
		FirstName: "jim",
		LastName:  "jenkins",
	}
	w2 := UserRegister(&body, t, TheTestServer.GetRouter())
	assert.Equal(t, w2.Code, http.StatusInternalServerError)
	var rr models.ErrorResult
	if err := json.Unmarshal(w2.Body.Bytes(), &rr); err != nil {
		t.Fatalf("invalid response JSON: %v", err)
	}
	assert.Equal(t, rr.Error == "Couldn't register user, the handle is used", true)
}

// register email uniqueness
func TestRegisterNewUsersWithWithDuplicateEmail(t *testing.T) {
	body := models.LoginRequest{
		Email:     t.Name() + "@example.com",
		Password:  "testpass",
		Handle:    "newuser1",
		FirstName: "george",
		LastName:  "jenkins",
	}
	w1 := UserRegister(&body, t, TheTestServer.GetRouter())
	assert.Equal(t, w1.Code, http.StatusOK)

	body = models.LoginRequest{
		Email:     t.Name() + "@example.com",
		Password:  "testpass",
		Handle:    "newuser2",
		FirstName: "jim",
		LastName:  "jenkins",
	}
	w2 := UserRegister(&body, t, TheTestServer.GetRouter())
	assert.Equal(t, w2.Code, http.StatusInternalServerError)
	var rr models.ErrorResult
	if err := json.Unmarshal(w2.Body.Bytes(), &rr); err != nil {
		t.Fatalf("invalid response JSON: %v", err)
	}
	assert.Equal(t, rr.Error == "Couldn't register user, the email is used", true)
}

// check for weird handle
func TestRegisterNewUserWithStrangeHandle(t *testing.T) {
	body := models.LoginRequest{
		Email:     "newuser+" + t.Name() + "@example.com",
		Password:  "testpass",
		Handle:    "asd 89 ! =  88sd&@asd.com",
		FirstName: "george",
		LastName:  "glooney",
	}
	w := UserRegister(&body, t, TheTestServer.GetRouter())
	assert.Equal(t, w.Code, http.StatusBadRequest)
	var rr models.ErrorResult
	if err := json.Unmarshal(w.Body.Bytes(), &rr); err != nil {
		t.Fatalf("invalid response JSON: %v", err)
	}
	assert.Equal(t, rr.Details == "Invalid handle provided", true)
}

// check for empty handle!
func TestRegisterNewUserWithEmptyHandle(t *testing.T) {
	body := models.LoginRequest{
		Email:     "newuser+" + t.Name() + "@example.com",
		Password:  "testpass",
		Handle:    " ",
		FirstName: "george",
		LastName:  "glooney",
	}
	w := UserRegister(&body, t, TheTestServer.GetRouter())
	assert.Equal(t, w.Code, http.StatusBadRequest)
	var rr models.ErrorResult
	if err := json.Unmarshal(w.Body.Bytes(), &rr); err != nil {
		t.Fatalf("invalid response JSON: %v", err)
	}
	assert.Equal(t, rr.Details == "No handle provided", true)
}

// login tests
func TestUserLoginNoPayload(t *testing.T) {
	w := UserLoginRaw("", "", t, TheTestServer.GetRouter())
	assert.Equal(t, w.Code, http.StatusBadRequest)
	var rr models.ErrorResult
	if err := json.Unmarshal(w.Body.Bytes(), &rr); err != nil {
		t.Fatalf("invalid response JSON: %v", err)
	}
	assert.Equal(t, rr.Details == "No password provided", true)
}

func TestUserLoginNoPassword(t *testing.T) {
	w := UserLoginRaw("gclkaze@gmail.com", "", t, TheTestServer.GetRouter())
	assert.Equal(t, w.Code, http.StatusBadRequest)
	var rr models.ErrorResult
	if err := json.Unmarshal(w.Body.Bytes(), &rr); err != nil {
		t.Fatalf("invalid response JSON: %v", err)
	}
	assert.Equal(t, rr.Details == "No password provided", true)
}

func TestLoginIncorrectPassword(t *testing.T) {
	w := UserLoginRaw("gclkaze@gmail.com", "wrongpass", t, TheTestServer.GetRouter())
	assert.Equal(t, w.Code, http.StatusUnauthorized)
}

func TestUnknownUserAttemptsToLogin(t *testing.T) {
	w := UserLoginRaw("unknown+"+t.Name()+"@example.com", "nopass", t, TheTestServer.GetRouter())
	assert.Equal(t, w.Code, http.StatusUnauthorized)
}

// refresh
func TestUserRefreshNoPayload(t *testing.T) {
	w := UserRefresh("", t, TheTestServer.GetRouter())
	assert.Equal(t, w.Code, http.StatusUnauthorized)
}

func TestUserRefreshInvalidRefreshToken(t *testing.T) {
	w := UserRefresh("this-is-not-a-token", t, TheTestServer.GetRouter())
	assert.Equal(t, w.Code, http.StatusUnauthorized)
}

func TestUserRefreshEmptyRefreshToken(t *testing.T) {
	w := UserRefresh("", t, TheTestServer.GetRouter())
	assert.Equal(t, w.Code, http.StatusUnauthorized)
}

func TestUserRefreshSuccessfull(t *testing.T) {
	// login to get refresh token
	lr := UserLogin(t)
	if lr == nil || lr.Value.RefreshToken == "" {
		t.Fatalf("couldn't login to obtain refresh token")
	}

	w := UserRefresh(lr.Value.RefreshToken, t, TheTestServer.GetRouter())
	assert.Equal(t, w.Code, http.StatusOK)
}

func TestUnknownUserAttemptsToRefreshToken(t *testing.T) {
	// craft a refresh token for non-existent user id
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"sub": 999999, "type": "refresh"})
	signed, _ := token.SignedString([]byte(TheTestServer.GetBackend().GetJWTSecret()))

	w := UserRefresh(signed, t, TheTestServer.GetRouter())
	assert.Equal(t, w.Code, http.StatusUnauthorized)
}

func TestUserRegistrationEmailNotCorrectlyFormed(t *testing.T) {

	body := models.LoginRequest{
		Email:     "new*!user+" + t.Name() + "@ex?amp.le.com.au",
		Password:  "testpass",
		Handle:    "hn",
		FirstName: "Fn",
		LastName:  "Ln",
	}
	w := UserRegister(&body, t, TheTestServer.GetRouter())
	assert.Equal(t, w.Code, http.StatusBadRequest)
	var rr models.ErrorResult
	if err := json.Unmarshal(w.Body.Bytes(), &rr); err != nil {
		t.Fatalf("invalid response JSON: %v", err)
	}
	assert.Equal(t, rr.Details == "Invalid email form", true)
}

func TestUserRegistrationNoPassword(t *testing.T) {
	body := models.LoginRequest{
		Email:     "newuser" + t.Name() + "@example.com",
		Password:  "   ",
		Handle:    "hn",
		FirstName: "Fn",
		LastName:  "Ln",
	}
	w := UserRegister(&body, t, TheTestServer.GetRouter())
	assert.Equal(t, w.Code, http.StatusBadRequest)
	var rr models.ErrorResult
	if err := json.Unmarshal(w.Body.Bytes(), &rr); err != nil {
		t.Fatalf("invalid response JSON: %v", err)
	}
	assert.Equal(t, rr.Details == "No password provided", true)
}

func TestUserRegistrationNoEmail(t *testing.T) {
	body := models.LoginRequest{
		Email:     "    ",
		Password:  "thisisapass",
		Handle:    "hn",
		FirstName: "Fn",
		LastName:  "Ln",
	}
	w := UserRegister(&body, t, TheTestServer.GetRouter())
	assert.Equal(t, w.Code, http.StatusBadRequest)
	var rr models.ErrorResult
	if err := json.Unmarshal(w.Body.Bytes(), &rr); err != nil {
		t.Fatalf("invalid response JSON: %v", err)
	}
	assert.Equal(t, rr.Details == "No email provided", true)
}

func TestUserLoginEmailNotCorrectlyFormed(t *testing.T) {
	body := models.LoginRequest{
		Email:    "new*!user+" + t.Name() + "@ex?amp.le.com.au",
		Password: "testpass",
	}
	w := UserLoginRaw(body.Email, body.Password, t, TheTestServer.GetRouter())
	assert.Equal(t, w.Code, http.StatusBadRequest)
	var rr models.ErrorResult
	if err := json.Unmarshal(w.Body.Bytes(), &rr); err != nil {
		t.Fatalf("invalid response JSON: %v", err)
	}
	assert.Equal(t, rr.Details == "Invalid email form", true)
}

func TestUserLoginEmptyEmail(t *testing.T) {
	body := models.LoginRequest{
		Email:    "   ",
		Password: "testpass",
	}
	w := UserLoginRaw(body.Email, body.Password, t, TheTestServer.GetRouter())
	assert.Equal(t, w.Code, http.StatusBadRequest)
	var rr models.ErrorResult
	if err := json.Unmarshal(w.Body.Bytes(), &rr); err != nil {
		t.Fatalf("invalid response JSON: %v", err)
	}
	assert.Equal(t, rr.Details == "No email provided", true)
}

func TestUserLoginEmptyPassword(t *testing.T) {
	body := models.LoginRequest{
		Email:    "user@example.com",
		Password: " ",
	}
	w := UserLoginRaw(body.Email, body.Password, t, TheTestServer.GetRouter())
	assert.Equal(t, w.Code, http.StatusBadRequest)
	var rr models.ErrorResult
	if err := json.Unmarshal(w.Body.Bytes(), &rr); err != nil {
		t.Fatalf("invalid response JSON: %v", err)
	}
	assert.Equal(t, rr.Details == "No password provided", true)
}

func TestRegisterWithBannedUser(t *testing.T) {
	t.Skip("TODO: implement - requires ban/unban flow to be defined for registration")
}

func TestUnbanUserThenRegister(t *testing.T) {
	t.Skip("TODO: implement - re-enable when supervise unban API usage is desired in tests")
}

func TestBanUserWithConsecutiveLogin(t *testing.T) {
	t.Skip("TODO: implement - no automatic ban on consecutive failed logins currently")
}

func TestUnbanUserWithLogin(t *testing.T) {
	t.Skip("TODO: implement - requires supervise unban behavior to be tested via admin endpoints")
}

func TestUserRefreshBannedUser(t *testing.T) {
	t.Skip("TODO: implement - requires banned-user behavior to be specified for refresh tokens")
}
