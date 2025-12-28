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
	"fmt"
	"net/http"

	"github.com/gclkaze/evamodulerepositoryserver/internal/services"
	"github.com/gclkaze/evamodulerepositoryserver/pkg/utils"
	"github.com/gin-gonic/gin"
)

type SuperviseHandler struct {
	service *services.UserService
}

func NewSuperviseHandler(service *services.UserService) *SuperviseHandler {
	return &SuperviseHandler{service: service}
}

func (u *SuperviseHandler) BanUser(c *gin.Context) {
	initiator := c.GetUint("userId")
	userID := c.Param("userId")
	userIDUint, err := utils.StringToUint(userID)

	if err != nil {
		c.JSON(http.StatusBadRequest, utils.Err(err, "Invalid User ID format"))
		return
	}

	if initiator == userIDUint {
		c.JSON(http.StatusInternalServerError, utils.ErrWithSimpleMessage("Admin User cannot ban himself"))
		return
	}
	err = u.service.BanUser(userIDUint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Err(err, fmt.Sprintf("Admin User with id %d couldn't ban user %d", initiator, userIDUint)))
		return
	}
	c.JSON(http.StatusOK, utils.SimpleOkMessage("User was banned successfully"))
}

func (u *SuperviseHandler) UnbanUser(c *gin.Context) {
	initiator := c.GetUint("userId")
	userID := c.Param("userId")
	userIDUint, err := utils.StringToUint(userID)

	if err != nil {
		c.JSON(http.StatusBadRequest, utils.Err(err, "Invalid User ID format"))
		return
	}

	err = u.service.UnbanUser(userIDUint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Err(err, fmt.Sprintf("Admin User with id %d couldn't unban user %d", initiator, userIDUint)))
		return
	}
	c.JSON(http.StatusOK, utils.SimpleOkMessage("User was unbanned successfully"))
}
