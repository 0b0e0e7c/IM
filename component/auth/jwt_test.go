package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestGenerateToken(t *testing.T) {
	tokenString, err := GenerateToken(1, "testuser")
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)
}

func TestParseToken_Expired(t *testing.T) {
	// 创建一个过期的 token
	now := time.Now()
	claims := Claims{
		UserID:   1,
		Username: "testuser",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(-1 * time.Hour)), // 设置为1小时前过期
			IssuedAt:  jwt.NewNumericDate(now.Add(-2 * time.Hour)), // 设置为2小时前签发
			Issuer:    "IM_service",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	assert.NoError(t, err)

	// 解析 token
	parsedClaims, err := ParseToken(tokenString)
	assert.Error(t, err)
	assert.Nil(t, parsedClaims)
	assert.ErrorIs(t, err, ErrInvalidToken)
}

func TestParseToken_Valid(t *testing.T) {
	tokenString, err := GenerateToken(1, "testuser")
	assert.NoError(t, err)

	parsedClaims, err := ParseToken(tokenString)
	assert.NoError(t, err)
	assert.NotNil(t, parsedClaims)
	assert.Equal(t, int64(1), parsedClaims.UserID)
	assert.Equal(t, "testuser", parsedClaims.Username)
}

func TestErrEqual(t *testing.T) {
	assert.Equal(t, ErrInvalidToken, jwt.ErrTokenExpired)
}
