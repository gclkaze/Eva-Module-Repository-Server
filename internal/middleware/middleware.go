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
