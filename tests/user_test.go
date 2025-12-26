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
		Handle:    "hn",
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
		Email:    "gclkaze@gmail.com",
		Password: "thisisapass",
	}
	w := UserRegister(&body, t, TheTestServer.GetRouter())
	assert.Equal(t, w.Code, http.StatusOK)
}

// login tests
func TestUserLoginNoPayload(t *testing.T) {
	w := UserLoginRaw("", "", t, TheTestServer.GetRouter())
	assert.Equal(t, w.Code, http.StatusBadRequest)
}

func TestUserLoginNoPassword(t *testing.T) {
	w := UserLoginRaw("gclkaze@gmail.com", "", t, TheTestServer.GetRouter())
	assert.Equal(t, w.Code, http.StatusUnauthorized)
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
	assert.Equal(t, w.Code, http.StatusBadRequest)
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

// The tests below were previously placeholders. Re-adding as skipped TODOs
// because behavior (auto-ban thresholds, email format validation) isn't
// enforced in the current codebase and requires a spec before implementing.
func TestUserRegistrationEmailNotCorrectlyFormed(t *testing.T) {
	t.Skip("TODO: implement once email validation rules are defined")
}

func TestUserRegistrationNoPassword(t *testing.T) {
	t.Skip("TODO: implement - define expected behavior for empty password registration")
}

func TestRegistedWithBannedUser(t *testing.T) {
	t.Skip("TODO: implement - requires ban/unban flow to be defined for registration")
}

func TestUnbanUserThenRegister(t *testing.T) {
	t.Skip("TODO: implement - re-enable when supervise unban API usage is desired in tests")
}

func TestUserLoginEmailNotCorrectlyFormed(t *testing.T) {
	t.Skip("TODO: implement once email validation rules are defined")
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
