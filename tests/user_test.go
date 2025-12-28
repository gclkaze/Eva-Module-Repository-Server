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
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/gclkaze/evamodulerepositoryserver/internal/models"
	"github.com/gclkaze/evamodulerepositoryserver/pkg/utils"
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
	//login as admin
	adminLogin := AdminLogin(t)

	email := t.Name() + "@example.com"
	pwd := "thisisapass"
	//register a new user
	body := models.LoginRequest{
		Email:     email,
		Password:  pwd,
		Handle:    strings.ToLower(t.Name()),
		FirstName: t.Name(),
		LastName:  t.Name(),
	}
	w := UserRegister(&body, t, TheTestServer.GetRouter())
	assert.Equal(t, w.Code, http.StatusOK)

	u, err := TheTestServer.GetBackend().GetUserService().FindByEmail(email)
	assert.Equal(t, err == nil, true)
	assert.Equal(t, u != nil, true)

	//admin bans user
	wBan := AdminBansUser(&adminLogin.Value, u.ID, t)
	var banResult models.RequestResult[models.EmptyRequestResult]
	if err := json.Unmarshal(wBan.Body.Bytes(), &banResult); err != nil {
		t.Fatalf("invalid response JSON: %v", err)
	}
	assert.Equal(t, banResult.Message == "User was banned successfully", true)

	//user registers again -> Invalid credentials
	wFail := UserRegister(&body, t, TheTestServer.GetRouter())
	assert.Equal(t, wFail.Code, http.StatusUnauthorized)
	var rr models.ErrorResult
	if err := json.Unmarshal(wFail.Body.Bytes(), &rr); err != nil {
		t.Fatalf("invalid response JSON: %v", err)
	}
	assert.Equal(t, rr.Error == "Invalid credentials", true)
}

func TestUnbanUserThenRegister(t *testing.T) {
	//login as admin
	adminLogin := AdminLogin(t)

	email := t.Name() + "@example.com"
	pwd := "thisisapass"
	//register a new user
	body := models.LoginRequest{
		Email:     email,
		Password:  pwd,
		Handle:    strings.ToLower(t.Name()),
		FirstName: t.Name(),
		LastName:  t.Name(),
	}
	w := UserRegister(&body, t, TheTestServer.GetRouter())
	assert.Equal(t, w.Code, http.StatusOK)

	u, err := TheTestServer.GetBackend().GetUserService().FindByEmail(email)
	assert.Equal(t, err == nil, true)
	assert.Equal(t, u != nil, true)

	//admin bans user
	wBan := AdminBansUser(&adminLogin.Value, u.ID, t)
	var banResult models.RequestResult[models.EmptyRequestResult]
	if err := json.Unmarshal(wBan.Body.Bytes(), &banResult); err != nil {
		t.Fatalf("invalid response JSON: %v", err)
	}
	assert.Equal(t, banResult.Message == "User was banned successfully", true)

	//user registers again -> Invalid credentials
	wFail := UserRegister(&body, t, TheTestServer.GetRouter())
	assert.Equal(t, wFail.Code, http.StatusUnauthorized)
	var rr models.ErrorResult
	if err := json.Unmarshal(wFail.Body.Bytes(), &rr); err != nil {
		t.Fatalf("invalid response JSON: %v", err)
	}
	assert.Equal(t, rr.Error == "Invalid credentials", true)

	//admin unbans user
	wUnBan := AdminUnBansUser(&adminLogin.Value, u.ID, t)
	var unbanResult models.RequestResult[models.EmptyRequestResult]
	if err := json.Unmarshal(wUnBan.Body.Bytes(), &unbanResult); err != nil {
		t.Fatalf("invalid response JSON: %v", err)
	}
	assert.Equal(t, unbanResult.Message == "User was unbanned successfully", true)

	//user registers again -> "Couldn't register user, the email is used"
	wSuccess := UserRegister(&body, t, TheTestServer.GetRouter())

	assert.Equal(t, wSuccess.Code, http.StatusInternalServerError)
	var errResult models.ErrorResult
	if err := json.Unmarshal(wSuccess.Body.Bytes(), &errResult); err != nil {
		t.Fatalf("invalid response JSON: %v", err)
	}
	assert.Equal(t, errResult.Error == "Couldn't register user, the email is used", true)
}

func TestBanUserWithConsecutiveLogin(t *testing.T) {
	//login as admin
	adminLogin := AdminLogin(t)

	email := t.Name() + "@example.com"
	pwd := "thisisapass"
	name := t.Name()
	//register a new user
	body := models.LoginRequest{
		Email:     email,
		Password:  pwd,
		Handle:    strings.ToLower(name)[:len(name)/2],
		FirstName: name,
		LastName:  name,
	}
	w := UserRegister(&body, t, TheTestServer.GetRouter())
	assert.Equal(t, w.Code, http.StatusOK)

	u, err := TheTestServer.GetBackend().GetUserService().FindByEmail(email)
	assert.Equal(t, err == nil, true)
	assert.Equal(t, u != nil, true)

	//admin bans user
	wBan := AdminBansUser(&adminLogin.Value, u.ID, t)
	var banResult models.RequestResult[models.EmptyRequestResult]
	if err := json.Unmarshal(wBan.Body.Bytes(), &banResult); err != nil {
		t.Fatalf("invalid response JSON: %v", err)
	}
	assert.Equal(t, banResult.Message == "User was banned successfully", true)

	//user logs in again -> Invalid credentials
	wFail := UserLoginRaw(email, pwd, t, TheTestServer.GetRouter())
	assert.Equal(t, wFail.Code, http.StatusUnauthorized)
	var rr models.ErrorResult
	if err := json.Unmarshal(wFail.Body.Bytes(), &rr); err != nil {
		t.Fatalf("invalid response JSON: %v", err)
	}
	assert.Equal(t, rr.Error == "Invalid credentials", true)
}

func TestUnbanUserWithLogin(t *testing.T) {
	//login as admin
	adminLogin := AdminLogin(t)

	email := t.Name() + "@example.com"
	pwd := "thisisapass"
	//register a new user
	body := models.LoginRequest{
		Email:     email,
		Password:  pwd,
		Handle:    strings.ToLower(t.Name()),
		FirstName: t.Name(),
		LastName:  t.Name(),
	}
	w := UserRegister(&body, t, TheTestServer.GetRouter())
	assert.Equal(t, w.Code, http.StatusOK)

	u, err := TheTestServer.GetBackend().GetUserService().FindByEmail(email)
	assert.Equal(t, err == nil, true)
	assert.Equal(t, u != nil, true)

	//admin bans user
	wBan := AdminBansUser(&adminLogin.Value, u.ID, t)
	var banResult models.RequestResult[models.EmptyRequestResult]
	if err := json.Unmarshal(wBan.Body.Bytes(), &banResult); err != nil {
		t.Fatalf("invalid response JSON: %v", err)
	}
	assert.Equal(t, banResult.Message == "User was banned successfully", true)

	//user registers again -> Invalid credentials
	wFail := UserLoginRaw(email, pwd, t, TheTestServer.GetRouter())
	assert.Equal(t, wFail.Code, http.StatusUnauthorized)
	var rr models.ErrorResult
	if err := json.Unmarshal(wFail.Body.Bytes(), &rr); err != nil {
		t.Fatalf("invalid response JSON: %v", err)
	}
	assert.Equal(t, rr.Error == "Invalid credentials", true)

	//admin unbans user
	wUnBan := AdminUnBansUser(&adminLogin.Value, u.ID, t)
	var unbanResult models.RequestResult[models.EmptyRequestResult]
	if err := json.Unmarshal(wUnBan.Body.Bytes(), &unbanResult); err != nil {
		t.Fatalf("invalid response JSON: %v", err)
	}
	assert.Equal(t, unbanResult.Message == "User was unbanned successfully", true)

	wSuccess := UserLoginRaw(email, pwd, t, TheTestServer.GetRouter())
	//user logs in again -> Status OK
	assert.Equal(t, wSuccess.Code, http.StatusOK)
	var rrOK models.RequestResult[models.LoginResponse]
	if err := json.Unmarshal(wSuccess.Body.Bytes(), &rrOK); err != nil {
		t.Fatalf("invalid response JSON: %v", err)
	}
	assert.Equal(t, (rrOK.Value.AccessToken != "" && rrOK.Value.RefreshToken != ""), true)
}

func TestUserRefreshBannedUser(t *testing.T) {
	//login as admin
	adminLogin := AdminLogin(t)

	email := t.Name() + "@example.com"
	pwd := "thisisapass"
	//register a new user
	body := models.LoginRequest{
		Email:     email,
		Password:  pwd,
		Handle:    strings.ToLower(t.Name()),
		FirstName: t.Name(),
		LastName:  t.Name(),
	}
	w := UserRegister(&body, t, TheTestServer.GetRouter())
	assert.Equal(t, w.Code, http.StatusOK)
	var initialTokens models.RequestResult[models.LoginResponse]
	if err := json.Unmarshal(w.Body.Bytes(), &initialTokens); err != nil {
		t.Fatalf("invalid response JSON: %v", err)
	}

	u, err := TheTestServer.GetBackend().GetUserService().FindByEmail(email)
	assert.Equal(t, err == nil, true)
	assert.Equal(t, u != nil, true)

	//admin bans user
	wBan := AdminBansUser(&adminLogin.Value, u.ID, t)
	var banResult models.RequestResult[models.EmptyRequestResult]
	if err := json.Unmarshal(wBan.Body.Bytes(), &banResult); err != nil {
		t.Fatalf("invalid response JSON: %v", err)
	}
	assert.Equal(t, banResult.Message == "User was banned successfully", true)

	//user refreshes -> Invalid credentials
	wFail := UserRefresh(initialTokens.Value.RefreshToken, t, TheTestServer.GetRouter())
	assert.Equal(t, wFail.Code, http.StatusUnauthorized)
	var rr models.ErrorResult
	if err := json.Unmarshal(wFail.Body.Bytes(), &rr); err != nil {
		t.Fatalf("invalid response JSON: %v", err)
	}
	assert.Equal(t, rr.Error == "Invalid credentials", true)

	//admin unbans user
	wUnBan := AdminUnBansUser(&adminLogin.Value, u.ID, t)
	var unbanResult models.RequestResult[models.EmptyRequestResult]
	if err := json.Unmarshal(wUnBan.Body.Bytes(), &unbanResult); err != nil {
		t.Fatalf("invalid response JSON: %v", err)
	}
	assert.Equal(t, unbanResult.Message == "User was unbanned successfully", true)

	//user refreshes again in an attempt to get tokens -> OK!
	wSuccess := UserRefresh(initialTokens.Value.RefreshToken, t, TheTestServer.GetRouter())

	assert.Equal(t, wSuccess.Code, http.StatusOK)
	var rrOK models.RequestResult[models.LoginResponse]
	if err := json.Unmarshal(wSuccess.Body.Bytes(), &rrOK); err != nil {
		t.Fatalf("invalid response JSON: %v", err)
	}
	assert.Equal(t, (rrOK.Value.AccessToken != "" && rrOK.Value.RefreshToken != ""), true)
}

// admin bans user that does not exist
func TestAdminBansUserThatDoesNotExist(t *testing.T) {
	//login as admin
	adminLogin := AdminLogin(t)
	u, _ := TheTestServer.GetBackend().GetUserService().GetFirstWithRole(models.Admin.String())
	//admin bans user
	unknownUserID := utils.GetRandomUintRange(100, 10000)
	wBan := AdminBansUser(&adminLogin.Value, unknownUserID, t)
	var banResult models.ErrorResult
	if err := json.Unmarshal(wBan.Body.Bytes(), &banResult); err != nil {
		t.Fatalf("invalid response JSON: %v", err)
	}
	assert.Equal(t, banResult.Details == fmt.Sprintf("unknown user with id %d", unknownUserID), true)
	assert.Equal(t, banResult.Error == fmt.Sprintf("Admin User with id %d couldn't ban user %d", u.ID, unknownUserID), true)
}

// admin unbans user that does not exist
func TestAdminUnBansUserThatDoesNotExist(t *testing.T) {
	//login as admin
	adminLogin := AdminLogin(t)
	u, _ := TheTestServer.GetBackend().GetUserService().GetFirstWithRole(models.Admin.String())
	//admin bans user
	unknownUserID := utils.GetRandomUintRange(100, 10000)
	wBan := AdminUnBansUser(&adminLogin.Value, unknownUserID, t)
	var unbanResult models.ErrorResult
	if err := json.Unmarshal(wBan.Body.Bytes(), &unbanResult); err != nil {
		t.Fatalf("invalid response JSON: %v", err)
	}
	assert.Equal(t, unbanResult.Details == fmt.Sprintf("unknown user with id %d", unknownUserID), true)
	assert.Equal(t, unbanResult.Error == fmt.Sprintf("Admin User with id %d couldn't unban user %d", u.ID, unknownUserID), true)
}

func TestAdminBansHimself(t *testing.T) {
	//login as admin
	adminLogin := AdminLogin(t)
	u, _ := TheTestServer.GetBackend().GetUserService().GetFirstWithRole(models.Admin.String())
	//admin attempts to ban himself
	wBan := AdminBansUser(&adminLogin.Value, u.ID, t)
	var banResult models.ErrorResult
	if err := json.Unmarshal(wBan.Body.Bytes(), &banResult); err != nil {
		t.Fatalf("invalid response JSON: %v", err)
	}
	assert.Equal(t, banResult.Error == "Admin User cannot ban himself", true)
}
