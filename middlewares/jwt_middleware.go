package middlewares

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"quiz-project-book-api/config"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		// **Periksa blacklist token**
		var tokenID string
		err := config.DB.QueryRow("SELECT token FROM blacklisted_tokens WHERE token = $1", tokenString).Scan(&tokenID)
		if err == nil {
			// Token ditemukan di blacklist, tolak akses
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is blacklisted"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				// return nil, jwt.NewValidationError("Invalid signing method", jwt.ValidationErrorSignatureInvalid)
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				// return nil, fmt.Errorf("invalid signing method")
			}
			return []byte(os.Getenv("JWT_SECRET_KEY")), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}
		// Menyimpan username dari token ke dalam Gin Context
		c.Set("username", claims["username"])

		c.Next()
	}
}
