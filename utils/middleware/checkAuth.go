package middleware

import (
	"net/http"
	"simple-api/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func CheckAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("Authorization") == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
		}

		tokenString := strings.Split(c.GetHeader("Authorization"), " ")

		if len(tokenString) != 2 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
		}

		if tokenString[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
		}

		if tokenString == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
		}

		claims, err := utils.AccessParams{}.VerifyToken(tokenString[1])

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
		} else {
			c.Set("user", claims)

		}

		c.Next()
	}
}
