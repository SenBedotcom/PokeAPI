package services

import (
	"errors"
	"os"
	"time"

	"pokemon-api/models"
	"pokemon-api/utils"

	"github.com/dgrijalva/jwt-go"
)

// AuthService struct manages authentication actions like register and login
type AuthService struct {
	UserModel *models.UserModel
}

// NewAuthService creates a new instance of AuthService
func NewAuthService(userModel *models.UserModel) *AuthService {
	return &AuthService{UserModel: userModel}
}

// Register registers a new user with a hashed password
func (s *AuthService) Register(username, password string) error {
	// Check if the user already exists
	existingUser, err := s.UserModel.FindByUsername(username)
	if err == nil && existingUser.Username != "" {
		return errors.New("username already exists")
	}

	// Hash the password
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return err
	}

	// Create and save the new user
	user := models.User{
		Username: username,
		Password: hashedPassword,
	}
	return s.UserModel.CreateUser(user)
}

// Login checks the credentials and returns a JWT if valid
func (s *AuthService) Login(username, password string) (string, error) {
	// Retrieve the user by username
	user, err := s.UserModel.FindByUsername(username)
	if err != nil {
		return "", errors.New("user not found")
	}

	// Check if the password matches the hashed password
	if !utils.CheckPasswordHash(password, user.Password) {
		return "", errors.New("invalid credentials")
	}

	// Generate a JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	})

	// Sign the token with a secret key
	secretKey := os.Getenv("JWT_SECRET_KEY")
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
