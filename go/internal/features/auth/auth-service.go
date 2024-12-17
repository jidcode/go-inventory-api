package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/ventry/internal/domain/config"
	"github.com/ventry/internal/domain/models"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo   *AuthRepository
	config *config.Variables
}

func NewAuthService(userRepo *AuthRepository, cfg *config.Variables) *AuthService {
	return &AuthService{repo: userRepo, config: cfg}
}

func (service *AuthService) RegisterUser(username, email, password string) (*models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
		Role:     "user",
	}

	err = service.repo.CreateUser(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (service *AuthService) LoginUser(email, password string) (string, error) {
	user, err := service.repo.GetUserByEmail(email)
	if err != nil {
		return "", errors.New("user not found")
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// Generate JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID.String(),
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(service.config.JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (service *AuthService) ValidateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(service.config.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}

func (service *AuthService) IsTokenExpired(tokenString string) bool {
	token, err := service.ValidateToken(tokenString)
	if err != nil {
		return true
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return true
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return true
	}

	return time.Now().Unix() > int64(exp)
}

func (service *AuthService) GetUserFromToken(tokenString string) (*models.User, error) {
	token, err := service.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	userID, err := uuid.Parse(claims["user_id"].(string))
	if err != nil {
		return nil, err
	}

	return service.repo.GetUserByID(userID)
}
