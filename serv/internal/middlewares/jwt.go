package middlewares

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
)

// CustomClaims 自定义的JWT声明
type CustomClaims struct {
	ClientID string `json:"client_id"`
	jwt.StandardClaims
}

//var SecretKey = []byte("Lr2E52Y0XPleQi5w277b7w0bH9W7ia4wewxGszr5QH0=")

func JwtMiddleWare(secretString string) gin.HandlerFunc {
	var SecretKey = []byte(secretString)
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			c.Abort()
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return SecretKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// 在请求上下文中存储解析后的声明，以便后续处理函数使用
		claims, ok := token.Claims.(*CustomClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		c.Set("claims", claims)
	}
}
