package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/0b0e0e7c/chat/service/user-service/pb/user"

	"github.com/gin-gonic/gin"
)

func JWTMiddleware(userRPCClient user.UserServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Bearer token is required"})
			c.Abort()
			return
		}

		// 调用 RPC 方法
		resp, err := userRPCClient.ValidateJWT(context.Background(), &user.ValidateRequest{Token: token})
		if err != nil || !resp.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// 将用户信息存储在上下文中
		c.Set("userID", resp.UserId)
		c.Next()
	}
}
