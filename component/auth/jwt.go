package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/zeromicro/go-zero/core/logx"
)

var jwtSecret = []byte("user_service_secret")

type Claims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func GenerateToken(userID int64, username string) (string, error) {
	// now := time.Now()
	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			// ExpiresAt: jwt.NewNumericDate(now.Add(24 * time.Hour)),
			// IssuedAt:  jwt.NewNumericDate(now),
			Issuer: "user-service",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ParseToken parses a JWT token and returns the claims
func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}

func Expire(c *Claims) bool {
	return time.Now().After(c.ExpiresAt.Time)
}

func RefreshToken(userID int64, username string) (string, error) {
	return GenerateToken(userID, username)
}

func ValidateToken(tokenString string) (bool, int64, error) {
	c, err := ParseToken(tokenString)
	if c != nil {
		logx.Info("token.uid: ", c.UserID)
		return true, c.UserID, nil
	}
	return false, 0, err
}
