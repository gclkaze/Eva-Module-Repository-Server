package handlers

import (
	"net/http"

	"github.com/gclkaze/evamodulerepositoryserver/internal/models"
	"github.com/gclkaze/evamodulerepositoryserver/internal/services"
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
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid payload",
			"result":  false,
			"details": err.Error(),
		})
		return
	}

	var user *models.UserAccount
	user, err := h.service.Authenticate(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Invalid crendentials",
			"result":  false,
			"details": err.Error()})
		return
	}

	//need to create the user
	if user == nil {
		userID, createUserErr := h.userService.CreateUser(req.Handle, req.FirstName, req.LastName, req.Email, req.Password, true)
		if createUserErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Couldn't create user",
				"result":  false,
				"details": createUserErr.Error()})
			return
		}
		userResult, findUserErr := h.userService.FindUserByID(userID)
		if findUserErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Couldn't find user ID",
				"details": findUserErr.Error(),
				"result":  false,
			})
			return
		}
		user = userResult
	}

	access, refresh, err := h.service.GenerateTokens(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"details": err.Error(),
			"error":   "token creation failed",
			"result":  false,
		})
		return
	}

	c.JSON(http.StatusOK, models.LoginResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid payload",
			"result":  false,
			"details": err.Error(),
		})
		return
	}

	user, err := h.service.Authenticate(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Invalid crendentials",
			"result":  false,
			"details": err.Error()})
		return
	}

	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"details": "unregistered user",
			"result":  false,
			"error":   "you need to register first!"})
		return
	}
	access, refresh, err := h.service.GenerateTokens(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"details": err.Error(),
			"result":  false,
			"error":   "token creation failed"})
		return
	}

	c.JSON(http.StatusOK, models.LoginResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	})
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var req models.RefreshRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid payload",
			"result":  false,
			"details": err.Error(),
		})
		return
	}

	user, err := h.service.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"details": err.Error(),
			"error":   "invalid refresh token",
			"result":  false})
		return
	}

	access, refresh, err := h.service.GenerateTokens(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"details": err.Error(),
			"result":  false,
			"error":   "failed to generate tokens"})
		return
	}

	c.JSON(http.StatusOK, models.LoginResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	})
}
