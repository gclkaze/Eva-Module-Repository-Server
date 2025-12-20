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
	initiator := c.GetUint("userID")
	userID := c.Param("userId")
	userIDUint, err := utils.StringToUint(userID)

	if err != nil {
		c.JSON(http.StatusBadRequest, utils.Err(err, "Invalid User ID format"))
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
	initiator := c.GetUint("userID")
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
