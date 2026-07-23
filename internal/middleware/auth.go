package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const APIKeyHeader = "X-API-Key"

func AuthMiddleware(expectedKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.GetHeader(APIKeyHeader)
		if key != expectedKey || expectedKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		c.Next()
	}
}
