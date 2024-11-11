package middleware

import (
	"fmt"
	"net/http"
	"pokemon-api/config"
	"pokemon-api/utils"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func RequireAuth(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" || !utils.VerifyJWT(token) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}
	c.Next()
}

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the token from the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			c.Abort() // Prevent further processing
			return
		}

		// Extract the token from the "Bearer {token}" format
		tokenString := strings.Split(authHeader, " ")[1]
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token missing"})
			c.Abort()
			return
		}

		// Parse the JWT token using your custom parsing function
		token, err := parseToken(tokenString)
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// If the token is valid, proceed to the next handler
		c.Next()
	}
}

func WrapGinHandlerToHTTP(ginHandler gin.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Create a Gin context from the HTTP request and response
		c, _ := gin.CreateTestContext(w)
		c.Request = r
		ginHandler(c)
	}
}

// parseToken parses and validates the JWT token
func parseToken(tokenString string) (*jwt.Token, error) {
	// You can replace this with your secret or a more complex JWT setup
	secretKey := config.AppConfig.JwtSecret // Assuming you store this in the config

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the token method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(secretKey), nil
	})

	return token, err
}
