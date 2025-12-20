package handlers

import (
	"net/http"

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

	var user *models.UserAccount
	user, err := h.service.Authenticate(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, utils.Err(err, "Invalid crendentials"))
		return
	}

	if user != nil && user.IsBanned {
		c.JSON(http.StatusUnauthorized, utils.Err(err, "Invalid crendentials"))
		return
	}

	//need to create the user
	if user == nil {
		userID, createUserErr := h.userService.CreateUser(req.Handle, req.FirstName, req.LastName, req.Email, req.Password, true)
		if createUserErr != nil {
			c.JSON(http.StatusInternalServerError, utils.Err(err, "Couldn't create user"))
			return
		}
		userResult, findUserErr := h.userService.FindUserByID(userID)
		if findUserErr != nil {
			c.JSON(http.StatusInternalServerError, utils.Err(findUserErr, "Couldn't find user ID"))
			return
		}
		user = userResult
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

	user, err := h.service.Authenticate(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, utils.Err(err, "Invalid crendentials"))
		return
	}
	if user != nil && user.IsBanned {
		c.JSON(http.StatusUnauthorized, utils.Err(err, "Invalid crendentials"))
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
		c.JSON(http.StatusUnauthorized, utils.Err(err, "Invalid crendentials"))
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
