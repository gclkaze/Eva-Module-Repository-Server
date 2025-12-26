// Package middleware contains auth code
package middleware

import (
	"net/http"
	"strings"

	"github.com/gclkaze/evamodulerepositoryserver/internal/models"
	"github.com/gclkaze/evamodulerepositoryserver/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthMiddleWare struct {
	userService *services.UserService
}

func NewAuthMiddleware(userService *services.UserService) *AuthMiddleWare {
	return &AuthMiddleWare{userService: userService}
}

func (r *AuthMiddleWare) AuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")

		if !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(header, "Bearer ")

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		claims := token.Claims.(jwt.MapClaims)

		id, ok := claims["sub"].(float64)
		if !ok {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Attach user ID to context
		c.Set("userId", uint(id))

		c.Next()
	}
}

func (r *AuthMiddleWare) PreAuthorize(condition func(c *gin.Context) bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !condition(c) {
			c.AbortWithStatusJSON(403, gin.H{"error": "Forbidden"})
			return
		}
		c.Next()
	}
}

func (r *AuthMiddleWare) HasRole(role string) func(c *gin.Context) bool {
	return func(c *gin.Context) bool {
		roles, _ := c.Get("roles")
		for _, r := range roles.([]string) {
			if r == role {
				return true
			}
		}
		return false
	}
}

func (r *AuthMiddleWare) HasPermissions(requiredPermissions []models.UserPermissionTypeDef) func(c *gin.Context) bool {
	return func(c *gin.Context) bool {
		userID := c.GetUint("userId")
		perms, err := r.userService.GetUserPermissions(userID)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return false
		}

		for i := range requiredPermissions {
			found := false
			for j := range perms {
				if perms[j].Value == requiredPermissions[i].String() {
					found = true
					break
				}
			}
			if !found {
				c.AbortWithStatus(http.StatusUnauthorized)
				return false
			}
		}
		return true
	}
}

func (r *AuthMiddleWare) MaxBodySize(limit int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(
			c.Writer,
			c.Request.Body,
			limit,
		)
		c.Next()
	}
}
