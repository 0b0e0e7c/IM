package handler

import (
	"context"
	"net/http"

	"github.com/0b0e0e7c/chat/service/user-service/pb/user"

	"google.golang.org/grpc/status"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context, client user.UserServiceClient) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := client.Register(context.Background(), &user.RegisterRequest{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "failed",
				"msg":    err.Error(),
			})
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "failed",
			"msg":    st.Message(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       resp.UserId,
		"username": resp.Username,
		"status":   "success",
		"msg":      "register success",
	})
}

func Login(c *gin.Context, client user.UserServiceClient) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := client.Login(context.Background(), &user.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "failed",
				"msg":    err.Error(),
			})
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "failed",
			"msg":    st.Message(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"username": resp.Username,
		"user_id":  resp.UserId,
		"status":   "success",
		"token":    resp.Token,
	})
}

func ValidateJWT(c *gin.Context, client user.UserServiceClient) {
	var req struct {
		Token string `json:"token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := client.ValidateJWT(context.Background(), &user.ValidateRequest{
		Token: req.Token,
	})

	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "failed",
				"msg":    err.Error(),
			})
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "failed",
			"msg":    st.Message(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"valid":   resp.Valid,
		"user_id": resp.UserId,
	})
}
