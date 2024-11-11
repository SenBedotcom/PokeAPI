// main.go
package main

import (
	"log"
	"net/http"
	"pokemon-api/config"
	"pokemon-api/controllers"
	"pokemon-api/middleware"
	"pokemon-api/models"
	"pokemon-api/services"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/mux"
)

func main() {
	// Load configuration
	config.LoadConfig()

	// Initialize Cache Service
	cacheService := services.NewCacheService()

	// Initialize Pokémon Service
	pokemonService := services.NewPokemonService(cacheService)

	// Initialize Auth Service with an empty UserModel
	userModel := &models.UserModel{}
	authService := services.NewAuthService(userModel)

	// Initialize Controllers
	authController := controllers.NewAuthController(*authService)
	pokemonController := controllers.NewPokemonController(*pokemonService, *cacheService)

	// Create a new router
	router := mux.NewRouter()

	// Define routes for Auth
	router.HandleFunc("/register", authController.Register).Methods("POST")
	router.HandleFunc("/login", authController.Login).Methods("POST")

	// Protected routes with JWT middleware (using Gin's context)
	router.HandleFunc("/pokemon/{name}", func(w http.ResponseWriter, r *http.Request) {
		// Create a Gin context from the HTTP request and response
		c, _ := gin.CreateTestContext(w)
		c.Request = r
		middleware.JWTAuth()(c)               // Applying JWT middleware here
		pokemonController.GetPokemonByName(c) // Your handler
	}).Methods("GET")

	router.HandleFunc("/pokemon/{name}/ability", func(w http.ResponseWriter, r *http.Request) {
		// Create a Gin context from the HTTP request and response
		c, _ := gin.CreateTestContext(w)
		c.Request = r
		middleware.JWTAuth()(c)                // Applying JWT middleware here
		pokemonController.GetPokemonAbility(c) // Your handler
	}).Methods("GET")

	router.HandleFunc("/pokemon/random", func(w http.ResponseWriter, r *http.Request) {
		// Create a Gin context from the HTTP request and response
		c, _ := gin.CreateTestContext(w)
		c.Request = r
		middleware.JWTAuth()(c)               // Applying JWT middleware here
		pokemonController.GetRandomPokemon(c) // Your handler
	}).Methods("GET")

	// Start the server
	log.Printf("Server starting on port 8080...")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
