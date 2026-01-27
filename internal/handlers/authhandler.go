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

package handlers

import (
	"net/http"
	"strings"

	"github.com/gclkaze/evamodulerepositoryserver/internal/models"
	"github.com/gclkaze/evamodulerepositoryserver/internal/services"
	"github.com/gclkaze/evamodulerepositoryserver/pkg/utils"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service     *services.AuthService
	userService *services.UserService
}

func NewAuthHandler(service *services.AuthService, userService *services.UserService) *AuthHandler {
	return &AuthHandler{service: service, userService: userService}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.Err(err, "Invalid payload"))
		return
	}

	pwd := strings.TrimSpace(req.Password)
	if pwd == "" {
		c.JSON(http.StatusBadRequest, utils.ErrWithSimpleMessage("No password provided"))
		return
	}

	email := strings.TrimSpace(req.Email)
	if email == "" {
		c.JSON(http.StatusBadRequest, utils.ErrWithSimpleMessage("No email provided"))
		return
	}

	if !utils.IsValidEmail(email) {
		c.JSON(http.StatusBadRequest, utils.ErrWithSimpleMessage("Invalid email form"))
		return
	}

	errPwd := utils.ValidatePassword(pwd, email)
	if errPwd != nil {
		c.JSON(http.StatusBadRequest, utils.ErrWithSimpleMessage(errPwd.Error()))
		return
	}

	firstName := strings.TrimSpace(req.FirstName)
	if firstName == "" {
		c.JSON(http.StatusBadRequest, utils.ErrWithSimpleMessage("No first name provided"))
		return
	}

	if !utils.IsValidName(firstName) {
		c.JSON(http.StatusBadRequest, utils.ErrWithSimpleMessage("Invalid first name provided"))
		return
	}

	lastName := strings.TrimSpace(req.LastName)
	if lastName == "" {
		c.JSON(http.StatusBadRequest, utils.ErrWithSimpleMessage("No last name provided"))
		return
	}

	if !utils.IsValidName(lastName) {
		c.JSON(http.StatusBadRequest, utils.ErrWithSimpleMessage("Invalid last name provided"))
		return
	}

	handle := strings.TrimSpace(req.Handle)
	if handle == "" {
		c.JSON(http.StatusBadRequest, utils.ErrWithSimpleMessage("No handle provided"))
		return
	}

	if !utils.IsValidHandle(handle) {
		c.JSON(http.StatusBadRequest, utils.ErrWithSimpleMessage("Invalid handle provided"))
		return
	}

	var user *models.UserAccount
	user, err := h.service.Authenticate(email, pwd)
	if err != nil {
		c.JSON(http.StatusUnauthorized, utils.Err(err, "Invalid credentials"))
		return
	}

	if user != nil && user.IsBanned {
		c.JSON(http.StatusUnauthorized, utils.Err(err, "Invalid credentials"))
		return
	}

	//need to create the user
	if user == nil {
		userID, createUserErr := h.userService.CreateUser(handle, firstName, lastName, email, pwd, true)
		if createUserErr != nil {
			c.JSON(http.StatusInternalServerError, utils.ErrWithSimpleMessage("Couldn't register user, "+createUserErr.Error()))
			return
		}
		userResult, findUserErr := h.userService.FindUserByID(userID)
		if findUserErr != nil {
			c.JSON(http.StatusInternalServerError, utils.ErrWithSimpleMessage("Couldn't register user, "+findUserErr.Error()))
			return
		}
		user = userResult
	} else {
		c.JSON(http.StatusInternalServerError, utils.ErrWithSimpleMessage("Couldn't register user, the email is used"))
		return
	}

	access, refresh, err := h.service.GenerateTokens(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Err(err, "token creation failed"))
		return
	}

	c.JSON(http.StatusOK, utils.OkWithMessage(models.LoginResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	}, "Registration was successfull"))
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.Err(err, "Invalid payload"))
		return
	}

	pwd := strings.TrimSpace(req.Password)
	if pwd == "" {
		c.JSON(http.StatusBadRequest, utils.ErrWithSimpleMessage("No password provided"))
		return
	}

	email := strings.TrimSpace(req.Email)
	if email == "" {
		c.JSON(http.StatusBadRequest, utils.ErrWithSimpleMessage("No email provided"))
		return
	}

	if !utils.IsValidEmail(email) {
		c.JSON(http.StatusBadRequest, utils.ErrWithSimpleMessage("Invalid email form"))
		return
	}

	err := utils.IsValidPassword(pwd)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.Err(err, "Invalid password provided"))
		return
	}

	user, err := h.service.Authenticate(email, pwd)
	if err != nil {
		c.JSON(http.StatusUnauthorized, utils.Err(err, "Invalid credentials"))
		return
	}
	if user != nil && user.IsBanned {
		c.JSON(http.StatusUnauthorized, utils.Err(err, "Invalid credentials"))
		return
	}

	if user == nil {
		c.JSON(http.StatusUnauthorized, utils.Err(nil, "unregistered user"))
		return
	}
	access, refresh, err := h.service.GenerateTokens(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Err(err, "token creation failed"))
		return
	}

	c.JSON(http.StatusOK, utils.OkWithMessage(models.LoginResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	}, "Login was successfull"))
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var req models.RefreshRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.Err(err, "Invalid payload"))
		return
	}

	user, err := h.service.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, utils.Err(err, "invalid refresh token"))
		return
	}
	if user != nil && user.IsBanned {
		c.JSON(http.StatusUnauthorized, utils.Err(err, "Invalid credentials"))
		return
	}

	if user == nil {
		c.JSON(http.StatusUnauthorized, utils.Err(err, "Invalid credentials"))
		return
	}
	access, refresh, err := h.service.GenerateTokens(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Err(err, "failed to generate tokens"))
		return
	}

	c.JSON(http.StatusOK, utils.OkWithMessage(models.LoginResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	}, "Refresh was successfull"))
}
