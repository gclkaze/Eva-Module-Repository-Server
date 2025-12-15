package services

import (
	"errors"
	"time"

	"github.com/gclkaze/evamodulerepositoryserver/internal/models"
	"github.com/gclkaze/evamodulerepositoryserver/internal/repositories"
	"github.com/gclkaze/evamodulerepositoryserver/pkg/logger"
	"github.com/gclkaze/evamodulerepositoryserver/pkg/runtime"
	"github.com/golang-jwt/jwt/v5"
	"github.com/magiconair/properties"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	p         *properties.Properties
	logger    logger.ILogger
	jwtSecret string

	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration

	usersRepo repositories.UserAccountRepository
	devRepo   repositories.DeveloperRepository
}

func NewAuthService(usersRepo repositories.UserAccountRepository,
	devRepo repositories.DeveloperRepository, p *properties.Properties) *AuthService {
	auth := &AuthService{p: p}
	auth.jwtSecret = p.GetString("jwt_secret", "")
	auth.logger = runtime.CreateLogger(p)
	auth.accessTokenTTL = p.GetDuration("access_token_ttl", 15) * time.Hour
	auth.refreshTokenTTL = p.GetDuration("refresh_token_ttl", 7*24) * time.Hour
	auth.usersRepo = usersRepo
	auth.devRepo = devRepo
	return auth
}

func (s AuthService) GetJWTSecret() string {
	return s.jwtSecret
}

func (s *AuthService) GenerateTokens(user *models.UserAccount) (string, string, error) {

	// ACCESS TOKEN
	accessClaims := jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"exp":   time.Now().Add(s.accessTokenTTL).Unix(),
	}

	access := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err := access.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", "", err
	}

	// REFRESH TOKEN
	refreshClaims := jwt.MapClaims{
		"sub":  user.ID,
		"type": "refresh",
		"exp":  time.Now().Add(s.refreshTokenTTL).Unix(),
	}

	refresh := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err := refresh.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *AuthService) Authenticate(email, password string) (*models.UserAccount, error) {
	user, err := s.usersRepo.FindByEmail(email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if user == nil {
		return nil, nil
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	return user, nil
}

func (s *AuthService) ValidateRefreshToken(tokenString string) (*models.UserAccount, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid refresh token")
	}

	claims := token.Claims.(jwt.MapClaims)

	// must be a refresh token
	if claims["type"] != "refresh" {
		return nil, errors.New("invalid token type")
	}

	userID, ok := claims["sub"].(float64)
	if !ok {
		return nil, errors.New("invalid token payload")
	}

	// Find user in DB
	user, err := s.usersRepo.FindByID(uint(userID))
	if err != nil {
		return nil, errors.New("user no longer exists")
	}

	return user, nil
}
