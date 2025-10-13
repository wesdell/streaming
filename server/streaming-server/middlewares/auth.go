package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/wesdell/streaming/server/streaming-server/utils"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := utils.GetToken(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "no token provided"})
			c.Abort()
			return
		}

		claims, err := utils.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("userId", claims.UserId)
		c.Set("role", claims.Role)
		c.Next()
	}
}
