package auth

import (
	"crypto/ed25519"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/zeromicro/go-zero/core/logx"
)

const (
	TokenExpireDuration = time.Hour * 24
	Issuer              = "chat_service"
)

var (
	ErrTokenExpired = jwt.ErrTokenExpired
	publicKey       ed25519.PublicKey
	privateKey      ed25519.PrivateKey
)

type Claims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func init() {
	publicKey, privateKey, _ = ed25519.GenerateKey(nil)
}

func GenerateToken(userID int64, username string) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(TokenExpireDuration)),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
	return token.SignedString(privateKey)
}

// ParseToken parses a JWT token and returns the claims
func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodEd25519); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}

func RefreshToken(userID int64, username string) (string, error) {
	return GenerateToken(userID, username)
}

func ValidateToken(tokenString string) (valid bool, uid int64, err error) {
	c, err := ParseToken(tokenString)
	if c != nil {
		logx.Info("token.uid: ", c.UserID)
		return true, c.UserID, nil
	}
	return false, 0, err
}
