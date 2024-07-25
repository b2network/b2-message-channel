package utils

import (
	"bsquared.network/b2-message-channel-serv/internal/middlewares"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ValidateToken(c *gin.Context) string {
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Claims not found in context"})
		return ""
	}

	clientID := claims.(*middlewares.CustomClaims).ClientID
	return clientID
}
