package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/yourusername/book-management-api/models"
	"github.com/yourusername/book-management-api/repository"
)

type Claims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func AuthMiddleware(userRepo *repository.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Status:  http.StatusUnauthorized,
				Message: "Authorization header required",
			})
			c.Abort()
			return
		}

		// Check if it's Bearer token (JWT)
		if strings.HasPrefix(authHeader, "Bearer ") {
			token := strings.TrimPrefix(authHeader, "Bearer ")
			userID, username, err := VerifyToken(token)
			if err != nil {
				c.JSON(http.StatusUnauthorized, models.ErrorResponse{
					Status:  http.StatusUnauthorized,
					Message: "Invalid token",
				})
				c.Abort()
				return
			}
			c.Set("user_id", userID)
			c.Set("username", username)
			c.Next()
			return
		}

		// Check if it's Basic Auth
		if strings.HasPrefix(authHeader, "Basic ") {
			username, password, ok := c.Request.BasicAuth()
			if !ok {
				c.JSON(http.StatusUnauthorized, models.ErrorResponse{
					Status:  http.StatusUnauthorized,
					Message: "Invalid basic auth",
				})
				c.Abort()
				return
			}

			user, err := userRepo.VerifyPassword(username, password)
			if err != nil {
				c.JSON(http.StatusUnauthorized, models.ErrorResponse{
					Status:  http.StatusUnauthorized,
					Message: "Invalid credentials",
				})
				c.Abort()
				return
			}

			c.Set("user_id", user.ID)
			c.Set("username", user.Username)
			c.Next()
			return
		}

		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Status:  http.StatusUnauthorized,
			Message: "Invalid authorization format",
		})
		c.Abort()
	}
}

func GenerateToken(userID int, username string) (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-secret-key-change-in-production"
	}

	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(tokenString string) (int, string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-secret-key-change-in-production"
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return 0, "", err
	}

	if !token.Valid {
		return 0, "", fmt.Errorf("invalid token")
	}

	return claims.UserID, claims.Username, nil
}
