package controllers

import (
	"encoding/json"
	"net/http"
	"pokemon-api/models"
	"pokemon-api/services"
	"pokemon-api/utils"
)

type AuthController struct {
	AuthService services.AuthService
}

// NewAuthController is a constructor for AuthController
func NewAuthController(authService services.AuthService) *AuthController {
	return &AuthController{AuthService: authService}
}

// Register handles user registration
func (ac *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	var user models.User

	// Decode the request body into the UserModel
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Hash the user's password
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	// Set the hashed password to the UserModel
	user.Password = hashedPassword

	// Call the AuthService to create a new user
	err = ac.AuthService.Register(user.Username, user.Password)
	if err != nil {
		http.Error(w, "Error registering user", http.StatusInternalServerError)
		return
	}

	// Respond with a success message
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}

// Login handles user login
func (ac *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	var user models.User

	// Decode the request body into the UserModel
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate the user credentials by comparing the hashed password
	existingUser, err := ac.AuthService.Login(user.Username, user.Password)
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Generate a JWT token for the logged-in user
	token, err := utils.GenerateJWT(existingUser)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	// Respond with the JWT token
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
