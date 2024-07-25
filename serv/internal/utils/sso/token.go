package sso

import (
	"bsquared.network/b2-message-channel-serv/internal/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
	"net/http"
	"time"
)

type Wallet struct {
	WalletAddress   string `json:"wallet_address"`
	ScaAddress      string `json:"sca_address"`
	ParticleAddress string `json:"particle_address"`
}

type User struct {
	Name    string   `json:"name,omitempty"`
	Wallets []Wallet `json:"wallets,omitempty"`
}

type UserClaims struct {
	User           User   `json:"user,omitempty"`
	Sub            string `json:"sub"`
	IsRefreshToken uint8  `json:"is_refresh_token,omitempty"`
	jwt.RegisteredClaims
}

type JwtToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func GenerateToken(secretKey string, user User, sub string, isRefreshToken uint8) (string, error) {
	expiredAt := time.Now().Add(7 * 24 * time.Hour)
	if isRefreshToken == 1 {
		expiredAt = time.Now().Add(14 * 24 * time.Hour)
	}
	numericExpiredAt := jwt.NewNumericDate(expiredAt)
	var userClaims UserClaims
	if isRefreshToken == 0 {
		userClaims = UserClaims{
			User: user,
			Sub:  sub,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: numericExpiredAt,
				Issuer:    "b2network",
			},
		}
	} else {
		userClaims = UserClaims{
			Sub:            sub,
			IsRefreshToken: isRefreshToken,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: numericExpiredAt,
				Issuer:    "b2network",
			},
		}
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaims)
	tokenString, err := t.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, err
}

func VerifyToken(secretKey string, tokenString string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}
	if userClaims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		return userClaims, nil
	} else {
		return nil, errors.New("token invald")
	}
}

func ValidateToken(c *gin.Context) string {
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Claims not found in context"})
		return ""
	}

	clientID := claims.(*middlewares.CustomClaims).ClientID
	return clientID
}
