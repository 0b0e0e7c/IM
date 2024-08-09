package main

import (
	"context"
	"net/http"

	"github.com/0b0e0e7c/IM/user-service/cmd/rpc/pb/user"
	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context, client user.UserClient) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := client.Register(context.Background(), &user.UserRequest{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       resp.Id,
		"username": resp.Username,
	})
}

func Login(c *gin.Context, client user.UserClient) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := client.Login(context.Background(), &user.UserRequest{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       resp.Id,
		"username": resp.Username,
	})
}
